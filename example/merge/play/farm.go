package play

import (
	"github.com/hiank/think/run"
	"golang.org/x/exp/slices"
)

const (
	ErrBlockNotEmpty = run.Err("play: block must empty")
)

type Position struct {
	X, Y int
}

func (pos Position) Equal(p2 Position) bool {
	return pos.X == p2.X && pos.Y == p2.Y
}

type farm struct {
	cache [][]*block           //[y][x]*block
	ltops [][]Position         //[layer][]Position
	free  []*block             //free blocks (without item, but useable)
	mres  map[uint64]*resource //map[res key]*resource
	// mfruit     map[int32]*resource    //map[item id]*resource; 保存果实数据
	// mres       map[ItemType]*resource //
	iset       Itemset
	unlockcode Unlockcode //地块解锁信息
}

// NewFarm
// pss [y][x]*Preset 地图配置信息. 数值为nil，表明块不可用. (*Preset).ItemId == -1，表明初始无物品
// NOTE: pss 与用户特有数据无关，完全配置于地图. 初始化时，会根据地块开放信息填充初始数据，已解锁的块中数据以dset中相关数据为准
func NewFarm(iset Itemset, dset Dataset, pss [][]*Preset) Farm {
	mst := make(map[uint32]State)
	for _, st := range dset.GetStates() {
		mst[st.GetSite()] = st
	}
	cache, ltops, useableCnt := unmarshalPresets(pss, iset)
	for _, dist := range dset.GetDists() {
		placeItem(iset.Get(dist.GetId()), dist.GetDistcodes(), mst, cache, ltops)
	}
	//
	f := &farm{
		cache: cache,
		ltops: ltops,
		free:  make([]*block, 0, useableCnt),
		iset:  iset,
		// mst:        mst, //make(map[uint32]State),
		mres:       make(map[uint64]*resource),
		unlockcode: Unlockcode(dset.GetUnlockcode()),
	}
	///init free
	///init mres

	return f
}

// func (f *farm) pre

func (f *farm) DoMove(from, end uint32) (suc bool) {
	var bf, be *block
	var found bool
	execute(
		func() bool { bf, found = f.siteToBlock(Sitecode(from)); return found && !bf.Empty() },
		func() bool { return !bf.Filler.Able(EBunmoveable) },
		func() bool { be, found = f.siteToBlock(Sitecode(end)); return found },
		func() bool { suc = f.merge(bf, be); return !suc },
		func() bool { suc = f.join(bf, be); return !suc },
		func() bool { suc = f.move(bf, be); return !suc },
	)
	return
}

func (f *farm) DoClick(end uint32) (suc bool) {
	// var pe Position
	var be *block
	var found bool
	execute(
		// func() bool { pe, suc = f.position(Sitecode(end)); return suc },
		// func() bool { be = f.cache[pe.X][pe.Y]; return true },
		func() bool { be, found = f.siteToBlock(Sitecode(end)); return found && !be.Empty() },
		func() bool { suc = f.pick(be); return !suc },  //尝试拾取
		func() bool { suc = f.build(be); return !suc }, //尝试修建
		func() bool { suc = f.mine(be); return !suc },  //尝试开采
		func() bool { suc = f.open(be); return !suc },  //尝试打开
		func() bool { suc = f.drop(be); return !suc },  //尝试掉落(开启后可能无足够地块全部放置，清理后点击继续放置)
		func() bool { suc = f.fly(be); return !suc },   //尝试飞到目标物品(完全体材料飞到对应城堡)
		func() bool { suc = f.cook(be); return !suc },  //尝试制作甜品
	)
	return
}

func (f *farm) DoEradicate(end uint32) (suc bool) {
	if b, found := f.siteToBlock(Sitecode(end)); found && !b.Empty() {
		suc = f.eradicate(b)
	}
	return
}

func (f *farm) position(sc Sitecode) (pos Position, loaded bool) {
	if layer := int(sc.Layer()); layer < len(f.ltops) {
		if idx, ps := sc.Index(), f.ltops[layer]; idx < len(ps) {
			pos, loaded = ps[idx], true
		}
	}
	return
}

func (f *farm) siteToBlock(sc Sitecode) (b *block, found bool) {
	pos, found := f.position(sc)
	if found {
		b = f.cache[pos.X][pos.Y]
	}
	return
}

// func (f *farm) fillerFromSite(sc Sitecode)

// move 移动
func (f *farm) move(from, to *block) bool {
	return true
}

// merge 合并(完全体水晶或者相同的物品)
func (f *farm) merge(s, d *block) (suc bool) {
	execute(
		func() bool { return !s.Empty() && !d.Empty() },    //地块不能为空
		func() bool { return !(s.X == d.X && s.Y == d.Y) }, //源位置和目标位置不能时同一个
		func() bool { return !d.Filler.Able(EBfinal) },     //目标非完全体
		func() bool { suc = s.Filler.GetItem().GetId() == d.Filler.GetItem().GetId(); return !suc },
		func() bool { suc = s.Filler.Able(EBuniversal); return !suc }, //源位置是一个万能合并物品
	)
	if suc {
		layer, itemId := int(d.Site.Layer()), int(d.Filler.GetItem().GetId())
		arr, cnt := filterLine(f.cache[d.Y], layer, itemId, encodeLinecode(Linecode(d.X), 1)) //, make(map[int][]Linecode)
		totalCnt, cache := cnt+1, make(map[int][]Linecode)
		for y := d.Y + 1; y < len(f.cache) && len(arr) > 0; y++ {
			arr, cnt = filterLines(f.lineX(y, s.Position), layer, itemId, arr...)
			totalCnt, cache[y] = totalCnt+cnt, arr
		}
		arr = cache[d.Y]
		for y := d.Y - 1; y >= 0 && len(arr) > 0; y-- {
			arr, cnt = filterLines(f.lineX(y, s.Position), layer, itemId, arr...)
			totalCnt, cache[y] = totalCnt+cnt, arr
		}
		if suc = totalCnt > 2; suc {
			remain, oldItemId, itemId := totalCnt%5, d.Filler.GetItem().GetId(), d.Filler.GetItem().GetNextItemId()
			ncnt, remain := totalCnt/5+remain/3, remain%3
			///delete all old item in blocks
			for y, arr := range cache {
				f.removeLines(f.cache[y], arr...)
			}
			f.clear(s) //clear source block
			///make new item blocks
			bs := nearestEmptyBlocks(f.free, d.Position, ncnt+remain)
			f.generate(itemId, bs[:ncnt]...) //generate merged item
			f.generate(int(oldItemId), bs[ncnt:]...)
			///暂时不考虑奖励，后续可能会考虑合并奖励掉落
		}
	}
	return
}

func (f *farm) lineX(y int, ignore Position) (line []*block) {
	line = f.cache[y]
	if y == ignore.Y && ignore.X < len(line) {
		line = slices.Clone(line)
		line[ignore.X] = nil //
	}
	return
}

func (f *farm) removeLines(line []*block, lcs ...Linecode) {
	for _, lc := range lcs {
		for idx, max := lc.Min(), lc.Min()+lc.Count(); idx < max; idx++ {
			f.clear(line[idx])
		}
	}
}

// func (f *farm)

// join 升级(完全体通用升级材料或者对应的建筑)或工人修建
func (f *farm) join(s, d *block) bool {
	return false
}

// pick 点击响应，拾取(入背包)
func (f *farm) pick(b *block) (suc bool) {
	if suc = b.Filler.Able(EBpickable); suc {

	}
	return
}

// build 点击响应，修建
func (f *farm) build(b *block) (suc bool) {
	// b.Item.GetStateBitag().
	if suc = b.Filler.Able(EBwaitbuild); suc {
		///
		// f.
	}
	return
}

// mine 点击响应，开采
func (f *farm) mine(b *block) (suc bool) {
	if suc = b.Filler.Able(EBwaitmine); suc {

	}
	return
}

// open 点击响应，开启(宝箱袋子等未开启状态下)(先生成奖励，然后执行掉落'drop')
func (f *farm) open(b *block) (suc bool) {
	if suc = b.Filler.Able(EBopenable); suc {
		//open here
	}
	return
}

// drop 点击响应，掉落(箱子袋子等开启后，城堡奖励cd完成，开采cd完成，农作物cd完成 等)
func (f *farm) drop(b *block) (suc bool) {

	return
}

// fly 点击响应，已知飞到城堡中(完全体对应建筑)，潜艇活动中飞到潜艇
func (f *farm) fly(b *block) bool {
	return false
}

// cook 点击响应，制作食物
func (f *farm) cook(b *block) bool {
	return false
}

// eradicate 铲除指定块中的物品(包括清理状态数据). 铲子的专属处理
func (f *farm) eradicate(b *block) (suc bool) {
	if suc = b.Filler.Able(EBeradicable); suc {
		f.clear(b)
	}
	return
}

// generate new item. the block must be empty
func (f *farm) generate(itemId int, bs ...*block) {
	if slices.IndexFunc(bs, func(b *block) (disable bool) {
		return b == nil || b.Filler != nil || f.unlockcode.Unlocked(b.PID)
	}) != -1 {
		panic(ErrBlockNotEmpty)
	}
	// it := f.iset.Get(int32(itemId))
	///
	// for _, b := range bs {
	///
	// b.Item = convertoDataItem(it, initialState(it))
	// if i := slices.Index(f.free, b); i != -1 {
	// 	f.free = slices.Delete(f.free, i, i+1)
	// }
	// }
}

// clear 清理块中的物品
func (f *farm) clear(b *block) {
	// b.item = nil
	// b.Clear()
	//clear state here
	// slices.Index()
	// slices.Clip()
	f.free = append(f.free, b) //add to free array
}

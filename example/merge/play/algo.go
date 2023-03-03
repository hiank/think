package play

import (
	"math"

	"github.com/hiank/think/exp/easy"
	"github.com/hiank/think/run"
	"golang.org/x/exp/slices"
)

const (
	ErrUnsupportedItemType = run.Err("unsupported item type")

	resourceKeyCover uint64 = (uint64(math.MaxUint32) << 32) & uint64(math.MaxUint8)
)

func encodeStateBitag(sb StateBitag, v uint, sbc StateBitagCover) StateBitag {
	// easy.ResetBit[]()
	return easy.ResetBit(sb, StateBitag(v), uint(sbc.Low()), uint(sbc.Count()))
}

func convertoResourceKey(rc Resourcecode) uint64 {
	return resourceKeyCover & uint64(rc)
}

func farmToDataset(f *farm) Dataset {
	return &farmDataset{f: f}
}

// convertoDataState convert proto State to *state
func convertoDataState(st State) *state {
	return &state{
		site:   st.GetSite(),
		Bitag:  StateBitag(st.GetBitag()),
		cutime: st.GetCutime(),
		ex:     st.GetEx(),
		awards: st.GetAwards(),
	}
}

// func convertoDataResource(res Resource) *resource {
// 	return nil
// }

// func convertoDataBackpack(bp Backpack) *backpack {
// 	return nil
// }

// func withState(st *state)

func convertoFiller(it Item, st State) *filler {
	var eb, final EasyBitag = EasyBitag(it.GetType()), EBfinal
	switch it.GetType() {
	case ITcastle: //城堡
		eb = EBhold4block | EBkeystate | EBkeybuilding | EBkeymining
	case ITrootmine: //母矿
		eb = EBhold4block | EBkeystate | EBkeybuilding | EBkeymining
	case ITtimeworker: //计时工人: 完全体可修建
		final = EBkeybuilding | EBkeystate
	case ITcrop: //农作物
		final = EBkeymining | EBkeystate
	case ITcastlematerial: //建筑(用于升级城堡)
		if it.GetBuiltDuration() > 0 { //需要修建
			eb = EBkeybuilding | EBkeystate
		}
	case ITmine: //矿
		final = EBkeymining | EBkeystate
	case ITlocalmine: //固定位矿(未解锁时填充矿)
		eb = EBsurround | EBkeymining
	case ITcrystal: //水晶
	case ITcook: //厨子
		final = EBkeystate
	case ITcoin: //金币
		eb = EBpickable
	case ITdiamond: //钻石
		eb = EBpickable
	case ITenergy: //闪电
		eb = EBpickable
	case ITmagic: //魔法棒，用于开始新地块
		eb = EBpickable
	case ITfruit: //果实，用于厨子制作甜品
		eb = EBpickable
	case ITbox: //箱子，袋子
		eb = EBopenable //此标识包含 EBwithstate
	case ITanimal: //动物
	case ITfree: //自由物品，除了移动(或回收)无法操作，相当于饰品
	default:
		panic(ErrUnsupportedItemType)
	}
	if it.GetNextItemId() == 0 {
		eb |= final //部分物品最终体包含更多属性
	}
	if it.IsUnique() {
		eb |= EBunique
	}
	if it.IsEradicable() {
		eb |= EBeradicable
	}
	// if st == nil {
	// 	st = initialState(it)
	// }
	f := &filler{
		item:      it,
		EasyBitag: eb,
		// state:     st,
	}

	return f
}

// initialState for given Item.
// some Items do not require State, return nil
func initialState(it Item) (st *state) {
	return &state{Non: true} //暂定没有初始状态，后续完善
}

func nearestEmptyBlocks(free []*block, dst Position, wantCnt int) (out []*block) {
	out = slices.Clone(free)
	slices.SortFunc(out, func(a, b *block) bool {
		return math.Abs(float64(a.X-dst.X))+math.Abs(float64(a.Y-dst.Y)) < math.Abs(float64(b.X-dst.X))+math.Abs(float64(b.Y-dst.Y))
	})
	if len(out) > wantCnt {
		out = out[:wantCnt]
	}
	return
}

func filterLines(line []*block, layer, itemId int, lcs ...Linecode) (out []Linecode, totalCnt int) {
	tarr, min, max := make([]Linecode, 0, 8), -1, -1
	for _, tlc := range lcs {
		if min != -1 && tlc.Min() > min && tlc.Min() < max {
			continue
		}
		if ttarr, cnt := filterLine(line, layer, itemId, tlc); len(ttarr) > 0 {
			lm, rm := ttarr[0].Min(), ttarr[len(ttarr)-1].Min()
			if min == -1 || min > lm {
				min = lm
			}
			if max == -1 || max < rm {
				max = rm
			}
			tarr, totalCnt = append(tarr, ttarr...), totalCnt+cnt
		}
	}
	return slices.Clip(tarr), totalCnt
}

func filterLine(line []*block, layer, id int, lc Linecode) (out []Linecode, totalCnt int) {
	min, max, out := lc.Min(), lc.Min()+lc.Count(), make([]Linecode, 0, 8)
	for i := min; i >= 0 && line[i] != nil && line[i].Lineable(layer, id); i-- {
		min = i //find min index
	}
	for i := max; i < len(line) && line[i] != nil && line[i].Lineable(layer, id); i++ {
		max = i //find max index
	}
	idx, cnt := -1, 0
	for i, b := range line[min : max+1] {
		if b != nil && b.Lineable(layer, id) {
			if idx == -1 {
				idx = min + i
			}
			cnt++
		} else if cnt != 0 {
			out, totalCnt = append(out, encodeLinecode(Linecode(idx), Linecode(cnt))), totalCnt+cnt
			cnt, idx = 0, -1
		}
	}
	if cnt != 0 {
		out, totalCnt = append(out, encodeLinecode(Linecode(idx), Linecode(cnt))), totalCnt+cnt
	}
	return slices.Clip(out), totalCnt
}

func clipAndSortTwodimensionalSlice[T any](s [][]T, sortFunc func(a, b T) bool) [][]T {
	for i, arr := range s {
		slices.SortFunc(arr, sortFunc)
		s[i] = slices.Clip(arr)
	}
	return slices.Clip(s)
}

func rangeTwodimensionalSlice[T any](s [][]T, f func(ia, ib int, v T)) {
	for ia, arr := range s {
		for ib, v := range arr {
			f(ia, ib, v)
		}
	}
}

func execute(fs ...func() bool) {
	for _, f := range fs {
		if !f() {
			break
		}
	}
}

// unmarshalPresets unmarshal presets to block map and layer-positions
// bmap: [y][x]*block
// lps: [layer][]Postion
// pss: [y][x]*Preset
// NOTE: 预设物品一定是初始状态
func unmarshalPresets(pss [][]*Preset, iset Itemset) (bs [][]*block, lps [][]Position, useableCnt int) {
	bs, lps = make([][]*block, len(pss)), make([][]Position, 0, 8)
	for i := range bs {
		bs[i] = make([]*block, len(pss[i]))
	}
	rangeTwodimensionalSlice(pss, func(iy, ix int, ps *Preset) {
		if ps != nil {
			if cnt := ps.Layer + 1 - len(lps); cnt > 0 {
				for i := 0; i < cnt; i++ {
					lps = append(lps, make([]Position, 0, 512))
				}
			}
			it := iset.Get(int32(ps.ItemId))
			b := &block{Position: Position{X: ix, Y: iy}, PID: ps.Plot, Site: newSitecode(ps.Layer, len(lps[ps.Layer])), Filler: convertoFiller(it, nil)}
			lps[ps.Layer], bs[iy][ix] = append(lps[ps.Layer], b.Position), b
			useableCnt++
		}
	})
	lps = clipAndSortTwodimensionalSlice(lps, func(p1, p2 Position) (less bool) {
		if less = p1.Y < p2.Y; !less {
			less = p1.Y == p2.Y && p1.X < p2.X
		}
		return
	})
	return
}

// placeItem 放置物品，初始化时，将物品放置到存储的位置
func placeItem(prefab Item, dists []uint64, mst map[uint32]State, cache [][]*block, lps [][]Position) {
	for _, v := range dists {
		dc := Distcode(v)
		layer := int(dc.Layer())
		ps := lps[layer]
		for _, idx := range dc.Indexes() {
			pos := ps[idx]
			if b := cache[pos.Y][pos.X]; b != nil {
				var dataSt *state
				if st, ok := mst[uint32(newSitecode(layer, idx))]; ok {
					dataSt = convertoDataState(st)
				}
				b.Filler = convertoFiller(prefab, dataSt)
			}
		}
	}
}

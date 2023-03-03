package play

import (
	"testing"

	"gotest.tools/v3/assert"
)

var tmpPss = [][]*Preset{
	{nil, nil, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 36, Plot: 6}, {Layer: 3, ItemId: 36, Plot: 6}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}},
	{nil, {Layer: 3, ItemId: 62, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 36, Plot: 6}, {Layer: 3, ItemId: 36, Plot: 6}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}},
	{{Layer: 3, ItemId: 1, Plot: 6}, {Layer: 3, ItemId: 1, Plot: 6}, {Layer: 3, ItemId: 62, Plot: 6}, {Layer: 3, ItemId: 91, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 81, Plot: 7}, {Layer: 3, ItemId: 81, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}},
	{{Layer: 3, ItemId: 1, Plot: 6}, {Layer: 3, ItemId: 1, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 62, Plot: 6}, {Layer: 3, ItemId: 81, Plot: 7}, {Layer: 3, ItemId: 81, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 94, Plot: 7}},
	{nil, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 12, Plot: 6}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 106, Plot: 7}, {Layer: 3, ItemId: 106, Plot: 7}},
	{nil, nil, nil, nil, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 62, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 106, Plot: 7}, {Layer: 3, ItemId: 106, Plot: 7}},
	{nil, nil, nil, nil, {Layer: 3, ItemId: 12, Plot: 6}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, nil, nil, {Layer: 3, ItemId: 71, Plot: 7}},
	{nil, nil, nil, nil, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 62, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, nil, nil, nil},
	{nil, nil, nil, nil, {Layer: 3, ItemId: 70, Plot: 6}, {Layer: 3, ItemId: 71, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, nil, nil, nil},
	{nil, nil, nil, nil, nil, {Layer: 3, ItemId: 62, Plot: 7}, {Layer: 3, ItemId: 71, Plot: 7}, nil, nil, nil},
}

func cloneTmpPss() [][]*Preset {
	pss := make([][]*Preset, len(tmpPss))
	for i, ps := range tmpPss {
		tp := make([]*Preset, len(ps))
		pss[i] = tp
		for a, p := range ps {
			if p != nil {
				tp[a] = new(Preset)
				*tp[a] = *p
			}
		}
	}
	return pss
}

var tmpDists = []Itemdist{
	&itemdist{},
}

// type tmpDist struct {
// 	id    int32
// 	codes []uint64
// }

// func (td *tmpDist) GetId() int32 {
// 	return td.id
// }

// func (td *tmpDist) GetDistcodes() []uint32 {
// 	return td.codes
// }

type tmpItem struct {
	ib  EasyBitag
	t   ItemType
	id  int32
	sbt StateBitag
	Item
}

func (ti *tmpItem) GetStateBitag() StateBitag {
	return ti.sbt
}

func (ti *tmpItem) GetId() int32 {
	return ti.id
}

func (ti *tmpItem) GetType() ItemType {
	return ti.t
}

func (ti *tmpItem) IsUnique() bool {
	return false
}
func (ti *tmpItem) IsEradicable() bool {
	return true
}

// func (ti *tmpItem) GetBitag() EasyBitag {
// 	return ti.ib
// }

func (ti *tmpItem) GetNextItemId() (nid int) {
	nid = int(ti.id)
	if nid > 0 {
		nid += 1
	}
	return nid
}

func TestFilterLines(t *testing.T) {
	line := []*block{
		// {Position: Position{X: 0, Y: 10}, LID: 2, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// {Position: Position{X: 1, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// {Position: Position{X: 2, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// nil,
		// {Position: Position{X: 4, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// {Position: Position{X: 5, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// {Position: Position{X: 6, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 22}, &state{Non: true})},
		// {Position: Position{X: 7, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// nil,
		// {Position: Position{X: 9, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// {Position: Position{X: 10, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		// {Position: Position{X: 11, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 23}, &state{Non: true})},
		// {Position: Position{X: 12, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 25}, &state{Non: true})},
		// {Position: Position{X: 13, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 22}, &state{Non: true})},
		// {Position: Position{X: 14, Y: 11}, LID: 3, Filler: convertoDataItem(&tmpItem{id: 21}, &state{Non: true})},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
	arr, cnt := filterLine(line, 3, 21, encodeLinecode(Linecode(2), Linecode(3)))
	assert.Equal(t, len(arr), 2)
	assert.Equal(t, cnt, 4)
	assert.Equal(t, arr[0], (Linecode(1)<<16)|Linecode(2))
	assert.Equal(t, arr[1], (Linecode(4)<<16)|Linecode(2))

	arr, cnt = filterLine(line, 3, 21, encodeLinecode(Linecode(2), Linecode(10)))
	assert.Equal(t, len(arr), 4)
	assert.Equal(t, cnt, 7)

	assert.Equal(t, arr[0], (Linecode(1)<<16)|Linecode(2))
	assert.Equal(t, arr[1], (Linecode(4)<<16)|Linecode(2))
	assert.Equal(t, arr[2], (Linecode(7)<<16)|Linecode(1))
	assert.Equal(t, arr[3], (Linecode(9)<<16)|Linecode(2))

	arr, cnt = filterLines(line, 3, 21, arr...)
	assert.Equal(t, len(arr), 4)
	assert.Equal(t, cnt, 7)

	arr, cnt = filterLines(line, 3, 21, encodeLinecode(Linecode(3), Linecode(2)), encodeLinecode(Linecode(9), Linecode(7)), encodeLinecode(Linecode(19), Linecode(2)))
	assert.DeepEqual(t, arr, []Linecode{
		encodeLinecode(Linecode(4), Linecode(2)),
		encodeLinecode(Linecode(9), Linecode(2)),
		encodeLinecode(Linecode(14), Linecode(1)),
	})
	assert.Equal(t, cnt, 5)
}

func TestNearestEmptyBlocks(t *testing.T) {
	// out = slices.Clone(f.free)
	// slices.SortFunc(out, func(a, b *block) bool {
	// 	return math.Abs(float64(a.X-dst.X))+math.Abs(float64(a.Y-dst.Y)) < math.Abs(float64(b.X-dst.X))+math.Abs(float64(b.Y-dst.Y))
	// })
	// if len(out) > wantCnt {
	// 	out = out[:wantCnt]
	// }
	// return
	free := []*block{
		// {Position: Position{X: 0, Y: 10}, LID: 2},
		// {Position: Position{X: 1, Y: 11}, LID: 3},
		// {Position: Position{X: 2, Y: 11}, LID: 3},
		// {Position: Position{X: 4, Y: 11}, LID: 3},
		// {Position: Position{X: 5, Y: 11}, LID: 3},
		// {Position: Position{X: 6, Y: 11}, LID: 3},
		// {Position: Position{X: 7, Y: 11}, LID: 3},
	}
	bs := nearestEmptyBlocks(free, Position{X: 3, Y: 3}, 3)
	assert.Equal(t, len(bs), 3)
	assert.DeepEqual(t, bs[0].Position, Position{X: 2, Y: 11})
	assert.DeepEqual(t, bs[1].Position, Position{X: 4, Y: 11})
	assert.DeepEqual(t, bs[2].Position, Position{X: 0, Y: 10})
	assert.Equal(t, len(free), 7, "不会改变原slice数据")

	bs = nearestEmptyBlocks(free[:1], Position{X: 3, Y: 3}, 3)
	assert.Equal(t, len(bs), 1)
	assert.DeepEqual(t, bs[0].Position, Position{X: 0, Y: 10})

	bs = nearestEmptyBlocks(nil, Position{X: 3, Y: 3}, 3)
	assert.Equal(t, len(bs), 0)
}

func TestClipAndSortTwodimensionalSlice(t *testing.T) {
	// for i, arr := range s {
	// 	slices.SortFunc(arr, sortFunc)
	// 	s[i] = slices.Clip(arr)
	// }
	// return slices.Clip(s)
	arr := make([][]int, 0, 16)
	for i := 0; i < 4; i++ {
		tmp := make([]int, 0, 8)
		for a := 0; a < i; a++ {
			tmp = append(tmp, a)
		}
		arr = append(arr, tmp)
	}
	assert.Equal(t, cap(arr), 16)

	arr = clipAndSortTwodimensionalSlice(arr, func(a, b int) bool { return a < b })
	assert.Equal(t, cap(arr), 4)

	for i, tmp := range arr {
		assert.Equal(t, cap(tmp), i)
	}
}

func TestExecute(t *testing.T) {
	num := 0
	execute(func() bool {
		// num++
		num |= 1 << 0
		return true
	}, func() bool {
		num |= 1 << 1
		return false
	}, func() bool {
		num |= 1 << 2
		return true
	})
	assert.Equal(t, num, (1<<1)|(1<<0))
}

// // unmarshalPresets unmarshal presets to block map and layer-positions
// // bmap: [y][x]*block
// // lps: [layer][]Postion
// // pss: [y][x]*Preset
// func TestUnmarshalPresets(t *testing.T) {
// 	pss := cloneTmpPss()
// 	bs, lps, cnt := unmarshalPresets(pss, func(itemId int) Item {
// 		// return newItemPrefab(&tmpItem{id: int32(itemId)}), false
// 		return &tmpItem{id: int32(itemId)}
// 	})
// 	tmpCnt := 0
// 	for _, ps := range pss {
// 		for _, p := range ps {
// 			if p != nil {
// 				tmpCnt++
// 			}
// 		}
// 	}
// 	assert.Equal(t, tmpCnt, cnt)
// 	tmpCnt = 0
// 	for _, arr := range bs {
// 		for _, b := range arr {
// 			if b != nil {
// 				tmpCnt++
// 			}
// 		}
// 	}
// 	assert.Equal(t, tmpCnt, cnt)
// 	assert.Equal(t, 4, len(lps), "最大layer + 1")
// 	tmpCnt = 0
// 	for _, ps := range lps {
// 		tmpCnt += len(ps)
// 	}
// 	assert.Equal(t, tmpCnt, cnt)

// 	for y, ps := range pss {
// 		for x, p := range ps {
// 			if p == nil {
// 				assert.Assert(t, bs[y][x] == nil)
// 			} else {
// 				it := bs[y][x]
// 				assert.Equal(t, p.ItemId, int(it.Item.GetId()))
// 				assert.Equal(t, p.Layer, it.LID)
// 				assert.Equal(t, p.Plot, it.PID)
// 			}
// 		}
// 	}
// }

func TestInitialState(t *testing.T) {

}

func TestDeployItem(t *testing.T) {
	// for _, v := range dists {
	// 	dc := Distcode(v)
	// 	ps := lps[dc.Layer()]
	// 	for _, idx := range dc.Indexes() {
	// 		pos := ps[idx]
	// 		if b := cache[pos.Y][pos.X]; b != nil && b.Hold(prefab.Clone()) {
	// 			sc := uint32(newSitecode(int(dc.Layer()), idx))
	// 			if st, ok := mst[sc]; ok {
	// 				b.GetItem().SetState(st)
	// 			}
	// 		}
	// 	}
	// }
}

func TestStateBitagCover(t *testing.T) {
	// SBCbuildingdone
	equals := [][]uint8{
		{8, 3},
		{11, 1},
		{12, 4},
		{16, 4},
		{20, 1},
		{21, 2},
	}
	covers := []StateBitagCover{
		SBCcomponentcount,
		SBCbuildingdone,
		SBCremaintimes,
		SBCstarcount,
		SBCtakeupworker,
		SBCnextworktype,
	}
	for i, cover := range covers {
		equal := equals[i]
		assert.Equal(t, cover.Low(), equal[0])
		assert.Equal(t, cover.Count(), equal[1])
	}
}

func TestStateBitag(t *testing.T) {
	var sb StateBitag
	sb = encodeStateBitag(sb, uint(ABwaitbuild), SBCactionbitag) //easy.ResetBit(sb, StateBitag(ABwaitbuild), uint(SBCactionbitag.Low()), uint(SBCactionbitag.Count()))
	assert.Equal(t, sb.GetActionBitag(), ABwaitbuild)

	sb = encodeStateBitag(sb, uint(ABmining), SBCactionbitag)
	assert.Equal(t, sb.GetActionBitag(), ABmining)

	sb = encodeStateBitag(sb, 3, SBCcomponentcount)
	assert.Equal(t, sb.GetActionBitag(), ABmining)
	assert.Equal(t, sb.GetCoverValue(SBCcomponentcount), uint(3))

	sb = encodeStateBitag(sb, 1, SBCcomponentcount)
	assert.Equal(t, sb.GetCoverValue(SBCcomponentcount), uint(1))

	//
	sbcs := []StateBitagCover{
		SBCbuildingdone,
		SBCremaintimes,
		SBCstarcount,
		SBCtakeupworker,
		SBCnextworktype,
	}
	for _, sbc := range sbcs {
		assert.Equal(t, sb.GetCoverValue(sbc), uint(0))
	}
}

func TestConvertoDataItem(t *testing.T) {
	t.Run("panic nil Item", func(t *testing.T) {
		defer func(t *testing.T) {
			r := recover()
			assert.Assert(t, r != nil)
		}(t)
		// convertoDataItem(nil, &state{Non: true})
	})
	// convertoDataItem(nil,nil)
	// ts := []ItemType{
	// 	ITcastle,         //城堡
	// 	ITrootmine,       //母矿
	// 	ITtimeworker,     //计时工人: 完全体可修建
	// 	ITcrop,           //农作物
	// 	ITcastlematerial, //城堡建材(建筑)：用于升级城堡
	// 	ITmine,           //矿
	// 	ITlocalmine,      //固定位矿(未解锁时填充矿)
	// 	ITcrystal,        //水晶
	// 	ITcook,           //厨子
	// 	ITcoin,           //金币
	// 	ITdiamond,        //钻石
	// 	ITenergy,         //闪电
	// 	ITmagic,          //魔法棒，用于开始新地块
	// 	ITfruit,          //果实，用于厨子制作甜品
	// 	ITbox,            //箱子，袋子
	// 	ITanimal,         //动物
	// 	ITfree,           //自由物品，除了移动(或回收)无法操作，相当于饰品
	// }
	// for _, t := range ts {
	// 	it := convertoDataItem(&tmpItem{t: t}, &state{Non: true})
	// 	it.Able()
	// }
}

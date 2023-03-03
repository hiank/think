package play

type block struct {
	Position
	// LID  int //Layer id (layer index)
	PID  int //plot id
	Site Sitecode
	// Locked bool  //locked before unlock by use magic
	Filler *filler //nil when non item
	///应该持有此item 的状态，并且有对此状态的'增删改查'接口对象. 另一种思路是，item 包含状态的变换处理
	// State *state
}

// // func newBlock(pos Position, )

// //Hold 持有给定的item. 如果原来已持有某个item, 需要先'Clear'清除旧item
// func (b *block) Hold(it *item) (suc bool) {
// 	if suc = b.item == nil; suc {
// 		b.item = it
// 	}
// 	return
// }

// //RemoveItem remove current item and it's state
// //
// func (b *block) RemoveItem() (oldItem *item, oldState State) {
// 	return
// }

// // func (b *block) Clear() {
// // 	b.item = nil
// // }

// func (b *block) GetItem() *item {
// 	return b.item
// }

// Empty 是否为空
func (b *block) Empty() bool {
	return b.Filler.EasyBitag == EBnon
}

// Lineable 是否可连线(用于合并)
func (b *block) Lineable(layer, itemId int) (able bool) {
	// if layer == b.LID && b.Item != nil
	// if b.Item != nil && b.Item.GetId() == int32(itemId)
	// return b.LID == layer && b.item != nil && int(b.item.GetId()) == itemId
	return int(b.Site.Layer()) == layer && b.Filler != nil && b.Filler.GetItem().GetId() == int32(itemId) && !b.Filler.Able(EBfinal) //b.mergeable()
}

// // mergeable check wether the block's item could merge
// func (b *block) mergeable() (able bool) {
// 	// if !b.Able(IBfinal) {
// 	// 	///必须是完成修建状态

// 	// }
// 	return
// }

// filler 填充物，包含item及物品
type filler struct {
	EasyBitag
	item  Item
	state *state
}

func (f *filler) GetItem() Item {
	return f.item
}

func (f *filler) SetState(cfg State) {
	///
	// if st == nil {
	if !f.Able(EBkeystate) {
		f.state = &state{Non: true}
		return
	}
	// }
	//
	// var st *state
	// if cfg == nil {
	// 	st = initialState(f.item)
	// } else {
	// 	st = convertoDataState(cfg)
	// }
	st := convertoDataState(cfg)
	////
	////
	f.state = st
}

// func (f *filler) GetState() *state {
// 	return f.state
// }

// //Update 更新，主要是
// func (f *filler) Update() {

// }

type item struct {
	Item
	// ItemBitag
	EasyBitag
	state *state
}

func (it *item) GetStateBitag() StateBitag {
	return StateBitag(it.state.GetBitag())
}

// // func newItem()

// func (it *item) MergeAble() (able bool) {
// 	if !it.GetBitag().Able(IBfinal) {

// 	}
// 	return false
// }

type state struct {
	Non bool //无状态，如果物品无需记录状态，则此值为true. 如果使用nil可能会不好判断是否已经处理过状态数据(非最终版本)
	// empty bool
	// State
	site uint32
	// bitag  uint32
	Bitag  StateBitag
	cutime int64
	ex     int32 //城堡记录经验值；计时工人记录剩余工作时间
	awards []int32
}

func (st *state) GetSite() uint32 {
	return st.site
}

func (st *state) GetBitag() uint32 {
	return uint32(st.Bitag)
}

func (st *state) GetCutime() int64 {
	return st.cutime
}

func (st *state) GetEx() int32 {
	return st.ex
}

func (st *state) GetAwards() []int32 {
	return st.awards
}

// func (st *state) Bitag() StateBitag {
// 	return StateBitag(st.bitag)
// }

type backpack struct {
	// mres  map[ItemType]int //map[res type]count
	// mitem map[int]int      //map[itemId]count
	m map[int]int //map[itemId]count
	// m map[ItemType]int //map[res type]max
}

func (bp *backpack) GetCount(id int) (cnt int, owned bool) {
	///
	cnt, owned = bp.m[id]
	return
}

type resource struct {
	// Resourcecode
}

func (res *resource) GetCode() uint64 {
	return 0
}
func (res *resource) GetId() int32 {
	return 0
}
func (res *resource) GetTimestamp() int64 {
	return 0
}

type itemdist struct {
	id    int32
	codes []Distcode
}

func (dist *itemdist) GetId() int32 {
	return dist.id
}

func (dist *itemdist) GetDistcodes() []uint64 {
	// return slices.Clip(dist.codes)
	return nil
}

type farmDataset struct {
	f *farm
}

func (fd *farmDataset) GetDists() []Itemdist {
	m := make(map[int32]*itemdist)
	// for _, b := range fd.
	// fd.f.cache
	rangeTwodimensionalSlice(fd.f.cache, func(_, _ int, b *block) {
		//
		if !b.Empty() {
			dist, ok := m[b.Filler.GetItem().GetId()]
			if !ok {
				dist = &itemdist{id: b.Filler.GetItem().GetId(), codes: make([]Distcode, 0, 64)}
				m[dist.GetId()] = dist
			}
			// dist.codes = append(dist.codes, b.Position)
		}
	})
	return nil
}
func (fd *farmDataset) GetStates() []State {
	return nil
}

// plots unlock info
func (fd *farmDataset) GetUnlockcode() uint64 {
	return 0
}
func (fd *farmDataset) GetResources() []Resource {
	return nil
}

// func (res *resource) GetType() ItemType {
// 	return ITundefined
// }

// func (res *resource) GetLimit() int {
// 	return 0
// }

// func (res *resource) GetCount() int {
// 	return 0
// }

// func (res *resource) GetTimestamp() int64 {
// 	return 0
// }

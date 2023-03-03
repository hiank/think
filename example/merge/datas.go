package merge

import (
	"github.com/hiank/think/exp/convert"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// type Position struct {
// 	X, Y int
// }

// func (pos Position) Equal(other Position) bool {
// 	return other.X == pos.X && other.Y == pos.Y
// }

type block struct {
	Position       //[readonly] {X, Y int}
	Sitecode       //[readonly] layer id and index in layer
	Plot     uint8 //[readonly] plot id
	Filler   filler
}

type filler struct {
	EasyBitag
	item  Item
	state *state
}

// Empty determine if the filler is empty
func (f filler) Empty() bool {
	return f.EasyBitag == EBempty
}

type state struct {
	Non    bool //no need
	site   uint32
	bitag  StateBitag
	cutime int64
	ex     int32
	awards []int32
}

func (st *state) GetSite() uint32 {
	return st.site
}
func (st *state) GetBitag() uint32 {
	return uint32(st.bitag)
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

type farmDataset struct {
	f *farm
}

func (fd *farmDataset) GetDists() []Itemdist {
	mdist := make(map[int32]map[uint16]Distcode)
	// rangeTwodimensional(fd.f.tdblock, func(_, _ int, b *block) {
	// 	if fd.f.unlockcode.Unlocked(b.Plot) && !b.Filler.Empty() {
	fd.rangeWorkingBlocks(func(b *block) {
		dc := convertoDistcode(b.Sitecode)
		base := dc.Base()
		m, ok := mdist[b.Filler.item.GetId()]
		if !ok {
			m = make(map[uint16]Distcode)
			mdist[b.Filler.item.GetId()] = m
		} else if lastdc, ok := m[base]; ok {
			dc |= lastdc
		}
		m[base] = dc
	})
	dists := make([]Itemdist, 0, len(mdist))
	converter := convert.ConverterFunc[Distcode, uint64](func(v Distcode) uint64 { return uint64(v) })
	for id, m := range mdist {
		dists = append(dists, Itemdist{
			Id:    id,
			Dists: convert.Slice[Distcode, uint64](maps.Values(m), converter),
		})
	}
	return dists
}

// GetStates get states from *farm
func (fd *farmDataset) GetStates() []State {
	sts := make([]State, 0, 1024)
	// rangeTwodimensional(fd.f.tdblock, func(_, _ int, b *block) {
	fd.rangeWorkingBlocks(func(b *block) {
		if !b.Filler.state.Non {
			sts = append(sts, b.Filler.state)
		}
	})
	return slices.Clip(sts)
}

// plots unlock info
func (fd *farmDataset) GetUnlockcode() uint64 {
	///
	return uint64(fd.f.unlockcode)
}

// GetResources resources count
func (fd *farmDataset) GetResources() []Resource {
	return maps.Values(fd.f.mres)
}

func (fd *farmDataset) rangeWorkingBlocks(f func(*block)) {
	///
	rangeTwodimensional(fd.f.blocks, func(_, _ int, b *block) {
		if fd.f.unlockcode.Unlocked(b.Plot) && !b.Filler.Empty() {
			f(b)
		}
	})
}

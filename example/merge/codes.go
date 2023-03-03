package merge

import (
	"math"

	"github.com/hiank/think/exp/easy"
	"github.com/hiank/think/run"
	"golang.org/x/exp/slices"
)

const (
	ErrInvalidValue = run.Err("invalid value (number exceed the limit)")

	distRecordCount      = 48
	siteIndexLimit  uint = (1+math.MaxUint8)*distRecordCount - 1
)

// Sitecode for save site info
// 0-7bit layer id
// 8-31bit index in layer
type Sitecode uint32

func encodeSitecode(layer uint8, idx uint) Sitecode {
	if idx > siteIndexLimit {
		panic(ErrInvalidValue)
	}
	return Sitecode(layer) | (Sitecode(idx) << 8)
}

// Layer layer id (index)
func (sc Sitecode) Layer() uint8 {
	return uint8(sc)
}

// Index index in layer
func (sc Sitecode) Index() uint {
	return uint(sc >> 8)
}

type Distcode uint64

func encodeBaseDistcode(layer, baseIdx uint8) Distcode {
	return (Distcode(layer) << 56) | (Distcode(baseIdx) << 48)
}

func (dc Distcode) Layer() uint8 {
	return uint8(dc >> 56)
}

func (dc Distcode) baseIndex() uint8 {
	return uint8(dc >> 48)
}

// Base baseindex & layer
func (dc Distcode) Base() uint16 {
	return uint16(dc >> 48)
}

func (dc Distcode) Indexes() []int {
	baseIdx, out := int(dc.baseIndex())*distRecordCount, make([]int, 0, distRecordCount)
	for i := 0; i < distRecordCount; i++ {
		if dc&(1<<i) != 0 {
			out = append(out, baseIdx+i)
		}
	}
	return slices.Clip(out)
}

type Position uint32

func encodePosition(x, y uint16) Position {
	return (Position(y) << 16) | Position(x)
}

func (pos Position) X() uint16 {
	return uint16(pos)
}

func (pos Position) Y() uint16 {
	return uint16(pos >> 16)
}

type Resourcecode uint64

func encodeResourcecode(it ItemType, limit uint8, cnt uint) Resourcecode {
	return Resourcecode(it) | (Resourcecode(limit) << 8) | (Resourcecode(cnt) << 32)
}

// Type low 8bit. ItemType
func (rc Resourcecode) Type() ItemType {
	return ItemType(rc)
}

// Limit 最大值限制
// 0 means no limit
func (rc Resourcecode) Limit() uint8 {
	v := easy.BitValue(rc, 8, 8)
	return uint8(v)
}

func (rc Resourcecode) Count() uint {
	return uint(easy.BitValue(rc, 32, 32))
}

type Unlockcode uint64

func (uc Unlockcode) Unlocked(plot uint8) bool {
	///
	return uc&(1<<plot) != 0
}

package play

import (
	"math"

	"github.com/hiank/think/exp/easy"
	"github.com/hiank/think/run"
	"golang.org/x/exp/slices"
)

const (
	ErrExceedMaxUnit16    = run.Err("play: exceed max uint16 limit")
	ErrExceedMaxSiteIndex = run.Err("paly: exceed max site index limit")

	siteBitLayer   int = 24
	siteIndexCount int = 24
	maxSiteIndex   int = 1<<siteIndexCount - 1

	distBitBaseIndex int = 48
	distBitLayer     int = 56
	distIndexCount   int = 48
)

type Linecode uint32

func encodeLinecode[T ~uint32](min, count T) T {
	if min > math.MaxUint16 || count > math.MaxUint16 {
		panic(ErrExceedMaxUnit16)
	}
	return (min << 16) | count
}

func (lc Linecode) Min() int {
	return int(lc >> 16)
}

func (lc Linecode) Count() int {
	return int(lc & math.MaxUint16)
}

type Sitecode uint32

func newSitecode(layer, index int) Sitecode {
	if layer > math.MaxUint8 || layer < 0 {
		panic(ErrExceedMaxUnit16)
	}
	if index > maxSiteIndex || index < 0 {
		panic(ErrExceedMaxSiteIndex)
	}
	return Sitecode(layer<<uint32(siteBitLayer) | index)
}

// Layer index (max is 15)
func (sc Sitecode) Layer() uint8 {
	return uint8((sc >> siteBitLayer) & math.MaxUint8)
}

// BaseIndex
func (sc Sitecode) Index() int {
	return maxSiteIndex & int(sc)
}

type Distcode uint64

// Layer index (max is 15)
func (dc Distcode) Layer() uint8 {
	return uint8((dc >> distBitLayer) & math.MaxUint8)
}

// BaseIndex
func (dc Distcode) BaseIndex() (baseIdx int) {
	return int((dc>>distBitBaseIndex)&math.MaxUint8) * distIndexCount
}

// Indexes tagged indexes
func (dc Distcode) Indexes() []int {
	baseIdx, out := dc.BaseIndex(), make([]int, 0, distIndexCount)
	for i := 0; i < distIndexCount; i++ {
		if int(dc>>i)&1 == 1 {
			out = append(out, i+baseIdx)
		}
	}
	return slices.Clip(out)
}

type Unlockcode uint64

func (uc Unlockcode) Unlocked(plotId int) bool {
	var one uint64 = 1
	return ((one << uint64(plotId)) & uint64(uc)) != 0
}

type Resourcecode uint64

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

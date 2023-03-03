package merge

import (
	"context"
	"io"
	"time"

	"github.com/hiank/think/run"
)

type farm struct {
	blocks     [][]*block          //[y][x]*block
	ltops      [][]Position        //[layer][]Position
	idles      []*block            //闲置的所有可用地块(已解锁&无物品)
	mres       map[uint64]Resource //map[res key]Resource
	unlockcode Unlockcode          //地块解锁信息

	tsset *run.Timestampset //timestampset
	iset  Itemset           //item config
	io.Closer
}

func New(ctx context.Context, iset Itemset, lastDs Dataset, mp [][]Preset) (Farm, Dataset) {
	ctx, closer := run.StartHealthyMonitoring(ctx)
	f := &farm{
		iset:   iset,
		tsset:  run.NewTimestampset(ctx, time.Second),
		Closer: closer,
	}

	return f, &farmDataset{f: f}
}

func (f *farm) DoMove(from, end uint32) bool {
	return false
}

func (f *farm) DoClick(end uint32) bool {
	return false
}

func (f *farm) DoEradicate(end uint32) bool {
	return false
}

func (f *farm) DoUnlock(plot uint8) bool {
	return false
}

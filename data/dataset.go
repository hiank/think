package dset

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hiank/think/data/db"
	"k8s.io/klog/v2"
)

func defaultOptions() options {
	return options{}
}

type liteDB struct {
	memory db.IClient //NOTE: memory database
	disk   db.IClient //NOTE: disk database ()
}

func (lc *liteDB) HGet(hkey *DBKey, fkey string) (parser db.IParser, err error) {
	useMem, err := hkey.InTag(TagUseMemory), fmt.Errorf("invalid HashKey %v", hkey)
	if useMem {
		//load data from memory database
		if parser, err = lc.memory.HGet(hkey.Result, fkey); err == nil {
			return
		}
	}
	if hkey.InTag(TagUseDisk) && (lc.disk != nil) {
		//load data from disk database (when load from memory db failed)
		if parser, err = lc.disk.HGet(hkey.Result, fkey); (err == nil) && useMem {
			var val string
			if val, err = parser.Result(); err == nil {
				//save data to memory db (TagUseMemory)
				err = lc.memory.HSet(hkey.Result, fkey, val)
			}
		}
	}
	return
}

func (lc *liteDB) HSet(hkey *DBKey, values ...interface{}) (err error) {
	if hkey.InTag(TagUseMemory) {
		if err = lc.memory.HSet(hkey.Result, values...); err != nil {
			klog.Warning("set data to memory database failed: ", err)
		}
	}
	if (lc.disk != nil) && hkey.InTag(TagUseDisk) {
		if err = lc.disk.HSet(hkey.Result, values...); err != nil {
			klog.Warning("set data to disk database failed: ", err)
		}
	}
	return
}

type liteSet struct {
	ctx  context.Context
	opts options
	*liteDB
}

func NewDataset(ctx context.Context, opts ...Option) IDataset {
	dopts := defaultOptions()
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	if dopts.memoryDB == nil {
		panic("memory database must be set")
	}
	ls := &liteSet{ctx: ctx, opts: dopts, liteDB: &liteDB{dopts.memoryDB, dopts.diskDB}}
	return ls
}

func (lh *liteSet) GetGamer(uid uint64) (gamer IGamer, err error) {
	parser, err := lh.HGet(&DBKey{Tag: TagUseMemory | TagUseDisk, Result: hkeyGamer}, strconv.FormatUint(uid, 10))
	if err == nil {
		gamer = lh.opts.buildGamer()
		err = parser.Scan(gamer)
	}
	return
}

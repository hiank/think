package data

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hiank/think/data/db"
	"k8s.io/klog/v2"
)

func pushError(err, ex error) error {
	if err == nil {
		return ex
	}
	if ex != nil {
		err = fmt.Errorf("%s&&%s", err.Error(), ex.Error())
	}
	return err
}

func decode(k string) (kt KeyTag, baseKey string, err error) {
	var val int
	if strs := regexp.MustCompile(ktRegexp).FindStringSubmatch(k); len(strs) < 2 {
		err = fmt.Errorf("invalid key(%s): non KeyTag information(%v)", k, strs)
	} else if len(strs[0]) == len(k) {
		err = fmt.Errorf("invalid key(%s): non base key information", k)
	} else if val, err = strconv.Atoi(strs[1]); err == nil {
		if kt = KeyTag(val); KTMix.equal(kt) {
			baseKey = k[len(strs[0]):] //strs[0] is "[`kt`@KT]"
		} else {
			kt, err = 0, fmt.Errorf("invalid key(%s): overflow keytag(%d)", k, val)
		}
	}
	return
}

//robustDB encapsulation of client
type robustDB struct {
	store db.KvDB
}

func (rd *robustDB) Set(k string, v interface{}) error {
	_, k, _ = decode(k)
	return rd.store.Set(k, v)
}

// Get retrieves the value for the given key.
func (rd *robustDB) Get(k string, v interface{}) (found bool, err error) {
	_, k, _ = decode(k)
	return rd.store.Get(k, v)
}

// Delete deletes the stored value for the given key.
func (rd *robustDB) Delete(k string) error {
	_, k, _ = decode(k)
	return rd.store.Delete(k)
}

func (rd *robustDB) Close() (err error) {
	return rd.store.Close()
}

type mixDB struct {
	mstore map[KeyTag]db.KvDB
}

func (md *mixDB) decode(k string) (kt KeyTag, baseKey string, err error) {
	if kt, baseKey, err = decode(k); err == nil {
		var lastKt uint8
		ktcacheKey := fmt.Sprintf("KTCACHE@%s", baseKey)
		found, gerr := md.mstore[KTMem].Get(ktcacheKey, &lastKt)
		if !found {
			if err = md.mstore[KTMem].Set(ktcacheKey, uint8(kt)); err != nil {
				kt, baseKey = 0, ""
			}
		} else if gerr != nil {
			//found but contians error
			kt, baseKey, err = 0, "", gerr
		} else if lastKt != uint8(kt) {
			//found but inconsistent
			kt, baseKey, err = 0, "", fmt.Errorf("inconsistent with cached KeyTag(cached:%d): %s", lastKt, k)
		}
		if kt == KTMix {
			klog.Warning("using mixed mode is prone to data inconsistency and should be avoided")
		}
	}
	return
}

func (md *mixDB) Set(k string, v interface{}) error {
	kt, k, err := md.decode(k)
	if err != nil {
		return err
	}
	for mkt, store := range md.mstore {
		if kt.equal(mkt) {
			err = pushError(err, store.Set(k, v))
		}
	}
	return err
}

// Get retrieves the value for the given key.
func (md *mixDB) Get(k string, v interface{}) (found bool, err error) {
	kt, k, err := md.decode(k)
	if err != nil {
		return
	}
	//
	if kt.equal(KTMem) {
		if found, err = md.mstore[KTMem].Get(k, v); found {
			return
		}
	}
	if kt.equal(KTDisk) {
		if found, err = md.mstore[KTDisk].Get(k, v); found && err == nil {
			if kt.equal(KTMem) {
				if serr := md.mstore[KTMem].Set(k, v); serr != nil {
					klog.Warning("failed to set value to memory store:", serr)
				}
			}
		}
	}
	return
}

// Delete deletes the stored value for the given key.
func (md *mixDB) Delete(k string) error {
	kt, k, err := md.decode(k)
	if err != nil {
		return err
	}
	for mkt, store := range md.mstore {
		if kt.equal(mkt) {
			err = pushError(err, store.Delete(k))
		}
	}
	md.mstore[KTMem].Delete(fmt.Sprintf("KTCACHE@%s", k))
	return err
}

func (md *mixDB) Close() (err error) {
	for _, store := range md.mstore {
		err = pushError(err, store.Close())
	}
	return
}

type liteSet struct {
	store db.KvDB
}

func (ls *liteSet) KvDB() db.KvDB {
	return ls.store
}

//NewDataset create a new Dataset
//NOTE: at least one k-v database is required
func NewDataset(mstore map[KeyTag]db.KvDB) Dataset {
	var store db.KvDB
	switch len(mstore) {
	case 1:
		for _, val := range mstore {
			store = &robustDB{store: val}
		}
	case 2:
		store = &mixDB{mstore: mstore}
	default:
		panic("at least one k-v database is required")
	}
	return &liteSet{store: store}
}
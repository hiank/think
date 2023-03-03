package testitem

import "github.com/hiank/think/example/merge/play"

type Itemset int

func (ts Itemset) Get(id int32) play.Item {
	return itemCacheInstance()[id]
}

func (ts Itemset) GetInType(t play.ItemType) []play.Item {
	cache, out := itemCacheInstance(), make([]play.Item, 0, 8)
	for _, it := range cache {
		if it.GetType() == t {
			out = append(out, it)
		}
	}
	return out
}

func (ts Itemset) GetRelated(tmp play.Item) []play.Item {
	tmpIt, out := tmp.(*item), make([]play.Item, 0, 8)
	if tmpIt.mergeId == 0 {
		return out
	}
	for _, it := range itemCacheInstance() {
		if it.mergeId == tmpIt.mergeId {
			out = append(out, it)
		}
	}
	return out
}

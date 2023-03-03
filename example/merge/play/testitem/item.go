package testitem

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/hiank/think/example/merge/play"
	// "github.com/hiank/think/example/merge/farm"
)

var (
	itemCache map[int32]*item = make(map[int32]*item)
	loadOnce  sync.Once
)

type itemSetJson struct {
	Items []*item `json:"items"`
}

func itemCacheInstance() map[int32]*item {
	loadOnce.Do(func() {
		path := "../testitem/item.json"
		itemJson := &itemSetJson{}
		buf, err := os.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(buf, itemJson)
		}
		if err != nil {
			panic(err)
		}
		for _, it := range itemJson.Items {
			itemCache[int32(it.GetId())] = it
		}
		mergeId, items := 1, itemJson.Items
		for len(items) > 0 {
			idx := len(items) - 1
			it := items[idx]
			items = items[:idx]
			if it == nil || it.GetNextItemId() == 0 || it.mergeId != 0 {
				continue
			}
			tmpIts, mid := make([]*item, 0, 8), mergeId
			for it != nil {
				if it.mergeId != 0 {
					mid = it.mergeId
				}
				tmpIts = append(tmpIts, it)
				it = itemCache[int32(it.GetNextItemId())]
			}
			for _, it := range tmpIts {
				it.mergeId = mid
			}
			if mid == mergeId {
				mergeId++
			}
		}
	})
	return itemCache
}

type item struct {
	Id                 int32  `json:"id"`
	Type               int    `json:"type"`
	Nextid             int    `json:"nextid"`
	Unique             bool   `json:"unique"`
	Castlebind         int    `json:"castlebind"`
	MaxMineTimes       int    `json:"maxMineTimes"`
	LatestMineDuration int    `json:"latestMineDuration"`
	MaxRewardTimes     int    `json:"maxRewardTimes"`
	BuiltDuration      int    `json:"builtDuration"`
	Interval           int    `json:"interval"`
	ExpAddition        int    `json:"expAddition"`
	WorkerDuration     int    `json:"workerDuration"`
	ClearAward         int    `json:"clearAward"`
	FixedAwards        int    `json:"fixedAwards"`
	ClickAwards        string `json:"clickAwards"`
	RemoveAwards       string `json:"removeAwards"`
	Count              int    `json:"count"`
	TransformId        int    `json:"transformId"`
	EradicateAble      bool   `json:"eradicateAble"`
	mergeId            int    //NOTE: use to GetRelated
}

// GetId item id
func (it *item) GetId() int32 {
	return it.Id
}

// GetType item type
func (it *item) GetType() (t play.ItemType) {
	switch it.Type {
	case 1:
		t = play.ITcastle
	case 2:
		t = play.ITcastlematerial
	case 3:
		t = play.ITtimeworker
	case 4:
		t = play.ITcrop
	case 5:
		t = play.ITfruit
	case 6:
		t = play.ITmine
	case 7:
		t = play.ITcrystal
	case 8:
		t = play.ITrootmine
	case 9:
		t = play.ITcook
	case 10:
		t = play.ITbox
	case 11:
		t = play.ITlocalmine
	case 12:
		t = play.ITcoin
	case 13:
		t = play.ITdiamond
	case 14:
		t = play.ITenergy
	case 15:
		t = play.ITmagic
	default:
		panic(fmt.Errorf("cannot support the item type: %d", it.Type))
	}
	return
}

// GetMaxMineTimes max mine times
func (it *item) GetMaxMineTimes() int {
	return it.MaxMineTimes
}

// GetLatestMineDuraion latest mine duration (s)
func (it *item) GetLatestMineDuration() int { return it.LatestMineDuration }

// GetLightningNeed lightning need
func (it *item) GetLightningNeed() int { return 5 }

// GetMaxRewardTimes max reward times
func (it *item) GetMaxRewardTimes() int { return it.MaxRewardTimes }

// GetInterval interval between towice reward
// crop or castle reward need the value
func (it *item) GetInterval() int { return it.Interval }

// GetResCount resource count (add value for click the resource type item)
func (it *item) GetResCount() int { return it.Count }

func (it *item) strToIds(str string) []int {
	arr, out := strings.Split(str, "|"), make([]int, 0, 8)
	for _, str := range arr {
		pair := strings.Split(str, ":")
		if id, err := strconv.Atoi(pair[0]); err == nil {
			if cnt, err := strconv.Atoi(pair[1]); err == nil {
				for ; cnt > 0; cnt-- {
					out = append(out, id)
				}
			}
		}
	}
	return out
}

// GetFinalAwards final awards (the item disappeared, then the awards appear)
func (it *item) GetFinalAwards() []int {
	return it.strToIds(it.RemoveAwards)
}

// GetClickAwards awards from click the item
func (it *item) GetClickAwards() []int {
	return it.strToIds(it.ClickAwards)
}

// GetBuiltDuration built the building (castle, rootmine also) duration
func (it *item) GetBuiltDuration() int { return it.BuiltDuration }

// GetTargetCastleId castle id for building
func (it *item) GetTargetCastleId() int { return it.Castlebind }

// GetExpaddForCastle exp value for building to upgrade the castle
func (it *item) GetExpaddForCastle() int { return it.ExpAddition }

// GetMaxExpCurrentStar max exp value for current star
func (it *item) GetMaxExpCurrentStar() int { return 1000 }

// GetNextItemId next item id by merger
func (it *item) GetNextItemId() int { return it.Nextid }

// GetWorkerDuration time worker's max work time (s)
func (it *item) GetWorkerDuration() int { return it.WorkerDuration }

// IsUnique whether the item is unique. when it is, the item would appear in game only once (one)
func (it *item) IsUnique() bool { return it.Unique }

// GetTransformId when the item is unique, and it appeard now, all base item change to this id's item
func (it *item) GetTransformId() int { return it.TransformId }

// IsEradicateAble whether then item could be eradicated
func (it *item) IsEradicable() bool { return it.EradicateAble }

// func GetAllItemIds() []int32 {
// 	m := itemCacheInstance()
// 	out := make([]int32, 0, len(m))
// 	for id := range m {
// 		out = append(out, id)
// 	}
// 	return out
// }

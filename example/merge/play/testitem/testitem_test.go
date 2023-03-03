package testitem

import (
	"testing"

	// "github.com/hiank/think/example/merge/farm"
	"github.com/hiank/think/example/merge/play"
	"gotest.tools/v3/assert"
)

func TestJsonLoad(t *testing.T) {
	val := itemCacheInstance()
	assert.Assert(t, len(val) > 0)
}

var (
	testIdFarmType map[int32]play.ItemType = map[int32]play.ItemType{
		1:   play.ITcastle,         //farm.ItemCastle,
		2:   play.ITcastlematerial, //farm.ItemBuilding,
		3:   play.ITcastlematerial, //farm.ItemBuilding,
		4:   play.ITcastlematerial, //farm.ItemBuilding,
		5:   play.ITcastlematerial,
		10:  play.ITtimeworker,
		11:  play.ITtimeworker,
		12:  play.ITtimeworker,
		13:  play.ITtimeworker,
		20:  play.ITcrop,
		21:  play.ITcrop,
		22:  play.ITcrop,
		23:  play.ITcrop,
		24:  play.ITfruit,
		30:  play.ITmine, //farm.ItemMergeMine | farm.ItemEradicable,
		31:  play.ITmine,
		32:  play.ITmine,
		33:  play.ITmine,
		34:  play.ITmine,
		35:  play.ITmine,
		36:  play.ITrootmine,
		40:  play.ITcrystal,
		41:  play.ITcrystal,
		42:  play.ITcrystal,
		43:  play.ITcrystal,
		50:  play.ITcook,
		51:  play.ITcook,
		52:  play.ITcook,
		53:  play.ITcook,
		60:  play.ITbox,
		61:  play.ITbox,
		62:  play.ITbox,
		70:  play.ITlocalmine,
		71:  play.ITlocalmine,
		81:  play.ITcastle,
		82:  play.ITcastlematerial,
		83:  play.ITcastlematerial,
		84:  play.ITcastlematerial,
		85:  play.ITcastlematerial,
		90:  play.ITcoin,
		91:  play.ITcoin,
		92:  play.ITcoin,
		93:  play.ITdiamond,
		94:  play.ITdiamond,
		95:  play.ITdiamond,
		96:  play.ITmagic,
		97:  play.ITmagic,
		100: play.ITmine, //farm.ItemMergeMine | farm.ItemEradicable,
		101: play.ITmine,
		102: play.ITmine,
		103: play.ITmine,
		104: play.ITmine,
		105: play.ITmine,
		106: play.ITrootmine,
	}
)

func TestItem(t *testing.T) {
	val := itemCacheInstance()
	t.Run("check type", func(t *testing.T) {
		for _, item := range val {
			assert.Equal(t, item.GetType(), testIdFarmType[item.GetId()], "error id:%d", item.GetId())
		}
	})
	t.Run("tmplate", func(t *testing.T) {
		//GetMaxMineTimes max mine times
		it := val[84]
		assert.Equal(t, it.GetMaxMineTimes(), 0)
		assert.Equal(t, it.GetLatestMineDuration(), 0)
		assert.Equal(t, it.GetLightningNeed(), 5, "测试数据固定配置")
		assert.Equal(t, it.GetMaxRewardTimes(), 0)
		assert.Equal(t, it.GetInterval(), 0)
		assert.Equal(t, it.GetResCount(), 0)
		assert.DeepEqual(t, it.GetFinalAwards(), []int{})
		assert.DeepEqual(t, it.GetClickAwards(), []int{})
		assert.Equal(t, it.GetBuiltDuration(), 40)
		assert.Equal(t, it.GetTargetCastleId(), 81)
		assert.Equal(t, it.GetExpaddForCastle(), 10)
		assert.Equal(t, it.GetMaxExpCurrentStar(), 1000, "测试数据固定配置")
		assert.Equal(t, it.GetNextItemId(), 85)
		assert.Equal(t, it.GetWorkerDuration(), 0)
		assert.Assert(t, !it.IsUnique())
		assert.Equal(t, it.GetTransformId(), 0)
		assert.Assert(t, !it.IsEradicable())
	})
}

func TestItemSet(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		var itemSet Itemset
		it := itemSet.Get(1)
		assert.Equal(t, it.GetType(), testIdFarmType[1])
		assert.Equal(t, Itemset(1).Get(1), it)
	})
	t.Run("GetInType", func(t *testing.T) {
		var itemSet Itemset
		its := itemSet.GetInType(play.ITcastle)
		assert.Equal(t, len(its), 2)
		its = itemSet.GetInType(play.ITcrop)
		assert.Equal(t, len(its), 4)
	})
	t.Run("GetRelated", func(t *testing.T) {
		var itemSet Itemset
		m := map[int32]int{
			2:  4,
			13: 4,
			21: 4,
			31: 2,
			40: 4,
			53: 4,
			60: 3,
			85: 4,
			90: 3,
			96: 2,
		}
		for id, cnt := range m {
			its := itemSet.GetRelated(itemSet.Get(id))
			assert.Equal(t, len(its), cnt)
		}
	})
}

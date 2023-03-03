package play

//Preset 预设数据，初始时，每个可用地块都会有个预设物品(可为空)
type Preset struct {
	ItemId int //预设物品id (需要注意多block占用的物品不能跨plot占用.多block占用 其余位置不设id)
	Layer  int //所处layer index, 同一个layer index上才可以合并
	Plot   int //所处plot id
}

type Farm interface {
	// Clone() Farm
	// Md5() string
	DoMove(from, end uint32) bool
	DoClick(end uint32) bool
	DoEradicate(end uint32) bool
}

//Resource 操作需要消耗的数值(体力，金币，钻石等)
type Resource interface {
	GetCode() uint64     //Resourcecode
	GetId() int32        //(option)果实需要记录id
	GetTimestamp() int64 //(option)记录的时刻
}

// type Typeset struct {
// 	Limit, Count int
// }

// //Backpack 背包，保存物品数量(例如果实数量)
// //NOTE: 图鉴对后端程序是没有意义的，不考虑
// type Backpack interface {
// 	//Count for item with given id
// 	//owned (true) means the gamer at least once owned this item
// 	GetCount(itemId int) (cnt int, owned bool)
// }

type State interface {
	GetSite() uint32
	GetBitag() uint32
	GetCutime() int64
	GetEx() int32 //城堡记录经验值；计时工人记录剩余工作时间
	GetAwards() []int32
}

type ItemType uint8

const (
	//以下为具体物品
	ITundefined      ItemType = iota
	ITcastle                  //城堡
	ITrootmine                //母矿
	ITtimeworker              //计时工人: 完全体可修建
	ITcrop                    //农作物
	ITcastlematerial          //城堡建材(建筑)：用于升级城堡
	ITmine                    //矿
	ITlocalmine               //固定位矿(未解锁时填充矿)
	ITcrystal                 //水晶
	ITcook                    //厨子
	ITcoin                    //金币
	ITdiamond                 //钻石
	ITenergy                  //闪电
	ITmagic                   //魔法棒，用于开始新地块
	ITworker                  //工人：一般用于Resource，暂时没有具体物品可以添加工人值
	ITfruit                   //果实，用于厨子制作甜品
	ITbox                     //箱子，袋子
	ITanimal                  //动物
	ITfree                    //自由物品，除了移动(或回收)无法操作，相当于饰品
)

//Item game item config (in excel)
type Item interface {
	//GetId item id
	GetId() int32
	//GetType item type
	//NOTE: return value should be play.ItemBitag
	GetType() ItemType
	//GetMaxMineTimes max mine times
	GetMaxMineTimes() int
	//GetLatestMineDuraion latest mine duration (s)
	GetLatestMineDuration() int
	//GetLightningNeed lightning need
	GetLightningNeed() int
	//GetMaxRewardTimes max reward times
	GetMaxRewardTimes() int
	//GetInterval interval between towice reward
	//crop or castle reward need the value
	GetInterval() int
	//GetResCount resource count (add value for click the resource type item)
	GetResCount() int
	//GetFinalAwards final awards (the item disappeared, then the awards appear)
	GetFinalAwards() []int
	//GetClickAwards awards from click the item
	GetClickAwards() []int
	//GetBuiltDuration built the building (castle, rootmine also) duration
	GetBuiltDuration() int
	//GetTargetCastleId castle id for building
	GetTargetCastleId() int
	//GetExpaddForCastle exp value for building to upgrade the castle
	GetExpaddForCastle() int
	//GetMaxExpCurrentStar max exp value for current star
	GetMaxExpCurrentStar() int
	//GetNextItemId next item id by merger (0 means no next item)
	GetNextItemId() int
	//GetWorkerDuration time worker's max work time (s)
	GetWorkerDuration() int
	//IsUnique whether the item is unique. when it is, the item would appear in game only once (one)
	IsUnique() bool
	//GetTransformId when the item is unique, and it appeard now, all base item change to this id's item
	GetTransformId() int
	//IsEradicateAble whether then item could be eradicated
	IsEradicable() bool
}

type Itemset interface {
	//Get get IItem by given item id
	Get(id int32) Item
	//GetInType get all Item by given type
	GetInType(t ItemType) []Item
	//GetRelated get all related Item by given Item
	//example: a merger to b, b merger to c, c is final item. the method would return [a, b, c] by given a or b or c
	GetRelated(Item) []Item
}

type Itemdist interface {
	//GetId item id
	GetId() int32
	//GetIdcodes item distribution
	GetDistcodes() []uint64
}

// type ResData interface {
// 	GetId() int32
// 	GetCount() int32
// }

type Dataset interface {
	GetDists() []Itemdist
	GetStates() []State
	//plots unlock info
	GetUnlockcode() uint64
	GetResources() []Resource
}

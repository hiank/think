package play

import (
	"github.com/hiank/think/exp/easy"
	"github.com/hiank/think/run"
)

const (
	ErrExceedLimit = run.Err("exceed the limit")
)

// //DoneType
// type Operable uint8

// func (oa Operable) Moveable() bool {
// 	return oa&OAmove == OAmove
// }

// func (oa Operable) Clickable() bool {
// 	return oa&OAclick == OAclick
// }

// func (oa Operable) Eradicable() bool {
// 	return oa&OAeradicate == OAeradicate
// }

// const (
// 	OAmove      Operable = 1 << 0 //移动
// 	OAclick     Operable = 1 << 1 //点击
// 	OAeradicate Operable = 1 << 2 //铲除
// 	OAsqueezed  Operable = 1 << 3 //排挤(移到其它空位上)
// 	// OAbebuilt  Operable = 1 << 1 //被修建(建筑，城堡，矿山(最后一次开采后))//点击后，弹出移除操作界面
// 	// OAmined    Operable = 1 << 2 //被开采(类矿)//点击后，弹出移除操作界面
// 	// OApickedup Operable = 1 << 3 //被拾取(资源)//点击后，飞到资源栏，并改变资源数量
// 	// OAopen     Operable = 1 << 4 //打开(宝箱，袋子等)//点击后，弹出操作界面
// 	// OAharvest  Operable = 1 << 5 //收获(植物成熟后，类矿开采后，宝箱开启但未能完全释放奖励)//点击后，物品飞到空闲块中
// 	// OAshow     Operable = 1 << 6 //显示信息(工人，升级道具，城堡完全体之前)//点击后，弹出信息界面(工人剩余工时，道具操作说明，城堡)
// )

type StateBitagCover uint16

func (sbc StateBitagCover) Low() uint8 {
	return uint8(sbc)
}

func (sbc StateBitagCover) Count() uint8 {
	return uint8(sbc >> 8)
}

type StateBitag uint32

func (sb StateBitag) GetActionBitag() ActionBitag {
	return ActionBitag(sb)
}

func (sb StateBitag) GetCoverValue(sbc StateBitagCover) uint {
	return uint(easy.BitValue(sb, uint(sbc.Low()), uint(sbc.Count())))
}

// func (sb StateBitag) GetCoverValue(StateBitagCover) uint32 {

// }

// func stateBitagCover(low, cnt int) (cover StateBitag) {
// 	// cover = (1 << high) - 1
// 	if low+cnt > 32 {
// 		panic(ErrExceedLimit)
// 	}
// 	cover = ((1 << cnt) - 1) << low
// 	return
// }

const (
	// // SBnon StateBitag = 0 //空状态
	// SBcovercomponentcount StateBitag = ((1 << 3) - 1) << 8  //拥有的组件数量：城堡，需要4个组件用于激活
	// SBcoverbuildingdone   StateBitag = 1 << 11              //执行了修建：城堡需要执行修建完成激活
	// SBcoverremaintimes    StateBitag = ((1 << 4) - 1) << 12 //剩余次数：矿/母矿/农作物等，需要记录剩余执行次数
	// SBcoverstarcount      StateBitag = ((1 << 4) - 1) << 16 //星数：母矿/城堡的等级
	// SBcovertakeupworker   StateBitag = 1 << 20              //占用工人(通用工人)
	// SBcovernextworktype   StateBitag = ((1 << 2) - 1) << 21 //下一次执行的工作类型(3:修建 1:开采 2:收获 0:没有)：部分物品可执行的操作会转换(比如农作物最后一次收获后，转为开采)
	// SBcover
	SBCcountoffset    StateBitagCover = 8
	SBCactionbitag    StateBitagCover = 0 | (8 << SBCcountoffset)  //(0-7bit)ActionBitag数值
	SBCcomponentcount StateBitagCover = 8 | (3 << SBCcountoffset)  //(8-10bit)拥有的组件数量：城堡，需要4个组件用于激活
	SBCbuildingdone   StateBitagCover = 11 | (1 << SBCcountoffset) //(11bit)执行了修建：城堡需要执行修建完成激活
	SBCremaintimes    StateBitagCover = 12 | (4 << SBCcountoffset) //(12-15bit)剩余次数：矿/母矿/农作物等，需要记录剩余执行次数
	SBCstarcount      StateBitagCover = 16 | (4 << SBCcountoffset) //(16-19bit)星数：母矿/城堡的等级
	SBCtakeupworker   StateBitagCover = 20 | (1 << SBCcountoffset) //(20bit)占用工人(通用工人)
	SBCnextworktype   StateBitagCover = 21 | (2 << SBCcountoffset) //(21-22bit)下一次执行的工作类型(3:修建 1:开采 2:收获 0:没有)：部分物品可执行的操作会转换(比如农作物最后一次收获后，转为开采)
)

type ActionBitag uint8

const (
	ABnon         ActionBitag = 0                        //无动作：动作完成后，也置此标识
	ABwaiting     ActionBitag = 1 << 0                   //等待操作中
	ABtiming      ActionBitag = 1 << 1                   //计时中
	ABtagmine     ActionBitag = 1 << 2                   //开采标识
	ABtagbuild    ActionBitag = 1 << 3                   //修建标识
	ABtagharvest  ActionBitag = 1 << 4                   //收获标识
	ABwaitbuild   ActionBitag = ABwaiting | ABtagbuild   //等待修建
	ABbuilding    ActionBitag = ABtiming | ABtagbuild    //建造中
	ABwaitmine    ActionBitag = ABwaiting | ABtagmine    //等待开采
	ABmining      ActionBitag = ABtiming | ABtagmine     //开采中
	ABwaitharvest ActionBitag = ABwaiting | ABtagharvest //等待收获
	ABharvesting  ActionBitag = ABtiming | ABtagharvest  //收获中
)

// EasyBitag 方便快速确定特性(很多物品有相同特性)
type EasyBitag uint32

func (eb EasyBitag) Able(want EasyBitag) (suc bool) {
	wantType := ItemType(want)
	if suc = (want & eb) == want; suc && (wantType != ITundefined) {
		suc = ItemType(eb) == wantType
	}
	return
}

// // Universal 万能的，完全体水晶
// func (eb EasyBitag) Universal() bool {
// 	return eb.Able(EBfinal | EasyBitag(ITcrystal))
// }

const (
	EBnon EasyBitag = 0 //无物品：ItemType会填充到低8bit，有物品时一定不为0
	// EBmarkbuilding EasyBitag = 1 << 0	//标记修建：用于转换
	// EBmarkmining   EasyBitag = 1 << 1
	// EBmarkstate    EasyBitag = 1 << 2
	//固定属性，有则有无则无
	EBunique       EasyBitag = 1 << 8                                //全局唯一: 某个item id 只能存在唯一个
	EBeradicable   EasyBitag = 1 << 9                                //可被铲除: 被铲子移除
	EBpickable     EasyBitag = 1 << 10                               //可被拾取：加入背包
	EBunsqueezable EasyBitag = 1 << 11                               //不可挤压: 可移动物品无法放置到此处(无法重新位移，多块占用及无法移动的物品拥有此属性)
	EBunmoveable   EasyBitag = (1 << 12) | EBunsqueezable            //不可移动: 原始矿等(某些饰品?预留这个属性)
	EBsurround     EasyBitag = (1 << 13) | EBunmoveable | EBkeystate //围挡：参考 原始矿，必须移除后才能露出里面的资源(不可移动，否则就没有意义了)
	EBfinal        EasyBitag = 1 << 14                               //完全体：由此标识，一定不能被合并
	EBhold4block   EasyBitag = (1 << 15) | EBunsqueezable | EBunique //占用4地块
	EBopenable     EasyBitag = 1<<16 | EBkeystate                    //可开启标识：可开启一定需要状态信息(是否已开启，开启了一半)
	EBkeystate     EasyBitag = 1 << 17                               //需要状态信息
	EBkeymining    EasyBitag = 1 << 18                               //开采：母矿/城堡/完全体矿/固定矿/完全体农作物
	EBkeybuilding  EasyBitag = 1 << 19                               //修建：母矿/城堡/城堡材料
	EBuniversal    EasyBitag = EBfinal | EasyBitag(ITcrystal)        //万能的：完全体水晶，务必使用 Able 方法检测

	//状态，不同地块中的相同type物品可能有不同状态
	EBwaiting     EasyBitag = 1 << 20                  //等待操作中
	EBtiming      EasyBitag = 1 << 21                  //计时中
	EBtagmine     EasyBitag = 1 << 22                  //开采标识
	EBtagbuild    EasyBitag = 1 << 23                  //修建标识
	EBtagharvest  EasyBitag = 1 << 24                  //收获标识
	EBwaitbuild   EasyBitag = EBwaiting | EBtagbuild   //等待修建
	EBbuilding    EasyBitag = EBtiming | EBtagbuild    //建造中
	EBwaitmine    EasyBitag = EBwaiting | EBtagmine    //等待开采
	EBmining      EasyBitag = EBtiming | EBtagmine     //开采中
	EBwaitharvest EasyBitag = EBwaiting | EBtagharvest //等待收获
	EBharvesting  EasyBitag = EBtiming | EBtagharvest  //收获中

	// EBmarkstate    EasyBitag = 1 << 25 //标记状态需求
	// EBmarkbuilding EasyBitag = 1 << 26 //标记修建需求
	// EBmarkmining   EasyBitag = 1 << 27 //标记开采需求
	// //具体类型，整合了ItemType
	// EBcastle         EasyBitag = EasyBitag(ITcastle) | EBhold4block   //城堡
	// EBrootmine       EasyBitag = EasyBitag(ITrootmine) | EBhold4block //母矿
	// EBtimeworker     EasyBitag = EasyBitag(ITtimeworker)              //计时工人: 完全体可修建
	// EBcrop           EasyBitag = EasyBitag(ITcrop) |                    //农作物
	// EBcastlematerial EasyBitag = EasyBitag(ITcastlematerial)          //建筑(用于升级城堡)
	// EBmine           EasyBitag = EasyBitag(ITmine)                    //矿
	// EBlocalmine      EasyBitag = EasyBitag(ITlocalmine)               //固定位矿(未解锁时填充矿)
	// EBcrystal        EasyBitag = EasyBitag(ITcrystal)                 //水晶
	// EBcook           EasyBitag = EasyBitag(ITcook)                    //厨子
	// EBcoin           EasyBitag = EasyBitag(ITcoin)                    //金币
	// EBdiamond        EasyBitag = EasyBitag(ITdiamond)                 //钻石
	// EBenergy         EasyBitag = EasyBitag(ITenergy)                  //闪电
	// EBmagic          EasyBitag = EasyBitag(ITmagic)                   //魔法棒，用于开始新地块
	// EBfruit          EasyBitag = EasyBitag(ITfruit)                   //果实，用于厨子制作甜品
	// EBbox            EasyBitag = EasyBitag(ITbox)                     //箱子，袋子
	// EBanimal         EasyBitag = EasyBitag(ITanimal)                  //动物
	// EBfree           EasyBitag = EasyBitag(ITfree)                    //自由物品，除了移动(或回收)无法操作，相当于饰品

	// EBuniversal EasyBitag = EB
	// EBkindcastlematerial EasyBitag = 1 << 8                               //城堡材料：用于升级城堡，需要状态数据(是否已激活)
	// EBkindcastle         EasyBitag = (1 << 9) | EBhold4block              //城堡类：会有多种状态，需要特别处理
	// EBkindrootmine       EasyBitag = 1 << 10                              //母矿类：会有多种状态，需要特别处理

	// //非固定属性，受 开采/修建 影响
	// EBcastlematerialactive EasyBitag = 1 << 11 //城堡材料激活标识：有些材料需要修建完成才能使用
	// EBwantopen
	// EBopenable EasyBitag = 1 << 8 //可被开启：城堡/宝箱袋子/矿，收获产出
	// // IBbuildworker              = 1 << 11
	// // IBsuperworker              = (1 << 12) | IBbuildworker
	// // IBkindcastle   ItemBitag = (1 << 13) | IBhold4block | IBfinal //城堡类：城堡有多种状态，需要特别标记
	// // IBkindrootmine ItemBitag = (1 << 14) | IBhold4block | IBfinal //母矿类：母矿有多种状态，需要特别标记
	// EBkindbox EasyBitag = (1 << 15) | EBopenable //宝箱宝袋类：因为矿/城堡也可能有打开状态，需要特别标识

)

// type ItemBitag uint

// // //Moveable 是否可移动
// // func (ib ItemBitag) Moveable() bool {
// // 	return (ib & IBunmoveable) != IBunmoveable
// // }

// // //Eradicable 是否可被移除
// // func (ib ItemBitag) Eradicable() bool {
// // 	return (ib & IBeradicable) == IBeradicable
// // }

// func (ib ItemBitag) Able(want ItemBitag) bool {
// 	return (ib & want) == want
// }

// func (ib ItemBitag) Type() ItemType {
// 	return ItemType(ib)
// }

// func (ib ItemBitag) ActionBitag() ActionBitag {
// 	return ActionBitag(ib >> 8)
// }

// func ResetItemBitag(ib ItemBitag, low, cnt int) ItemBitag {
// 	return run.ResetBit(ib, 0, low, cnt)
// }

// const (
// 	IBnon ItemBitag = 0 //无任何特性
// 	// ibtypecover ItemBitag = math.MaxUint8 //0-7bit: 用于记录ItemType值
// 	// IBkind

// 	// //以下为特性
// 	// // IBpickable ItemBitag = 1 << 0 //可被拾取(可入背包): 包括金币，钻石，闪电(用于解锁物品)，魔法(用于解锁地块)，各类果实等
// 	// // IBunsqueezable ItemBitag = 1 << 1                               //不可挤压: 可移动物品无法放置到此处(无法重新位移，多块占用及无法移动的物品拥有此属性)
// 	// // IBunique       ItemBitag = 1 << 2                               //全局唯一: 某个item id 只能存在唯一个
// 	// // IBunmoveable   ItemBitag = (1 << 3) | IBunsqueezable            //不可移动: 原始矿等(某些饰品?预留这个属性)
// 	// // IBeradicable   ItemBitag = 1 << 4                               //可被铲除: 被铲子移除
// 	// // IBfinal      ItemBitag = 1 << 5                               //完全体(无法再合并)，大部分的附加操作需要是最终体
// 	// IBbuildable  ItemBitag = 1 << 6                               //可被修建: 修建建筑，通用工人/计时工人都可以执行修建
// 	// IBmineable   ItemBitag = (1 << 7) | IBfinal                   //可被开采: 通用工人，可被开采的一定是完全体
// 	// IBopenable   ItemBitag = 1 << 8                               //可被打开: 宝箱，袋子等
// 	// IBhold4block ItemBitag = (1 << 9) | IBunsqueezable | IBunique //占用4个块: 城堡，母矿等. 已知此类物品都是全局唯一的
// 	// // IBsurround   ItemBitag = (1 << 10) | IBunmoveable             //围挡：参考 原始矿，必须移除后才能露出里面的资源(不可移动，否则就没有意义了)
// 	// IBworkable ItemBitag = 1 << 11 //可工作(最终体计时工人)
// 	// // IBworker ItemBitag =

// 	//固定属性，有则有无则无
// 	IBunique         ItemBitag = 1 << 0                               //全局唯一: 某个item id 只能存在唯一个
// 	IBeradicable     ItemBitag = 1 << 1                               //可被铲除: 被铲子移除
// 	IBpickable       ItemBitag = 1 << 2                               //可被拾取：加入背包
// 	IBunsqueezable   ItemBitag = 1 << 3                               //不可挤压: 可移动物品无法放置到此处(无法重新位移，多块占用及无法移动的物品拥有此属性)
// 	IBunmoveable     ItemBitag = (1 << 4) | IBunsqueezable            //不可移动: 原始矿等(某些饰品?预留这个属性)
// 	IBsurround       ItemBitag = (1 << 5) | IBunmoveable              //围挡：参考 原始矿，必须移除后才能露出里面的资源(不可移动，否则就没有意义了)
// 	IBfinal          ItemBitag = 1 << 6                               //完全体：由此标识，一定不能被合并
// 	IBhold4block     ItemBitag = (1 << 7) | IBunsqueezable | IBunique //占用4地块
// 	IBcastlematerial ItemBitag = 1 << 10                              //城堡材料：用于升级城堡，需要状态数据(是否已激活)
// 	IBkindcastle     ItemBitag = 1 << 11                              //城堡类：会有多种状态，需要特别处理
// 	IBkindrootmine   ItemBitag = 1 << 12                              //母矿类：会有多种状态，需要特别处理

// 	//非固定属性，受 开采/修建 影响
// 	IBcastlematerialactive ItemBitag = 1 << 11 //城堡材料激活标识：有些材料需要修建完成才能使用
// 	IBwantopen
// 	IBopenable ItemBitag = 1 << 8 //可被开启：城堡/宝箱袋子/矿，收获产出
// 	// IBbuildworker              = 1 << 11
// 	// IBsuperworker              = (1 << 12) | IBbuildworker
// 	// IBkindcastle   ItemBitag = (1 << 13) | IBhold4block | IBfinal //城堡类：城堡有多种状态，需要特别标记
// 	// IBkindrootmine ItemBitag = (1 << 14) | IBhold4block | IBfinal //母矿类：母矿有多种状态，需要特别标记
// 	IBkindbox ItemBitag = (1 << 15) | IBopenable //宝箱宝袋类：因为矿/城堡也可能有打开状态，需要特别标识

// 	// IBbuildworkable ItemBitag = IBbuildworker | IBfinal //可执行修建：计时工人
// 	// IBcastlematerialusable ItemBitag = IBcastlematerial |
// 	// IBfinalcrop     ItemBitag = IBkindcrop | IBneedstate | IBfinal
// 	// IBfinalmine     ItemBitag = IBkindmine | IBneedstate | IBfinal
// 	// IBfinalbuilding ItemBitag = IBkindbuilding | IBneedstate | IBfinal
// 	// IBlocalmine     ItemBitag = IBfinalmine | IBsurround
// 	// IBfinal
// 	// IBworker ItemBitag = 1 << 10	//工人()
// 	// IBisbuilding   ItemBitag = 1 << 9  //是建筑
// 	// IBismine       ItemBitag = 1 << 10 //是矿
// 	// IBfocustag ItemBitag = 1 << 1 //标记激活：
// 	// IBtagbuild ItemBitag = 1 << 2 //建造标志
// 	// IBtagmine  ItemBitag = 1 << 2 //开采标志
// 	// IBtagopen  ItemBitag = 1 << 2 //开启标志

// 	// IBcastle   ItemBitag = IBhold4block | IBfinal
// 	// IBrootmine ItemBitag = IBhold4block | IBfinal	//母矿

// 	// IBbuildworkable ItemBitag = IBexecutor | IBtagbuild | IBable | IBfinal //可执行修建
// 	// IBmineworkable  ItemBitag = IBexecutor | IBtagmine | IBable | IBfinal  //可执行开采：目前只有系统的通用工人可以执行开采
// 	// IBbebuiltable   ItemBitag = IBtagbuild | IBable                        //可被修建
// 	// IBminedable     ItemBitag = IBtagmine | IBable                         //可被开采
// 	// IB
// 	//以下为具体物品
// 	// IBkindcastle     ItemBitag =             //城堡
// 	// IBkindrootmine   ItemBitag = (1 << 12) | IBhold4block | IBmineable         //母矿
// 	// IBkindtimeworker ItemBitag = 1 << 13                                       //计时工人: 完全体可修建
// 	// IBkindcrop       ItemBitag = 1 << 14                                       //农作物
// 	// IBkindbuilding   ItemBitag = 1 << 15                                       //建筑(用于升级城堡)
// 	// IBkindmine       ItemBitag = 1 << 16                                       //矿
// 	// IBkindlocalmine  ItemBitag = (1 << 16) | IBsurround | IBmineable | IBfinal //原始矿(未解锁时填充的矿)，概念上也是矿，围挡属性的矿，并且一定是最终体
// 	// IBkindanimal     ItemBitag = 1 << 18                                       //动物
// 	// IBkindcrystal    ItemBitag = 1 << 19                                       //水晶
// 	// IBkindcook       ItemBitag = 1 << 20                                       //厨子
// 	// IBkindcoin       ItemBitag = (1 << 21) | IBpickable                        //金币
// 	// IBkinddiamond    ItemBitag = (1 << 22) | IBpickable                        //钻石
// 	// IBkindenergy     ItemBitag = (1 << 23) | IBpickable                        //闪电
// 	// IBkindmagic      ItemBitag = (1 << 24) | IBpickable                        //魔法棒，用于开始新地块
// 	// IBkindfruit      ItemBitag = (1 << 25) | IBpickable | IBfinal              //果实，用于厨子制作甜品
// 	// IBkindbox        ItemBitag = (1 << 26) | IBopenable                        //箱子，袋子
// 	// IBkindfree       ItemBitag = 1 << 27                                       //自由物品，除了移动(或回收)无法操作，相当于饰品
// 	// out.COMMON = (1 << 26);      //NOTE: 通用物品 不参与地图上的合并或者建造操作
// 	// out.LIGHTNING_RECOVER = 1 << 27; //NOTE: 体力恢复，用于将上次一次体力变化时刻(当时的时间戳)记录在ResStatus中
// 	// out.HOME_ITEM = (1 << 28) | out.MERGER;      //NOTE: 家园地图合成物品
// 	// out.LIGHTNING_ITEM = 1 << 29;      //NOTE: 体力使用道具 指定时间内无限使用体力
// )

// type StateBitag uint

// const (
// 	ABnotready     StateBitag = 1 << 0                //未准备好(主要用于合并判断，有这个标志表明当前无法合并)
// 	SBwait         StateBitag = (1 << 1) | SBnotready //等待操作，比如等待开采，等待收获，等待修建等
// 	SBtiming       StateBitag = (1 << 2) | SBnotready //定时中(修建)
// 	SBactionbuild  StateBitag = 1 << 3                //修建，配合SBwait/SBtiming 使用
// 	SBactionmine   StateBitag = 1 << 4                //开采，配合SBwait/SBtiming 使用
// 	SBactionreward StateBitag = 1 << 5                //收获(农作物，箱子，袋子，城堡，矿等)，配合SBwait 使用
// )

// type item struct {
// 	IItem
// 	bit  int      //NOTE: type value
// 	dist sync.Map //NOTE: item distribution
// }

// func newItem(cfg IItem) *item {
// 	it := &item{IItem: cfg, bit: cfg.GetType()}
// 	if cfg.GetNextItemId() == 0 {
// 		it.bit |= ITFinal
// 	}
// 	if cfg.IsUnique() {
// 		it.bit |= ITUnique
// 	}
// 	if cfg.IsEradicable() {
// 		it.bit |= ItemEradicable
// 	}
// 	return it
// }

// //GetType rewrite IItem's GetType method
// func (it *item) GetType() int {
// 	return it.bit
// }

// //distSlice get distribution slice
// func (it *item) distSlice() []idcode {
// 	out := make([]idcode, 0, 16)
// 	it.dist.Range(func(key, value interface{}) bool {
// 		out = append(out, key.(idcode))
// 		return true
// 	})
// 	return out
// }

// func (it *item) place(id idcode) {
// 	_, ok := it.dist.LoadOrStore(id, byte(1))
// 	if !ok {
// 		//NOTE: data is added here. should do something under
// 	}
// }

// func (it *item) delete(id idcode) {
// 	_, ok := it.dist.LoadAndDelete(id)
// 	if ok {
// 		//NOTE: data is deleted here. should do something under
// 	}
// }

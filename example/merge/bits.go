package merge

type EasyBitag uint32

const (
	EBempty EasyBitag = 0

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
)

type StateBitag uint32

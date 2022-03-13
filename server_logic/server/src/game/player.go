package game

import (
	"fmt"
	"os"
	"server/csvs"
)

const (
	TASK_STATE_INIT   = 0
	TASK_STATE_DOING  = 1
	TASK_STATE_FINISH = 2
)

type Player struct {
	ModPlayer     *ModPlayer
	ModIcon       *ModIcon
	ModCard       *ModCard
	ModUniqueTask *ModUniqueTask
	ModRole       *ModRole
	ModBag        *ModBag
	ModWeapon     *ModWeapon
	ModRelics     *ModRelics
	ModCook       *ModCook
	ModHome       *ModHome
	ModPool       *ModPool
	ModMap        *ModMap

	modManage map[string]interface{}
}

func NewTestPlayer() *Player {
	player := new(Player)
	player.ModPlayer = new(ModPlayer)
	player.ModIcon = new(ModIcon)
	player.ModIcon.IconInfo = make(map[int]*Icon)
	player.ModCard = new(ModCard)
	player.ModCard.CardInfo = make(map[int]*Card)
	player.ModUniqueTask = new(ModUniqueTask)
	player.ModUniqueTask.MyTaskInfo = make(map[int]*TaskInfo)
	//player.ModUniqueTask.Locker = new(sync.RWMutex)
	player.ModRole = new(ModRole)
	player.ModRole.RoleInfo = make(map[int]*RoleInfo)
	player.ModBag = new(ModBag)
	player.ModBag.BagInfo = make(map[int]*ItemInfo)
	player.ModWeapon = new(ModWeapon)
	player.ModWeapon.WeaponInfo = make(map[int]*Weapon)
	player.ModRelics = new(ModRelics)
	player.ModRelics.RelicsInfo = make(map[int]*Relics)
	player.ModCook = new(ModCook)
	player.ModCook.CookInfo = make(map[int]*Cook)
	player.ModHome = new(ModHome)
	player.ModHome.HomeItemIdInfo = make(map[int]*HomeItemId)
	player.ModPool = new(ModPool)
	player.ModPool.UpPoolInfo = new(PoolInfo)
	player.ModMap = new(ModMap)
	player.ModMap.InitData()
	//****************************************
	player.ModPlayer.PlayerLevel = 1
	player.ModPlayer.Name = "旅行者"
	player.ModPlayer.WorldLevel = 1
	player.ModPlayer.WorldLevelNow = 1
	//****************************************
	player.ModPlayer.UserId = 10000666
	player.InitData()
	return player
}

func (self *Player) InitData() {
	path := GetServer().Config.LocalSavePath
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return
		}
	}

	selfPath := path + fmt.Sprintf("/%d", self.ModPlayer.UserId)
	_, err = os.Stat(selfPath)
	if err != nil {
		err = os.Mkdir(selfPath, os.ModePerm)
		if err != nil {
			return
		}
		self.ModPlayer.SaveData(selfPath+"/player.json")
	}else{
		self.ModPlayer.LoadData(selfPath+"/player.json")
	}
}

//对外接口
func (self *Player) RecvSetIcon(iconId int) {
	self.ModPlayer.SetIcon(iconId, self)
}

func (self *Player) RecvSetCard(cardId int) {
	self.ModPlayer.SetCard(cardId, self)
}

func (self *Player) RecvSetName(name string) {
	self.ModPlayer.SetName(name, self)
}

func (self *Player) RecvSetSign(sign string) {
	self.ModPlayer.SetSign(sign, self)
}

func (self *Player) ReduceWorldLevel() {
	self.ModPlayer.ReduceWorldLevel(self)
}

func (self *Player) ReturnWorldLevel() {
	self.ModPlayer.ReturnWorldLevel(self)
}

func (self *Player) SetBirth(birth int) {
	self.ModPlayer.SetBirth(birth, self)
}

func (self *Player) SetShowCard(showCard []int) {
	self.ModPlayer.SetShowCard(showCard, self)
}

func (self *Player) SetShowTeam(showRole []int) {
	self.ModPlayer.SetShowTeam(showRole, self)
}

func (self *Player) SetHideShowTeam(isHide int) {
	self.ModPlayer.SetHideShowTeam(isHide, self)
}

func (self *Player) SetEventState(state int) {
	//self.ModMap.SetEventState(state, self)
}

func (self *Player) Run() {
	fmt.Println("从0开始写原神服务器------测试工具v0.7")
	fmt.Println("作者:B站------golang大海葵")
	fmt.Println("模拟用户创建成功OK------开始测试")
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	for {
		fmt.Println(self.ModPlayer.Name, ",欢迎来到提瓦特大陆,请选择功能：1基础信息2背包3角色(八重神子UP池)4地图5圣遗物6角色7武器8存储数据")
		var modChoose int
		fmt.Scan(&modChoose)
		switch modChoose {
		case 1:
			self.HandleBase()
		case 2:
			self.HandleBag()
		case 3:
			self.HandlePool()
		case 4:
			self.HandleMap()
		case 5:
			self.HandleRelics()
		case 6:
			self.HandleRole()
		case 7:
			self.HandleWeapon()
		case 8:
			self.ModPlayer.SaveData(fmt.Sprintf("./save/%d/player.json",self.ModPlayer.UserId))
		}
	}
}

//基础信息
func (self *Player) HandleBase() {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作：0返回1查询信息2设置名字3设置签名4头像5名片6设置生日")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBaseGetInfo()
		case 2:
			self.HandleBagSetName()
		case 3:
			self.HandleBagSetSign()
		case 4:
			self.HandleBagSetIcon()
		case 5:
			self.HandleBagSetCard()
		case 6:
			self.HandleBagSetBirth()
		}
	}
}

func (self *Player) HandleBaseGetInfo() {
	fmt.Println("名字:", self.ModPlayer.Name)
	fmt.Println("等级:", self.ModPlayer.PlayerLevel)
	fmt.Println("大世界等级:", self.ModPlayer.WorldLevelNow)
	if self.ModPlayer.Sign == "" {
		fmt.Println("签名:", "未设置")
	} else {
		fmt.Println("签名:", self.ModPlayer.Sign)
	}

	if self.ModPlayer.Icon == 0 {
		fmt.Println("头像:", "未设置")
	} else {
		fmt.Println("头像:", csvs.GetItemConfig(self.ModPlayer.Icon), self.ModPlayer.Icon)
	}

	if self.ModPlayer.Card == 0 {
		fmt.Println("名片:", "未设置")
	} else {
		fmt.Println("名片:", csvs.GetItemConfig(self.ModPlayer.Card), self.ModPlayer.Card)
	}

	if self.ModPlayer.Birth == 0 {
		fmt.Println("生日:", "未设置")
	} else {
		fmt.Println("生日:", self.ModPlayer.Birth/100, "月", self.ModPlayer.Birth%100, "日")
	}
}

func (self *Player) HandleBagSetName() {
	fmt.Println("请输入名字:")
	var name string
	fmt.Scan(&name)
	self.RecvSetName(name)
}

func (self *Player) HandleBagSetSign() {
	fmt.Println("请输入签名:")
	var sign string
	fmt.Scan(&sign)
	self.RecvSetSign(sign)
}

func (self *Player) HandleBagSetIcon() {
	for {
		fmt.Println("当前处于基础信息--头像界面,请选择操作：0返回1查询头像背包2设置头像")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBagSetIconGetInfo()
		case 2:
			self.HandleBagSetIconSet()
		}
	}
}

func (self *Player) HandleBagSetIconGetInfo() {
	fmt.Println("当前拥有头像如下:")
	for _, v := range self.ModIcon.IconInfo {
		config := csvs.GetItemConfig(v.IconId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (self *Player) HandleBagSetIconSet() {
	fmt.Println("请输入头像id:")
	var icon int
	fmt.Scan(&icon)
	self.RecvSetIcon(icon)
}

func (self *Player) HandleBagSetCard() {
	for {
		fmt.Println("当前处于基础信息--名片界面,请选择操作：0返回1查询名片背包2设置名片")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBagSetCardGetInfo()
		case 2:
			self.HandleBagSetCardSet()
		}
	}
}

func (self *Player) HandleBagSetCardGetInfo() {
	fmt.Println("当前拥有名片如下:")
	for _, v := range self.ModCard.CardInfo {
		config := csvs.GetItemConfig(v.CardId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (self *Player) HandleBagSetCardSet() {
	fmt.Println("请输入名片id:")
	var card int
	fmt.Scan(&card)
	self.RecvSetCard(card)
}

func (self *Player) HandleBagSetBirth() {
	if self.ModPlayer.Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}
	fmt.Println("生日只能设置一次，请慎重填写,输入月:")
	var month, day int
	fmt.Scan(&month)
	fmt.Println("请输入日:")
	fmt.Scan(&day)
	self.ModPlayer.SetBirth(month*100+day, self)
}

//背包
func (self *Player) HandleBag() {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作：0返回1增加物品2扣除物品3使用物品4升级七天神像(风)")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBagAddItem()
		case 2:
			self.HandleBagRemoveItem()
		case 3:
			self.HandleBagUseItem()
		case 4:
			self.HandleBagWindStatue()
		}
	}
}

//抽卡
func (self *Player) HandlePool() {
	for {
		fmt.Println("当前处于模拟抽卡界面,请选择操作：0返回1角色信息2十连抽(入包)3单抽(可选次数,入包)" +
			"4五星爆点测试5十连多黄测试6视频原版函数(30秒)7单抽(仓检版,独宠一人)8单抽(仓检版,雨露均沾)")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.ModRole.HandleSendRoleInfo(self)
		case 2:
			self.ModPool.HandleUpPoolTen(self)
		case 3:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolSingle(times, self)
		case 4:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolTimesTest(times)
		case 5:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolFiveTest(times)
		case 6:
			self.ModPool.DoUpPool()
		case 7:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolSingleCheck1(times, self)
		case 8:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolSingleCheck2(times, self)
		}
	}
}

func (self *Player) HandleBagAddItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.ModBag.AddItem(itemId, int64(itemNum), self)
}

func (self *Player) HandleBagRemoveItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.ModBag.RemoveItemToBag(itemId, int64(itemNum), self)
}

func (self *Player) HandleBagUseItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.ModBag.UseItem(itemId, int64(itemNum), self)
}

func (self *Player) HandleBagWindStatue() {
	fmt.Println("开始升级七天神像")
	self.ModMap.UpStatue(1, self)
	self.ModRole.CalHpPool()
}

//地图
func (self *Player) HandleMap() {
	fmt.Println("向着星辰与深渊,欢迎来到冒险家协会！")
	for {
		fmt.Println("请选择互动地图1蒙德2璃月1001深入风龙废墟2001无妄引咎密宫")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		default:
			self.HandleMapIn(action)
		}
	}
}

func (self *Player) HandleMapIn(mapId int) {

	config := csvs.ConfigMapMap[mapId]
	if config == nil {
		fmt.Println("无法识别的地图")
		return
	}
	self.ModMap.RefreshByPlayer(mapId)
	for {
		self.ModMap.GetEventList(config)
		fmt.Println("请选择触发事件Id(0返回)")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		default:
			eventConfig := csvs.ConfigMapEventMap[action]
			if eventConfig == nil {
				fmt.Println("无法识别的事件")
				break
			}
			self.ModMap.SetEventState(mapId, eventConfig.EventId, csvs.EVENT_END, self)
		}
	}
}

func (self *Player) HandleRelics() {
	for {
		fmt.Println("当前处于圣遗物界面，选择功能0返回1强化测试2满级圣遗物3极品头测试")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.ModRelics.RelicsUp(self)
		case 2:
			self.ModRelics.RelicsTop(self)
		case 3:
			self.ModRelics.RelicsTestBest(self)
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func (self *Player) HandleRole() {
	for {
		fmt.Println("当前处于角色界面，选择功能0返回1查询2穿戴圣遗物3卸下圣遗物4穿戴武器5卸下武器")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.ModRole.HandleSendRoleInfo(self)
		case 2:
			self.HandleWearRelics()
		case 3:
			self.HandleTakeOffRelics()
		case 4:
			self.HandleWearWeapon()
		case 5:
			self.HandleTakeOffWeapon()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func (self *Player) HandleWearRelics() {
	for {
		fmt.Println("输入操作的目标英雄Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := self.ModRole.RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("英雄不存在")
			continue
		}

		RoleInfo.ShowInfo(self)
		fmt.Println("输入需要穿戴的圣遗物key:,0返回")
		var relicsKey int
		fmt.Scan(&relicsKey)
		if relicsKey == 0 {
			return
		}
		relics := self.ModRelics.RelicsInfo[relicsKey]
		if relics == nil {
			fmt.Println("圣遗物不存在")
			continue
		}
		self.ModRole.WearRelics(RoleInfo, relics, self)
	}
}

func (self *Player) HandleTakeOffRelics() {
	for {
		fmt.Println("输入操作的目标英雄Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := self.ModRole.RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("英雄不存在")
			continue
		}

		RoleInfo.ShowInfo(self)
		fmt.Println("输入需要卸下的圣遗物key:,0返回")
		var relicsKey int
		fmt.Scan(&relicsKey)
		if relicsKey == 0 {
			return
		}
		relics := self.ModRelics.RelicsInfo[relicsKey]
		if relics == nil {
			fmt.Println("圣遗物不存在")
			continue
		}
		self.ModRole.TakeOffRelics(RoleInfo, relics, self)
	}
}

func (self *Player) HandleWeapon() {
	for {
		fmt.Println("当前处于武器界面，选择功能0返回1强化测试2突破测试3精炼测试")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleWeaponUp()
		case 2:
			self.HandleWeaponStarUp()
		case 3:
			self.HandleWeaponRefineUp()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func (self *Player) HandleWeaponUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range self.ModWeapon.WeaponInfo {
			fmt.Println(fmt.Sprintf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d",
				v.KeyId, v.Level, v.StarLevel, v.RefineLevel))
		}
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		self.ModWeapon.WeaponUp(weaponKeyId, self)
	}
}

func (self *Player) HandleWeaponStarUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range self.ModWeapon.WeaponInfo {
			fmt.Println(fmt.Sprintf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d",
				v.KeyId, v.Level, v.StarLevel, v.RefineLevel))
		}
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		self.ModWeapon.WeaponUpStar(weaponKeyId, self)
	}
}

func (self *Player) HandleWeaponRefineUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range self.ModWeapon.WeaponInfo {
			fmt.Println(fmt.Sprintf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d",
				v.KeyId, v.Level, v.StarLevel, v.RefineLevel))
		}
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		for {
			fmt.Println("输入作为材料的武器keyId:,0返回")
			var weaponTargetKeyId int
			fmt.Scan(&weaponTargetKeyId)
			if weaponTargetKeyId == 0 {
				return
			}
			self.ModWeapon.WeaponUpRefine(weaponKeyId, weaponTargetKeyId, self)
		}
	}
}

func (self *Player) HandleWearWeapon() {
	for {
		fmt.Println("输入操作的目标英雄Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := self.ModRole.RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("英雄不存在")
			continue
		}

		RoleInfo.ShowInfo(self)
		fmt.Println("输入需要穿戴的武器key:,0返回")
		var weaponKey int
		fmt.Scan(&weaponKey)
		if weaponKey == 0 {
			return
		}
		weaponInfo := self.ModWeapon.WeaponInfo[weaponKey]
		if weaponInfo == nil {
			fmt.Println("武器不存在")
			continue
		}
		self.ModRole.WearWeapon(RoleInfo, weaponInfo, self)
		RoleInfo.ShowInfo(self)
	}
}

func (self *Player) HandleTakeOffWeapon() {
	for {
		fmt.Println("输入操作的目标英雄Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := self.ModRole.RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("英雄不存在")
			continue
		}

		RoleInfo.ShowInfo(self)
		fmt.Println("输入需要卸下的武器key:,0返回")
		var weaponKey int
		fmt.Scan(&weaponKey)
		if weaponKey == 0 {
			return
		}
		weapon := self.ModWeapon.WeaponInfo[weaponKey]
		if weapon == nil {
			fmt.Println("武器不存在")
			continue
		}
		self.ModRole.TakeOffWeapon(RoleInfo, weapon, self)
		RoleInfo.ShowInfo(self)
	}
}

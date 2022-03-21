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

const (
	MOD_PLAYER     = "player"
	MOD_ICON       = "icon"
	MOD_CARD       = "card"
	MOD_UNIQUETASK = "uniquetask"
	MOD_ROLE       = "role"
	MOD_BAG        = "bag"
	MOD_WEAPON     = "weapon "
	MOD_RELICS     = "relics"
	MOD_COOK       = "cook"
	MOD_HOME       = "home"
	MOD_POOL       = "pool"
	MOD_MAP        = "map"
)

type ModBase interface {
	LoadData(player *Player)
	SaveData()
	InitData()
}

type Player struct {
	UserId    int64
	modManage map[string]ModBase
	localPath string
}

func NewTestPlayer(userid int64) *Player {
	//***************泛型架构***************************
	player := new(Player)
	player.UserId = userid
	player.modManage = map[string]ModBase{
		MOD_PLAYER:     new(ModPlayer),
		MOD_ICON:       new(ModIcon),
		MOD_CARD:       new(ModCard),
		MOD_UNIQUETASK: new(ModUniqueTask),
		MOD_ROLE:       new(ModRole),
		MOD_BAG:        new(ModBag),
		MOD_WEAPON:     new(ModWeapon),
		MOD_RELICS:     new(ModRelics),
		MOD_COOK:       new(ModCook),
		MOD_HOME:       new(ModHome),
		MOD_POOL:       new(ModPool),
		MOD_MAP:        new(ModMap),
	}
	player.InitData()
	player.InitMod()
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
	self.localPath = path + fmt.Sprintf("/%d", self.UserId)
	_, err = os.Stat(self.localPath)
	if err != nil {
		err = os.Mkdir(self.localPath, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func (self *Player) InitMod() {
	for _, v := range self.modManage {
		v.LoadData(self)
	}
}

//对外接口
func (self *Player) RecvSetIcon(iconId int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetIcon(iconId)
}

func (self *Player) RecvSetCard(cardId int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetCard(cardId)
}

func (self *Player) RecvSetName(name string) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetName(name)
}

func (self *Player) RecvSetSign(sign string) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetSign(sign)
}

func (self *Player) ReduceWorldLevel() {
	self.GetMod(MOD_PLAYER).(*ModPlayer).ReduceWorldLevel()
}

func (self *Player) ReturnWorldLevel() {
	self.GetMod(MOD_PLAYER).(*ModPlayer).ReturnWorldLevel()
}

func (self *Player) SetBirth(birth int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetBirth(birth)
}

func (self *Player) SetShowCard(showCard []int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetShowCard(showCard, self)
}

func (self *Player) SetShowTeam(showRole []int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetShowTeam(showRole, self)
}

func (self *Player) SetHideShowTeam(isHide int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetHideShowTeam(isHide, self)
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
		fmt.Println(self.GetMod(MOD_PLAYER).(*ModPlayer).Name, ",欢迎来到提瓦特大陆,请选择功能：1基础信息2背包3角色(八重神子UP池)4地图5圣遗物6角色7武器8存储数据")
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
			for _,v:=range self.modManage{
				v.SaveData()
			}
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
	fmt.Println("名字:", self.GetMod(MOD_PLAYER).(*ModPlayer).Name)
	fmt.Println("等级:", self.GetMod(MOD_PLAYER).(*ModPlayer).PlayerLevel)
	fmt.Println("大世界等级:", self.GetMod(MOD_PLAYER).(*ModPlayer).WorldLevelNow)
	if self.GetMod(MOD_PLAYER).(*ModPlayer).Sign == "" {
		fmt.Println("签名:", "未设置")
	} else {
		fmt.Println("签名:", self.GetMod(MOD_PLAYER).(*ModPlayer).Sign)
	}

	if self.GetMod(MOD_PLAYER).(*ModPlayer).Icon == 0 {
		fmt.Println("头像:", "未设置")
	} else {
		fmt.Println("头像:", csvs.GetItemConfig(self.GetMod(MOD_PLAYER).(*ModPlayer).Icon), self.GetMod(MOD_PLAYER).(*ModPlayer).Icon)
	}

	if self.GetMod(MOD_PLAYER).(*ModPlayer).Card == 0 {
		fmt.Println("名片:", "未设置")
	} else {
		fmt.Println("名片:", csvs.GetItemConfig(self.GetMod(MOD_PLAYER).(*ModPlayer).Card), self.GetMod(MOD_PLAYER).(*ModPlayer).Card)
	}

	if self.GetMod(MOD_PLAYER).(*ModPlayer).Birth == 0 {
		fmt.Println("生日:", "未设置")
	} else {
		fmt.Println("生日:", self.GetMod(MOD_PLAYER).(*ModPlayer).Birth/100, "月", self.GetMod(MOD_PLAYER).(*ModPlayer).Birth%100, "日")
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
	for _, v := range self.GetModIcon().IconInfo {
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
	for _, v := range self.GetModCard().CardInfo {
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
	if self.GetMod(MOD_PLAYER).(*ModPlayer).Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}
	fmt.Println("生日只能设置一次，请慎重填写,输入月:")
	var month, day int
	fmt.Scan(&month)
	fmt.Println("请输入日:")
	fmt.Scan(&day)
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetBirth(month*100 + day)
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
			self.GetModRole().HandleSendRoleInfo(self)
		case 2:
			self.GetModPool().HandleUpPoolTen(self)
		case 3:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.GetModPool().HandleUpPoolSingle(times, self)
		case 4:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.GetModPool().HandleUpPoolTimesTest(times)
		case 5:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.GetModPool().HandleUpPoolFiveTest(times)
		case 6:
			self.GetModPool().DoUpPool()
		case 7:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.GetModPool().HandleUpPoolSingleCheck1(times, self)
		case 8:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.GetModPool().HandleUpPoolSingleCheck2(times, self)
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
	self.GetModBag().AddItem(itemId, int64(itemNum))
}

func (self *Player) HandleBagRemoveItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.GetModBag().RemoveItemToBag(itemId, int64(itemNum))
}

func (self *Player) HandleBagUseItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.GetModBag().UseItem(itemId, int64(itemNum))
}

func (self *Player) HandleBagWindStatue() {
	fmt.Println("开始升级七天神像")
	self.GetModMap().UpStatue(1)
	self.GetModRole().CalHpPool()
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
	self.GetModMap().RefreshByPlayer(mapId)
	for {
		self.GetModMap().GetEventList(config)
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
			self.GetModMap().SetEventState(mapId, eventConfig.EventId, csvs.EVENT_END, self)
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
			self.GetModRelics().RelicsUp(self)
		case 2:
			self.GetModRelics().RelicsTop(self)
		case 3:
			self.GetModRelics().RelicsTestBest(self)
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
			self.GetModRole().HandleSendRoleInfo(self)
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

		RoleInfo := self.GetModRole().RoleInfo[roleId]
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
		relics := self.GetModRelics().RelicsInfo[relicsKey]
		if relics == nil {
			fmt.Println("圣遗物不存在")
			continue
		}
		self.GetModRole().WearRelics(RoleInfo, relics, self)
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

		RoleInfo := self.GetModRole().RoleInfo[roleId]
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
		relics := self.GetModRelics().RelicsInfo[relicsKey]
		if relics == nil {
			fmt.Println("圣遗物不存在")
			continue
		}
		self.GetModRole().TakeOffRelics(RoleInfo, relics, self)
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
		for _, v := range self.GetModWeapon().WeaponInfo {
			fmt.Println(fmt.Sprintf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d",
				v.KeyId, v.Level, v.StarLevel, v.RefineLevel))
		}
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		self.GetModWeapon().WeaponUp(weaponKeyId, self)
	}
}

func (self *Player) HandleWeaponStarUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range self.GetModWeapon().WeaponInfo {
			fmt.Println(fmt.Sprintf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d",
				v.KeyId, v.Level, v.StarLevel, v.RefineLevel))
		}
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		self.GetModWeapon().WeaponUpStar(weaponKeyId, self)
	}
}

func (self *Player) HandleWeaponRefineUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range self.GetModWeapon().WeaponInfo {
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
			self.GetModWeapon().WeaponUpRefine(weaponKeyId, weaponTargetKeyId, self)
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

		RoleInfo := self.GetModRole().RoleInfo[roleId]
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
		weaponInfo := self.GetModWeapon().WeaponInfo[weaponKey]
		if weaponInfo == nil {
			fmt.Println("武器不存在")
			continue
		}
		self.GetModRole().WearWeapon(RoleInfo, weaponInfo, self)
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

		RoleInfo := self.GetModRole().RoleInfo[roleId]
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
		weapon := self.GetModWeapon().WeaponInfo[weaponKey]
		if weapon == nil {
			fmt.Println("武器不存在")
			continue
		}
		self.GetModRole().TakeOffWeapon(RoleInfo, weapon, self)
		RoleInfo.ShowInfo(self)
	}
}

func (self *Player) GetMod(modName string) ModBase {
	return self.modManage[modName]
}

func (self *Player) GetModPlayer() *ModPlayer {
	return self.modManage[MOD_PLAYER].(*ModPlayer)
}

func (self *Player) GetModIcon() *ModIcon {
	return self.modManage[MOD_ICON].(*ModIcon)
}

func (self *Player) GetModCard() *ModCard {
	return self.modManage[MOD_CARD].(*ModCard)
}

func (self *Player) GetModUniqueTask() *ModUniqueTask {
	return self.modManage[MOD_UNIQUETASK].(*ModUniqueTask)
}

func (self *Player) GetModRole() *ModRole {
	return self.modManage[MOD_ROLE].(*ModRole)
}

func (self *Player) GetModBag() *ModBag {
	return self.modManage[MOD_BAG].(*ModBag)
}

func (self *Player) GetModWeapon() *ModWeapon {
	return self.modManage[MOD_WEAPON].(*ModWeapon)
}

func (self *Player) GetModRelics() *ModRelics {
	return self.modManage[MOD_RELICS].(*ModRelics)
}

func (self *Player) GetModCook() *ModCook {
	return self.modManage[MOD_COOK].(*ModCook)
}

func (self *Player) GetModHome() *ModHome {
	return self.modManage[MOD_HOME].(*ModHome)
}

func (self *Player) GetModPool() *ModPool {
	return self.modManage[MOD_POOL].(*ModPool)
}

func (self *Player) GetModMap() *ModMap {
	return self.modManage[MOD_MAP].(*ModMap)
}
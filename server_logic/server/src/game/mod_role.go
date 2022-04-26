package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server/csvs"
	"time"
)

type RoleInfo struct {
	RoleId     int
	GetTimes   int
	RelicsInfo []int
	WeaponInfo int
}

type ModRole struct {
	RoleInfo  map[int]*RoleInfo
	HpPool    int
	HpCalTime int64

	player *Player
	path   string
}

func (self *ModRole) IsHasRole(roleId int) bool {
	return true
}

func (self *ModRole) GetRoleLevel(roleId int) int {
	return 80
}

func (self *ModRole) AddItem(roleId int, num int64) {
	config := csvs.GetRoleConfig(roleId)
	if config == nil {
		fmt.Println("配置不存在roleId:", roleId)
		return
	}
	for i := 0; i < int(num); i++ {
		_, ok := self.RoleInfo[roleId]
		if !ok {
			data := new(RoleInfo)
			data.RoleId = roleId
			data.GetTimes = 1
			self.RoleInfo[roleId] = data
		} else {
			//判断实际获得东西
			self.RoleInfo[roleId].GetTimes++
			if self.RoleInfo[roleId].GetTimes >= csvs.ADD_ROLE_TIME_NORMAL_MIN &&
				self.RoleInfo[roleId].GetTimes <= csvs.ADD_ROLE_TIME_NORMAL_MAX {
				self.player.GetModBag().AddItemToBag(config.Stuff, config.StuffNum)
				self.player.GetModBag().AddItemToBag(config.StuffItem, config.StuffItemNum)
			} else {
				self.player.GetModBag().AddItemToBag(config.MaxStuffItem, config.MaxStuffItemNum)
			}
		}
	}
	itemConfig := csvs.GetItemConfig(roleId)
	if itemConfig != nil {
		fmt.Println("获得角色", itemConfig.ItemName, "次数", roleId, "------", self.RoleInfo[roleId].GetTimes, "次")
	}
	self.player.GetModIcon().CheckGetIcon(roleId)
	self.player.GetModCard().CheckGetCard(roleId, 10)
}

func (self *ModRole) HandleSendRoleInfo(player *Player) {
	fmt.Println(fmt.Sprintf("当前拥有角色信息如下:"))
	for _, v := range self.RoleInfo {
		v.SendRoleInfo(player)
	}
}

func (self *RoleInfo) SendRoleInfo(player *Player) {
	fmt.Println(fmt.Sprintf("%s:,Id:%d,累计获得次数:%d", csvs.GetItemName(self.RoleId), self.RoleId, self.GetTimes))
	self.ShowInfo(player)
}

func (self *ModRole) GetRoleInfoForPoolCheck() (map[int]int, map[int]int) {
	fiveInfo := make(map[int]int)
	fourInfo := make(map[int]int)

	for _, v := range self.RoleInfo {
		roleConfig := csvs.GetRoleConfig(v.RoleId)
		if roleConfig == nil {
			continue
		}
		if roleConfig.Star == 5 {
			fiveInfo[roleConfig.RoleId] = v.GetTimes
		} else if roleConfig.Star == 4 {
			fourInfo[roleConfig.RoleId] = v.GetTimes
		}
	}
	return fiveInfo, fourInfo
}

func (self *ModRole) CalHpPool() () {
	if self.HpCalTime == 0 {
		self.HpCalTime = time.Now().Unix()
	}
	calTime := time.Now().Unix() - self.HpCalTime
	self.HpPool += int(calTime) * 10
	self.HpCalTime = time.Now().Unix()
	fmt.Println("当前血池回复量:", self.HpPool)
}

//把圣遗物穿在角色身上
func (self *ModRole) WearRelics(roleInfo *RoleInfo, relics *Relics, player *Player) {
	relicsConfig := csvs.GetRelicsConfig(relics.RelicsId)
	if relicsConfig == nil {
		return
	}
	self.CheckRelicsPos(roleInfo, relicsConfig.Pos)
	if relicsConfig.Pos < 0 || relicsConfig.Pos > len(roleInfo.RelicsInfo) {
		return
	}

	oldRelicsKeyId := roleInfo.RelicsInfo[relicsConfig.Pos-1]
	if oldRelicsKeyId > 0 {
		oldRelics := player.GetModRelics().RelicsInfo[oldRelicsKeyId]
		if oldRelics != nil {
			oldRelics.RoleId = 0
		}
		roleInfo.RelicsInfo[relicsConfig.Pos-1] = 0
	}

	oldRoleId := relics.RoleId
	if oldRoleId > 0 {
		oldRole := player.GetModRole().RoleInfo[oldRoleId]
		if oldRole != nil {
			oldRole.RelicsInfo[relicsConfig.Pos-1] = 0
		}
		relics.RoleId = 0
	}

	roleInfo.RelicsInfo[relicsConfig.Pos-1] = relics.KeyId
	relics.RoleId = roleInfo.RoleId

	if oldRelicsKeyId > 0 && oldRoleId > 0 {
		oldRelics := player.GetModRelics().RelicsInfo[oldRelicsKeyId]
		oldRole := player.GetModRole().RoleInfo[oldRoleId]
		if oldRelics != nil && oldRole != nil {
			self.WearRelics(oldRole, oldRelics, player)
		}
	}

	roleInfo.ShowInfo(player)
}

func (self *ModRole) CheckRelicsPos(roleInfo *RoleInfo, pos int) {
	nowSize := len(roleInfo.RelicsInfo)
	needAdd := pos - nowSize

	for i := 0; i < needAdd; i++ {
		roleInfo.RelicsInfo = append(roleInfo.RelicsInfo, 0)
	}
}

func (self *RoleInfo) ShowInfo(player *Player) {
	fmt.Println(fmt.Sprintf("当前角色:%s,角色ID:%d", csvs.GetItemName(self.RoleId), self.RoleId))

	weaponNow := player.GetModWeapon().WeaponInfo[self.WeaponInfo]
	if weaponNow == nil {
		fmt.Println(fmt.Sprintf("武器:未穿戴"))
	} else {
		fmt.Println(fmt.Sprintf("武器:%s,key:%d", csvs.GetItemName(weaponNow.WeaponId), self.WeaponInfo))
	}

	suitMap := make(map[int]int)
	for _, v := range self.RelicsInfo {
		relicsNow := player.GetModRelics().RelicsInfo[v]
		if relicsNow == nil {
			fmt.Println(fmt.Sprintf("未穿戴"))
			continue
		}
		fmt.Println(fmt.Sprintf("%s,key:%d", csvs.GetItemName(relicsNow.RelicsId), v))
		relicsNowConfig := csvs.GetRelicsConfig(relicsNow.RelicsId)
		if relicsNowConfig != nil {
			suitMap[relicsNowConfig.Type]++
		}
	}

	suitSkill := make([]int, 0)
	for suit, num := range suitMap {
		for _, config := range csvs.ConfigRelicsSuitMap[suit] {
			if num >= config.Num {
				suitSkill = append(suitSkill, config.SuitSkill)
			}
		}
	}
	for _, v := range suitSkill {
		fmt.Println(fmt.Sprintf("激活套装效果:%d", v))
	}

}

func (self *ModRole) TakeOffRelics(roleInfo *RoleInfo, relics *Relics, player *Player) {
	relicsConfig := csvs.GetRelicsConfig(relics.RelicsId)
	if relicsConfig == nil {
		return
	}
	self.CheckRelicsPos(roleInfo, relicsConfig.Pos)
	if relicsConfig.Pos < 0 || relicsConfig.Pos > len(roleInfo.RelicsInfo) {
		return
	}
	if roleInfo.RelicsInfo[relicsConfig.Pos-1] != relics.KeyId {
		fmt.Println(fmt.Sprintf("当前角色没有穿戴这个物品"))
		return
	}

	roleInfo.RelicsInfo[relicsConfig.Pos-1] = 0
	relics.RoleId = 0
	roleInfo.ShowInfo(player)
}

func (self *ModRole) WearWeapon(roleInfo *RoleInfo, weapon *Weapon, player *Player) {
	weaponConfig := csvs.GetWeaponConfig(weapon.WeaponId)
	if weaponConfig == nil {
		fmt.Println("数据异常，武器配置不存在")
		return
	}

	//先判断武器和角色是否匹配
	roleConfig := csvs.GetRoleConfig(roleInfo.RoleId)
	if roleConfig.Type != weaponConfig.Type {
		fmt.Println("武器和角色不匹配")
		return
	}

	oldWeaponKey := 0
	if roleInfo.WeaponInfo > 0 {
		oldWeaponKey = roleInfo.WeaponInfo
		roleInfo.WeaponInfo = 0
		oldWeapon := player.GetModWeapon().WeaponInfo[oldWeaponKey]
		if oldWeapon != nil {
			oldWeapon.RoleId = 0
		}
	}

	oldRoleId := 0
	if weapon.RoleId > 0 {
		oldRoleId = weapon.RoleId
		weapon.RoleId = 0
		oldRole := player.GetModRole().RoleInfo[oldRoleId]
		if oldRole != nil {
			oldRole.WeaponInfo = 0
		}
	}

	roleInfo.WeaponInfo = weapon.KeyId
	weapon.RoleId = roleInfo.RoleId

	if roleInfo.WeaponInfo > 0 && weapon.RoleId > 0 {
		oldWeapon := player.GetModWeapon().WeaponInfo[oldWeaponKey]
		oldRole := player.GetModRole().RoleInfo[oldRoleId]
		if oldWeapon != nil && oldRole != nil {
			self.WearWeapon(oldRole, oldWeapon, player)
		}
	}
}

func (self *ModRole) TakeOffWeapon(roleInfo *RoleInfo, weapon *Weapon, player *Player) {
	weaponConfig := csvs.GetWeaponConfig(weapon.WeaponId)
	if weaponConfig == nil {
		fmt.Println("数据异常，武器配置不存在")
		return
	}
	if roleInfo.WeaponInfo != weapon.KeyId {
		fmt.Println("角色没有装备这把武器")
		return
	}
	//根据位置看是否身上有对应圣遗物
	roleInfo.WeaponInfo = 0
	weapon.RoleId = 0
}

func (self *ModRole) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModRole) LoadData(player *Player) {

	self.player = player
	self.path = self.player.localPath + "/role.json"

	configFile, err := ioutil.ReadFile(self.path)
	if err != nil {
		self.InitData()
		return
	}
	err = json.Unmarshal(configFile, &self)
	if err != nil {
		self.InitData()
		return
	}

	if self.RoleInfo == nil {
		self.RoleInfo = make(map[int]*RoleInfo)
	}
	return
}

func (self *ModRole) InitData() {
	if self.RoleInfo == nil {
		self.RoleInfo = make(map[int]*RoleInfo)
	}
}

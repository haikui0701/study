package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server/csvs"
	"time"
)

type ShowRole struct {
	RoleId    int
	RoleLevel int
}
type ModPlayer struct {
	UserId         int       `json:"userid" `  //唯一id
	Icon           int         //头像   新增icon模块
	Card           int         //名片   新增card模块
	Name           string      //名字   新增banword模块
	Sign           string      //签名
	PlayerLevel    int         //等级
	PlayerExp      int         //阅历(经验)
	WorldLevel     int         //大世界等级
	WorldLevelNow  int         //大世界等级(当前)
	WorldLevelCool int64       //操作大世界等级的冷却时间
	Birth          int         //生日
	ShowTeam       []*ShowRole //展示阵容
	HideShowTeam   int         //隐藏开关
	ShowCard       []int       //展示名片
	//看不见的字段
	Prohibit int //封禁状态
	IsGM     int //GM账号标志
}

func (self *ModPlayer) SetIcon(iconId int, player *Player) {
	if !player.ModIcon.IsHasIcon(iconId) {
		//通知客户端，操作非法
		fmt.Println("没有头像:", iconId)
		return
	}

	player.ModPlayer.Icon = iconId
	fmt.Println("变更头像为:", csvs.GetItemName(iconId), player.ModPlayer.Icon, )
}

func (self *ModPlayer) SetCard(cardId int, player *Player) {
	if !player.ModCard.IsHasCard(cardId) {
		//通知客户端，操作非法
		return
	}

	player.ModPlayer.Card = cardId
	fmt.Println("当前名片", player.ModPlayer.Card)
}

func (self *ModPlayer) SetName(name string, playerTest interface{}) {

	player := playerTest.(*ModPlayer)

	if GetManageBanWord().IsBanWord(name) {
		return
	}
	print(player)
	//player.ModPlayer.Name = name
	//fmt.Println("设置成功,名字变更为:", player.ModPlayer.Name)
}

func (self *ModPlayer) SetSign(sign string, player *Player) {
	if GetManageBanWord().IsBanWord(sign) {
		return
	}

	player.ModPlayer.Sign = sign
	fmt.Println("设置成功,签名变更为:", player.ModPlayer.Sign)
}

func (self *ModPlayer) AddExp(exp int, player *Player) {
	self.PlayerExp += exp
	for {
		config := csvs.GetNowLevelConfig(self.PlayerLevel)
		if config == nil {
			break
		}
		if config.PlayerExp == 0 {
			break
		}
		//是否完成任务
		if config.ChapterId > 0 && !player.ModUniqueTask.IsTaskFinish(config.ChapterId) {
			break
		}
		if self.PlayerExp >= config.PlayerExp {
			self.PlayerLevel += 1
			self.PlayerExp -= config.PlayerExp
		} else {
			break
		}
	}
	fmt.Println("当前等级:", self.PlayerLevel, "---当前经验：", self.PlayerExp)
}

func (self *ModPlayer) ReduceWorldLevel(player *Player) {
	if self.WorldLevel < csvs.REDUCE_WORLD_LEVEL_START {
		fmt.Println("操作失败:, ---当前世界等级：", self.WorldLevel)
		return
	}

	if self.WorldLevel-self.WorldLevelNow >= csvs.REDUCE_WORLD_LEVEL_MAX {
		fmt.Println("操作失败:, ---当前世界等级：", self.WorldLevel, "---真实世界等级：", self.WorldLevelNow)
		return
	}

	if time.Now().Unix() < self.WorldLevelCool {
		fmt.Println("操作失败:, ---冷却中")
		return
	}

	self.WorldLevelNow -= 1
	self.WorldLevelCool = time.Now().Unix() + csvs.REDUCE_WORLD_LEVEL_COOL_TIME
	fmt.Println("操作成功:, ---当前世界等级：", self.WorldLevel, "---真实世界等级：", self.WorldLevelNow)
	return
}

func (self *ModPlayer) ReturnWorldLevel(player *Player) {
	if self.WorldLevelNow == self.WorldLevel {
		fmt.Println("操作失败:, ---当前世界等级：", self.WorldLevel, "---真实世界等级：", self.WorldLevelNow)
		return
	}

	if time.Now().Unix() < self.WorldLevelCool {
		fmt.Println("操作失败:, ---冷却中")
		return
	}

	self.WorldLevelNow += 1
	self.WorldLevelCool = time.Now().Unix() + csvs.REDUCE_WORLD_LEVEL_COOL_TIME
	fmt.Println("操作成功:, ---当前世界等级：", self.WorldLevel, "---真实世界等级：", self.WorldLevelNow)
	return
}

func (self *ModPlayer) SetBirth(birth int, player *Player) {
	if self.Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}

	month := birth / 100
	day := birth % 100

	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		if day <= 0 || day > 31 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	case 4, 6, 9, 11:
		if day <= 0 || day > 30 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	case 2:
		if day <= 0 || day > 29 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	default:
		fmt.Println("没有", month, "月！")
		return
	}

	self.Birth = birth
	fmt.Println("设置成功，生日为:", month, "月", day, "日")

	if self.IsBirthDay() {
		fmt.Println("今天是你的生日，生日快乐！")
	} else {
		fmt.Println("期待你生日的到来!")
	}

}

func (self *ModPlayer) IsBirthDay() bool {
	month := time.Now().Month()
	day := time.Now().Day()
	if int(month) == self.Birth/100 && day == self.Birth%100 {
		return true
	}
	return false
}

func (self *ModPlayer) SetShowCard(showCard []int, player *Player) {

	if len(showCard) > csvs.SHOW_SIZE {
		return
	}

	cardExist := make(map[int]int)
	newList := make([]int, 0)
	for _, cardId := range showCard {
		_, ok := cardExist[cardId]
		if ok {
			continue
		}
		if !player.ModCard.IsHasCard(cardId) {
			continue
		}
		newList = append(newList, cardId)
		cardExist[cardId] = 1
	}
	self.ShowCard = newList
	fmt.Println(self.ShowCard)
}

func (self *ModPlayer) SetShowTeam(showRole []int, player *Player) {
	if len(showRole) > csvs.SHOW_SIZE {
		fmt.Println("消息结构错误")
		return
	}

	roleExist := make(map[int]int)
	newList := make([]*ShowRole, 0)
	for _, roleId := range showRole {
		_, ok := roleExist[roleId]
		if ok {
			continue
		}
		if !player.ModRole.IsHasRole(roleId) {
			continue
		}
		showRole := new(ShowRole)
		showRole.RoleId = roleId
		showRole.RoleLevel = player.ModRole.GetRoleLevel(roleId)
		newList = append(newList, showRole)
		roleExist[roleId] = 1
	}
	self.ShowTeam = newList
	fmt.Println(self.ShowCard)
}

func (self *ModPlayer) SetHideShowTeam(isHide int, player *Player) {
	if isHide != csvs.LOGIC_FALSE && isHide != csvs.LOGIC_TRUE {
		return
	}
	self.HideShowTeam = isHide
}

func (self *ModPlayer) SetProhibit(prohibit int) {
	self.Prohibit = prohibit
}

func (self *ModPlayer) SetIsGM(isGm int) {
	self.IsGM = isGm
}

func (self *ModPlayer) IsCanEnter() bool {
	return int64(self.Prohibit) < time.Now().Unix()
}

func (self *ModPlayer) GetWorldLevelNow() int {
	return self.WorldLevelNow
}

func (self *ModPlayer) RelicsUp(player *Player) {
	player.ModRelics.RelicsUp(player)
}

func (self *ModPlayer) SaveData(path string) {
	self.ShowCard = append(self.ShowCard, 1)
	self.ShowCard = append(self.ShowCard, 2)
	self.ShowCard = append(self.ShowCard, 3)
	self.ShowCard = append(self.ShowCard, 4)
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModPlayer) LoadData(path string) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error")
		return
	}
	err = json.Unmarshal(configFile, &self)
	if err != nil {
		fmt.Println("error")
		return
	}
	return
}
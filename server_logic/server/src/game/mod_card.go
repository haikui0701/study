package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server/csvs"
)

type Card struct {
	CardId int
}

type ModCard struct {
	CardInfo map[int]*Card

	player *Player
	path   string
}

func (self *ModCard) IsHasCard(cardId int) bool {
	_, ok := self.CardInfo[cardId]
	return ok
}

func (self *ModCard) AddItem(itemId int,friendliness int) {
	_, ok := self.CardInfo[itemId]
	if ok {
		fmt.Println("已存在名片：", itemId)
		return
	}
	config := csvs.GetCardConfig(itemId)
	if config == nil {
		fmt.Println("非法名片：", itemId)
		return
	}
	if friendliness<config.Friendliness{
		fmt.Println("好感度不足：", itemId)
		return
	}

	self.CardInfo[itemId] = &Card{CardId: itemId}
	fmt.Println("获得名片：", itemId)
}

func (self *ModCard) CheckGetCard(roleId int,friendliness int) {
	config:=csvs.GetCardConfigByRoleId(roleId)
	if config==nil{
		return
	}
	self.AddItem(config.CardId,friendliness)
}

func (self *ModCard) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModCard) LoadData(player *Player) {

	self.player = player
	self.path = self.player.localPath + "/card.json"

	configFile, err := ioutil.ReadFile(self.path)
	if err != nil {
		fmt.Println("error")
		return
	}
	err = json.Unmarshal(configFile, &self)
	if err != nil {
		self.InitData()
		return
	}

	if self.CardInfo == nil {
		self.CardInfo = make(map[int]*Card)
	}
	return
}

func (self *ModCard) InitData() {

}
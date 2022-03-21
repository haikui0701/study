package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server/csvs"
)

type Cook struct {
	CookId int
}

type ModCook struct {
	CookInfo map[int]*Cook

	player *Player
	path   string
}

func (self *ModCook) AddItem(itemId int) {
	_, ok := self.CookInfo[itemId]
	if ok {
		fmt.Println("已习得：", csvs.GetItemName(itemId))
		return
	}
	config := csvs.GetCookConfig(itemId)
	if config == nil {
		fmt.Println("没有这个烹饪技能：", csvs.GetItemName(itemId))
		return
	}
	self.CookInfo[itemId] = &Cook{CookId: itemId}
	fmt.Println("学会烹饪：", itemId)
}

func (self *ModCook) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModCook) LoadData(player *Player) {

	self.player = player
	self.path = self.player.localPath + "/cook.json"

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

	if self.CookInfo == nil {
		self.CookInfo = make(map[int]*Cook)
	}
	return
}

func (self *ModCook) InitData() {

}
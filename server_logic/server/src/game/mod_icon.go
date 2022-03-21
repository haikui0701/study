package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server/csvs"
)

type Icon struct {
	IconId int
}

type ModIcon struct {
	IconInfo map[int]*Icon

	player *Player
	path   string
}

func (self *ModIcon) IsHasIcon(iconId int) bool {
	_, ok := self.IconInfo[iconId]
	return ok
}

func (self *ModIcon) AddItem(itemId int) {
	_, ok := self.IconInfo[itemId]
	if ok {
		fmt.Println("已存在头像：", itemId)
		return
	}
	config := csvs.GetIconConfig(itemId)
	if config == nil {
		fmt.Println("非法头像：", itemId)
		return
	}
	self.IconInfo[itemId] = &Icon{IconId: itemId}
	fmt.Println("获得头像：", itemId)
}

func (self *ModIcon) CheckGetIcon(roleId int) {
	config:=csvs.GetIconConfigByRoleId(roleId)
	if config==nil{
		return
	}
	self.AddItem(config.IconId)
}

func (self *ModIcon) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModIcon) LoadData(player *Player) {

	self.player=player
	self.path=self.player.localPath+"/icon.json"

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

	if self.IconInfo==nil{
		self.IconInfo=make(map[int]*Icon,0)
	}

	return
}

func (self *ModIcon) InitData() {

}
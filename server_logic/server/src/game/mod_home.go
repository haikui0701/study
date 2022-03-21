package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server/csvs"
)

type HomeItemId struct {
	HomeItemId  int
	HomeItemNum int64
	KeyId       int
}

type ModHome struct {
	HomeItemIdInfo map[int]*HomeItemId

	player *Player
	path   string
}

func (self *ModHome) AddItem(itemId int, num int64) {
	_, ok := self.HomeItemIdInfo[itemId]
	if ok {
		self.HomeItemIdInfo[itemId].HomeItemNum += num
	} else {
		self.HomeItemIdInfo[itemId] = &HomeItemId{HomeItemId: itemId, HomeItemNum: num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("获得家具物品", config.ItemName, "----数量：", num, "----当前数量：", self.HomeItemIdInfo[itemId].HomeItemNum)
	}
}

func (self *ModHome) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModHome) LoadData(player *Player) {

	self.player = player
	self.path = self.player.localPath + "/home.json"

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

	if self.HomeItemIdInfo == nil {
		self.HomeItemIdInfo = make(map[int]*HomeItemId)
	}
	return
}

func (self *ModHome) InitData() {

}

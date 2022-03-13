package game

import (
	"fmt"
	"server/csvs"
)

type HomeItemId struct {
	HomeItemId  int
	HomeItemNum int64
	KeyId       int
}

type ModHome struct {
	HomeItemIdInfo map[int]*HomeItemId
	//UseHomeItemIdInfo map[int]*HomeItemId
	//Map
}

func (self *ModHome) AddItem(itemId int, num int64, player *Player) {
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

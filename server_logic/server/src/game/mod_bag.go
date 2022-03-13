package game

import (
	"fmt"
	"server/csvs"
)

type ItemInfo struct {
	ItemId  int
	ItemNum int64
}

type ModBag struct {
	BagInfo map[int]*ItemInfo
}

func (self *ModBag) AddItem(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}
	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	self.AddItemToBag(itemId, num)
	case csvs.ITEMTYPE_ROLE:
		player.ModRole.AddItem(itemId, num, player)
	case csvs.ITEMTYPE_ICON:
		player.ModIcon.AddItem(itemId)
	case csvs.ITEMTYPE_CARD:
		player.ModCard.AddItem(itemId, 12)
	case csvs.ITEMTYPE_WEAPON:
		player.ModWeapon.AddItem(itemId, num)
	case csvs.ITEMTYPE_RELICS:
		player.ModRelics.AddItem(itemId, num)
	case csvs.ITEMTYPE_COOK:
		player.ModCook.AddItem(itemId)
	case csvs.ITEMTYPE_HOME_ITEM:
		player.ModHome.AddItem(itemId, num, player)
	default: //同普通
		self.AddItemToBag(itemId, num)
	}
}

func (self *ModBag) AddItemToBag(itemId int, num int64) {
	_, ok := self.BagInfo[itemId]
	if ok {
		self.BagInfo[itemId].ItemNum += num
	} else {
		self.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("获得物品", config.ItemName, "----数量：", num, "----当前数量：", self.BagInfo[itemId].ItemNum)
	}

}

func (self *ModBag) RemoveItem(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	switch itemConfig.SortType {
	case csvs.ITEMTYPE_NORMAL:
		self.RemoveItemToBagGM(itemId, num)
	default: //同普通
		//self.AddItemToBag(itemId, 1)
	}
}

func (self *ModBag) RemoveItemToBagGM(itemId int, num int64) {
	_, ok := self.BagInfo[itemId]
	if ok {
		self.BagInfo[itemId].ItemNum -= num
	} else {
		self.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("扣除物品", config.ItemName, "----数量：", num, "----当前数量：", self.BagInfo[itemId].ItemNum)
	}
}

func (self *ModBag) RemoveItemToBag(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	self.AddItemToBag(itemId, num)
	case csvs.ITEMTYPE_ROLE:
		fmt.Println("此物品无法扣除")
		return
	case csvs.ITEMTYPE_ICON:
		fmt.Println("此物品无法扣除")
		return
	case csvs.ITEMTYPE_CARD:
		fmt.Println("此物品无法扣除")
		return
	default: //同普通
	}

	if !self.HasEnoughItem(itemId, num) {
		config := csvs.GetItemConfig(itemId)
		if config != nil {
			nowNum := int64(0)
			_, ok := self.BagInfo[itemId]
			if ok {
				nowNum = self.BagInfo[itemId].ItemNum
			}
			fmt.Println(config.ItemName, "数量不足", "----当前数量：", nowNum)
		}
		return
	}

	_, ok := self.BagInfo[itemId]
	if ok {
		self.BagInfo[itemId].ItemNum -= num
	} else {
		self.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	fmt.Println("扣除物品", itemConfig.ItemName, "----数量：", num, "----当前数量：", self.BagInfo[itemId].ItemNum)
}

func (self *ModBag) HasEnoughItem(itemId int, num int64) bool {
	if itemId == 0 {
		return true
	}
	_, ok := self.BagInfo[itemId]
	if !ok {
		return false
	} else if self.BagInfo[itemId].ItemNum < num {
		return false
	}
	return true
}

func (self *ModBag) UseItem(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	if !self.HasEnoughItem(itemId, num) {
		config := csvs.GetItemConfig(itemId)
		if config != nil {
			nowNum := int64(0)
			_, ok := self.BagInfo[itemId]
			if ok {
				nowNum = self.BagInfo[itemId].ItemNum
			}
			fmt.Println(config.ItemName, "数量不足", "----当前数量：", nowNum)
		}
		return
	}

	switch itemConfig.SortType {
	case csvs.ITEMTYPE_COOKBOOK:
		self.UseCookBook(itemId, num, player)
	case csvs.ITEMTYPE_FOOD:
		//给英雄加属性
	default: //同普通
		fmt.Println(itemId, "此物品无法使用")
		return
	}
}

func (self *ModBag) UseCookBook(itemId int, num int64, player *Player) {
	cookBookConfig := csvs.GetCookBookConfig(itemId)
	if cookBookConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}
	self.RemoveItem(itemId, num, player)
	self.AddItem(cookBookConfig.Reward, num, player)
}

func (self *ModBag) GetItemNum(itemId int, player *Player) int64 {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		return 0
	}
	_, ok := self.BagInfo[itemId]
	if !ok {
		return 0
	}
	return self.BagInfo[itemId].ItemNum
}

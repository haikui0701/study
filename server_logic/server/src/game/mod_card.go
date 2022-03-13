package game

import (
	"fmt"
	"server/csvs"
)

type Card struct {
	CardId int
}

type ModCard struct {
	CardInfo map[int]*Card
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
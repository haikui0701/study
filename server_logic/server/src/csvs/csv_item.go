package csvs

import (
	"server/utils"
)

const (
	ITEMTYPE_NORMAL    = 1
	ITEMTYPE_ROLE      = 2
	ITEMTYPE_ICON      = 3
	ITEMTYPE_CARD      = 4
	ITEMTYPE_WEAPON    = 6
	ITEMTYPE_RELICS    = 7
	ITEMTYPE_COOKBOOK  = 8
	ITEMTYPE_COOK      = 9
	ITEMTYPE_FOOD      = 10
	ITEMTYPE_HOME_ITEM = 11
)

type ConfigItem struct {
	ItemId   int    `json:"ItemId"`
	SortType int    `json:"SortType"`
	ItemName string `json:"ItemName"`
}

var (
	ConfigItemMap map[int]*ConfigItem
)

func init() {
	ConfigItemMap = make(map[int]*ConfigItem)
	utils.GetCsvUtilMgr().LoadCsv("Item", &ConfigItemMap)
	return
}

func GetItemConfig(itemId int) *ConfigItem {
	return ConfigItemMap[itemId]
}

func GetItemName(itemId int) string {
	config := GetItemConfig(itemId)
	if config == nil {
		return ""
	}
	return config.ItemName
}

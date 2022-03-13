package csvs

import "server/utils"

type ConfigDrop struct {
	DropId int `json:"DropId"`
	Weight int `json:"Weight"`
	Result int `json:"Result"`
	IsEnd  int `json:"IsEnd"`
}

type ConfigDropItem struct {
	DropId     int `json:"DropId"`
	DropType   int `json:"DropType"`
	Weight     int `json:"Weight"`
	ItemId     int `json:"ItemId"`
	ItemNumMin int `json:"ItemNumMin"`
	ItemNumMax int `json:"ItemNumMax"`
	WorldAdd   int `json:"WorldAdd"`
}

var (
	ConfigDropSlice     []*ConfigDrop
	ConfigDropItemSlice []*ConfigDropItem
)

func init() {
	ConfigDropSlice = make([]*ConfigDrop, 0)
	utils.GetCsvUtilMgr().LoadCsv("Drop", &ConfigDropSlice)

	ConfigDropItemSlice = make([]*ConfigDropItem, 0)
	utils.GetCsvUtilMgr().LoadCsv("DropItem", &ConfigDropItemSlice)
	return
}

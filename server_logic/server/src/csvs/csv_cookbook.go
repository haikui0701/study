package csvs

import (
	"server/utils"
)

type ConfigCookBook struct {
	CookBookId int `json:"CookBookId"`
	Reward     int `json:"Reward"`
}

var (
	ConfigCookBookMap map[int]*ConfigCookBook
)

func init() {
	ConfigCookBookMap = make(map[int]*ConfigCookBook)
	utils.GetCsvUtilMgr().LoadCsv("CookBook", &ConfigCookBookMap)
	return
}

func GetCookBookConfig(cookBookId int) *ConfigCookBook {
	return ConfigCookBookMap[cookBookId]
}

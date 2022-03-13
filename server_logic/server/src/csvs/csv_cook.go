package csvs

import "server/utils"

type ConfigCook struct {
	CookId int `json:"CookId"`
}

var (
	ConfigCookMap         map[int]*ConfigCook
)

func init() {
	ConfigCookMap = make(map[int]*ConfigCook)
	utils.GetCsvUtilMgr().LoadCsv("Cook", &ConfigCookMap)
	return
}

func GetCookConfig(cookId int) *ConfigCook {
	return ConfigCookMap[cookId]
}

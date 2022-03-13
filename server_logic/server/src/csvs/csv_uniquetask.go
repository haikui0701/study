package csvs

import (
	"server/utils"
)

type ConfigUniqueTask struct {
	TaskId    int `json:"TaskId"`
	SortType  int `json:"SortType"`
	OpenLevel int `json:"OpenLevel"`
	TaskType  int `json:"TaskType"`
	Condition int `json:"Condition"`
}

var (
	ConfigUniqueTaskMap map[int]*ConfigUniqueTask
)

func init() {
	ConfigUniqueTaskMap=make(map[int]*ConfigUniqueTask)
	utils.GetCsvUtilMgr().LoadCsv("UniqueTask", &ConfigUniqueTaskMap)
	return
}


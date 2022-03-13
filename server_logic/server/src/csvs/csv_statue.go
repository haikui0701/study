package csvs

import "server/utils"

type ConfigStatue struct {
	StatueId   int   `json:"StatueId"`
	Level      int   `json:"Level"`
	CostItem   int   `json:"CostItem"`
	CostNum    int64 `json:"CostNum"`
	RewardItem []int `json:"RewardItem"`
	RewardNum  []int `json:"RewardNum"`
}

var (
	ConfigStatueSlice []*ConfigStatue
)

func init() {
	ConfigStatueSlice = make([]*ConfigStatue, 0)
	utils.GetCsvUtilMgr().LoadCsv("Statue", &ConfigStatueSlice)
	return
}

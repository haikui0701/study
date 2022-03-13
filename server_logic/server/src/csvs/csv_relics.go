package csvs

import "server/utils"

type ConfigRelics struct {
	RelicsId      int `json:"RelicsId"`
	Type          int `json:"Type"`
	Pos           int `json:"Pos"`
	Star          int `json:"Star"`
	MainGroup     int `json:"MainGroup"`
	OtherGroup    int `json:"OtherGroup"`
	OtherGroupNum int `json:"OtherGroupNum"`
}

type ConfigRelicsEntry struct {
	Id        int    `json:"Id"`
	Group     int    `json:"Group"`
	AttrType  int    `json:"AttrType"`
	AttrValue int    `json:"AttrValue"`
	Weight    int    `json:"Weight"`
	AttrName  string `json:"AttrName"`
}

type ConfigRelicsLevel struct {
	EntryId   int    `json:"EntryId"`
	Level     int    `json:"Level"`
	AttrType  int    `json:"AttrType"`
	AttrName  string `json:"AttrName"`
	AttrValue int    `json:"AttrValue"`
	NeedExp   int    `json:"NeedExp"`
}

type ConfigRelicsSuit struct {
	Type        int    `json:"Type"`        //套装类型
	Num         int    `json:"Num"`         //需要件数
	SuitSkill   int    `json:"SuitSkill"`   //套装技能
	SkillString string `json:"SkillString"` //套装描述
}

var (
	ConfigRelicsMap        map[int]*ConfigRelics
	ConfigRelicsEntryMap   map[int]*ConfigRelicsEntry
	ConfigRelicsLevelSlice []*ConfigRelicsLevel
	ConfigRelicsSuitSlice  []*ConfigRelicsSuit
)

func init() {
	ConfigRelicsMap = make(map[int]*ConfigRelics)
	utils.GetCsvUtilMgr().LoadCsv("Relics", &ConfigRelicsMap)

	ConfigRelicsEntryMap = make(map[int]*ConfigRelicsEntry)
	utils.GetCsvUtilMgr().LoadCsv("RelicsEntry", &ConfigRelicsEntryMap)

	utils.GetCsvUtilMgr().LoadCsv("RelicsLevel", &ConfigRelicsLevelSlice)
	utils.GetCsvUtilMgr().LoadCsv("RelicsSuit", &ConfigRelicsSuitSlice)
	return
}

func GetRelicsConfig(relicsId int) *ConfigRelics {
	return ConfigRelicsMap[relicsId]
}

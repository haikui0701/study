package csvs

import "server/utils"

type ConfigMap struct {
	MapId   int    `json:"MapId"`
	MapName string `json:"MapName"`
	MapType int    `json:"MapType"`
}

type ConfigMapEvent struct {
	EventId        int    `json:"EventId"`
	EventType      int    `json:"EventType"`
	Name           string `json:"Name"`
	RefreshType    int    `json:"RefreshType"`
	EventDrop      int    `json:"EventDrop"`
	EventDropTimes int    `json:"EventDropTimes"`
	MapId          int    `json:"MapId"`
	CostItem       int    `json:"CostItem"`
	CostNum        int64  `json:"CostNum"`
}

var (
	ConfigMapMap      map[int]*ConfigMap
	ConfigMapEventMap map[int]*ConfigMapEvent
)

func init() {
	ConfigMapMap = make(map[int]*ConfigMap)
	utils.GetCsvUtilMgr().LoadCsv("Map", &ConfigMapMap)

	ConfigMapEventMap = make(map[int]*ConfigMapEvent)
	utils.GetCsvUtilMgr().LoadCsv("MapEvent", &ConfigMapEventMap)
	return
}

func GetMapName(mapId int) string {
	_, ok := ConfigMapMap[mapId]
	if !ok {
		return ""
	}
	return ConfigMapMap[mapId].MapName
}

func GetEventName(eventId int) string {
	_, ok := ConfigMapEventMap[eventId]
	if !ok {
		return ""
	}
	return ConfigMapEventMap[eventId].Name
}

func GetEventConfig(eventId int) *ConfigMapEvent {
	return ConfigMapEventMap[eventId]
}

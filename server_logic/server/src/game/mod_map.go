package game

import (
	"fmt"
	"math/rand"
	"server/csvs"
	"time"
)

type Map struct {
	MapId     int
	EventInfo map[int]*Event
}

type Event struct {
	EventId       int
	State         int
	NextResetTime int64
}

type StatueInfo struct {
	StatueId int
	Level    int
	ItemInfo map[int]*ItemInfo
}

type ModMap struct {
	MapInfo map[int]*Map
	Statue  map[int]*StatueInfo
}

func (self *ModMap) InitData() {
	self.MapInfo = make(map[int]*Map)
	self.Statue = make(map[int]*StatueInfo)

	for _, v := range csvs.ConfigMapMap {
		_, ok := self.MapInfo[v.MapId]
		if !ok {
			self.MapInfo[v.MapId] = self.NewMapInfo(v.MapId)
		}
	}

	for _, v := range csvs.ConfigMapEventMap {
		_, ok := self.MapInfo[v.MapId]
		if !ok {
			continue
		}
		_, ok = self.MapInfo[v.MapId].EventInfo[v.EventId]
		if !ok {
			self.MapInfo[v.MapId].EventInfo[v.EventId] = new(Event)
			self.MapInfo[v.MapId].EventInfo[v.EventId].EventId = v.EventId
			self.MapInfo[v.MapId].EventInfo[v.EventId].State = csvs.EVENT_START
		}
	}
}

func (self *ModMap) NewMapInfo(mapId int) *Map {
	mapInfo := new(Map)
	mapInfo.MapId = mapId
	mapInfo.EventInfo = make(map[int]*Event)
	return mapInfo
}

func (self *ModMap) GetEventList(config *csvs.ConfigMap) {
	_, ok := self.MapInfo[config.MapId]
	if !ok {
		return
	}
	for _, v := range self.MapInfo[config.MapId].EventInfo {
		self.CheckRefresh(v)
		lastTime := v.NextResetTime - time.Now().Unix()
		noticeTime := ""
		if lastTime <= 0 {
			noticeTime = "已刷新"
		} else {
			noticeTime = fmt.Sprintf("%d秒后刷新", lastTime)
		}
		fmt.Println(fmt.Sprintf("事件Id:%d,名字:%s,状态:%d,%s", v.EventId, csvs.GetEventName(v.EventId), v.State, noticeTime))
	}
}

func (self *ModMap) SetEventState(mapId int, eventId int, state int, player *Player) {
	_, ok := self.MapInfo[mapId]
	if !ok {
		fmt.Println("地图不存在")
		return
	}
	_, ok = self.MapInfo[mapId].EventInfo[eventId]
	if !ok {
		fmt.Println("事件不存在")
		return
	}
	if self.MapInfo[mapId].EventInfo[eventId].State >= state {
		fmt.Println("状态异常")
		return
	}
	eventConfig := csvs.GetEventConfig(self.MapInfo[mapId].EventInfo[eventId].EventId)
	if eventConfig == nil {
		return
	}
	configMap := csvs.ConfigMapMap[mapId]
	if configMap == nil {
		return
	}
	if !player.ModBag.HasEnoughItem(eventConfig.CostItem,eventConfig.CostNum){
		fmt.Println(fmt.Sprintf("%s不足!",csvs.GetItemName(eventConfig.CostItem)))
		return
	}
	if configMap.MapType == csvs.REFRESH_PLAYER && eventConfig.EventType == csvs.EVENT_TYPE_REWARD {
		for _, v := range self.MapInfo[mapId].EventInfo {
			eventConfigNow := csvs.GetEventConfig(v.EventId)
			if eventConfigNow == nil {
				continue
			}
			if eventConfigNow.EventType!= csvs.EVENT_TYPE_NORMAL{
				continue
			}
			if v.EventId == eventId {
				continue
			}
			if v.State != csvs.EVENT_END {
				fmt.Println("有事件尚未完成:", v.EventId)
				return
			}
		}
	}

	self.MapInfo[mapId].EventInfo[eventId].State = state
	if state == csvs.EVENT_FINISH {
		fmt.Println("事件完成")
	}
	if state == csvs.EVENT_END {
		for i:=0;i<eventConfig.EventDropTimes;i++{
			config := csvs.GetDropItemGroupNew(eventConfig.EventDrop)
			for _, v := range config {
				randNum := rand.Intn(csvs.PERCENT_ALL)
				if randNum < v.Weight {
					randAll := v.ItemNumMax - v.ItemNumMin + 1
					itemNum := rand.Intn(randAll) + v.ItemNumMin
					worldLevel := player.ModPlayer.GetWorldLevelNow()
					if worldLevel > 0 {
						itemNum = itemNum * (csvs.PERCENT_ALL + worldLevel*v.WorldAdd) / csvs.PERCENT_ALL
					}
					player.ModBag.AddItem(v.ItemId, int64(itemNum), player)
				}
			}
		}
		fmt.Println("事件领取")
	}
	if state > 0 {
		switch eventConfig.RefreshType {
		case csvs.MAP_REFRESH_SELF:
			self.MapInfo[mapId].EventInfo[eventId].NextResetTime = time.Now().Unix() + csvs.MAP_REFRESH_SELF_TIME
		}
	}
}

func (self *ModMap) RefreshDay() {
	for _, v := range self.MapInfo {
		for _, v := range self.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_DAY {
				continue
			}
			v.State = csvs.EVENT_START
		}
	}
}

func (self *ModMap) RefreshWeek() {
	for _, v := range self.MapInfo {
		for _, v := range self.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_WEEK {
				continue
			}
			v.State = csvs.EVENT_START
		}
	}
}

func (self *ModMap) RefreshSelf() {
	for _, v := range self.MapInfo {
		for _, v := range self.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_SELF {
				continue
			}
			if time.Now().Unix() <= v.NextResetTime {
				v.State = csvs.EVENT_START
			}
		}
	}
}

func (self *ModMap) CheckRefresh(event *Event) {
	if event.NextResetTime > time.Now().Unix() {
		return
	}
	eventConfig := csvs.GetEventConfig(event.EventId)
	if eventConfig == nil {
		return
	}
	switch eventConfig.RefreshType {
	case csvs.MAP_REFRESH_DAY:
		count := time.Now().Unix() / csvs.MAP_REFRESH_DAY_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_DAY_TIME
	case csvs.MAP_REFRESH_WEEK:
		count := time.Now().Unix() / csvs.MAP_REFRESH_WEEK_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_WEEK_TIME
	case csvs.MAP_REFRESH_SELF:
	case csvs.MAP_REFRESH_CANT:
		return
	}
	event.State = csvs.EVENT_START
}

func (self *ModMap) RefreshByPlayer(mapId int) {
	config := csvs.ConfigMapMap[mapId]
	if config == nil {
		return
	}
	if config.MapType != csvs.REFRESH_PLAYER {
		return
	}
	_, ok := self.MapInfo[config.MapId]
	if !ok {
		return
	}
	for _, v := range self.MapInfo[config.MapId].EventInfo {
		v.State = csvs.EVENT_START
	}
}

func (self *ModMap) NewStatue(statueId int) *StatueInfo {
	data := new(StatueInfo)
	data.StatueId = statueId
	data.Level = 0
	data.ItemInfo = make(map[int]*ItemInfo)
	return data
}
func (self *ModMap) UpStatue(statueId int, player *Player) {
	_, ok := self.Statue[statueId]
	if !ok {
		self.Statue[statueId] = self.NewStatue(statueId)
	}
	info, ok := self.Statue[statueId]
	if !ok {
		return
	}
	nextLevel := info.Level + 1
	nextConfig := csvs.GetStatueConfig(statueId, nextLevel)
	if nextConfig == nil {
		return
	}

	_, okNow := info.ItemInfo[nextConfig.CostItem]
	nowNum := int64(0)
	if okNow {
		nowNum = info.ItemInfo[nextConfig.CostItem].ItemNum
	}
	needNum := nextConfig.CostNum - nowNum

	if !player.ModBag.HasEnoughItem(nextConfig.CostItem, needNum) {
		num := player.ModBag.GetItemNum(nextConfig.CostItem, player)
		if num <= 0 {
			fmt.Println(fmt.Sprintf("神像升级物品不足"))
			return
		}
		_, okItem := info.ItemInfo[nextConfig.CostItem]
		if !okItem {
			info.ItemInfo[nextConfig.CostItem] = new(ItemInfo)
			info.ItemInfo[nextConfig.CostItem].ItemId = nextConfig.CostItem
			info.ItemInfo[nextConfig.CostItem].ItemNum = 0
		}
		_, okItem = info.ItemInfo[nextConfig.CostItem]
		if !okItem {
			return
		}
		info.ItemInfo[nextConfig.CostItem].ItemNum += num
		player.ModBag.RemoveItemToBag(nextConfig.CostItem, num, player)
		fmt.Println(fmt.Sprintf("神像升级,提交物品%d，数量%d，当前数量%d", nextConfig.CostItem, num, info.ItemInfo[nextConfig.CostItem].ItemNum))

	} else {
		player.ModBag.RemoveItemToBag(nextConfig.CostItem, needNum, player)
		info.Level++
		info.ItemInfo = make(map[int]*ItemInfo)
		fmt.Println(fmt.Sprintf("神像升级成功,神像:%d，当前等级:%d", info.StatueId, info.Level))
	}
}

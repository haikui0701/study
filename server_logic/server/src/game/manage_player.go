package game

import (
	"golang.org/x/net/websocket"
	"sync"
)

var managePlayer *ManagePlayer

type ManagePlayer struct {
	Players map[int64]*Player
	lock    *sync.RWMutex
}

func GetManagePlayer() *ManagePlayer {
	if managePlayer == nil {
		managePlayer = new(ManagePlayer)
		managePlayer.Players = make(map[int64]*Player)
		managePlayer.lock = new(sync.RWMutex)
	}
	return managePlayer
}

func (self *ManagePlayer) GetPlayer(uid int64) *Player {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.Players[uid]
}

func (self *ManagePlayer) CreatePlayer(uid int64) *Player {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.Players[uid]
	if ok {
		return self.Players[uid]
	}
	self.Players[uid] = NewTestPlayer(uid)
	return self.Players[uid]
}

func (self *ManagePlayer) PlayerLoginIn(ws *websocket.Conn, userId int64) *Player {
	playerInfo := self.GetPlayer(userId)
	if playerInfo == nil {
		playerInfo = self.CreatePlayer(userId)
	}
	playerInfo.ws = ws
	return playerInfo
}

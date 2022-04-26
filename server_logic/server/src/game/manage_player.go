package game

import (
	"fmt"
	"golang.org/x/net/websocket"
	"sync"
	"time"
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

func (self *ManagePlayer) PlayerLoginIn(ws *websocket.Conn, userId int64) *Player {
	self.lock.Lock()
	defer self.lock.Unlock()

	playerInfo, ok := self.Players[userId]
	if ok {
		//顶号
		if player.ws != ws {
			oldWs := player.ws
			playerInfo.ws = ws
			playerInfo.exitTime = 0
			if oldWs != nil {
				oldWs.Write([]byte("账号在别处登陆"))
				oldWs.Close()
			}
		}
	}

	playerInfo = NewTestPlayer(ws, userId)
	self.Players[userId] = playerInfo

	return playerInfo
}

func (self *ManagePlayer) PlayerClose(ws *websocket.Conn, userId int64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	playerInfo, ok := self.Players[userId]
	if ok {
		//顶号
		if playerInfo.ws == ws {
			playerInfo.ws = nil
			playerInfo.exitTime = time.Now().Unix() + 10
			fmt.Println("websocket连接断开,ws设置为空")
		}
	}
	return
}

func (self *ManagePlayer) Run() {
	ticker := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-ticker.C:
			self.CheckPlayerOff()
		}
	}
}

func (self *ManagePlayer) CheckPlayerOff() {
	self.lock.Lock()
	defer self.lock.Unlock()
	for k, v := range self.Players {
		if v.exitTime > time.Now().Unix() {
			fmt.Println("内存中清除角色:", v.UserId)
			delete(self.Players, k)
		}
	}
}

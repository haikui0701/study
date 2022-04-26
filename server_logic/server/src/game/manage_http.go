package game

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
	"time"
)

var manageHttp *ManageHttp

type ManageHttp struct {
}

func GetManageHttp() *ManageHttp {
	if manageHttp == nil {
		manageHttp = new(ManageHttp)
	}
	return manageHttp
}

func (self *ManageHttp) InitData() {
	http.Handle("/", websocket.Handler(self.WebsocketHandler))

	http.HandleFunc("/correctname", self.CorrectName)
}

func (self *ManageHttp) CorrectName(w http.ResponseWriter, r *http.Request) {
	player.GetModPlayer().Name = "修改任意名字"
}

func (self *ManageHttp) WebsocketHandler(ws *websocket.Conn) {
	defer ws.Close()

	var player *Player

	fmt.Println("服务器连接成功")

	for {

		var msg []byte

		ws.SetReadDeadline(time.Now().Add(3 * time.Second))
		err := websocket.Message.Receive(ws, &msg)
		fmt.Println(err)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				continue
			}
			if player != nil {
				//存档
				GetManagePlayer().PlayerClose(ws, player.UserId)
			}
			break
		}
		fmt.Println(string(msg))

		if player == nil {
			var loginMsg MsgLogin
			msgErr := json.Unmarshal(msg, &loginMsg)
			if msgErr == nil {
				player = GetManagePlayer().PlayerLoginIn(ws, loginMsg.UserId)
				go player.LogicRun()
			}
		}else{
			player.SendLogic(msg)
		}
	}
	return
}

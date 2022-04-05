package game

import (
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
			break
		}
		fmt.Println(string(msg))

		if player == nil {
			player = NewTestPlayer(10000666)
		}
	}

	return
}

package models

import (
	"github.com/gorilla/websocket"
	"logManager/common"
)

var (
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan Message)
)

func init() {
	go handlerMessage()
}

//广播发送至页面
func handlerMessage() {

	for {
		msg := <-Broadcast
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				common.Logger.Error("client writeJson err : %v ", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}

}

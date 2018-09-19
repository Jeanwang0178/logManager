package utils

import (
	"github.com/gorilla/websocket"
	"logManager/models"
)

var (
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan models.Message)
)

func init() {
	go handlerMessage()
}

//广播发送至页面
func handlerMessage() {

	for {
		msg := <-Broadcast
		Logger.Info("client len ", len(Clients))
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				Logger.Info("client writeJson err : %v ", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}

}

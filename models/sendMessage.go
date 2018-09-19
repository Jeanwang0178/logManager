package models

import (
	"github.com/beego/bee/logger"
	"github.com/gorilla/websocket"
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
		beeLogger.Log.Info("client len " + string(len(Clients)))
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				beeLogger.Log.Errorf("client writeJson err : %v ", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}

}

package models

import (
	"github.com/gorilla/websocket"
	"logManager/src/common"
	"math/rand"
)

type BroadCast struct {
	Gid     uint64
	msgChan chan string
}

var (
	Clients      = make(map[*websocket.Conn]BroadCast)
	BroadCastMap = make(map[string]BroadCast)
)

func init() {
	//go handlerMessage()
}

func NewBroadCast() BroadCast {
	broadCast := &BroadCast{
		Gid: uint64(rand.Int63()),
	}
	msgChan := make(chan string)
	broadCast.msgChan = msgChan
	return *broadCast
}

//广播发送至页面
func HandlerMessage(msgKey string, socketConn websocket.Conn) {
	broad := BroadCastMap[msgKey]
	for {
		msg := <-broad.msgChan
		err := socketConn.WriteJSON(msg)
		if err != nil {
			common.Logger.Info("client writeJson err : %v ", err)
			socketConn.Close()
			delete(Clients, &socketConn)
			delete(BroadCastMap, msgKey) //删除key
			break
		}
	}

}

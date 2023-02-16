package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
	"github.com/theone-daxia/chat-demo/config"
	"github.com/theone-daxia/chat-demo/pkg/e"
)

func (manager *ClientManager) Start() {
	for {
		logging.Println("<---监听管道通信--->")
		select {
		case conn := <-manager.Register: // 建立连接
			logging.Printf("建立新连接: %v", conn.ID)
			manager.Clients[conn.ID] = conn
			WriteMessage(conn, websocket.TextMessage, e.WebsocketSuccess, "已连接至服务器", "")
		case conn := <-manager.Unregister: // 断开连接
			logging.Printf("连接断开: %v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				WriteMessage(conn, websocket.TextMessage, e.WebsocketEnd, "连接已断开", "")
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-manager.Broadcast:
			message := broadcast.Message
			sendId := broadcast.Client.SendID
			flag := false // 默认对方不在线
			for id, conn := range manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(manager.Clients, conn.ID)
				}
			}

			var read uint = 0 // 0-未读 1-已读
			if flag {
				logging.Println("对方在线应答")
				WriteMessage(broadcast.Client, websocket.TextMessage, e.WebsocketOnlineReply, "对方在线", "")
				read = 1 // 这里简单认为对方在线就已读
			} else {
				logging.Println("对方不在线")
				WriteMessage(broadcast.Client, websocket.TextMessage, e.WebsocketOfflineReply, "对方不在线", "")
			}

			id := broadcast.Client.ID
			err := InsertMsg(config.MongoDBName, id, string(message), read, int64(3*month))
			if err != nil {
				fmt.Println("InsertMsg error: ", err)
			}
		}
	}
}

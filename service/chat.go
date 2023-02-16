package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
	"github.com/theone-daxia/chat-demo/cache"
	"github.com/theone-daxia/chat-demo/config"
	"github.com/theone-daxia/chat-demo/pkg/e"
	"net/http"
	"time"
)

const month = 60 * 60 * 24 * 30 // 30天算一个月

// SendMsg 发送的消息
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// ReplyMsg 回复的消息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// Client 用户类
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// Broadcast 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// ClientManager 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// Message 信息转JSON (包括：发送者、接收者、内容)
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func createID(uid, toUid string) string {
	return uid + "->" + toUid
}

func WsHandler(c *gin.Context) {
	uid := c.Query("uid")     // 自己的uid
	toUid := c.Query("toUid") // 对方的uid

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 创建一个用户实例
	client := &Client{
		ID:     createID(uid, toUid),
		SendID: createID(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 注册到用户管理上
	Manager.Register <- client

	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		err := c.Socket.ReadJSON(&sendMsg) // 读取json格式
		if err != nil {
			fmt.Println("数据格式不正确")
			break
		}

		if sendMsg.Type == 1 { // 1给2发消息
			r1, _ := cache.RedisClient.Get(c.ID).Result()     // 1->2
			r2, _ := cache.RedisClient.Get(c.SendID).Result() // 2->1
			if r1 > "3" && r2 == "" {                         // 限制单聊
				code := e.WebsocketLimit
				WriteMessage(c, websocket.TextMessage, code, e.GetMsg(code), "")
				// 防止骚扰，未建立连接刷新过期时间一个月
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30).Result()
				continue
			} else {
				cache.RedisClient.Incr(c.ID)
				// 防止过快"分手"，建立连接三个月过期
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30*3).Result()
			}
			logging.Println(c.ID, "发送消息", sendMsg.Content)
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 { // 拉取历史消息
			historyList := FindMany(config.MongoDBName, c.ID, c.SendID, 10)
			if len(historyList) > 10 {
				historyList = historyList[:10]
			} else if len(historyList) == 0 {
				WriteMessage(c, websocket.TextMessage, e.WebsocketEnd, "到底了", "")
				continue
			}
			for _, his := range historyList {
				WriteMessage(c, websocket.TextMessage, 0, his.Msg, his.From)
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				WriteMessage(c, websocket.CloseMessage, 0, "", "")
				return
			}
			logging.Println(c.ID, "接受消息:", string(message))
			content := fmt.Sprintf("%s", string(message))
			WriteMessage(c, websocket.TextMessage, e.WebsocketSuccessMessage, content, "")
		}
	}
}

func WriteMessage(client *Client, msgType int, code int, content string, from string) {
	if msgType == websocket.CloseMessage {
		_ = client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	replyMsg := ReplyMsg{Content: content}
	if code != 0 {
		replyMsg.Code = code
	}
	if from != "" {
		replyMsg.From = from
	}
	msg, _ := json.Marshal(replyMsg)
	_ = client.Socket.WriteMessage(msgType, msg)
	return
}

package Chat

import (
	"MyGO-IM/DB"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var H *Hub

func InitHub() {
	H = NewHub()
	H.ClientPool.New = func() any {
		H.mu.Lock()
		defer H.mu.Unlock()
		c := new(Client)
		return c
	}
	go H.Run()
	log.Println("WS集中管理器Hub创建成功")
}
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		//上线
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.Name] = client
			h.mu.Unlock()

		//下线
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.Name]; ok {
				OfflineHandler(client)
			}
			h.mu.Unlock()

		//广播
		case msg := <-h.Broadcast:
			data, _ := json.Marshal(msg)
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- data:
				default:
					// 发送缓冲区满，关闭连接
					close(client.Send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// ReadPump 持续读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // 最大消息 512KB
	// 设置心跳检测（Pong 响应）
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {

		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取错误: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println(err)
			continue
		}

		msg.FromName = c.Name
		msg.FromID = c.ID
		msg.Timestamp = time.Now().Unix()
		DB.MsgNumPlus()
		// 根据消息类型路由
		switch msg.Type {
		case "private":
			c.SendPrivate(msg)
		case "group":

			c.SendGroup(msg)
		case "system":
			c.Hub.Broadcast <- msg
		}
	}
}

// WritePump 持续向客户端写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second) // 心跳间隔
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func OfflineHandler(client *Client) {
	close(client.Close)                     //关闭通道，通知子协程死亡
	client.MQCh.Close()                     //关闭AMQP信道，防止出现僵尸消费者
	DB.SetOffline(client.Name)              //删除Redis缓存中在线记录
	DB.SetLastOnline(client.ID)             //设置下线时间
	delete(client.Hub.Clients, client.Name) //从Hub中删除对应Client连接
	close(client.Send)                      //关闭发送消息发送通道
	H.ClientPool.Put(client)                //将client实例放回池
}

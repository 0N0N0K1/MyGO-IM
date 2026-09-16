package MyGOWS

import (
	"MyGO-IM/DB"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"log"
	"sync"
)

var H *Hub
var SnowID *snowflake.Node

// InitHub 初始化并启动 Hub
func InitHub() {
	H = NewHub()

	H.ClientPool.New = func() any {
		H.mu.Lock()
		defer H.mu.Unlock()
		c := new(Client)
		return c
	}
	H.ServerMsgPool = sync.Pool{
		New: func() any {
			H.mu.Lock()
			defer H.mu.Unlock()
			c := new(ServerMessage)
			return c
		}}
	H.ClientMsgPool = sync.Pool{
		New: func() any {
			H.mu.Lock()
			defer H.mu.Unlock()
			c := new(ClientMessage)
			return c
		}}
	SnowID, _ = snowflake.NewNode(1)
	go H.Run()
	log.Println("WS集中管理器Hub创建成功")
}

// NewHub 返回一个全局 Hub 实例
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan ServerMessage, 256),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		//上线
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			OnlineHandler(client)
			log.Printf("用户%s上线", client.Name)

		//下线
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				OfflineHandler(client)
			}
			h.mu.Unlock()
			log.Printf("用户%s下线", client.Name)
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

// OfflineHandler 下线后的处理
func OfflineHandler(client *Client) {
	close(client.Close)                   //关闭通道，通知子协程死亡
	client.MQCh.Close()                   //关闭AMQP信道，防止出现僵尸消费者
	DB.SetOffline(client.ID)              //删除Redis缓存中在线记录
	DB.SetLastOnline(client.ID)           //设置下线时间
	delete(client.Hub.Clients, client.ID) //从Hub中删除对应Client连接
	close(client.Send)                    //关闭发送消息发送通道
	H.ClientPool.Put(client)              //将client实例放回池
}

// OnlineHandler 注册上线后的处理
func OnlineHandler(client *Client) {
	//Redis进行缓存上线记录
	DB.SetOnline(client.ID)
	//开辟 读+写goroutine、消费消息的goroutine
	go client.WritePump()
	go client.ReadPump()
	go client.ConsumeMyQueue()
	go client.ConfirmHandler()
}

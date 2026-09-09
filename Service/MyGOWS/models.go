package MyGOWS

import (
	"MyGO-IM/Conf"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"net/http"
	"sync"
)

type Message struct {
	Type      string `json:"type"` // 消息类型 "private" | "system" | "group"
	FromName  string `json:"from_name"`
	FromID    uint   `json:"form_id"` // 发送者ID
	ToName    string `json:"to_name"`
	ToID      uint   `json:"to_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"` // 发送时间戳
}
type Client struct {
	Close chan struct{}
	ID    uint
	Name  string
	Conn  *websocket.Conn  // WS连接
	Queue amqp091.Queue    // 客户端持有的队列
	MQCh  *amqp091.Channel // 客户端持有的AMQP信道
	Send  chan []byte      // 发送消息队列
	Hub   *Hub
}

// Hub 管理所有客户端连接
type Hub struct {
	Clients    map[string]*Client // 在线客户端集合
	Register   chan *Client       // 连接注册Chan
	Unregister chan *Client       // 连接断开Chan
	Broadcast  chan Message       // 广播消息Chan
	mu         sync.RWMutex       // 保护 Clients 集合
	ClientPool sync.Pool
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  Conf.Conf.W.ReadBufferSize,
	WriteBufferSize: Conf.Conf.W.WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应校验域名
	},
}

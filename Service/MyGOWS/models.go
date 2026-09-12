package MyGOWS

import (
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"sync"
)

const (
	private = "private"
	system  = "system"
	group   = "group"
)

type Message struct {
	FromName string `json:"from_name"`
	FromID   uint   `json:"form_id"`
	ToName   string `json:"to_name"`
	ToID     uint   `json:"to_id"`
	Method   string `json:"method"` //  "private" | "system" | "group"

	Type      string `json:"type"` //TODO 扩展消息类型
	Payload   string `json:"Payload"`
	Timestamp int64  `json:"timestamp"`

	Seq   uint  `json:"seq"`
	MsgID int64 `json:"msg_id"`
}
type Client struct {
	Close    chan struct{}
	ID       uint
	Name     string
	Confirm  chan amqp091.Confirmation
	Conn     *websocket.Conn  // WS连接
	Queue    amqp091.Queue    // 客户端持有的队列
	MQCh     *amqp091.Channel // 客户端持有的AMQP信道
	Send     chan []byte      // 发送消息队列
	Hub      *Hub
	SendFunc func(msg *Message)
}

// Hub 管理所有客户端连接
type Hub struct {
	Clients    map[uint]*Client // 在线客户端集合
	Register   chan *Client     // 连接注册Chan
	Unregister chan *Client     // 连接断开Chan
	Broadcast  chan Message     // 广播消息Chan
	mu         sync.RWMutex     // 保护 Clients 集合
	ClientPool sync.Pool
	MsgPool    sync.Pool
}

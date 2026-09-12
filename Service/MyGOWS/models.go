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

	Seq   uint64 `json:"seq"`
	MsgID int64  `json:"msg_id"`
}
type Client struct {
	Close       chan struct{}
	ID          uint
	Name        string
	Confirm     chan amqp091.Confirmation
	SendReady   chan struct{}    //得到confirm并更新seq后通知可以发送下一条消息的Chan
	AckReady    chan bool        //消息拉取Send完成（成功/失败）后通知consume Ack/Nack的Chan
	ExpertSeq   uint             //本次希望consume并send的消息seq
	DisorderMag []Message        // 存放乱序到达的Msg
	Conn        *websocket.Conn  // WS连接
	Queue       amqp091.Queue    // 客户端持有的队列
	MQCh        *amqp091.Channel // 客户端持有的AMQP信道
	Send        chan []byte      // 发送消息队列
	Hub         *Hub
	SendFunc    func(msg *Message) bool
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

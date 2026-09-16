package MyGOWS

import (
	"github.com/gorilla/websocket"
	"sync"
)

const (
	private = "private"
	system  = "system"
	group   = "group"
)

type ServerACK struct {
	Cmd            string `json:"cmd"`
	ConversationId string `json:"conversation_id"`
	MsgID          int64  `json:"msg_id"`
	Resend         bool   `json:"resend"`
	Desc           string `json:"describe"`
	ToId           uint   `json:"to_id"`
	ReplyID        int64  `json:"reply_id"`
	Method         string `json:"method"`
}
type ClientMessage struct {
	Cmd            string `json:"cmd"`
	MsgId          int64  `json:"msg_id"`
	Method         string `json:"method"`
	ConversationId string `json:"conversation_id"`
	ToId           uint   `json:"to_id"`
	Type           string `json:"type"`
	Payload        string `json:"payload"`
	ReplyID        int64  `json:"reply_id"`
}

type ServerMessage struct {
	Cmd            string `json:"cmd"`
	FromID         uint   `json:"form_id"`
	ToID           uint   `json:"to_id"`
	Method         string `json:"method"` //  "private" | "system" | "group"
	ConversationId string `json:"conversation_id"`
	Type           string `json:"type"`
	Payload        string `json:"Payload"`
	Timestamp      int64  `json:"timestamp"`
	Seq            uint64 `json:"seq"`
	MsgID          int64  `json:"msg_id"`
	ReplyID        int64  `json:"reply_id"`
}

type Client struct {
	Close       chan struct{}
	ID          uint
	Name        string
	AckReady    chan struct{}   //消息拉取Send完成（成功/失败）后通知consume Ack/Nack的Chan
	ExpertSeq   uint            //本次希望consume并send的消息seq
	DisorderMag []ServerMessage // 存放乱序到达的Msg
	Conn        *websocket.Conn // WS连接
	Send        chan []byte     // 发送消息队列
	Hub         *Hub
	SendFunc    func(msg *ServerMessage) // 发送消息用的函数
}

// Hub 管理所有客户端连接
type Hub struct {
	Clients       map[uint]*Client   // 在线客户端集合
	Register      chan *Client       // 连接注册Chan
	Unregister    chan *Client       // 连接断开Chan
	Broadcast     chan ServerMessage // 广播消息Chan
	mu            sync.RWMutex       // 保护 Clients 集合
	ClientPool    sync.Pool          // Client 池
	ClientMsgPool sync.Pool          // ClientMessage 池
	ServerMsgPool sync.Pool          // ClientMessage 池
}

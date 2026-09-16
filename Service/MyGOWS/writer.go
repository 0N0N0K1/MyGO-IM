package MyGOWS

import (
	"MyGO-IM/DB"
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

// WritePump 持续向客户端写消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(27 * time.Second) // 心跳间隔
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case <-c.Close:
			return
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.WriteHandler(message)
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WriteHandler 向客户端发送消息的处理函数
func (c *Client) WriteHandler(message []byte) {
	msg := c.Hub.ServerMsgPool.Get().(*ServerMessage)
	defer c.Hub.ServerMsgPool.Put(msg)
	err := json.Unmarshal(message, msg)
	if err != nil {
		return
	}
	switch msg.Method {
	case "system":
		_ = c.Conn.WriteMessage(websocket.TextMessage, message)
		return
	case "group":
		readSeq := DB.GetGroupReadSeq(msg.ToID, c.ID)
		if readSeq > msg.Seq {
			return
		}
		err = c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
		DB.IncrGroupReadSeq(msg.ToID, c.ID, DB.GetGroupReadSeq(msg.ToID, c.ID)+1)

	case "private":
		readSeq, _ := DB.GetPrivateSeq(msg.ToID, msg.FromID)
		if readSeq > msg.Seq {
			return
		}
		err = c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
		DB.IncrPrivateReadSeq(msg.ToID, msg.FromID, msg.Seq+1)
	}

}

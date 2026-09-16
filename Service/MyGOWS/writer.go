package MyGOWS

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
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
		c.AckReady <- struct{}{}
		return
	}
	switch msg.Method {
	case "system":
		_ = c.Conn.WriteMessage(websocket.TextMessage, message)
		return
	case "group":
		if !Utils.Dedup(c.ID, msg.MsgID, msg.Seq) {
			log.Printf("repeat")
			c.AckReady <- struct{}{}
			return
		}
		readSeq := DB.GetGroupReadSeq(msg.ToID, c.ID)
		if readSeq > msg.Seq {
			c.AckReady <- struct{}{}
			return
		}
		err = c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			c.AckReady <- struct{}{}
			return
		}
		DB.IncrGroupReadSeq(msg.ToID, c.ID, DB.GetGroupReadSeq(msg.ToID, c.ID)+1)

	case "private":
		if !Utils.Dedup(msg.ToID, msg.MsgID, msg.Seq) {
			log.Printf("repeat")
			c.AckReady <- struct{}{}
			return
		}

		err = c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			c.AckReady <- struct{}{}
			return
		}
		DB.IncrPrivateReadSeq(msg.ToID, msg.FromID, msg.Seq+1)
	}
	c.AckReady <- struct{}{}

}

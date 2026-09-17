package MyGOHTTP

import (
	"MyGO-IM/DB"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetMessage 查询聊天记录
// 通过 query 指定
// 1.查询聊天对象 群/好友 ID
// 2.查询方式，双方消息/发送消息
func GetMessage(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	groupID := c.Query("group")
	friendID := c.Query("friend")
	page := c.Query("page")
	limit := c.Query("limit")
	mode := c.Query("mode")
	lm, _ := strconv.Atoi(limit)
	pg, _ := strconv.Atoi(page)
	if groupID != "" {
		gID, err := strconv.Atoi(groupID)
		if err != nil || gID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "ID invalid"})
			return
		}
		result, err := DB.QueryGroupMsg(int64(gID), actorID.(int64), pg, lm, DB.MsgQueryMode(mode))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "QueryGroupMemberMsg error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "messages": result})
		return
	}
	if friendID != "" {
		fID, err := strconv.Atoi(friendID)
		if err != nil || fID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "ID invalid"})
			return
		}
		result, err := DB.QueryPrivateMsg(int64(fID), actorID.(int64), pg, lm, DB.MsgQueryMode(mode))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "QueryGroupMemberMsg error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "messages": result})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1112, "error": "no appointed ID"})
	return
}

// PullMsgByCID 点开会话时拉取聊天记录
func PullMsgByCID(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	conversationID := c.Query("conversationID")
	limit := c.Query("limit")
	endSeq := c.Query("endSeq")
	lm, err := strconv.Atoi(limit)
	if err != nil || lm <= 0 {
		lm = 20
	}
	Seq, seqerr := strconv.ParseUint(endSeq, 10, 64)

	if []byte(conversationID)[0] == 'g' {
		if DB.QueryMemberByCID(actorID.(int64), conversationID).GroupID == 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "You aren' the member "})
			return
		}
		if seqerr != nil {
			Seq = DB.GetGroupWriterSeqByCID(conversationID)
		}
		result, err := DB.PullGroupMsg(conversationID, lm, Seq)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "QueryGroupMemberMsg error"})
			return
		}
		if DB.GetGroupReadSeqByCID(conversationID, actorID.(int64)) < Seq {
			DB.IncrGroupReadSeqByCID(conversationID, actorID.(int64), Seq)
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "messages": result})
		return
	}
	if []byte(conversationID)[0] == 'p' {
		frd, _ := DB.QueryFrd(actorID.(int64), conversationID, DB.ByCID)
		//是否已添加
		if len(frd) == 0 {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "对方不是你的好友"})
			return
		}
		readSeq, writeSeq := DB.GetPrivateSeqByCID(actorID.(int64), conversationID)
		if seqerr != nil {
			Seq = writeSeq
		}
		result, err := DB.PullPrivateMsg(conversationID, lm, Seq)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "QueryGroupMemberMsg error"})
			return
		}
		log.Println(readSeq, writeSeq)
		if readSeq < Seq {
			DB.IncrPrivateReadSeqByCID(actorID.(int64), conversationID, Seq-1)
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "messages": result})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1112, "error": "conversationID invalid"})
	return
}

package Service

import (
	"MyGO-IM/Chat"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func WSUpgrade(c *gin.Context) {
	//升级HTTP为WS
	conn, err := Chat.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": "102", "msg": "WS UP ERROR!"})
		c.Abort()
		return
	}
	userID, _ := c.Get("ID")
	userName, _ := c.Get("name")

	//初始化聊天
	if err = Chat.InitChat(userName.(string), userID.(uint), conn); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "102", "msg": "TRY AGAIN!"})
		c.Abort()
		return
	}

}

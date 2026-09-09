package MyGOHTTP

import (
	"MyGO-IM/Service/MyGOWS"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func WSUpgrade(c *gin.Context) {
	//升级HTTP为WS
	conn, err := MyGOWS.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": "102", "msg": "WS UP ERROR!"})
		c.Abort()
		return
	}
	actorID, _ := c.Get("actorID")
	actorName, _ := c.Get("actorName")

	//初始化聊天
	if err = MyGOWS.InitChat(actorName.(string), actorID.(uint), conn); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "102", "msg": "TRY AGAIN!"})
		c.Abort()
		return
	}

}

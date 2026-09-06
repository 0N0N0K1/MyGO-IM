package Service

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
)

type JWToken struct {
	Token string
}

func VerifyJWT(c *gin.Context) {
	var JWTRecv JWToken
	JWTRecv.Token = c.GetHeader("JWT")
	username, ok := Utils.VerifyJWT(JWTRecv.Token)
	nameParam := c.Param("name")
	if !ok || nameParam != username {
		c.JSON(http.StatusOK, gin.H{"code": "101", "msg": "NO Auth!"})
		c.Abort()
		return
	}
	user, err := DB.QueryUser(username, DB.ByName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "101", "msg": err.Error()})
		c.Abort()
		return
	}
	c.Set("ID", user.ID)
	c.Set("name", username)
	if ban, reason, t := DB.QueryIfBan(strconv.Itoa(int(user.ID))); ban {
		c.JSON(http.StatusOK, gin.H{"code": "101", "msg": "Banned!!!", "reason": reason, "ban-time": t})
		c.Abort()
		return
	}
	c.Next()
}

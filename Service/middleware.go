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

// VerifyJWT 检查JWT鉴权
func VerifyJWT(c *gin.Context) {
	var JWTRecv JWToken
	JWTRecv.Token = c.GetHeader("JWT")
	userID, ok := Utils.VerifyJWT(JWTRecv.Token)
	nameParam := c.Param("name")
	uid, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "Int2string: " + err.Error()})
		return
	}
	users, err := DB.QueryUser(uid, DB.ByID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "QueryUser: " + err.Error()})
		return
	}
	if !ok || nameParam != users[0].Name {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "NO Auth!"})
		c.Abort()
		return
	}
	c.Set("ID", uint(uid))
	c.Set("nickname", nameParam)
	if ban, reason, t := DB.QueryIfBan(userID); ban {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "Banned!!!", "reason": reason, "ban-time": t})
		c.Abort()
		return
	}
	c.Next()
}

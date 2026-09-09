package MyGOHTTP

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
	id := c.Param("ID")
	uid, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "NO Auth! Int2string: " + err.Error()})
		c.Abort()
		return
	}
	users, err := DB.QueryUser(uid, DB.ByID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "QueryUser: " + err.Error()})
		c.Abort()
		return
	}
	if !ok || id != userID {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "NO Auth!"})
		c.Abort()
		return
	}
	if ban, reason, t := DB.QueryIfBan(userID); ban {
		c.JSON(http.StatusOK, gin.H{"code": 101, "error": "Banned!!!", "reason": reason, "ban-time": t})
		c.Abort()
		return
	}
	c.Set("actorID", uint(uid))
	c.Set("actorName", users[0].Name)

	c.Next()
}
func GroupMiddlewarw(c *gin.Context) {
	gID := c.Param("gID")
	if gID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "No groupID input"})
		c.Abort()
		return
	}
	groupID, err := strconv.Atoi(gID)
	if err != nil || groupID <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "String2int error"})
		c.Abort()
		return
	}
	group := DB.QueryGroup(uint(groupID))
	c.Set("gID", group.ID)
	c.Set("gName", group.GroupName)
	c.Set("ownerName", group.OwnerName)
	c.Set("ownerID", group.OwnerID)
	c.Next()
}

package API

import (
	"github.com/gin-gonic/gin"
	"mygoim/Service"
	"net/http"
)

var UserEngine *gin.Engine
var AdminEngine *gin.Engine

func InitRouters() {
	UserEngine = gin.Default()
	AdminEngine = gin.Default()
	UserEngine.POST("/login", Service.Login)
	UserEngine.POST("/register", Service.Register)

	// 用户组API
	User := UserEngine.Group("/users/:name")
	User.Use(Service.VerifyJWT)
	{
		User.GET("/friends", Service.GetFriends)
		User.GET("/friends/:friendname")
		User.POST("/friends/:friendname", Service.AddFriend)

		User.POST("/groups/:groupname", Service.CreateGroup)
		User.PATCH("/groups/:groupname", Service.EnterGroup)
		User.GET("/groups/:groupname", Service.GetMembers)

		User.GET("/chat", Service.WSUpgrade)

		// user test API
		User.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": "999", "msg": "Hello!"})
		})
	}

	// 管理员API

	Admin := AdminEngine.Group("/admin")

	{
		Admin.GET("/chat", Service.WSUpgrade)
		Admin.GET("/basicinfo", Service.GetBasicInfo)
		Admin.GET("/message", Service.QueryMessage)

		UserManamge := Admin.Group("/users")
		{
			UserManamge.GET("/password", Service.ChangePassword)
			UserManamge.GET("/ban", Service.BanUser)
			UserManamge.GET("/info", Service.QueryUser)
		}
	}
}

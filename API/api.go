package API

import (
	"MyGO-IM/Service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var UserEngine *gin.Engine
var AdminEngine *gin.Engine

func InitRouters() {

	UserEngine = gin.Default()
	AdminEngine = gin.Default()

	UserEngine.POST("/login", Service.Login)
	UserEngine.POST("/register", Service.Register)
	UserEngine.POST("/send-code", Service.AuthCode)

	// 用户组API
	User := UserEngine.Group("/users/:name")
	User.Use(Service.VerifyJWT)
	{
		User.GET("/friends", Service.GetFriends)
		User.POST("/friends", Service.AddFriend)
		User.DELETE("/friends", Service.DeleteFriend)

		User.POST("/groupse", Service.CreateGroup)
		User.PATCH("/groups/:groupname", Service.EnterGroup)
		User.GET("/groups/:groupname", Service.GetMembers)

		User.GET("/chat", Service.WSUpgrade)

		//  test API
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
		// 用户管理API
		UserManamge := Admin.Group("/users")
		{
			UserManamge.GET("/password", Service.ChangePassword)
			UserManamge.GET("/ban", Service.BanUser)
			UserManamge.GET("/info", Service.QueryUser)
		}
	}
}

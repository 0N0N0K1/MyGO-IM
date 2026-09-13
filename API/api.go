package API

import (
	"MyGO-IM/Service/MyGOHTTP"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ServiceEngine 负责user服务的Engine
var ServiceEngine *gin.Engine

// AdminEngine 负责admin管理的Engine
var AdminEngine *gin.Engine

func InitRouters() {

	ServiceEngine = gin.New()
	AdminEngine = gin.Default()

	// 登录 API
	ServiceEngine.POST("/login", MyGOHTTP.Login)
	// 注册 API
	ServiceEngine.POST("/register", MyGOHTTP.Register)
	// 发送验证码 API
	ServiceEngine.POST("/send-code", MyGOHTTP.AuthCode)

	// 用户服务路由组
	//TODO 获取聊天未读数
	Users := ServiceEngine.Group("/users/:ID")
	Users.Use(MyGOHTTP.VerifyJWT)
	{
		// 账号注销 API
		Users.DELETE("/", MyGOHTTP.DeleteUser)
		// 更新用户密码
		Users.PATCH("/password", MyGOHTTP.ChangePassword)
		// 更新用户信息 API
		Users.PATCH("/info", MyGOHTTP.UpdateInfo)
		// 更新用户name API
		Users.PATCH("/nickname", MyGOHTTP.UpdateName)

		// 查询好友 API
		Users.GET("/friends", MyGOHTTP.GetFriends)
		// 添加好友 API
		Users.POST("/friends", MyGOHTTP.AddFriend)
		// 删除好友 API
		Users.DELETE("/friends", MyGOHTTP.DeleteFriend)
		// 升级为 WS 拉取消息 API
		Users.GET("/chat", MyGOHTTP.WSUpgrade)

		//  test connection API
		Users.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": "999", "msg": "Hello!"})
		})

		// 群聊路由组
		Groups := Users.Group("/groups")
		{
			// 创建群聊 API
			Groups.POST("/", MyGOHTTP.CreateGroup)

			GroupsWithMW := Groups.Group("/:gID")
			GroupsWithMW.Use(MyGOHTTP.GroupMiddleware)
			{
				// 销毁群聊 API
				GroupsWithMW.DELETE("/", MyGOHTTP.DropMyGroup)
				// 进入群聊 API
				GroupsWithMW.DELETE("/members", MyGOHTTP.ExitGroup)
				// 退出群聊 API
				GroupsWithMW.POST("/members", MyGOHTTP.EnterGroup)

				// 踢出成员 API
				GroupsWithMW.DELETE("/owner", MyGOHTTP.KickOutMember)
				// 禁言成员 API
				GroupsWithMW.PATCH("/owner", MyGOHTTP.DoNotSpeak)
				// 查询成员 API
				GroupsWithMW.GET("/:gID/members", MyGOHTTP.GetMembers)
			}
		}
	}

	// admin 路由组
	Admin := AdminEngine.Group("/admin")

	{
		// 发送系统广播 API
		Admin.GET("/broadcast", MyGOHTTP.WSUpgrade)
		// 获取服务情报统计 API
		Admin.GET("/intelligence", MyGOHTTP.GetBasicInfo)
		// 查询用户消息 API
		Admin.GET("/message", MyGOHTTP.QueryMessage)

		// 用户管理路由组
		UserManamge := Admin.Group("/users")
		{
			// 修改用户密码 API
			UserManamge.PATCH("/password", MyGOHTTP.ChangePassword)
			// 禁言用户 API
			UserManamge.POST("/ban", MyGOHTTP.BanUser)
			// 查询用户信息 API
			UserManamge.GET("/info", MyGOHTTP.QueryUser)
		}
	}
}

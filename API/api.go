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

	// 登录
	ServiceEngine.POST("/login", MyGOHTTP.Login)
	// 注册
	ServiceEngine.POST("/register", MyGOHTTP.Register)
	// 发送验证码
	ServiceEngine.POST("/send-code", MyGOHTTP.AuthCode)

	// 用户服务路由组
	Users := ServiceEngine.Group("/users/:ID")
	Users.Use(MyGOHTTP.VerifyJWT)
	{
		{
			// 账号注销
			Users.DELETE("/", MyGOHTTP.DeleteUser)
			// 更新用户密码
			Users.PATCH("/password", MyGOHTTP.ChangePassword)
			// 更新用户信息
			Users.PATCH("/info", MyGOHTTP.UpdateInfo)
			// 更新用户name
			Users.PATCH("/nickname", MyGOHTTP.UpdateName)
		}
		{
			// 查询好友
			Users.GET("/friends", MyGOHTTP.GetFriends)
			// 添加好友
			Users.POST("/friends", MyGOHTTP.AddFriend)
			// 删除好友
			Users.DELETE("/friends", MyGOHTTP.DeleteFriend)
		}
		{
			// 拉取系统通知
			Users.GET("/notices",MyGOHTTP.PullNotice)
			// 拉取聊天记录
			Users.GET("/message", MyGOHTTP.PullMsgByCID)
			// 查询聊天记录 <弃用>
			Users.GET("/<message>", MyGOHTTP.GetMessage)
		}
		// 升级为 WS 获得推送消息
		Users.GET("/chat", MyGOHTTP.WSUpgrade)

		//  test connection
		Users.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": "999", "msg": "Hello!"})
		})

		// 群聊路由组
		Groups := Users.Group("/groups")
		{
			// 创建群聊
			Groups.POST("/", MyGOHTTP.CreateGroup)
			// 查看加入的群聊
			Groups.GET("/", MyGOHTTP.GetGroups)

			GroupsWithMW := Groups.Group("/:gID")
			GroupsWithMW.Use(MyGOHTTP.GroupMiddleware)
			{
				// 销毁群聊
				GroupsWithMW.DELETE("/", MyGOHTTP.DropMyGroup)

				// 邀请进入群聊
				GroupsWithMW.POST("/invitation", MyGOHTTP.)

				// 退出群聊
				GroupsWithMW.DELETE("/members", MyGOHTTP.ExitGroup)
				// 进入群聊
				GroupsWithMW.POST("/members", MyGOHTTP.EnterGroup)
				// 查询群成员
				GroupsWithMW.GET("/members", MyGOHTTP.GetMembers)

				// 踢出成员
				GroupsWithMW.DELETE("/owner", MyGOHTTP.KickOutMember)
				// 禁言成员
				GroupsWithMW.PATCH("/owner", MyGOHTTP.DoNotSpeak)

			}
		}
	}

	// admin 路由组
	Admin := AdminEngine.Group("/admin")

	{
		// 发送系统广播
		Admin.GET("/broadcast", MyGOHTTP.WSUpgrade)
		// 获取服务情报统计
		Admin.GET("/intelligence", MyGOHTTP.GetBasicInfo)
		// 查询用户消息
		Admin.GET("/message", MyGOHTTP.QueryMessage)

		// 用户管理路由组
		UserManamge := Admin.Group("/users")
		{
			// 修改用户密码
			UserManamge.PATCH("/password", MyGOHTTP.ChangePassword)
			// 禁言用户
			UserManamge.POST("/ban", MyGOHTTP.BanUser)
			// 查询用户信息
			UserManamge.GET("/info", MyGOHTTP.QueryUser)
		}
	}
}

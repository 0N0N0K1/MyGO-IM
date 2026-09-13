package main

import (
	"MyGO-IM/API"
	"MyGO-IM/Conf"
	"MyGO-IM/DB"
	"MyGO-IM/Service/MyGOWS"
	"log"
)

// TODO 解决群聊获取/修改Read_seq与Writer_seq时的并发问题——————MySQL事务
func main() {

	Conf.LoadConfig()
	DB.InitMySQL()
	DB.InitRedis()
	MyGOWS.InitHub()
	MyGOWS.InitMQ()
	API.InitRouters()
	go func() {
		err := API.ServiceEngine.Run(Conf.UserAddr)
		if err != nil {
			log.Fatal("Engine.Run: ", err)
		}
	}()
	err := API.AdminEngine.Run(Conf.AdminAddr)
	if err != nil {
		log.Fatal("Engine.Run: ", err)
	}
}

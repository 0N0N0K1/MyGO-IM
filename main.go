package main

import (
	"MyGO-IM/API"
	"MyGO-IM/Chat"
	"MyGO-IM/Conf"
	"MyGO-IM/DB"
	"log"
)

func main() {
	Conf.LoadConfig()
	DB.InitMySQL()
	DB.InitRedis()
	Chat.InitHub()
	Chat.InitMQ()
	API.InitRouters()
	go func() {
		err := API.UserEngine.Run(Conf.UserAddr)
		if err != nil {
			log.Fatal("Engine.Run: ", err)
		}
	}()
	err := API.AdminEngine.Run(Conf.AdminAddr)
	if err != nil {
		log.Fatal("Engine.Run: ", err)
	}
}

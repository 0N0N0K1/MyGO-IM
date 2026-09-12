package main

import (
	"MyGO-IM/API"
	"MyGO-IM/Conf"
	"MyGO-IM/DB"
	WS2 "MyGO-IM/Service/MyGOWS"
	"log"
)

func main() {
	Conf.LoadConfig()
	DB.InitMySQL()
	DB.InitRedis()
	WS2.InitHub()
	WS2.InitMQ()
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

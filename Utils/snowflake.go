package Utils

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var SnowID *snowflake.Node

func InitSnowflake() {
	var err error
	SnowID, err = snowflake.NewNode(1)
	if err != nil {
		log.Fatal("snowflake初始化失败")
	}
}

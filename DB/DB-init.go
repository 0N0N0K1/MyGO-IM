package DB

import (
	"MyGO-IM/Conf"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

// MySQL 全局实例
var MySQL *gorm.DB

// RDB redis全局实例
var RDB *redis.Client

// InitRedis 初始化与MySQL连接
func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     Conf.Conf.R.Addr,
		Password: Conf.Conf.R.Password,
		DB:       Conf.Conf.R.DB,
		PoolSize: Conf.Conf.R.PoolSize,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis连接失败:", err)
	}
	RDB.HDel(context.TODO(), "online")
	log.Println("Redis连接成功:")

}

// InitMySQL 初始化与Redis连接
func InitMySQL() {
	var err error
	MySQL, err = gorm.Open(mysql.Open(Conf.MySQLDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("MySQL连接失败:", err)
	}
	sqlDB, err := MySQL.DB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MySQL连接成功")
	sqlDB.SetMaxIdleConns(Conf.Conf.M.MaxIdleConns)
	sqlDB.SetMaxOpenConns(Conf.Conf.M.MaxOpenConns)
	err = MySQL.AutoMigrate(
		&User{},
		&UserInfo{},
		&Group{},
		&UserUser{},
		&GroupUser{},
		&GroupMessage{},
		&PrivateMessage{},
	)
	if err != nil {
		log.Fatal("建表失败:", err)
	}

}

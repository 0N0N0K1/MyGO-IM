package Conf

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	M      MySQLConfig `mapstructure:"mysql"`
	R      RedisConfig `mapstructure:"redis"`
	J      JWTConfig   `mapstructure:"JWT"`
	W      WSConfig    `mapstructure:"WS"`
	MQ     MQConfig    `mapstructure:"rabbitMQ"`
	Web    WebConfig   `mapstructure:"web"`
	Worker Worker      `mapstructure:"worker"`
}
type Worker struct {
	Consumer int `mapstructure:"consumer"`
	Producer int `mapstructure:"producer"`
}
type WebConfig struct {
	Host      string `mapstructure:"host"`
	UserPort  string `mapstructure:"userport"`
	AdminPort string `mapstructure:"adminport"`
}
type MQConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
type WSConfig struct {
	ReadBufferSize  int `mapstructure:"ReadBufferSize"`
	WriteBufferSize int `mapstructure:"WriteBufferSize"`
}
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	TTL    int    `mapstructure:"ttl"`
}
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"db"`
	MaxIdleConns int    `mapstructure:"MaxIdleConns"`
	MaxOpenConns int    `mapstructure:"MaxOpenConns"`
}
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"PoolSize"`
}

var Conf Config
var MySQLDSN string
var MQURL string
var UserAddr string
var AdminAddr string

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./Conf")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("ReadInConfig error:", err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatal("Config Unmarshal error:", err)
	}

	MySQLDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Conf.M.User, Conf.M.Password, Conf.M.Host, Conf.M.Port, Conf.M.DB)
	MQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", Conf.MQ.User, Conf.MQ.Password, Conf.MQ.Host, Conf.MQ.Port)
	UserAddr = fmt.Sprintf("%s:%s", Conf.Web.Host, Conf.Web.UserPort)
	AdminAddr = fmt.Sprintf("%s:%s", Conf.Web.Host, Conf.Web.AdminPort)
	log.Println("读取配置成功")
}

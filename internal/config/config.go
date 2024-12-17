package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Config Conf

type Conf struct {
	GoogleClientId     string
	GoogleClientSecret string
	MySqlConnectString string
	RedisUrl           string
	RedisPassword      string
}

func init() {
	err := godotenv.Load("/Users/winter_wang/Project/chat/.env")
	if err != nil {
		panic(err)
	}
	Config = Conf{
		GoogleClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		MySqlConnectString: os.Getenv("MYSQL_CONNECT_STRING"),
		RedisUrl:           os.Getenv("REDIS_URL"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD")}
}

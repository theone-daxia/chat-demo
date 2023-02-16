package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strconv"
)

var (
	RedisClient *redis.Client
	RedisDB     string
	RedisAddr   string
	RedisPwd    string
	RedisDBName string
)

func init() {
	file, err := ini.Load("./config/config.ini") // 加载配置文件
	if err != nil {
		fmt.Println("redis ini load filed: ", err)
	}
	LoadRedis(file) // 读取配置文件
	Redis()         // redis 连接
}
func LoadRedis(file *ini.File) {
	RedisDB = file.Section("redis").Key("RedisDB").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPwd = file.Section("redis").Key("RedisPwd").String()
	RedisDBName = file.Section("redis").Key("RedisDBName").String()
}

func Redis() {
	db, _ := strconv.ParseUint(RedisDBName, 10, 64) // string to uint64
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   int(db),
	})
	_, err := client.Ping().Result()
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	RedisClient = client
}

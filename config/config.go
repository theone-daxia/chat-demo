package config

import (
	"context"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"github.com/theone-daxia/chat-demo/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
)

var (
	MongoDBClient *mongo.Client
	AppMode       string
	HttpPort      string

	DB         string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisDB     string
	RedisAddr   string
	RedisPwd    string
	RedisDBName string

	MongoDBName string
	MongoDBAddr string
	MongoDBPwd  string
	MongoDBPort string
)

func Init() {
	// 加载翻译功能
	if err := LoadLocales("config/locales/zh-cn.yaml"); err != nil {
		logging.Info(err)
		panic(err)
	}

	// 读取配置文件
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径：", err)
	}
	LoadServer(file)
	LoadMySQL(file)
	LoadMongoDB(file)

	// MySQL 连接
	dsn := DBUser + ":" + DBPassword +
		"@tcp(" + DBHost + ":" + DBPort +
		")/" + DBName + "?charset=utf8mb4&parseTime=true&loc=Local"
	model.DataBase(dsn)
	MongoDB() // mongodb 连接
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMySQL(file *ini.File) {
	DB = file.Section("mysql").Key("DB").String()
	DBHost = file.Section("mysql").Key("DBHost").String()
	DBPort = file.Section("mysql").Key("DBPort").String()
	DBUser = file.Section("mysql").Key("DBUser").String()
	DBPassword = file.Section("mysql").Key("DBPassword").String()
	DBName = file.Section("mysql").Key("DBName").String()
}

func LoadMongoDB(file *ini.File) {
	MongoDBName = file.Section("mongodb").Key("MongoDBName").String()
	MongoDBAddr = file.Section("mongodb").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("mongodb").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("mongodb").Key("MongoDBPort").String()
}

func MongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)
	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	logging.Info("mongodb connect successfully")
}

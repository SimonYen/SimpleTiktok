package database

import (
	"app/config"
	"app/models"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局变量，用于操作数据库
var Handler *gorm.DB
var RDB *redis.Client
var CTX context.Context
var MON *mongo.Client

func init() {
	db, err := gorm.Open(mysql.Open(config.Mysql.DSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败：%v", err))
	}
	//自动迁移
	db.AutoMigrate(&models.User{}, &models.Video{}, &models.Comment{})
	Handler = db.Debug()

	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	RDB = r
	CTX = ctx
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(c, options.Client().ApplyURI(config.Mongo.URI()))
	if err != nil {
		panic(err)
	}
	MON = client
}

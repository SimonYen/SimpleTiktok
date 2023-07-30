package database

import (
	"app/config"
	"app/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局变量，用于操作数据库
var Handler *gorm.DB

func init() {
	db, err := gorm.Open(mysql.Open(config.Mysql.DSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败：%v", err))
	}
	//自动迁移
	db.AutoMigrate(&models.User{}, &models.Video{})
	Handler = db
}

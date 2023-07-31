package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type server struct {
	Port int
	Host string
}

type mysql struct {
	Port     int
	Host     string
	User     string
	Password string
	Database string
}

// 全局变量，用于保存viper读取到的
var Server *server
var Mysql *mysql

func init() {
	Server = new(server)
	Mysql = new(mysql)
	viper.SetConfigName("app")  //设置配置文件名
	viper.SetConfigType("toml") //设置配置文件后缀
	viper.AddConfigPath(".")    //设置配置文件路径
	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		//直接退出整个程序，因为没有运行的必要了
		panic(fmt.Errorf("配置文件读写出错：%v", err))
	}
	//填入读取到的值
	Server.Port = viper.GetInt("server.port")
	Server.Host = viper.GetString("server.host")

	Mysql.Database = viper.GetString("mysql.database")
	Mysql.Host = viper.GetString("mysql.host")
	Mysql.Password = viper.GetString("mysql.password")
	Mysql.Port = viper.GetInt("mysql.port")
	Mysql.User = viper.GetString("mysql.user")
}

// Mysql配置struct返回dsn字符串，便于之后连接数据库
func (m *mysql) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.Database)
}

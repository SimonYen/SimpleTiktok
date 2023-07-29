package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// 全局变量，用于保存viper读取到的
var Conf struct {
	ServerPort int
}

func init() {
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
	Conf.ServerPort = viper.GetInt("server.port")
}

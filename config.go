package main

import (

	// "delete/utils"
	"fmt"
	"path"

	// "github.com/fsnotify/fsnotify"

	ulog "utils/ulog"

	"github.com/spf13/viper"
)

const Version = "21.390v"

// var filename string

func viperinit(filenames string) {
	fileType := path.Ext(filenames)
	viper.SetConfigName(filenames)    //设置配置文件的名字
	viper.AddConfigPath(".")          //添加配置文件所在的路径
	viper.SetConfigType(fileType[1:]) //设置配置文件类型，可选
	err := viper.ReadInConfig()       // 查找并读取配置文件
	if err != nil {                   // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
func readConfig(key string) (data string) {
	err := viper.ReadInConfig()
	if err != nil {
		ulog.Logs.Panic(err)
	}
	da := viper.Get(key)
	if da == nil {
		ulog.Logs.Panic("config file hasnot this key:\n")
	}
	data = da.(string)
	ulog.Logs.Println("readConfig:", key, data)
	return data
}
func writeConfig(key string, value interface{}, newfile string) {
	err := viper.ReadInConfig()
	if err != nil {
		ulog.Logs.Panic(err)
	}
	viper.Set(key, value)
	// now := time.Now()
	// dateString := fmt.Sprintf("%d-%d-%d %d:%d:%d",
	// 	now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	// viper.Set("date", dateString)
	err = viper.WriteConfigAs(newfile)
	if err != nil {
		ulog.Logs.Panic(err)
	}
}

// func viperwatch() {
// 	// 设置监听回调函数(非必须)
// 	viper.OnConfigChange(func(e fsnotify.Event) {
// 		fmt.Println("Config file changed:", e)
// 	})
// 	//开始监听
// 	viper.WatchConfig()
// }

package ulog

import (
	"log"
	"os"
)

var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题
)

func Logpackageinit() {
	logFile, err := os.OpenFile("./log/package.txt",
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	log.SetOutput(logFile)
	// 输出前缀
	log.SetPrefix("[log]")
	// log格式
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	// // log输出到文件
	// log.Println([]string{"你好", "golang日志 - log"})
}

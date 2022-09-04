package ulog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Logs *logrus.Logger

func Loginit(Filenames string, LinkName string, MaxAge int, RotationTime int) {
	logClient := logrus.New()
	Logs = logClient
	//禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	logClient.Out = src

	//////////////////////////////////////////////
	apiLogPath := Filenames + "info"
	logWriterinfo, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(LinkName),                                  // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(MaxAge)*time.Hour),             // 文件最大保存时间 7*24*time.Hour
		rotatelogs.WithRotationTime(time.Duration(RotationTime)*time.Hour), // 日志切割时间间隔
	)
	//////////////////////////////////////////////
	apiLogPath = Filenames + "debug"
	logWriterDebug, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(LinkName),                                  // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(MaxAge)*time.Hour),             // 文件最大保存时间 7*24*time.Hour
		rotatelogs.WithRotationTime(time.Duration(RotationTime)*time.Hour), // 日志切割时间间隔
	)
	//////////////////////////////////////////////
	apiLogPath = Filenames + "warn"
	logWriterwarn, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(LinkName),                                  // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(MaxAge)*time.Hour),             // 文件最大保存时间 7*24*time.Hour
		rotatelogs.WithRotationTime(time.Duration(RotationTime)*time.Hour), // 日志切割时间间隔
	)
	//////////////////////////////////////////////
	var n = strings.LastIndex(Filenames, "/")
	var Filenametrace string
	if n > 0 {
		Filenametrace = Filenames[:n] + "/trace" + Filenames[n:]
	}
	apiLogPath = Filenametrace + "trace"
	logWritertrace, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(LinkName),                      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(MaxAge)*time.Hour), // 文件最大保存时间 7*24*time.Hour
		rotatelogs.WithRotationSize(1024*1024*10),              // 日志切割时间间隔 10M
	)
	//////////////////////////////////////////////
	apiLogPath = Filenames
	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(LinkName),                                  // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(MaxAge)*time.Hour),             // 文件最大保存时间 7*24*time.Hour
		rotatelogs.WithRotationTime(time.Duration(RotationTime)*time.Hour), // 日志切割时间间隔
	)

	logClient.SetLevel(logrus.TraceLevel)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriterinfo,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriterDebug,
		logrus.WarnLevel:  logWriterwarn,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriterwarn,
		logrus.TraceLevel: logWritertrace,
	}
	//第二个参数默认不填 JSONFormatter&TextFormatter
	//	&logrus.TextFormatter{
	//		TimestampFormat:"2006-01-02 15:04:05",
	//	}
	lfHook := lfshook.NewHook(writeMap, new(MyFormatter))
	logClient.AddHook(lfHook)
	// logoutput()
}

//var diff_time int64
func stack() string {
	var buf [2 << 10]byte
	stackstrt := string(buf[:runtime.Stack(buf[:], false)])
	index := strings.Index(stackstrt, "main.main")
	if index < 0 {
		Logs.Warnln("stack() can not find main")
	}
	return "\n\t" + stackstrt
}
func Debugln(args ...interface{}) {
	Logs.Debugln(args...)
}
func Println(args ...interface{}) {
	Logs.Println(args...)
}
func Warnln(args ...interface{}) {
	Logs.Warnln(args, stack())
}
func Panicln(args ...interface{}) {
	Logs.Panic(args, stack())
}

///////////////////////////////////////////////////////////////////////

type MyFormatter struct{}

func (s *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

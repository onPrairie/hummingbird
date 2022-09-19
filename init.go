package main

import (
	"bufio"
	_ "embed"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	ulog "utils/ulog"
	utilsEx "utils/utilsEx"

	"github.com/vua/vfmt"
)

//go:embed title
var title string
var con DeleteConfig
var ticker *time.Ticker
var crtlC chan os.Signal

func initparm() {
	var dev = map[string]string{
		"appversion": version,
	}
	context := os.Expand(title, func(k string) string { return dev[k] })
	fmt.Println(context)
	// vfmt.Println("[更多请访问] @[https://onprairie.github.io::blue|underline]")
	title = ""
	data1 := readfile("Hconfig.xml")
	config_path(data1)
	var Maxage int
	if con.Jsinit.Log.Maxage != "" {
		agei, err := strconv.Atoi(con.Jsinit.Log.Maxage)
		if err != nil {
			out_warn(1, err)
			return
		}
		Maxage = agei
	} else {
		Maxage = -1
	}
	var rottime int
	if con.Jsinit.Log.RotationTime != "" {
		rottimes, err := strconv.Atoi(con.Jsinit.Log.RotationTime)
		if err != nil {
			out_warn(1, err)
			return
		}
		rottime = rottimes
	} else {
		rottime = 24
	}
	ulog.Loginit("log/log", "", Maxage, rottime)
	ulog.Println("start----------------------------------------->", time.Now())

	//初始化网络
	connect_network()
	//初始化
	initjsenv()
	//初始化文件
	go Filewatch()

	//初始化结束，开始配置定时器
	//********若在init中写死循环则以下代码无法工作，包括定时器************
	var tinterval = con.Jsinit.Interval.Value
	tinterval = strings.TrimSpace(tinterval)
	if tinterval != "" {
		runticker()
		var timeformat = tinterval[len(tinterval)-1]
		var dura time.Duration

		timeval, err := strconv.Atoi(tinterval[:len(tinterval)-1])
		if err != nil {
			ulog.Println(tinterval[:len(tinterval)-1])
			panic("atio to number failed")
		}
		switch timeformat {
		case 'h':
			dura = time.Hour * time.Duration(timeval)
		case 'm':
			dura = time.Minute * time.Duration(timeval)
		case 's':
			dura = time.Second * time.Duration(timeval)
		}

		ticker = time.NewTicker(dura)
	} else {
		ticker = time.NewTicker(time.Minute)
		ticker.Stop()
	}

	//tickers := time.NewTicker(time.Second * 30)
	//接受ctrl +c 推出
	crtlC = make(chan os.Signal, 1)
	signal.Notify(crtlC, os.Interrupt, os.Kill)
	//是否开启udp监测
	if con.Jsinit.Udp.Bindaddress != "" {
		go Udprecvele()
	}
	if con.Jsinit.Http.Bindaddress != "" {
		go HttpRequest()
	}
	if con.Jsinit.Tcp.Bindaddress != "" {
		go tcpserver()
	}
	//是否开启Http监测

}
func config_path(data1 string) {
	err := xml.Unmarshal([]byte(data1), &con)
	if err != nil {
		out_warn(1, err)
		return
	}
}
func readfile(paths string) (data string) {
	file, err := ioutil.ReadFile(paths)
	if err != nil {
		ulog.Logs.Panic(err)
	}
	return string(file)
}
func connect_network() {
	if con.Jsinit.Tcp.Connect != "" {
		conn, err := net.Dial("tcp", con.Jsinit.Tcp.Connect)
		if err != nil {
			out_warn(1, err)
			return
		}
		utilsEx.TcpConn = conn
		//conn.Close()
	}
	if con.Jsinit.Udp.Connect != "" {
		addr, err := net.ResolveUDPAddr("udp", con.Jsinit.Udp.Connect)
		if err != nil {
			out_warn(1, err)
			return
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			out_warn(1, err)
			return
		}
		utilsEx.UdpConn = conn
	}

}
func initjsenv() {
	len := len(con.Jsparams.Params)
	utilsEx.Paramsmp = make(map[string]interface{}, len)
	for i := 0; i < len; i++ {
		_, ok := utilsEx.Paramsmp[con.Jsparams.Params[i].Id]
		if ok == true {
			out_warn(1, "the Param key is exit")
		}
		utilsEx.Paramsmp[con.Jsparams.Params[i].Id] = con.Jsparams.Params[i].Arg
		// con.Jsparams.Params[i].Id = ""
		// con.Jsparams.Params[i].Arg = con.Jsparams.Params[i].Arg[:0]
	}
	con.Jsparams = nil

	utilsEx.JsInit(true)
	if con.Jsinit.DB.Conname != "" {
		utilsEx.OpenMysql(con.Jsinit.DB.Conname)
		utilsEx.RegisterJsParser("domysql")
		utilsEx.RegisterJsParser("domysqlselect")
	}
	//文件及目录操作
	utilsEx.RegisterJsParser("__filemove")
	utilsEx.RegisterJsParser("findfiles")
	utilsEx.RegisterJsParser("__copyfile")
	utilsEx.RegisterJsParser("__writefile")
	utilsEx.RegisterJsParser("__readfile")
	utilsEx.RegisterJsParser("__filestate")
	utilsEx.RegisterJsParser("__dirstate")
	utilsEx.RegisterJsParser("__mkdir")
	utilsEx.RegisterJsParser("__filerename")
	utilsEx.RegisterJsParser("__fileremove")

	utilsEx.RegisterJsParser("log")
	utilsEx.RegisterJsParser("__Getjsparamsbyid")
	utilsEx.RegisterJsParser("__RemoveBeforeHour")
	//网络相关 客户端
	utilsEx.RegisterJsParser("__HttpSend")
	utilsEx.RegisterJsParser("__tcpsend")
	utilsEx.RegisterJsParser("__tcpread")
	utilsEx.RegisterJsParser("__udpsend")
	utilsEx.RegisterJsParser("__udpread")
	//网络相关 服务端
	inbuilt_func()

	//程序资源监控
	utilsEx.RegisterJsParser("Getmemory")
	utilsEx.RegisterJsParser("Restartservice")
	utilsEx.RegisterJsParser("Restartprocess")
	//睡眠
	utilsEx.RegisterJsParser("sleep")

	//引入js文件
	filecontext := getbigcontext()
	utilsEx.JsRun(filecontext)
	utilsEx.JsParserbuffer("init")
}
func inbuilt_func() {
	utilsEx.RegisterJsParser("__Httpwrite")
	utilsEx.RegisterJsParser("__Udpwrite")
	utilsEx.RegisterJsParser("__Tcpwrite")
}

//js内部对象转化
//__objectConversion() 参数1,type :用以设置不同对象的函数：
//************以下服务端设置***************************
//type值为1：设置为http的write
//值为2：设置为udp的write
//值为3：设置为tcp的write
func objConversion() string {
	var filecontext string
	filecontext = `
	function __objectConversion(){
		if(arguments[0] == 1){
			var res = {}
			res = arguments[1][1]
			req =  arguments[1][0]
			eval(res.Funcwrite + "res.write=write")
			httprecv(req,res)
		}else if(arguments[0] == 2){
			var res = {}
			res = arguments[1][1]
			req =  arguments[1][0]
			eval(res.Funcwrite + "res.write=write")
			udprecv(req,res)
		}else if(arguments[0] == 3){
			var res = {}
			res = arguments[1][1]
			req =  arguments[1][0]
			eval(res.Funcwrite + "res.write=write")
			tcprecv(req,res)
		}
	}
	`
	return filecontext
}
func RegisterObjet() string {
	var filecontext string
	if con.Jsinit.DB.Conname != "" {
		filecontext += `
		function __Mysql(){
			this.exec = domysql
			this.select = domysqlselect
		}
		var Mysql = new __Mysql()
		`
	}
	if con.Jsinit.Tcp.Connect != "" {
		filecontext += `	
		function __Tcp(){
			this.write = __tcpsend_chain
			this.read = __tcpread
		}
		var tcp = new __Tcp()`
	}
	if con.Jsinit.Udp.Connect != "" {
		filecontext += `	
		function __Udp(){
			this.write = __udpsend_chain
			this.read = __udpread
		}
		var udp = new __Udp()`
	}
	return filecontext
}

//js链式转化
func chainConversion() string {
	var filecontext string
	filecontext = `
	function __udpsend_chain(){
		__udpsend(arguments[0])
		return this
	}
	function __tcpsend_chain(){
		__tcpsend(arguments[0])
		return this
	}
	`
	return filecontext
}

//js函数包装
// argnum不仅表示参数个数
func packkagefunction(funcname string, argnum int) string {
	var filecontext string
	if argnum == 1 {
		filecontext = `
		function ${0}(a){
			var t =  __${0}(a)
			if(t == undefined){
				return t
			}
			if(t.hasOwnProperty("Err") == true){
				throw JSON.stringify(t)
			}
			return t		
		}
		`
	} else if argnum == 2 {
		filecontext = `
		function ${0}(a,b){
			var t =  __${0}(a,b)
			if(t == undefined){
				return t
			}
			if(t.hasOwnProperty("Err") == true){
				throw JSON.stringify(t)
			}
			return t		
		}
		`
	} else if argnum == 3 {
		filecontext = `
		function ${0}(a,b,c){
			var t =  __${0}(a,b,c)
			if(t == undefined){
				return t
			}
			if(t.hasOwnProperty("Err") == true){
				throw JSON.stringify(t)
			}
			return t		
		}
		`
	} else {
		panic("err argnum")
	}
	var dev = map[string]string{
		"0": funcname,
	}
	context := os.Expand(filecontext, func(k string) string { return dev[k] })
	return context
}

//内部加载优先制，如果同时设置外部与内部，则采用内部加载，外部将无效
func getbigcontext() string {
	var filecontext string
	//内部加载
	if con.Jscode.Loadfromfile == "" {
		for i := 0; i < len(con.Jscode.Script); i++ {
			if con.Jscode.Script[i].Src != "" {
				index := strings.Index(con.Jscode.Script[i].Src, "http")
				if index == 0 {
					s, err := utilsEx.HttpSendfortext("get", con.Jscode.Script[i].Src, nil, nil)
					if err != nil {
						ulog.Warnln(err)
						return ""
					}
					filecontext += string(s)
				} else {
					file, err := ioutil.ReadFile(con.Jscode.Script[i].Src)
					if err != nil {
						ulog.Panicln(err)
					}
					filecontext += string(file)
				}

				con.Jscode.Script[i].Src = ""
			} else if strings.TrimSpace(con.Jscode.Script[i].Value) != "" {
				filecontext += con.Jscode.Script[i].Value
				ulog.Debugln("load in xml")
				con.Jscode.Script[i].Value = ""
			}
		}
	} else {
		//外部加载
		for i := 0; i < len(con.Jscode.Script); i++ {
			index := strings.Index(con.Jscode.Script[i].Src, "http")
			if index == 0 {
				s, err := utilsEx.HttpSendfortext("get", con.Jscode.Script[i].Src, nil, nil)
				if err != nil {
					ulog.Warnln(err)
					return ""
				}
				filecontext += string(s)
			} else {
				file, err := ioutil.ReadFile(con.Jscode.Script[i].Src)
				if err != nil {
					ulog.Panicln(err)
				}
				filecontext += string(file)
			}

			con.Jscode.Script[i].Src = ""
		}
		file, err := ioutil.ReadFile(con.Jscode.Loadfromfile)
		if err != nil {
			ulog.Panicln(err)
		}
		filecontext += string(file)
		watcherFile = con.Jscode.Loadfromfile
	}
	filecontext += packkagefunction("readfile", 1)
	filecontext += packkagefunction("filestate", 1)
	filecontext += packkagefunction("copyfile", 1)
	filecontext += packkagefunction("dirstate", 1)
	filecontext += packkagefunction("writefile", 2)
	filecontext += packkagefunction("filemove", 2)
	filecontext += packkagefunction("filerename", 2)
	filecontext += packkagefunction("fileremove", 1)
	filecontext += packkagefunction("mkdir", 1)
	filecontext += packkagefunction("Getjsparamsbyid", 1)
	filecontext += packkagefunction("RemoveBeforeHour", 2)
	filecontext += packkagefunction("HttpSend", 3)
	filecontext += objConversion()
	filecontext += RegisterObjet()
	filecontext += chainConversion()
	// ioutil.WriteFile("text.js", []byte(filecontext), 0777)
	return filecontext
}
func runticker() {
	ulog.Logs.Debugln("the ticker ------------>", time.Now())
	utilsEx.JsParserbuffer("Interval")
}

//参数一: 如果为1则忽略日志输出
//参数二: 报错信息
// 以后的参数：必为错误信息，信息信息往日志里记
func out_warn(args ...interface{}) {
	output := args[1]
	style := "#00ff00|bg#ff0000|bold"
	if args[0] == 1 {
		vfmt.Printf("@[WARN: %s::%s]\n", output, style)
	} else {
		ulog.Warnln(args...)
		vfmt.Printf("@[WARN: %s::%s]\n", output, style)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("按回车结束")
	_, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

}

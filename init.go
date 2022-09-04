package main

import (
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
	title = ""
	data1 := readfile("Hconfig.xml")
	config_path(data1)
	var Maxage int
	if con.Jsinit.Log.Maxage != "" {
		agei, err := strconv.Atoi(con.Jsinit.Log.Maxage)
		if err != nil {
			ulog.Warnln(err)
			return
		}
		Maxage = agei
	} else {
		Maxage = 24 * 30
	}
	var rottime int
	if con.Jsinit.Log.RotationTime != "" {
		rottimes, err := strconv.Atoi(con.Jsinit.Log.Maxage)
		if err != nil {
			ulog.Warnln(err)
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
		ulog.Logs.Panic(err)
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
			ulog.Warnln(err)
			return
		}
		utilsEx.TcpConn = conn
		//conn.Close()
	}
	if con.Jsinit.Udp.Connect != "" {
		addr, err := net.ResolveUDPAddr("udp", con.Jsinit.Udp.Connect)
		if err != nil {
			ulog.Warnln(err)
			return
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			ulog.Warnln(err)
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
			panic("mp is exit")
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
	utilsEx.RegisterJsParser("filemove")
	utilsEx.RegisterJsParser("findfiles")
	utilsEx.RegisterJsParser("copyfile")
	utilsEx.RegisterJsParser("writefile")
	utilsEx.RegisterJsParser("readfile")
	utilsEx.RegisterJsParser("filestate")
	utilsEx.RegisterJsParser("dirstate")
	utilsEx.RegisterJsParser("mkdir")

	utilsEx.RegisterJsParser("log")
	utilsEx.RegisterJsParser("Getjsparamsbyid")
	utilsEx.RegisterJsParser("RemoveBeforeHour")
	//网络相关
	utilsEx.RegisterJsParser("HttpSend")
	utilsEx.RegisterJsParser("__tcpsend")
	utilsEx.RegisterJsParser("__tcpread")
	utilsEx.RegisterJsParser("__udpsend")
	utilsEx.RegisterJsParser("__udpread")
	//服务端
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

//内部加载优先制，如果同时设置外部与内部，则采用内部加载，外部将无效
func getbigcontext() string {
	var filecontext string
	//内部加载
	if con.Jscode.Loadfromfile == "" {
		for i := 0; i < len(con.Jscode.Script); i++ {
			if con.Jscode.Script[i].Src != "" {
				file, err := ioutil.ReadFile(con.Jscode.Script[i].Src)
				if err != nil {
					ulog.Panicln(err)
				}
				filecontext += string(file)
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
			if con.Jscode.Script[i].Src != "" {
				file, err := ioutil.ReadFile(con.Jscode.Script[i].Src)
				if err != nil {
					ulog.Panicln(err)
				}
				filecontext += string(file)
				con.Jscode.Script[i].Src = ""
			}
		}
		file, err := ioutil.ReadFile(con.Jscode.Loadfromfile)
		if err != nil {
			ulog.Panicln(err)
		}
		filecontext += string(file)
		watcherFile = con.Jscode.Loadfromfile
	}

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

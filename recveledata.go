package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
	ulog "utils/ulog"
	utilsEx "utils/utilsEx"
)

type UdpRecv struct {
	Port    int
	Address string
	Data    string
}
type UdpConnect struct {
	Funcwrite string
}

func Udprecvele() {
	udp_addr, err := net.ResolveUDPAddr("udp", con.Jsinit.Udp.Bindaddress)
	if err != nil {
		ulog.Logs.Panic("err udp:", err)
	}
	listen, err := net.ListenUDP("udp", udp_addr)
	if err != nil {
		ulog.Logs.Panic("err udp:", err)
		return
	}

	// listen.SetReadBuffer(1000)
	defer listen.Close()
	for {
		data := make([]byte, 1024)
		n, addr, err := listen.ReadFromUDP(data[:])
		if err != nil {
			ulog.Logs.Panic("err udp:", err)
		}
		ulog.Logs.Debugln("recv data:", string(data[:n]))

		var req UdpRecv
		req.Port = udp_addr.Port
		req.Address = udp_addr.IP.String()
		req.Data = string(data[:n])

		var udbcon UdpConnect
		udbcon.Funcwrite = `
		function write(data){
			__Udpwrite(data)
		}
		`
		if runJsobjectionConversion(2, req, udbcon) == true {
			data := <-utilsEx.UdpData
			_, err = listen.WriteToUDP([]byte(data), addr)
			if err != nil {
				ulog.Logs.Panic("err udp:", err)
			}
		}

	}
}

type Req struct {
	Path       string
	Method     string
	RemoteAddr string
	Querys     map[string][]string
	Body       string
}
type Res struct {
	Funcwrite string
}

func (tgis *Req) writes(str string) {
	fmt.Println(str)
}

func HttpRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var req Req
		req.Path = r.URL.Path
		req.Method = r.Method
		req.RemoteAddr = r.RemoteAddr
		req.Querys = r.URL.Query()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ulog.Warnln("read request.Body failed", err)
			return
		}
		req.Body = string(b)

		var res Res
		res.Funcwrite = `
		function write(data){
			__Httpwrite(data)
		}
		`

		if runJsobjectionConversion(1, req, res) == true {
			data := <-utilsEx.HttpData
			w.Write([]byte(data))
		}

	})
	http.ListenAndServe(con.Jsinit.Http.Bindaddress, nil)
}

type TcpRecv struct {
	Address string
	Data    string
}
type TcpConnect struct {
	Funcwrite string
}

func tcpserver() {
	// 创建 listener 192.168.1.170
	var listener net.Listener
	var err error
	for true {
		time.Sleep(time.Second)
		listener, err = net.Listen("tcp", con.Jsinit.Tcp.Bindaddress)
		if err != nil {
			ulog.Warnln("Error listening", err)
			//终止程序
		} else {
			break
		}
	}

	// 监听并接受来自客户端的连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			ulog.Warnln("Error accepting", err)
			return // 终止程序
		}
		go doServerStuff(conn)
	}
}

func doServerStuff(conn net.Conn) {
	var tcprecv TcpRecv
	for {
		buf := make([]byte, 1024)
		// err := conn.SetReadDeadline(time.Now().Add(time.Second * 3)) // timeout
		// if err != nil {
		// 	log.Println("setReadDeadline failed:", err)
		// }
		len, err := conn.Read(buf)
		if err != nil {
			ulog.Println("connect is close", err.Error())
			return //终止程序
		}

		str := string(buf[:len])
		ulog.Println("tcpserver Received data", str)

		tcprecv.Address = conn.LocalAddr().String()

		tcprecv.Data = str

		var tcpcon TcpConnect
		tcpcon.Funcwrite = `
		function write(data){
			__Tcpwrite(data)
		}
		`
		if runJsobjectionConversion(3, tcprecv, tcpcon) == true {
			data := <-utilsEx.TcpData
			conn.Write([]byte(data))
		}

	}
}
func runJsobjectionConversion(types int, args ...interface{}) bool {
	//未执行转换函数
	var typestr = "__Objwrite_" + strconv.Itoa(types)
	utilsEx.JsSet(typestr, false)
	utilsEx.JsParserbuffer("__objectConversion", types, args)

	iswrite, err := utilsEx.JsGet(typestr).ToBoolean()
	if err != nil {
		ulog.Warnln("read request.Body failed", err)
		return false
	}
	return iswrite
}

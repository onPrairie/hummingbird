package utilsEx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
	"utils/ulog"

	"github.com/gookit/color"
	"github.com/robertkrimen/otto"
)

var vm *otto.Otto

var Paramsmp map[string]interface{}

var jscode string
var isBuffer bool

var HttpData chan string
var UdpData chan string
var TcpData chan string

type Cfileinfo struct {
	Name    string
	IsDir   bool
	ModTime time.Time
	Size    int64
}
type JsErr struct {
	Err      string
	Funcname string
}

//网络连接
var TcpConn net.Conn
var UdpConn net.Conn

//是否开启缓存
func JsInit(isbuf bool) {
	vm = otto.New()
	isBuffer = true
	HttpData = make(chan string, 10)
	UdpData = make(chan string, 10)
	TcpData = make(chan string, 10)
}
func RegisterJsParser(funcname string) {
	switch funcname {
	//****************************数据库mysql操作****************************
	//***************************************************************
	case "domysql":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			ulog.Debugln("sql:", call.Argument(0).String())
			ret, err := MysqlDb.Exec(call.Argument(0).String())
			if err != nil {
				ulog.Warnln(err)
				vals, _ := vm.ToValue(false)
				return vals
			}
			theID, err := ret.RowsAffected() // 新插入数据的id
			if err != nil {
				ulog.Warnln("get lastinsert ID failed, err:", err)
				return otto.Value{}
			}
			ulog.Debugln("RowsAffected: ", theID)
			vals, _ := vm.ToValue(true)
			return vals
		})
	case "domysqlselect":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			sqlStr := call.Argument(0).String()
			ulog.Debugln("sql:", call.Argument(0).String())

			rows, err := MysqlDb.Query(sqlStr)
			if err != nil {
				ulog.Warnln("query failed, err:%v", err)
				return otto.Value{}
			}
			//获取列名
			columns, _ := rows.Columns()

			//定义一个切片,长度是字段的个数,切片里面的元素类型是sql.RawBytes
			values := make([]sql.RawBytes, len(columns))
			//定义一个切片,元素类型是interface{} 接口
			scanArgs := make([]interface{}, len(values))
			for i := range values {
				//把sql.RawBytes类型的地址存进去了
				scanArgs[i] = &values[i]
			}
			//获取字段值
			var result []map[string]string
			for rows.Next() {
				res := make(map[string]string)
				rows.Scan(scanArgs...)
				for i, col := range values {
					res[columns[i]] = string(col)
				}
				result = append(result, res)
			}

			bytes, e := json.Marshal(result)
			if e != nil {
				ulog.Warnln("序列化失败", e)
				return otto.Value{}
			} else {
				jsonStr := string(bytes)
				ulog.Debugln("get result", jsonStr)
			}
			vals, _ := vm.ToValue(result)
			return vals
		})

	//****************************文件操作****************************
	//***************************************************************
	case "__filemove":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			oldpath := call.Argument(0).String()
			newpath := call.Argument(1).String()
			ulog.Debugln("filemove:", call.Argument(0).String(), call.Argument(1).String())
			err := Move(oldpath, newpath)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "filemove"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			return otto.Value{}
		})
	case "__copyfile":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			oldpath := call.Argument(0).String()
			newpath := call.Argument(1).String()
			ulog.Debugln("copyfile:", call.Argument(0).String(), call.Argument(1).String())
			err := Copy(oldpath, newpath)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "copyfile"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			return otto.Value{}
		})
	case "__writefile":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			filename := call.Argument(0).String()
			data := call.Argument(1).String()
			ulog.Debugln("copyfile:", call.Argument(0).String(), call.Argument(1).String())
			err := ioutil.WriteFile(filename, []byte(data), 0777)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "writefile"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			return otto.Value{}
		})
	case "findfiles":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			absolute_path := call.Argument(0).String()
			t := findfiles(absolute_path)
			result, _ := vm.ToValue(t)
			return result
		})
	case "__readfile":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			absolute_path := call.Argument(0).String()
			data, err := ioutil.ReadFile(absolute_path)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "readfile"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			result, _ := vm.ToValue(string(data))
			return result
		})
	case "__filestate":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			path := call.Argument(0).String()
			fileinfo, err := os.Stat(path)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "filestate"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			var info Cfileinfo
			info.Name = fileinfo.Name()
			info.IsDir = fileinfo.IsDir()
			info.ModTime = fileinfo.ModTime()
			//暂时无法统计目录大小
			if info.IsDir == true {
				info.Size = 0
			} else {
				info.Size = fileinfo.Size()
			}
			if err != nil {
				ulog.Warnln(err)
			}
			result, _ := vm.ToValue(info)
			return result
		})
	case "__dirstate":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			path := call.Argument(0).String()
			dirs, err := ioutil.ReadDir(path)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "filestate"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			var info []Cfileinfo
			info = make([]Cfileinfo, len(dirs))
			for i := 0; i < len(dirs); i++ {
				info[i].Name = dirs[i].Name()
				info[i].IsDir = dirs[i].IsDir()
				info[i].ModTime = dirs[i].ModTime()

				//暂时无法统计目录大小
				if info[i].IsDir == true {
					info[i].Size = 0
				} else {
					info[i].Size = dirs[i].Size()
				}
				if err != nil {
					var Serr JsErr
					Serr.Funcname = "filestate"
					Serr.Err = err.Error()
					result, _ := vm.ToValue(Serr)
					return result
				}
			}
			result, _ := vm.ToValue(info)
			return result
		})
	case "__mkdir":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			paths := call.Argument(0).String()
			err := os.MkdirAll(paths, 0777)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "mkdir"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			return otto.Value{}
		})
	case "__Getjsparamsbyid":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			key := call.Argument(0).String()
			ulog.Debugln("Getjsparamsbyid:", key)
			v, ok := Paramsmp[key]
			if ok == false {
				var Serr JsErr
				Serr.Funcname = "Getjsparamsbyid"
				Serr.Err = "mp is not exit"
				result, _ := vm.ToValue(Serr)
				return result
			}
			result, _ := vm.ToValue(v)
			return result
		})
	case "__RemoveBeforeHour":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			paths := call.Argument(0).String()
			hours := call.Argument(1).String()
			hour, err := strconv.Atoi(hours)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "Getjsparamsbyid"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			difftime := int64(hour) * 3600
			err = RemoveBefore(paths, difftime)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "Getjsparamsbyid"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			return otto.Value{}
		})
	case "__filerename":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			odlpath := call.Argument(0).String()
			newpath := call.Argument(1).String()
			err := os.Rename(odlpath, newpath)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "filerename"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			return otto.Value{}
		})
	case "__fileremove":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			path := call.Argument(0).String()
			err := os.RemoveAll(path)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "fileremove"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			return otto.Value{}
		})
	//****************************日志相关****************************
	//***************************************************************
	case "log":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			argls := call.ArgumentList
			ulog.Println("jscode out:", argls)
			return otto.Value{}
		})
	//****************************网络相关****************************
	//***************************************************************
	case "__HttpSend":
		//***  暂不支持header
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			types := call.Argument(0).String()
			Ledadddress := call.Argument(1).String()
			msg := call.Argument(2).String()
			if msg == "null" || msg == "undefined" {
				msg = ""
			}
			s, err := HttpSend(types, Ledadddress, []byte(msg), nil)
			if err != nil {
				var Serr JsErr
				Serr.Funcname = "fileremove"
				Serr.Err = err.Error()
				result, _ := vm.ToValue(Serr)
				return result
			}
			result, _ := vm.ToValue(s)
			return result
		})

	case "__Httpwrite":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			date := call.Argument(0).String()
			HttpData <- date
			JsSet("__Objwrite_1", true)
			return otto.Value{}
		})
	case "__Udpwrite":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			date := call.Argument(0).String()
			UdpData <- date
			JsSet("__Objwrite_2", true)
			return otto.Value{}
		})
	case "__Tcpwrite":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			date := call.Argument(0).String()
			TcpData <- date
			JsSet("__Objwrite_3", true)
			return otto.Value{}
		})
	case "__tcpsend":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			date := call.Argument(0).String()
			TcpConn.Write([]byte(date))
			return otto.Value{}
		})
	case "__tcpread":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			buf := make([]byte, 1024)
			n, err := TcpConn.Read(buf)
			if err != nil {
				fmt.Println(err)
			}
			s := string(buf[:n])
			result, _ := vm.ToValue(s)
			return result
		})
	case "__udpsend":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			date := call.Argument(0).String()
			UdpConn.Write([]byte(date))
			return otto.Value{}
		})
	case "__udpread":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			buf := make([]byte, 1024)
			n, err := UdpConn.Read(buf)
			if err != nil {
				fmt.Println(err)
			}
			s := string(buf[:n])
			result, _ := vm.ToValue(s)
			return result
		})
	//****************************系统相关****************************
	//***************************************************************
	case "Getmemory":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			programname := call.Argument(0).String()
			t := Getmemory(programname)
			ulog.Debugln("js build uilt-in", "programname:", t)
			result, _ := vm.ToValue(t)
			return result
		})
	case "Restartservice":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			servicename := call.Argument(0).String()
			Restartservice(servicename)
			return otto.Value{}
		})
	case "Restartprocess":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			servicename := call.Argument(0).String()
			Restartprocess(servicename)
			return otto.Value{}
		})
	case "sleep":
		vm.Set(funcname, func(call otto.FunctionCall) otto.Value {
			sec, err := call.Argument(0).ToInteger()
			if err != nil {
				ulog.Warnln(err)
				return otto.Value{}
			}
			time.Sleep(time.Duration(sec) * time.Millisecond)
			return otto.Value{}
		})
	default:
		panic("can not register JsParser")

	}

}
func JsRun(data string) {
	_, err := vm.Run(data)
	if err != nil {
		panic(err)
	}
}
func JsSet(key string, value interface{}) {
	vm.Set(key, value)
}
func JsGet(key string) otto.Value {
	value, err := vm.Get(key)
	if err != nil {
		panic(err)
	}
	return value
}

//缓存调用 必须调用JsRun
func JsParserbuffer(functionName string, args ...interface{}) (result string) {
	value, err := vm.Call(functionName, nil, args...)
	if err != nil {
		color.Red.Println("WARN:", err)
		ulog.Warnln(err)
		return ""
	}
	return value.String()
}
func JsParser(data string, functionName string, args ...interface{}) (result string) {
	if isBuffer == true {
		if jscode != data {
			_, err := vm.Run(data)
			if err != nil {
				panic(err)
			}
		}
		jscode = data
	} else {
		_, err := vm.Run(data)
		if err != nil {
			panic(err)
		}
	}

	value, err := vm.Call(functionName, nil, args...)
	if err != nil {
		ulog.Warnln(err)
		return ""
	}
	return value.String()
}

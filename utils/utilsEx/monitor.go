package utilsEx

import (
	"bytes"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	ulog "utils/ulog"

	"github.com/axgle/mahonia"
)

func Getmemory(programname string) (res int) {
	cmd := exec.Command("cmd")
	// cmd := exec.Command("powershell")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in //绑定输入
	var out bytes.Buffer
	cmd.Stdout = &out //绑定输出
	go func() {
		// start stop restart
		in.WriteString("wmic process where name='" + programname + "' get WorkingSetSize\n") //写入你的命令，可以有多行，"\n"表示回车
	}()
	err := cmd.Start()
	if err != nil {
		ulog.Warnln(err)
	}
	ulog.Logs.Println(cmd.Args)
	err = cmd.Wait()
	if err != nil {
		ulog.Logs.Printf("Command finished with error: %v", err)
	}
	//out.String() //
	rt := mahonia.NewDecoder("gbk").ConvertString(out.String()) //
	//log.Println(rt)
	arr := strings.Fields(rt)
	ulog.Logs.Println(arr)
	var ks int
	for i := len(arr) - 1; i > -1; i-- {
		if arr[i] == "WorkingSetSize" {
			ks, err = strconv.Atoi(arr[i+1])
			if err != nil {
				ulog.Logs.Warn(err)
				ks = -1
			}
			break
		}
	}
	var kskbyte int
	if ks == -1 {
		kskbyte = -1
	} else {
		kskbyte = ks / 1024
	}

	ulog.Logs.Println("get mem: bytes:", ks, "kb:", kskbyte)
	return kskbyte
}

func Restartservice(servicename string) {
	ulog.Logs.Println("start memmonitor!!!")
	var ostype = runtime.GOOS
	if ostype == "linux" {

	} else if ostype == "windows" {
		//停止服务 sc stop servicename
		//启动服务 sc start servicename
		Cmdexec("sc stop " + servicename + "\nsc start " + servicename)
	} else {
		ulog.Warnln("can supout this system!")
	}
}
func Restartprocess(servicename string) {
	ulog.Logs.Println("start memmonitor!!!")
	var ostype = runtime.GOOS
	if ostype == "linux" {

	} else if ostype == "windows" {
		dir, path := path.Split(servicename)
		dir = strings.Replace(dir, "/", "\\", -1)

		//杀死进程 TASKKILL /F /IM notepad++.exe /T
		//开启进程 start Notepad++.exe
		Cmdexec("TASKKILL /F /IM " + path + " /T\ncd " + dir + " \nstart " + path)
	} else {
		ulog.Warnln("can supout this system!")
	}
}

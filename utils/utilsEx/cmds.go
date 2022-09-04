package utilsEx

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"
	ulog "utils/ulog"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//50 130
// var limit = 100
// 	var bufferLen = 8000
func ModifySTringFromBuffer(sdata string) string {
	var newstr = ""
	var limit = 50
	var bufferLen = 1300 //缓冲区大小

	// var jstart int
	var jstartda int
	var j int = bufferLen
	var jlen = limit
	//times := len(sdata) / (bufferLen + limit) //余数一定比除数小
	for jlen < len(sdata) {
		for ; j < jlen; j++ {
			if sdata[j-1] == ']' && sdata[j] == ',' {
				newstr += sdata[jstartda : j+1]
				newstr += "\n"
				// j = len(newstr) - 1
				jstartda = j + 1
				j += bufferLen
				break
			}
			if j == jlen-1 {
				ulog.Logs.Warnln("ModifySTringFromBuffer", "can not find newstr", sdata[j-limit:])
			}
		}
		jlen = limit + j
	}

	newstr += sdata[jstartda:]
	return newstr
}
func DeleteExtraSpace(str string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	//  /* */：单行注释 /\*{1,2}[\s\S]*?\*/  [^\w]
	// //:多行注释 [\s\S]*?\n
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

//cmds 多条命令以\n间隔
func Cmdexec(cmds string) (res string) {
	var ostype = runtime.GOOS
	if ostype == "linux" {
		res = cmdexecforLinux(cmds)
	} else if ostype == "windows" {
		if cmds[len(cmds)-1] != '\n' {
			// ulog.Logs.Warnln("need enter!!:\\n")
			cmds += "\n"
		}
		res = cmdexecforWindows(cmds)
	} else {
		ulog.Warnln("can supout this system!")
	}
	return res
}
func cmdexecforWindows(cmds string) (res string) {
	cmd := exec.Command("cmd")
	// cmd := exec.Command("powershell")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in //绑定输入
	var out bytes.Buffer
	cmd.Stdout = &out //绑定输出
	cmd.Stderr = &out
	cmdsgbk, err := Utf8ToGbk([]byte(cmds))
	if err != nil {
		ulog.Logs.Warnln(err)
		return
	}
	go func() {
		// start stop restart
		_, err := in.WriteString(string(cmdsgbk)) //写入你的命令，可以有多行，"\n"表示回车
		if err != nil {
			ulog.Logs.Warnln("WriteString!!:\\n", err)
		}
	}()
	err = cmd.Start()
	if err != nil {
		ulog.Logs.Warnln(err)
		return
	}
	ulog.Logs.Println(cmd.Args)
	err = cmd.Wait()
	if err != nil {
		ulog.Logs.Warnf("Command finished with error: %v", err)
		return
	}
	resbyte, err1 := GbkToUtf8(out.Bytes())
	if err1 != nil {
		ulog.Logs.Warnln(err)
		return
	}
	res = string(resbyte)
	ulog.Logs.Debugln("cmds run:  ", res)
	return res
}

//注意：如果需要后台运行nohup，可参照此命令： nohup  ping www.baidu.com  >/dev/null 2>&1 &
func cmdexecforLinux(cmds string) (res string) {
	cmd := exec.Command("/bin/bash", "-c", cmds)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	ulog.Logs.Println(cmd.Args)
	err := cmd.Run()
	if err != nil {
		ulog.Warnln(err)
		return
	}
	res = out.String() //
	ulog.Logs.Println("cmds run:  ", res)
	return res
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

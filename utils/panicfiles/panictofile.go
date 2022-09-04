package panicfiles

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
)

var initforone = false
var stdErrFileHandler *os.File

const (
	kernel32dll = "kernel32.dll"
)

func getCurrentExeMd5Sum() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	filePath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	var md5Sum string
	fp, err := os.Open(filePath)
	if err != nil {
		return md5Sum, err
	}
	defer fp.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, fp); err != nil {
		return md5Sum, err
	}
	// hashInBytes := hash.Sum(nil)[:4] // only show 4 bytes
	hashInBytes := hash.Sum(nil)
	md5Sum = hex.EncodeToString(hashInBytes)
	return md5Sum, nil
}
func writejsondata(filename string) (con Config, t bool) {
	filename += Panicconfigname
	//var con Config
	con.Fileformat = ""
	summd5, err := getCurrentExeMd5Sum()
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "writejsondata"))
		return con, false
	}
	con.Program = summd5
	con.Version = version
	con.Panicfile = panicFile
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "getjsontodata"))
		return con, false
	}
	marshaldata, err := json.Marshal(&con)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "getjsontodata"))
		return con, false
	}
	var out bytes.Buffer
	err = json.Indent(&out, marshaldata, "", "\t")

	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(filename, out.Bytes(), 0777)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "getjsontodata"))
		return con, false
	}
	return con, false
}
func operatejsonRw(con Config) {
	var err error
	filename := panicFile[:12] + Panicconfigname
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "operatejsonRw"))

	}
	var conR Config
	err = json.Unmarshal(file, &conR)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "operatejsonRw"))
	}
	if con.Version != "" && con.Version != conR.Version {
		conR.Version = con.Version
	}
	if con.Program != "" {
		conR.Program = con.Program
	}

	marshaldata, err := json.Marshal(&conR)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "getjsontodata"))
		return
	}
	var out bytes.Buffer
	err = json.Indent(&out, marshaldata, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(filename, out.Bytes(), 0777)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "getjsontodata"))
	}
}
func getjsontodata(filename string) (con Config) {
	var err error
	file, err := ioutil.ReadFile(panicFile + Panicconfigname)
	if err != nil {
		perr, ok := (err).(*fs.PathError)
		if ok == false {
			log.Printf("err :%+v\n", errors.Wrap(errors.New("*fs.PathError failed"),
				"getjsontodata"))
			return
		}
		var strerr string
		if runtime.GOOS == "windows" {
			strerr = "The system cannot find the file specified."
		} else {
			strerr = "no such file or directory"
		}
		//no such file or directory
		if perr.Err.Error() == strerr {
			writejsondata(filename)
			return
		} else {
			log.Printf("err :%+v\n", errors.Wrap(err, "getjsontodata"))
		}
	}
	err = json.Unmarshal(file, &con)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "getjsontodata"))
	}
	return
}
func init0() {
	filet := panicFile[:12]
	filenameall := path.Dir(panicFile)
	os.Mkdir(filenameall, os.ModePerm)
	//////////////////////////////////////
	con := getjsontodata(panicFile)
	summd5, err := getCurrentExeMd5Sum()
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "custom message"))
		return
	}

	getfilenames()
	dir, err := ioutil.ReadDir(filet)
	if err != nil {
		log.Printf("err :%+v\n", errors.Wrap(err, "custom message"))
		return
	}
	if summd5 == con.Program {
		log.Println("Because it is the same executable, only the necessary files are cleaned up")
		for _, info := range dir {
			// if info.IsDir() {
			// 	log.Println("Ignore folder")
			// }
			if info.Size() == 0 {
				if runtime.GOOS == "windows" {
					os.RemoveAll(filet + "\\" + info.Name())
				} else {
					os.RemoveAll(filet + "/" + info.Name())
				}
			}
		}
	} else {
		for _, info := range dir {
			// if info.IsDir() {
			// 	log.Println("Ignore folder")
			// }
			if info.Name() != Panicconfigname {
				if runtime.GOOS == "windows" {
					os.RemoveAll(filet + "\\" + info.Name())
				} else {
					os.RemoveAll(filet + "/" + info.Name())
				}
			} else {
				var con Config
				con.Program = summd5
				operatejsonRw(con)
			}

		}
	}
	{
		var con Config
		con.Version = version
		operatejsonRw(con)
	}
}
func getfilenames() {
	//exeName := os.Args[0] //获取程序名称
	//
	//currentPath := filepath.ToSlash(exeName)
	//filenameall := path.Base(currentPath)
	now := time.Now()  //获取当前时间
	pid := os.Getpid() //获取进程ID

	time_str := now.Format("2006-01-02==15-04-05")        //设定时间格式
	fname := fmt.Sprintf("%d_%s_dump.log", pid, time_str) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
	panicFile += "Panic" + fname
	return
}
func TryE(callback func()) {
	if initforone == false {
		log.Println("you not call InitPanicFile or error!")
		return
	}
	errs := recover()
	if errs == nil {
		return
	}
	stdErrFileHandler.WriteString(fmt.Sprintf("%v\r\n========================\r\n%s\n",
		errs, string(debug.Stack())))
	if callback != nil {
		callback()
	}
}

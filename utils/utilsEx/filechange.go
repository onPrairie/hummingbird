package utilsEx

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	ulog "utils/ulog"
)

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func createFile(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

var base int = 0

//仅限于查询windows路径下的文件名，支持通配符
func findfiles(absolute_path string) (strls []string) {
	absolute_path = strings.Replace(absolute_path, "/", "\\", -1)
	fileanmes := Cmdexec("dir " + absolute_path + " /s /b")
	arr := strings.Fields(fileanmes)

	if base == 0 {
		for i := 0; i < len(arr); i++ {
			if arr[i] == "/b" {
				base = i
			}
		}
		base += 1
	}
	//除去最后无用数据
	strls = make([]string, len(arr)-base-1)
	for i := base; i < len(arr)-1; i++ {
		strls[i-base] = arr[i]
	}
	return strls
}
func Move(oldpath, newpath string) error {
	dir, _ := path.Split(newpath)
	err := createFile(dir)
	if err != nil {
		return err
	}
	var ostype = runtime.GOOS
	if ostype == "linux" {
		return os.Rename(oldpath, newpath)
	} else if ostype == "windows" {
		from, err := syscall.UTF16PtrFromString(oldpath)
		if err != nil {
			return err
		}
		to, err := syscall.UTF16PtrFromString(newpath)
		if err != nil {
			return err
		}
		return syscall.MoveFile(from, to) //windows API
	}
	// ulog.Warnln("can supout this system!")
	return errors.New("can supout this system!" + ostype)
}
func Copy(oldpath, newpath string) error {
	var ostype = runtime.GOOS
	if ostype == "windows" {
		if !isExist(oldpath) {
			return errors.New("can not find file!" + oldpath)
		}
		in, err := os.OpenFile(oldpath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		out, err := os.OpenFile(newpath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			ulog.Logs.Warnln("to ", err)
		}
		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		in.Close()
		out.Close()
		return nil
	} else {
		return errors.New("can supout this system!" + ostype)
	}

}

//paths 目录，此目录下文件全部清空
//diff_time 秒
func RemoveBefore(paths string, diff_time int64) error {
	now_time := time.Now().Unix() //当前时间，main使用Unix时间戳
	err := filepath.Walk(paths, func(paths string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		file_time := f.ModTime().Unix()
		if (now_time - file_time) > diff_time { //判断文件是否超过7天
			ulog.Logs.Printf("Delete file %v ", paths)
			dirs := isFile(paths)
			if dirs == true {
				os.RemoveAll(paths)
			}

		}
		return nil
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}
func isFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

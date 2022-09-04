// +build windows

package panicfiles

import (
	"log"
	"os"
	"runtime"
	"syscall"
)




func InitPanicFile() error {
	init0();
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	stdErrFileHandler = file //把文件句柄保存到全局变量，避免被GC回收
	kernel32 := syscall.NewLazyDLL(kernel32dll)
	setStdHandle := kernel32.NewProc("SetStdHandle")
	sh := syscall.STD_ERROR_HANDLE
	v, _, err := setStdHandle.Call(uintptr(sh), uintptr(file.Fd()))
	if v == 0 {
		return err
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})
	initforone = true;

	return nil
}

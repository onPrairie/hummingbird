package panicfiles

import (
	"log"
	"os"
	"runtime"
	"syscall"
)

func InitPanicFile() error {
	init0()
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	stdErrFileHandler = file
	if err != nil {
		println(err)
		return err
	}
	if runtime.GOARCH == "arm64" {
		// LINUX系统
		err = syscall.Dup3(int(file.Fd()), int(os.Stderr.Fd()), 0)
	} else {
		err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	}
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})
	initforone = true
	log.SetOutput(os.Stdout)
	return nil
}

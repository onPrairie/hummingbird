package main

import (
	"fmt"
	"log"
	"time"
	"utils/ulog"

	"github.com/fsnotify/fsnotify"
)

var JScodeDATA chan struct{} = make(chan struct{}, 1) //js代码内存，用于加速读取js代码
var watcherFile string
var timeonce bool

func Filewatch() {
	var path string
	if con.Jscode.Loadfromfile == "" {
		path = "Hconfig.xml"
	} else {
		path = con.Jscode.Loadfromfile
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if timeonce == false {
						time.AfterFunc(2*time.Second, func() {
							JScodeDATA <- struct{}{}
							timeonce = false
						})
						timeonce = true
					}
					fmt.Println("writing", time.Now())
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("Remove file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		ulog.Println(err)
	}
	<-done
}

package main

import (
	"fmt"
	ulog "utils/ulog"

	_ "utils/entry"

	utilsEx "utils/utilsEx"
)

const version = "202206.2884"

func main() {
	initparm()

	if ticker != nil {
		for {
			select {
			case <-ticker.C:
				runticker()
			case <-JScodeDATA:
				LoadJsfile()
			case s := <-crtlC:
				ulog.Logs.Println("=========>over", s)
				close(crtlC)
				return
			}
		}
	} else {
		for {
			select {
			case <-JScodeDATA:
				LoadJsfile()
			case s := <-crtlC:
				ulog.Logs.Println("=========>over", s)
				close(crtlC)
				return
			}
		}
	}

}
func LoadJsfile() {
	if con.Jscode.Loadfromfile == "" {
		data1 := readfile("Hconfig.xml")
		if data1 == "" {
			return
		}
		config_path(data1)
	}
	filecontext := getbigcontext()
	utilsEx.JsRun(filecontext)
	fmt.Println("JS run again")
}

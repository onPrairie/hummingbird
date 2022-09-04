package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	_ "utils/entry"
	// _ "github.com/go-sql-driver/mysql"
)

// func init() {
// 	utils.Loginit("log/log", "", 24*30, 24)
// 	utils.InitPanicFile()
// }
func Test_main(t *testing.T) {
	fmt.Println("version", version)

	//测试js 运行
	// utilsEx.Move("D:/ftp1/20201030162626760-1-黄色-鲁BY0751-大型货-非-D.jpg",
	// 	"D:/images/2020/10/20201030162626760-1-黄色-鲁BY0751-大型货-非-D.jpg")
	// // fmt.Println(err)

	fileInfo, err := os.Stat("D:/ftp")
	//filePath包含文件名
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fileInfo.Name())    //xxx.txt
	fmt.Println(fileInfo.IsDir())   //false  判断是否是目录
	fmt.Println(fileInfo.ModTime()) //2019-12-05 16:59:36.8832788 +0800 CST   文件的修改时间
	fmt.Println(fileInfo.Size())    //3097  文件大小（字节）
	fmt.Println(fileInfo.Mode())    // -rw-rw-rw-  读写属性
	fmt.Println(fileInfo.Sys())     //&{32 {2160608986 30778972} {2160608986 30778972} {1375605524 30780234} 0 3097}

	fi, err2 := ioutil.ReadDir("D:/ftp1")
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(fi[0].Name(), len(fi))
}

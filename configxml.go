package main

import "encoding/xml"

//js初始化相关
type CJsinit struct {
	// Value        string `xml:",cdata"`
	Interval CInterval `xml:"Interval"`
	DB       Database  `xml:"Database"`
	Udp      CUdp      `xml:"udp"`
	Http     CHttp     `xml:"http"`
	Tcp      CTcp      `xml:"tcp"`
	Log      CLog      `xml:"log"`
}

//js引用相关
type CJscode struct {
	Script       []CScript `xml:"script"`
	Loadfromfile string    `xml:"loadfromfile,attr"`
}
type Database struct {
	Conname string `xml:"conname"`
}
type CParams struct {
	Id  string   `xml:"id,attr"`
	Arg []string `xml:"arg"`
}
type CJsparams struct {
	Params []CParams `xml:"params"`
}

//定时器相关
type CInterval struct {
	Value string `xml:",cdata"`
}

//udp相关
type CUdp struct {
	Bindaddress string `xml:"bindaddress,attr"`
	Connect     string `xml:"connect,attr"`
}

//udp相关
type CHttp struct {
	Bindaddress string `xml:"bindaddress,attr"`
}
type CTcp struct {
	Bindaddress string `xml:"bindaddress,attr"`
	Connect     string `xml:"connect,attr"`
}

//log相关
type CLog struct {
	Maxage       string `xml:"maxage,attr"`
	RotationTime string `xml:"rotationTime,attr"`
}
type CScript struct {
	Src   string `xml:"src,attr"`
	Value string `xml:",cdata"`
}

//root
type DeleteConfig struct {
	Hummingbird xml.Name   `xml:"Hummingbird"`
	Version     string     `xml:"vesion,attr"`
	Jsparams    *CJsparams `xml:"jsparams"`
	Jsinit      CJsinit    `xml:"jsinit"`
	Jscode      CJscode    `xml:"jscode"`
}

// var Paramsmp map[string]interface{}

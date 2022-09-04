function init() {
	console.log(readfile("D:/1/1.txt")); 
}

//定时器触发
function Interval(){
	console.log("Intervals:",new Date())
	console.log(Getmemory("notepad++.exe"))
	//HttpSend("GET", "http://", null, null)
	// var res = domysqlselect("SELECT * FROM axisstandard");
	// console.log(res[0],res.length,res[10].AxisName,res[10].AxisType);
}

function httprecv(req,res){
	console.log(JSON.stringify(req),"-------"); 
	// console.log("222",obj.Path);
	res.write(req.Path)
}
//在获取udprecv回调接口时触发
function udprecv(req,res) {
	console.log("udp",JSON.stringify(req)); 
	res.write(req.Data)
};
function tcprecv(req,res){
	console.log("tcp",JSON.stringify(req)); 
	res.write(req.Data)
}

function filemovepics(jsonObj,rootpath){
	  var ftppath = "E:/ftp/" + jsonObj.vehicleD
		var fullpathD;
		var fullpathQ; 
		var fullpathS; 
		var fullpathC;
		var fullpathW;
		fullpathD = "E:/pics/" + rootpath + "/" + jsonObj.vehicleD.substr(0,4) + 
		"/" + jsonObj.vehicleD.substr(4,4) +"/"  + jsonObj.vehicleD
		 filemove(ftppath,fullpathD)

		if(jsonObj.hasOwnProperty('vehicleQ')){
			fullpathQ = "E:/pics/" + rootpath + "/" + jsonObj.vehicleQ.substr(0,4) + 
			"/" + jsonObj.vehicleQ.substr(4,4) +"/"  + jsonObj.vehicleQ
			ftppath = "E:/ftp/" + jsonObj.vehicleQ
			filemove(ftppath,fullpathQ)
		}

		if(jsonObj.hasOwnProperty('vehicleS')){
			fullpathS = "E:/pics/" + rootpath + "/" + jsonObj.vehicleS.substr(0,4) + 
			"/" + jsonObj.vehicleS.substr(4,4) +"/"  + jsonObj.vehicleS
			ftppath = "E:/ftp/" + jsonObj.vehicleS
			filemove(ftppath,fullpathS)
		}
		if(jsonObj.hasOwnProperty('vehicleC')){
			console.log( "1234",jsonObj.vehicleC)
			fullpathC = "E:/pics/" + rootpath + "/" + jsonObj.vehicleC.substr(0,4) + 
			"/" + jsonObj.vehicleC.substr(4,4) +"/"  + jsonObj.vehicleC
			ftppath = "E:/ftp/" + jsonObj.vehicleC
			filemove(ftppath,fullpathC)
		}
		 if(jsonObj.hasOwnProperty('vehicleW')){
			fullpathW = "E:/pics/" + rootpath + "/" + jsonObj.vehicleW.substr(0,4) + 
			"/" + jsonObj.vehicleW.substr(4,4) +"/"  + jsonObj.vehicleW
			ftppath = "E:/ftp/" + jsonObj.vehicleW
			filemove(ftppath,fullpathW)
		}
		updatepicadds(rootpath,fullpathD,jsonObj.id,fullpathQ,fullpathS,fullpathC,fullpathW)
}
function updatepicadds(rootpath,fullpathD,id,fullpathQ,fullpathS,fullpathC,fullpathW) {
	var nowDate = new Date()
	var tablename = "d_entry"+ formatDate(nowDate, "YYMM")

	var sql = "UPDATE "+tablename +" SET d_stationcode = '" + rootpath
		+ "',d_dpic='" +fullpathD + "',d_qpic='" + fullpathQ + "',d_wcpic='" + fullpathS +
		"',d_cpic='" + fullpathC + "',d_wpic='" + fullpathW +
		"'  WHERE d_jiancdh = '"+ id +"'"
	domysql(sql)
}

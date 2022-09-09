# Hummingbird
---

 - 介绍：

Hummingbird 是一套遵循 MIT协议的开源 JavaScript 运行引擎。设计的初衷是为了敏捷开发，尽量做到省时省力。

 - Hummingbird 可以做什么？

1. 搭建你所需的网络测试平台，集成tcp,http,udp客户端与服务端
2. 监控本地资源，分析程序内存大小，可对其重启操作。
3. 内置mysql驱动模块，通过内置api可以轻松调用。

 - 为什么用Hummingbird ？
 
  1. 一旦保存编辑的代码，即可快速编译运行，无需手动编译。
  2. 无需记住过多API,采用xml配置方式快速完成初始化。
  3. 避免回调地狱，改造Node.js的原有api,将异步函数改为同步。


Hummingbird只遵循一种原则：**简单易用**。

## 快速开始
**目标：** 向控制台输出hello world
1. 在Hummingbird程序路径下创建Hconfig.xml,内容如下：

```xml
<?xml version="1.0" encoding="utf-8"?>
<Hummingbird vesion="13">
    <jsinit>
    </jsinit>
    
    <jscode loadfromfile="">   
        <script>
            function init() {
                console.log("hello world")
            }
        </script>
    </jscode>

    <jsparams>
    </jsparams>
</Hummingbird>
```
2. 启动Hummingbird，在控制台你会看到打印字符串
![图1](https://img-blog.csdnimg.cn/91b6fcd828664415a4d45d02f6aa2d56.png)
	
## Hconfig.xml说明
###  Hconfig.xml的基础框架
	
在上面的标题下，贴的Hummingbird依赖的基础xml文件。下面介绍Hconfig.xm里面的内容：

第一行为程序xml版本号，编码默认utf-8
	
```bash
	<?xml version="1.0" encoding="utf-8"?> 
	
```
	
标签Hummingbird为根节点，vesion为版本号。

```bash
	<Hummingbird vesion="13"></Hummingbird>
```
	
Hummingbird标签有三个标签：分别为jsparams，jscode，jsinit。这三个标签建议不要缺少，因为他们是保证Hummingbird运行基础标签。
	

###  jsinit标签
用以配置js代码初始化工作，不需要写相关的js的代码。配置相关标签即可完成初始化。以下包含与初始化有关的子标签。
#### Interval定时器标签
1. 示例：配置定时器的时间：1秒在标签写上1s。

```bash
 <Interval>
    1s
</Interval>
```
2. 说明：标签内可写内容为数字加上时间单位，单位需要小写，仅支持以下单位：秒（s），分（m），时（h）。配置完成后，会触发回调Interval函数。

#### log日志标签
1. 示例：设置日志保存时间为48小时，分割时间为24小时
```bash
 <jsinit>
      <log maxage="48" rotationTime="24"/>
 </jsinit>
```
2. 说明：属性maxage为日志分割时间，单位48小时。属性RotationTime    若此项不配置，日志为24小时保存一次，永久保存。
#### Database连接数据库标签
1. 示例：连接本地mysql下的test库
```bash
 <Database>
        <conname>root:@(127.0.0.1:3306)/test</conname>
    </Database>
```

2. 说明：子级标签conname中写入数据库连接url，即可连接成功mysql。后续调用js函数访问数据库。

#### http标签
#### tcp标签
#### udp标签
### jscode标签

### jsparams标签
## API接口
### 系统相关与资源监控
#### Getmemory 获取内存大小
```javascript
Getmemory(programname)
```
参数programname:程序名

类型：String

---
返回值：返回某进程内存大小

类型：int

---
应用：获取进程内运行程序的物理内存大小,单位是KB,如果此进程不存在，返回值为-1。

限制：无法获取多个重名进程的内存大小。

示例：获取Notepad++内存大小
```javascript
function init() {
	var memsize = Getmemory("Notepad++.exe")
	console.log(memsize + "KB");
}
```
#### Restartprocess 重启进程
```javascript
Restartprocess(fullname)
```
参数fullname:程序的完整路径名，路径斜杠使用正斜杠

类型：String

---
返回值：无

---
示例：重启Notepad++进程 
```javascript
function init() {
	Restartprocess("C:/Program Files/Notepad++.exe")
}
```
#### Restartservice 重启服务

```javascript
Restartservice(servicename)
```

参数servicename:服务名

类型：String

---
返回值：无

---
应用：使用此函数重启服务需要管理员权限。若服务不存在则无法启动，服务未运行则会启动，服务运行则重启启动。

---
示例：重启MySQL服务
```javascript
function init() {
	Restartservice("MySQL")
}
```
#### filemove 移动文件

```javascript
filemove(src,dst)
```

参数src：原先文件的完整路径名，路径斜杠使用正斜杠

类型：String

参数dst：目的文件的完整路径名，路径斜杠使用正斜杠

类型：String

---
返回值：无

---
应用：从原先文件移动到目的文件所在位置，路径斜杠使用正斜杠。若文件名不存在或者目的文件存在，则无法移动。

---
示例：将Hummingbird图片从D盘1目录移动到2目录
```javascript
function init() {
	filemove("D:/1/Hummingbird.jpg","D:/2/NewHummingbird.jpg")
}
```
#### findfiles 获取文件名
```javascript
findfiles(filename)
```
参数filename:文件完整路径名，路径斜杠使用正斜杠

类型：String

---
返回值：返回匹配的文件名字符串数组

类型：Array

---
应用：获取文件路径下匹配到的文件名，支持通配符匹配

示例：查询以1开头的所有图片
```javascript
function init() {
	var pics =  findfiles("D:/ftp/1*.jpg")
	for (var i = 0; i < pics.length; i++) {
		console.log(pics[i]);
	}
}
```
#### copyfile 复制文件

```javascript
copyfile(src,dst)
```

参数src：原先文件的完整路径名，路径斜杠使用正斜杠

类型：String

参数dst：目的文件的完整路径名，路径斜杠使用正斜杠

类型：String

---
返回值：无

---
应用：从原先文件复制到目的文件所在位置，路径斜杠使用正斜杠。若文件名不存在或者目的文件存在，则无法复制。

---
示例：将Hummingbird图片从D盘1目录复制到2目录
```javascript
function init() {
	filemove("D:/1/Hummingbird.jpg","D:/2/NewHummingbird.jpg")
}
```
#### writefile 写入文件

```javascript
writefile(filename,data)
```

参数filename：文件的完整路径名，路径斜杠使用正斜杠

类型：String

参数data：写入文件的内容，编码为utf-8格式

类型：String

---
返回值：无

---
应用：将字符串写入文件中

限制：写入的文件名所在的路径必须存在。如果文件的内容已经存在，则会覆盖写入

---
示例：在D盘下写入一段内容
```javascript
function init() {
	writefile("D:/1/1.txt","hello Hummingbird")
}
```
#### readfile  读取文件

```javascript
writefile(filename)
```

参数filename：文件的完整路径名，路径斜杠使用正斜杠

类型：String

---
返回值：文件内容，编码为utf-8格式

类型：String

---
应用：读取存在的文件中的内容，以字符串返回

---
示例：在D盘下写入一段内容
```javascript
function init() {
	console.log(readfile("D:/1/1.txt")); 
}
```
#### filestate 获取文件信息
#### dirstate 获取目录下所有文件信息
#### mkdir  创建文件
#### filerename 重命名文件
#### remove 删除文件
#### RemoveBeforeHour 删除几小时之前的文件
#### sleep 睡眠多少毫秒

```javascript
sleep(Millisec)
```

参数Millisec：让程序自动暂停毫秒数

类型：int

---
返回值：无

---
应用：在js代码运行期间，通过调用此函数暂停一段时间。

---
示例：让程序暂停一秒钟
```javascript
function init() {
	sleep(1000); 
}
```

### 日志相关
#### log 将日志写入文件
#### consloe.log 将日志输出到console
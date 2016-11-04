```
{
    "url": "charles-proxy-tool",
    "time": "2016/07/28 19:52",
    "tag": "soft"
}
```

# 一、概述
## 1.1 关于Charles
> 一个HTTP代理服务器，HTTP监视器，反转代理服务器。它允许一个开发者查看所有连接互联网的HTTP通信。这些包括request， response现HTTP headers （包含cookies与caching信息）。

Charles有Win、Mac、Linux版本，**主要功能有：**

- 支持SSL代理。可以截取分析SSL的请求。
- 支持流量控制。可以模拟慢速网络以及等待时间（latency）较长的请求。
- 支持AJAX调试。可以自动将json或xml数据格式化，方便查看。
- 支持AMF调试。可以将Flash Remoting 或 Flex Remoting信息格式化，方便查看。
- 支持重发网络请求，方便后端调试。
- 支持修改网络请求参数。
- 支持网络请求的截获并动态修改。
- 检查HTML，CSS和RSS内容是否符合W3C标准。

## 1.2 Win下安装Charles
Charles用Java语言编写完成，安装Charles之前需要先确认本地是否有安装Java虚拟机。

### 1.2.1 Java虚拟机
先查看本地是否已安装java虚拟机，打开命令行窗口执行java -v看是否出现类似以下结果：
```
C:\Users\pbt>java -version
java version "1.8.0_74"
Java(TM) SE Runtime Environment (build 1.8.0_74-b02)
Java HotSpot(TM) 64-Bit Server VM (build 25.74-b02, mixed mode)
```
若未出现则表示系统未安装Java虚拟机，可直接网上搜索jdk下载并安装。安装完成后重复上面动作确认java虚拟机是否安装完成。

### 1.2.2 安装Charles
- 通过[官网](https://www.charlesproxy.com/download/)下载Charles并安装
- Charles是收费软件，未注册的情况下有一些限制，可下载破解补丁破解。

[charles.jar](/attachments/charles3.9.2-64位破解版.zip) 是Charles3.9.2 64位的补丁文件。


复制补丁文件charles.jar到软件安装根目录的lib目录下，覆盖原文件。启动Charles.exe，看到如下启动界面即表示成功。

![](/static/uploads/charles-start.png)

# 二、Charles简介
## 2.1 界面说明
先来看一下启动之后的界面，Window下启动之后浏览器访问的页面默认就会记录下来了，先打开一个HTTP站点（baidu采用https，建议开个http的站点。），下图是打开OSCHINA后效果图。

![](/static/uploads/charles-preview.jpg)

**先来对各按钮来个初步的认识：**

编号     | 说明
--- | ---
1 | 开启一个新的会话，点击会在15、16上新增一个Session 2的Tab，当前访问的会话将记录在此。
2    | 打开本地保存的会话记录
3     | 关闭当前会话
4     | 保存当前会话，格式如x.chls
5     | 清空当前会话
6     | 查找
7     | 开启记录或者关闭会话记录。只有开启的情况下Charles才会记录请求
8     | Start Throttling，模拟慢速网络。开启后网络请求将明显变慢。具体设置网速的地方：Proxy -> Throttling Settings.
9     | 开启断点，这样子满足条件的请求将会被拦截，我们可以对请求进行编辑后在执行。关于断点的配置：Proxy -> Breakpoints
10     | 重发选中请求，同样也可以在选中请求上右键 - Repeat 或者 Repeat Advanced
11     | 编辑选中请求，可以修改请求数据、Cookie等
12     | 校验请求返回，会请返回数据发送到w3上去校验。
13     | 工具，常用的断点、映射请求、Rewrite、黑名单等都可在此配置，这些工具也可以通过工具栏Tools找到
14     | 设置，常规设置、记录设置、代理设置等
15     | 结构化视图方式。针对一个个站点的方式来展示。
16     | 序列化视方式。针对一个个请求的方式来展示。
17     | 过滤请求，比如上图中只显示oschina的请求。
18     | 概述、可以看到HTTP请求的一些信息，比如，请求时间、返回时间、DNS时间、请求方式、是否压缩等等。
19     | 请求数据，包括请求头、COOKIE、请求参数等
20     | 返回数据，返回头、COOKIE、返回以及一些不同形式的展示方式。
上面编号的功能基本上都可以在菜单栏上找到，可以点击下各个菜单看看。

**熟悉上面的20个功能点之后已可以满足基本的抓包需求，请确保上面提到的点熟悉之后再查看下面章节！**

## 2.2 常用配置
一般PC端浏览器上的抓包并不需要什么配置，这里记录一下可能需要配置的地方。

功能|说明
--- | ---
关闭Window | 去掉Proxy -> Windows Proxy，这样子将不会抓取window系统上的包，抓手机包时为避免混淆可以去掉该功能。
代理设置   | Proxy -> Recording Settings，该界面用来设置记录的请求最大数以及最大内存。如果Charles提示内存不够，可设置该地方。

# 三、高级功能
## 3.1 移动端抓包
移动端抓包即抓取移动端设备里的请求，比如手机端。需要确保设备与电脑处于同一局域网。

### 3.1.1 Charles设置
Proxy -> Proxy Settings或者点击14号按钮里的Proxy Settings，设置端口信息，如下图：

![](/static/uploads/charles-proxy-settings.png)

这里设置端口为8888，端口没有被系统其他工具占用即可。

### 3.1.2 手机端设置
#### 3.1.2.1 查看电脑IP
此处为Win下查看方式，其他系统通过相应方式找到机器的IP。

Win+R 打开命令行，输入cmd， 在弹出的黑窗口输入 ipconfig ，看到类似如下信息：
```
C:\Users\pbt>ipconfig
Windows IP 配置
 
以太网适配器 USB NetWork:
   连接特定的 DNS 后缀 . . . . . . . : ""
   本地链接 IPv6 地址. . . . . . . . : fe80::d53f:f582:c00f:5859%2
   IPv4 地址 . . . . . . . . . . . . : 192.168.1.129
   子网掩码  . . . . . . . . . . . . : 255.255.255.0
   默认网关. . . . . . . . . . . . . : 192.168.1.1
```
可以看到电脑的IP为192.168.1.129。

#### 3.1.2.2 设置手机代理
手机设置 - WLAN- 链接的WLAN - 点击已链接的网络查看详情，将页面移动到最下面，设置代理IP和端口。IP为电脑的IP，端口为Charles里设置的端口，然后保存即可。

**说明：** 不同的手机设置代码的方式不同，但基本都是找到对应的网络在里面去设置。

![](/static/uploads/charles-android.png)

#### 3.1.2.3 开始抓包
此时通过手机网络浏览器或者相应的APP，如果是第一次链接，此时Charles会有如下对话框弹出，请点击Allow

![](/static/uploads/charles-allow.png)

后面的流程就和PC端一样，手机端发请求，Charles记录到请求，查看请求参数即响应结果等等。

## 3.2 篡改请求
篡改请求可以通过断点的开始来设置，可以直接设置在Proxy -> Settings设置，也可以选中要篡改的请求邮件点击Breakpoints。点击之后重新发请求，Charles收到请求后进入到这个页面

![](/static/uploads/chares-breakpoint.png)

我们可以在Edit Request里编辑请求的参数等信息， 编辑完成之后点击执行（Execute），然后收到返回后会同样进入一个这样子的返回页面，也可以对返回数据进行编辑，编辑完成之后再点击Execute。这样子发起请求者收到的将是经过篡改之后的数据。

此时查看Proxy -> Breakpoints里的设置信息将多处一条记录。

![](/static/uploads/chares-breakpoint-settings.png)

可以看到断点是针对请求和返回， 如果只想修改请求数据或者只修改返回数据则编辑一下，去掉对应的勾即可。

篡改请求可以帮助我们做一些特殊处理。比如扣费的时候需要传入用户ID，假设此时我们将用户ID调整为其他人的用户ID，查看请求结果。

## 3.3 重发请求
前面有提到过重发请求，这里可以直接在选中的请求上右键，Repeat是重复一次，Repeat Advanced可以进行高级设置。

![](/static/uploads/charles-repeat-advanced.png)

Iterations代表要执行多少次，Concurrency表示并发数。

我们可以通过重发请求来查看某些请求会不会产生异常处理，比如给用户送券或者增加余额，抓到请求后重复发送，然后查看这次请求是否有效。

## 3.4 映射请求到本地或者远程
有些时候我们可以将一组请求直接映射到本地或者映射到某个其他服务器上，可以通过Tools -> Map Remote 或者 Map Local来设置。比如我们将微信的请求映射到测试环境上，这样子访问生产环境将会映射到测试环境。

![](/static/uploads/charles-map-remote.png)

同样， 我们也可以将请求映射到本地。

映射功能可以让我们将请求转向到任何地方，比如，某个游戏助手需要授权，我们可以将授权接口的请求转向本地，通过本地的返回告诉请求者我授权通过了。当然这里并不是修改了就一定有效，客户端跟服务端也可以有一些签名校验的动作防止篡改。

## 3.5 模拟弱网络
弱网络可以帮助我们模拟网络不好的情况，通过8号按钮可以用来模拟弱网络，打开设置页面可以进行相应设置。

![](/static/uploads/charles-throttling-settings.png)

# 四、常见问题
## 4.1 安装后无法启动
Charles需要安装Java虚拟机，确保是否有安装。
## 4.2 为什么用着用着会自动关闭？
未注册的情况下每次只能使用30分钟，到时间后会自动关闭。需要注册或者破解Charles。
## 4.3 为什么会突然卡死，无法关掉？
这是Charles的BUG，网上提供操作方法：

- 随便抓取一个图片的请求包，选中此请求包
- 点击 Response - Raw。加载完成后内容框内会显示一些乱码。然后就不会卡死了。

## 4.4 Charles异常关闭后无法访问网页。
Charles异常关闭后浏览器的代理未取消，需要再次启动Charles即可访问网页。
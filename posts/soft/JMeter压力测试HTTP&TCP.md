```
{
    "url": "jmeter-http-tcp",
    "time": "2016/01/11 16:21",
    "tag": "常用软件,JMeter,压力测试"
}
```

# 一、概述
Apache JMeter是Apache组织开发的基于Java的压力测试工具。用于对软件做压力测试，它最初被设计用于Web应用测试但后来扩展到其他测试领域。 它可以用于测试静态和动态资源例如静态文件、Java 小服务程序、CGI 脚本、Java 对象、数据库， FTP 服务器， 等等。JMeter 可以用于对服务器、网络或对象模拟巨大的负载，来自不同压力类别下测试它们的强度和分析整体性能。另外，JMeter能够对应用程序做功能/回归测试，通过创建带有断言的脚本来验证你的程序返回了你期望的结果。为了最大限度的灵活性，JMeter允许使用正则表达式创建断言。

# 二、安装与启动
下载地址： http://jmeter.apache.org/download_jmeter.cgi 

JMeter基于Java开发，需要系统有安装JDK环境。解压后进入bin目录，点击`jmeter.bat`。若运行正常则会进入JMeter控制面板。界面如下：

![](/static/uploads/jmeter-start.png)

# 三、常规测试 - HTTP
## 3.1 添加线程组
测试计划 -> 添加 -> `Threads(User)` -> 线程组，添加后进入如下界面：

![](/static/uploads/jmeter-threads.png)


- **线程数：**表示将模拟多少个用户进行测试。
- **Ramp-Up Period(in seconds)：**线程启动间隔，所有线程将在这个时间内依次启动。
- **循环次数：**所有线程执行一次为一次循环。

默认为1个线程执行1次，可等请求调通后再修改此处。

## 3.2 添加采样器
样器可理解为针对前面创建的线程需要做什么事情，这里以添加HTTP请求为例。

线程组 -> 添加 -> `Sampler` -> `HTTP请求`，添加后进入如下界面：

![](/static/uploads/jmeter-HTTP-Sampler.png)

从界面上可以看到，分别可以对GET和POST进行压测。同时需要设置请求的域名或IP，端口，请求路径等等参数。同时可以上传文件、设置代理等。

## 3.3 添加监视器
监视器可以理解为针对结果的不同查看方式。JMeter里提供了多种结果表现形式。可通过 HTTP请求 -> 添加 -> 监视器 -> 察看结果树、聚合报告等等。

## 3.4 运行
切换到察看结果树页面，点击工具栏上的绿色执行按钮。如果请求没错会看到HTTP请求的绿色图标。

![](/static/uploads/jmeter-result.png)

查看请求和响应是否正确。确保无误后可以清空全部（工具栏扫帚按钮）并删掉该页面（留着后面清空比较慢），然后设置线程组中的线程属性。重新执行，在聚合报告页看压测结果，如下图：

![](/static/uploads/jmeter-report.png)

聚合报告各指标：

编号     | 说明
--- | ---
Label|每个 JMeter 的 element（例如 HTTP Request）都有一个 Name 属性，这里显示的就是 Name 属性的值
#Samples|表示你这次测试中一共发出了多少个请求，如果模拟10个用户，每个用户迭代10次，那么这里显示100
Average|平均响应时间——默认情况下是单个 Request 的平均响应时间，单位为毫秒。当使用了 Transaction Controller 时，也可以以Transaction 为单位显示平均响应时间
Median|中位数，也就是 50％ 用户的响应时间
90% Line|90％ 用户的响应时间
Min|最小响应时间
Max|最大响应时间
Error%|本次测试中出现错误的请求的数量/请求的总数
Throughput|吞吐量——默认情况下表示每秒完成的请求数（Request per Second），当使用了 Transaction Controller 时，也可以表示类似 LoadRunner 的 Transaction per Second 数
KB/Sec|每秒从服务器端接收到的数据量，相当于LoadRunner中的Throughput/Sec

# 四、TCP压测
官方TCP文档： https://wiki.apache.org/jmeter/UserManual/Reference/TcpSampler 

同HTTP请求流程，建立线程组后添加TCP取样器，然后根据实际情况配置TCP取样器。

![](/static/uploads/jmeter-TCP-Sampler.png)

TCPClient classname有三种设置：

- **TCPClientImpl：**文本数据
- **BinaryTCPClientImpl：**传输二进制数据，指定包结束符。
- **LengthPrefixedBinaryTCPClientImpl：**数据包中前2个字节为数据长度。可在bin/jmeter.properties配置文件中tcp.binarylength.prefix.length设置。
 
TCP交互数据包一般有两种协议，协议头返回数据长度或者指定结束符。这里以指定结束符为例传输数据为二进制。

- 1、设置TCPClient classname：`org.apache.jmeter.protocol.tcp.sampler.BinaryTCPClientImpl`
- 2、设置结束符End of line(EOL) byte value.
- 3、将数据内容转换为16进制写在要发送的文本区域。

然后就可以测试请求是否正常了。很多时候协议自带的协议可能不满足，就需要自己去写脚本实现了。

以上只是JMeter的一个简单使用，JMeter远远比这个强大，存在Linux下版本，也可以对数据库、FTP、REDIS等进行测试，对请求设置前置、后置处理等等，具体可根据自己的需要来扩展。
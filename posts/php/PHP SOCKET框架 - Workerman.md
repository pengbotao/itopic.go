```
{
    "url": "workerman",
    "time": "2016/08/13 16:04",
    "tag": "PHP"
}
```

# 一、概述
## 1.1 关于Workerman
workerman是一个高性能的PHP socket 服务器框架，workerman基于PHP多进程以及libevent事件轮询库，PHP开发者只要实现一两个接口，便可以开发出自己的网络应用，例如Rpc服务、聊天室服务器、手机游戏服务器等。

workerman的目标是让PHP开发者更容易的开发出基于socket的高性能的应用服务，而不用去了解PHP socket以及PHP多进程细节。 workerman本身是一个PHP多进程服务器框架，具有PHP进程管理以及socket通信的模块，所以不依赖php-fpm、nginx或者apache等这些容器便可以独立运行。

## 1.2 Workerman特性
- 支持HHVM
- 使用PHP开发
- 支持PHP多进程/多线程（多线程版本）
- 标准输入输出重定向
- 支持毫秒定时器
- 支持基于事件的异步编程
- 守护进程化
- 支持TCP/UDP
- 支持多端口监听
- 接口上支持各种应用层协议
- 支持libevent事件轮询库，支持高并发
- 支持服务平滑重启
- 支持PHP文件更新检测及自动加载
- 支持PHP长连接
- 支持以指定用户运行子进程
- 支持telnet远程控制
- 高性能

## 1.3 Workerman与Swoole
- swoole使用C编写，更加底层，以PHP扩展的方式调用。workerman是对PHP原生的socket的封装，类似一个装好的PHP代码库。
- workerman需要PHP安装有pcntl和posix扩展，基于libevent事件轮询。Swoole扩展是基于epoll高性能事件轮询，并且是多线程的。
- swoole更底层，提供的功能更多更强大，workerman更多的是socket服务。
- workerman文档上更清晰，相比感觉上更有态度。

# 二、Workerman功能介绍
## 2.1 Workerman安装
安装方式比较简单，直接下载文件即可。官网地址为：http://www.workerman.net/，文档地址：http://doc3.workerman.net/
## 2.2 Workerman通讯协议
传统PHP开发都是基于Web的，基本上都是HTTP协议，HTTP协议的解析处理都由WebServer独自承担了，所以开发者不会关心协议方面的事情。然而当我们需要基于非HTTP协议开发时，开发者就需要考虑协议的事情了。

WorkerMan目前已经支持HTTP、websocket、text协议(见附录说明)、frame协议(见附录说明)，ws协议(见附录说明)，需要基于这些协议通讯时可以直接使用，使用方法及时在初始化Worker时指定协议，例如

```
use Workerman\Worker;
// websocket://0.0.0.0:2345 表明用websocket协议监听2345端口
$websocket_worker = new Worker('websocket://0.0.0.0:2345');
 
// text协议
$text_worker = new Worker('text://0.0.0.0:2346');
 
// frame协议
$frame_worker = new Worker('frame://0.0.0.0:2347');
 
// tcp Worker，直接基于socket传输，不使用任何应用层协议
$tcp_worker = new Worker('tcp://0.0.0.0:2348');
 
// udp Worker，不使用任何应用层协议
$udp_worker = new Worker('udp://0.0.0.0:2349');
 
// unix domain Worker，不使用任何应用层协议
$unix_worker = new Worker('unix:///tmp/wm.sock');
```

## 2.3 自定义协议
在WorkerMan中如果要实现上面的协议，假设协议的名字叫JsonNL，所在项目为MyApp，则需要以下步骤

- 1、协议文件放到项目的Protocols文件夹，例如文件MyApp/Protocols/JsonNL.php
- 2、实现JsonNL类，以namespace Protocols;为命名空间，必须实现三个静态方法分别为 input、encode、decode

MyApp/Protocols/JsonNL.php的实现

```
namespace Protocols;
class JsonNL
{
    /**
     * 检查包的完整性
     * 如果能够得到包长，则返回包的在buffer中的长度，否则返回0继续等待数据
     * 如果协议有问题，则可以返回false，当前客户端连接会因此断开
     * @param string $buffer
     * @return int
     */
    public static function input($buffer)
    {
        // 获得换行字符"\n"位置
        $pos = strpos($buffer, "\n");
        // 没有换行符，无法得知包长，返回0继续等待数据
        if($pos === false)
        {
            return 0;
        }
        // 有换行符，返回当前包长（包含换行符）
        return $pos+1;
    }
 
    /**
     * 打包，当向客户端发送数据的时候会自动调用
     * @param string $buffer
     * @return string
     */
    public static function encode($buffer)
    {
        // json序列化，并加上换行符作为请求结束的标记
        return json_encode($buffer)."\n";
    }
 
    /**
     * 解包，当接收到的数据字节数等于input返回的值（大于0的值）自动调用
     * 并传递给onMessage回调函数的$data参数
     * @param string $buffer
     * @return string
     */
    public static function decode($buffer)
    {
        // 去掉换行，还原成数组
        return json_decode(trim($buffer), true);
    }
}
```
至此，JsonNL协议实现完毕，可以在MyApp项目中使用，使用方法例如下面。文件：MyApp\start.php

```
use Workerman\Worker;
require_once '/your/path/Workerman/Autoloader.php'
$json_worker = new Worker('JsonNL://0.0.0.0:1234');
$json_worker->onMessage = ...
...
```

# 三、常用示例
## 3.1 HttpServer应用
一个简单的HTTPServer实现方式
```
<?php
use Workerman\Worker;
require_once './Workerman/Autoloader.php';
 
// 创建一个Worker监听2345端口，使用http协议通讯
$http_worker = new Worker("http://0.0.0.0:2345");
 
// 启动4个进程对外提供服务
$http_worker->count = 4;
 
// 接收到浏览器发送的数据时回复hello world给浏览器
$http_worker->onMessage = function($connection, $data)
{
    // 向浏览器发送hello world
    $connection->send('hello world');
};
 
// 运行worker
Worker::runAll();
```
监听2345端口，并且启动4个进程来提供服务。执行性能轻松可上2w，并且比较稳定。但这个是完全没有业务逻辑，加上业务逻辑之后，C语言也可以写到只有300的QPS，后面可以测一下增加业务逻辑后的效果。

![](/static/uploads/workerman-http.png)

查看程序状态，内存上控制的也不错。event-loop为select，如果是libevent可能会更高。

```
[root@asm workerman]# php http.php status
Workerman[http.php] status
---------------------------------------GLOBAL STATUS--------------------------------------------
Workerman version:3.3.4          PHP version:7.0.5
start time:2016-08-13 15:14:18   run 0 days 0 hours  
load average: 1.1, 0.7, 0.43     event-loop:select
1 workers       4 processes
worker_name  exit_status     exit_count
none         0                0
---------------------------------------PROCESS STATUS-------------------------------------------
pid memory  listening           worker_name  connections total_request send_fail throw_exception
18640   2M      http://0.0.0.0:2345 none         24          910535         0         0             
18641   2M      http://0.0.0.0:2345 none         27          1024346        0         0             
18642   2M      http://0.0.0.0:2345 none         24          912041         0         0             
18643   2M      http://0.0.0.0:2345 none         25          949782         0         0
```

### 3.1.1 Workerman高性能原因
官方ab压测给出的结果可以到13w，关于高性能的原因官方给出以下几点说明：
#### 3.1.1.1 避免读取磁盘和反复编译
workerman运行过程中，单个进程生命周期内只要PHP文件被载入编译过一次，便会常驻内存，不会再去从磁盘读取或者再去编译。 workerman省去了重复的请求初始化、创建执行环境、词法解析、语法解析、编译生成opcode以及请求关闭等诸多耗时的步骤。 实际上workerman运行起来后便几乎没有磁盘IO及PHP文件编译开销，剩下的只是简单的计算过程，这使得workerman运行飞快。
#### 3.1.1.2 数据或者资源可以全局共享
workerman中多个请求是可以共享数据或者资源的，当前请求产生的全局变量或者类的静态成员在下一次请求中仍然有效。 这对于减少开销，提高运行效率非常有用。例如业务只要初始化一次数据库连接，那么全局都可以共享这个连接，这实现了真正意义上的数据库长连接。 从而不必每次用户请求都去创建新的连接，避免了连接时三次握手、连接后权限验证以及断开连接时四次挥手等耗时的交互过程。不仅数据库，像redis、 memcache等同样有效。少了这些开销和网络交互，使得workerman运行更快。
#### 3.1.1.3 没有多余的网络开销
传统PHP应用程序需要借助apache、nginx等容器才能对外提供网络服务，这就导致多了一层apache、nginx等容器到PHP进程的数据传输开销。 并且由于运行了apache或者nginx容器，这将进一步消耗服务器的资源。 workerman便没有这部分开销，这是由于workerman自身便是一个服务器容器具有PHP进程管理以及网络通讯功能， 完全不依赖于apache、nginx、php-fpm等这些容器便可以独立运行，所以性能更高。
#### 3.1.1.4 进程模型简单
workerman是多进程（也有多线程版本）的，可以充分利用服务器多核资源。并且workerman具有简单的进程模型，主进程只负责监控子进程，而每个子进程独自接受维护客户端的连接，独自读取连接上发来的数据，独自处理。 子进程间默认没有任何数据通讯，主进程和子进程之间只有一次信号通讯。简单的进程通讯模型使得workerman相比其它复杂的进程模型的软件更高效。
## 3.2 TCP Server应用
这里我们采用自定义协议，为方便测试还是chr(35)为分隔符，实际为前四个字节表示长度，后面跟一串JSON数据。server的写法如下：

```
<?php
require_once './Workerman/Autoloader.php';
use Workerman\Worker;
 
// 运行在主进程
$tcp_worker = new Worker("JsonInt://0.0.0.0:2347");
$tcp_worker->count=2;
// 赋值过程运行在主进程
$tcp_worker->onMessage = function($connection, $data)
{
    // 这部分运行在子进程
    $connection->send(json_encode($data));
};
Worker::runAll();
```

协议中的实现类似上面，在协议的decode方法中已经将数据转换为数组，onMessage回调时收到为PHP数组，这里encode后发给客户端。同样，也进行了一些压测：

![](/static/uploads/workerman-tcp.png)

测试机上效率上也还可以，而且比较稳定。

## 3.3 异步调用
workerman基于进程和子进程的方式来实现的，官方给出的异步的执行方式可参考文档：http://doc3.workerman.net/faq/async-task.html ，但和一般的异步线程相比还是逊色不少，workerman的主要侧重还是在socket通信方面。

# 四、小结
上面只列了http和tcp的示例，workerman还支持更多的协议，比如websocket，而且文档比较清晰，比较稳定。侧重点主要在socket通信这块，所以如果是实现聊天室，TCP长连也是PHPer的一种解决方案。
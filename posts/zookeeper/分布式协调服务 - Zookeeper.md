```
{
    "url": "zookeeper",
    "time": "2019/06/04 06:28",
    "tag": "zookeeper",
    "toc": "yes"
}
```



# 一、概述

`ZooKeeper`是一个分布式的，开源分布式应用程序协调服务，常用来做配置中心，将配置信息存放在`ZooKeeper`的节点上，通过事件监听机制实现配置更新时可以及时通知到监听客户端，通过强一致性确保分布式操作节点数据的可靠性和准确性。常见项目中有用到`ZooKeeper`的有：`Codis`、`Hbase`、`Kafka`、`Spark`、`Dubbo`等。常用的场景有：

## 1.1 配置中心

实现数据发布与订阅模型，将应用中用到的一些配置信息放到ZooKeeper的节点上，客户端可以从该节点读取信息，同时订阅节点的变化，当节点信息有变化时，ZooKeeper会重新通知订阅的客户端，从而实现配置信息的更新。

## 1.2 集群管理

通过`Znode`的特点和`Watcher`机制可以方便的实现集群资源的动态维护。比如负载均衡中节点的维护，当一个新节点启动时向`ZooKeeper`注册一个临时节点，当客户端需要访问时，先从`ZooKeeper`的节点上获取可服务列表，客户端本地根据需求在实现负载均衡的算法。如果出现意外，客户端会话失效，则创建的临时节点会被移除。

即可以借用`ZooKeeper`实现服务的注册与自动发现。

## 1.3 分布式锁

因为相同的节点只能创建一个，如果有多个服务节点创建同一个临时ZNode，则只会有一个服务节点能创建成功，可以用于全局订单号的生成，也可以理解成该服务节点获得了锁，其他服务可以注册监听，业务流程处理完成后可以删除ZNode节点，然后重复前面的过程去竞争创建Znode节点，从而实现排它锁。类似场景就可以用于有事务要求的地方。



由`ZooKeeper`的节点特性、事件监听、强一致性应该还有很多场景可以使用，但核心还是在于配置管理，在分布式服务中起到协调的作用。

# 二、安装与配置

官网：https://zookeeper.apache.org/releases.html

`ZooKeeper`通过`Java`开发，需要`JDK`环境支持，安装过程比较简单，从官网下载解压即可即可。

## 2.1 配置文件

配置文件位于`conf/`目录下，可以拷贝`zoo_sample.cfg`为`zoo.cfg`，调整下数据目录。

```
$ cat conf/zoo.cfg | grep -v ^#
tickTime=2000
initLimit=10
syncLimit=5
dataDir=/Users/peng/data/apache-zookeeper-3.6.2-1/data
clientPort=2181


server.1=127.0.0.1:2888:3888
```

启动脚本也实现好了，在`bin/zkServer.sh`中可以看到操作命令。

```
$ ./bin/zkServer.sh version
ZooKeeper JMX enabled by default
Using config: /Users/peng/data/apache-zookeeper-3.6.2-1/bin/../conf/zoo.cfg
Apache ZooKeeper, version 3.6.2- 09/04/2020 12:44 GMT

$ ./bin/zkServer.sh start
ZooKeeper JMX enabled by default
Using config: /Users/peng/data/apache-zookeeper-3.6.2-1/bin/../conf/zoo.cfg
Starting zookeeper ... STARTED

$ ./bin/zkServer.sh status
ZooKeeper JMX enabled by default
Using config: /Users/peng/data/apache-zookeeper-3.6.2-1/bin/../conf/zoo.cfg
Client port found: 2181. Client address: localhost. Client SSL: false.
Mode: standalone

$ ./bin/zkServer.sh stop
ZooKeeper JMX enabled by default
Using config: /Users/peng/data/apache-zookeeper-3.6.2-1/bin/../conf/zoo.cfg
Stopping zookeeper ... STOPPED
```

配置方式算是比较简单，和`Elasticsearch`差不多，文件参数说明：

| 参数       | 说明                                                         | 示例                             |
| ---------- | ------------------------------------------------------------ | -------------------------------- |
| tickTime   | 客户端与服务端的心跳时间间隔                                 | 2000，单位毫秒                   |
| initLimit  | `Follower`与`Leader`初始化连接超时时间                       | 10，表示10*tickTime              |
| syncLimit  | `Follower`与`Leader`之间心跳超时时间                         | 5，表示5*tickTime                |
| dataDir    | 数据目录                                                     |                                  |
| dataLogDir | 日志目录                                                     |                                  |
| clientPort | 客户端连接端口                                               | 2181                             |
|            |                                                              |                                  |
| server.1   | 服务器+端口配置，端口可自定义<br />- 2888：表示`Follower`与`Leader`之间的通信端口<br />- 3888：表示选举端口 | server.1=172.16.60.100:2888:3888 |
| server.2   |                                                              | server.1=172.16.60.101:2888:3888 |
| ...        |                                                              |                                  |

## 2.2 集群配置

集群的配置方式也比较简单，这里以本机示例，多机器也是一样。

第一步配置文件的调整：`dataDir`、`server.N`、`clientPort`

```
$ cat apache-zookeeper-3.6.2-1/conf/zoo.cfg | grep -v ^#
tickTime=2000
initLimit=10
syncLimit=5
dataDir=/Users/peng/data/apache-zookeeper-3.6.2-1/data
clientPort=2181


server.1=127.0.0.1:2888:3888
server.2=127.0.0.1:2898:3898
server.3=127.0.0.1:2808:3808

$ cat apache-zookeeper-3.6.2-2/conf/zoo.cfg | grep -v ^#
tickTime=2000
initLimit=10
syncLimit=5
dataDir=/Users/peng/data/apache-zookeeper-3.6.2-2/data
clientPort=2182


server.1=127.0.0.1:2888:3888
server.2=127.0.0.1:2898:3898
server.3=127.0.0.1:2808:3808


$ cat apache-zookeeper-3.6.2-3/conf/zoo.cfg | grep -v ^#
tickTime=2000
initLimit=10
syncLimit=5
dataDir=/Users/peng/data/apache-zookeeper-3.6.2-3/data
clientPort=2183


server.1=127.0.0.1:2888:3888
server.2=127.0.0.1:2898:3898
server.3=127.0.0.1:2808:3808
```

第二步：在对应的`data`目录下创建`myid`文件，写入对应的服务编号：

```
$ echo 1 > apache-zookeeper-3.6.2-1/data/myid
$ echo 2 > apache-zookeeper-3.6.2-2/data/myid
$ echo 3 > apache-zookeeper-3.6.2-3/data/myid
```

然后按照前面的方式启动服务即可：

```
$ ./apache-zookeeper-3.6.2-1/bin/zkServer.sh status
ZooKeeper JMX enabled by default
Using config: /Users/peng/data/apache-zookeeper-3.6.2-1/bin/../conf/zoo.cfg
Client port found: 2181. Client address: localhost. Client SSL: false.
Mode: follower

$ ./apache-zookeeper-3.6.2-2/bin/zkServer.sh status
ZooKeeper JMX enabled by default
Using config: /Users/peng/data/apache-zookeeper-3.6.2-2/bin/../conf/zoo.cfg
Client port found: 2182. Client address: localhost. Client SSL: false.
Mode: leader

$ ./apache-zookeeper-3.6.2-3/bin/zkServer.sh status
ZooKeeper JMX enabled by default
Using config: /Users/peng/data/apache-zookeeper-3.6.2-3/bin/../conf/zoo.cfg
Client port found: 2183. Client address: localhost. Client SSL: false.
Mode: follower
```

## 2.3 测试集群

连接到集群：

```
$ ./bin/zkCli.sh -server 127.0.0.1:2181,127.0.0.1:2182,127.0.0.1:2183

[zk: 127.0.0.1:2181,127.0.0.1:2182,127.0.0.1:2183(CONNECTED) 0] ls /
[demo, zookeeper]
[zk: 127.0.0.1:2181,127.0.0.1:2182,127.0.0.1:2183(CONNECTED) 1] ls /zookeeper
[config, quota]
[zk: 127.0.0.1:2181,127.0.0.1:2182,127.0.0.1:2183(CONNECTED) 2] get /zookeeper/config
server.1=127.0.0.1:2888:3888:participant
server.2=127.0.0.1:2898:3898:participant
server.3=127.0.0.1:2808:3808:participant
version=0
```

# 三、节点操作




---

- [1] [ZooKeeper典型应用场景一览](https://www.cnblogs.com/tommyli/p/3766189.html)
- [2] [【ZooKeeper】ZooKeeper安装及简单操作](https://www.cnblogs.com/h--d/p/10269869.html)




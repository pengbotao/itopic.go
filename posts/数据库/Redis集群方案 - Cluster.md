```
{
    "url": "redis-cluster",
    "time": "2019/11/21 05:53",
    "tag": "数据库,Redis",
    "public": "no"
}
```

# 一、下载安装

从Redis官方（`https://redis.io/download`）可以下载到Redis安装包，最新稳定版是6.0.6，安装过程也很简单。

```
$ wget http://download.redis.io/releases/redis-6.0.6.tar.gz
$ tar xzf redis-6.0.6.tar.gz
$ cd redis-6.0.6
$ make
$ sudo make PREFIX=/usr/local/server/redis6.0.6 install
```

会在`PREFIX`目录里创建bin目录以及客户端、服务端的命令工具。

```
pengbotao:redis6.0.6 peng$ pwd
/usr/local/server/redis6.0.6
pengbotao:redis6.0.6 peng$ ls -lh bin/
total 10312
-rwxr-xr-x  1 root  wheel   351K  8  8 05:55 redis-benchmark
-rwxr-xr-x  1 root  wheel   1.5M  8  8 05:55 redis-check-aof
-rwxr-xr-x  1 root  wheel   1.5M  8  8 05:55 redis-check-rdb
-rwxr-xr-x  1 root  wheel   312K  8  8 05:55 redis-cli
lrwxr-xr-x  1 root  wheel    12B  8  8 05:55 redis-sentinel -> redis-server
-rwxr-xr-x  1 root  wheel   1.5M  8  8 05:55 redis-server
```

源代码根目录中有`redis.conf`示例配置文件，utils目录中有一些脚本，`install_server.sh`可以将快速配置Redis的各个配置目录以及配置为服务开机启动。我们先拷贝一个配置文件到自己设定的目录

```
mkdir ~/data/redis/{6380,6381,6382,6383,6384,6385}
cp redis.conf ~/data/redis/6380/6380.conf
```

bind 127.0.0.1
port 6380
daemonize yes
pidfile /Users/peng/data/redis/6380/6380.pid
logfile "/Users/peng/logs/redis/6380.log"
dir /Users/peng/data/redis/6380/data/


尝试启动服务：`$ /usr/local/server/redis6.0.6/bin/redis-server ~/data/redis/6380/6380.conf`

如果启动报错可到日志文件中查看原因。

# 主从配置

增加重节点6381，调整对应配置文件，主要为指定主节点的地址。如果有密码，则指定对应的密码。

```
replicaof 127.0.0.1 6380
```

我们启动2台重节点，查看主库效果：

```
# Replication
role:master
connected_slaves:2
slave0:ip=127.0.0.1,port=6381,state=online,offset=428,lag=1
slave1:ip=127.0.0.1,port=6382,state=online,offset=428,lag=1
master_replid:4d6ec07400aa5d7d4e1af9ecb2cb00aed4b71220
```

- 整个上下线过程不需要重启主节点。
- 重节点上线后会自动同步主节点数据

## 同步逻辑


## 在线切换主从

slaveof no one

Redis Slaveof 命令可以将当前服务器转变为指定服务器的从属服务器(slave server)。

如果当前服务器已经是某个主服务器(master server)的从属服务器，那么执行 SLAVEOF host port 将使当前服务器停止对旧主服务器的同步，丢弃旧数据集，转而开始对新主服务器进行同步。

另外，对一个从属服务器执行命令 SLAVEOF NO ONE 将使得这个从属服务器关闭复制功能，并从从属服务器转变回主服务器，原来同步所得的数据集不会被丢弃。

利用『 SLAVEOF NO ONE 不会丢弃同步所得数据集』这个特性，可以在主服务器失败的时候，将从属服务器用作新的主服务器，从而实现无间断运行。
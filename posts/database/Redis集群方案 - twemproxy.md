```
{
    "url": "twemproxy",
    "time": "2016/07/29 16:13",
    "tag": "数据库,Redis,twemproxy"
}
```

# 一、概述
Twemproxy是Twitter的一个开源项目，用户实现reids、memcache的集群服务，可以有效的避免单点故障，不过twemproxy本身也是一个单点，需要其他服务来实现搞可用。使用twemproxy后后端缓存服务对用户透明，用户直连twemproxy即可。

**Twemproxy特性：**

- 轻量级、快速
- 保持长连接
- 减少了直接与缓存服务器连接的连接数量
- 使用 pipelining 处理请求和响应
- 支持代理到多台服务器上
- 同时支持多个服务器池
- 自动分片数据到多个服务器上
- 实现完整的 memcached 的 ASCII 和再分配协议
- 通过 yaml 文件配置服务器池
- 支持多个哈希模式，包括一致性哈希和分布
- 能够配置删除故障节点
- 可以通过端口监控状态
- 支持 linux, *bsd,os x 和 solaris

# 二、twemproxy安装
twemproxy需要autoconf2.64以上版本，开发机版本偏低重新安装。

## 2.1 安装autoconf
```
# wget http://ftp.gnu.org/gnu/autoconf/autoconf-2.69.tar.gz
# tar zxvf autoconf-2.69.tar.gz
# cd autoconf-2.69
# ./configure --prefix=/usr
# make
# make install
```
安装完成之后执行`autoreconf --version`，可以看到当前版本是2.69
```
# autoreconf --version
autoreconf (GNU Autoconf) 2.69
Copyright (C) 2012 Free Software Foundation, Inc.
License GPLv3+/Autoconf: GNU GPL version 3 or later
```
## 2.2 安装及启动
### 2.2.1 安装
```
# wget https://codeload.github.com/twitter/twemproxy/zip/master
# unzip master
# cd twemproxy-master
# autoreconf -fvi
```
开发机此时报错：
```
configure.ac:36: error: possibly undefined macro: AC_PROG_LIBTOOL
      If this token and others are legitimate, please use m4_pattern_allow.
      See the Autoconf documentation.
autoreconf: /usr/bin/autoconf failed with exit status: 1
```
提示需要安装libtool，
```
# yum install -y libtool libsysfs
# ./configure --prefix=/usr/local/twemproxy
# make
# make install
# cd /usr/local/twemproxy
```
### 2.2.2 配置
创建运行时目录和配置目录
```
# mkdir run conf
```
添加proxy配置文件
```
# vi conf/nutcracker.yml
alpha:
  listen: 127.0.0.1:22121
  hash: fnv1a_64
  distribution: ketama
  auto_eject_hosts: true
  redis: true
  server_retry_timeout: 30000
  server_failure_limit: 1
  servers:
   - 127.0.0.1:6379:1
```

**配置解读：**

**listen** 监听地址和端口（name:port 或者ip:port）,也可以用sock文件（/var/run/nutcracker.sock）的绝对路径

**hash** hash的函数名

- one_at_a_time
- md5
- crc16
- crc32 (crc32 implementation compatible with libmemcached)
- crc32a (correct crc32 implementation as per the spec) 
- fnv1_64
- fnv1a_64（默认配置） 
- fnv1_32
- fnv1a_32
- hsieh
- murmur
- jenkins

**distribution** 数据分配方式。

- ketama：一致性hash算法，根据server构造hash ring，为每个阶段分配hash范围它的优点是一个节点down后，整个集群re-hash，有部分key-range会跟之前的key-range重合，所以它只能合适做单纯的cache
- modula：根据key做hash值取模，根据结果分配到对应的server这种方式如果集群做re-hash，所有的key值都会目标错乱
- random：不管key值的hash结果是啥，随机选取一个server作为操作目标适合只读场景，需要配合数据加载？
- timeout：单位毫秒，等待到server建立连接的时间或者接收server相应过程的等待时间。默认是无限期等待，超时报错：SERVER_ERROR Connection timed out

**auto_eject_host** 当连接一个server失败次数超过server_failure_limit值时，是否把这个server驱逐出集群，默认是false
**redis** 使用redis还是memcached协议，默认false（即memcached）
**server_retry_timeout** 单位毫秒，当auto_eject_host打开后，重试被临时驱逐的server之前的等待时间
**server_failure_limit** 当auto_eject_host打开后，驱逐一个server之前重试次数
**servers** serverpool中包含的的server的地址、端口和权重的列表（name:port:weight or ip:port:weight）

目前开发机只有一个redis端口，配置完成之后检测配置文件是否校验通过：
```
# ./sbin/nutcracker -t
nutcracker: configuration file 'conf/nutcracker.yml' syntax is ok
```

### 2.2.3 启动twemproxy
```
# ./sbin/nutcracker -d -c /usr/local/twemproxy/conf/nutcracker.yml -p /usr/local/twemproxy/run/redisproxy.pid -o /usr/local/twemproxy/run/redisproxy.log
```

# 三、twemproxy使用
使用方式同redis客户端一样，只是将端口改成了22121，但需要注意的是有一些命令不支持，支持的命令可查看：https://github.com/twitter/twemproxy/blob/master/notes/redis.md 
```
# redis-cli -p 22121
# set cache:user:k1 "v1"
# get cache:user:k1
"v1"
```
我们设置可以`cache:user:k1`，设置后可以成功获取。同时我们去后端6379端口查看，数据已经写入。因为这里只有一台redis服务器。我们在开发服新增一个6380端口，用于模拟2台机器，启动后修改twemproxy配置文件并重启。然后在twemproxy中重新设置并查看。
```
127.0.0.1:22121> get cache:user:k1
(nil)
127.0.0.1:22121> set cache:user:k1 "v1.1"
OK
127.0.0.1:22121> get cache:user:k1
"v1.1"
 
127.0.0.1:6379> get cache:user:k1
"v1"
127.0.0.1:6380> get cache:user:k1
"v1.1"
```
可以看到重启后我们原有的KEY在6379端口上还存在， 但是因为新增机器hash值已变化，直接查找时找到6380端口上去了。此时两个端口上都有该KEY，但最新的是在6380上。

# 四、性能测试
经过一层代理后官方给出的极限情况性能下降20%，这里通过redis-benchmark进行一下简单的set、get压测，可以看到性能有一定下降。

## 4.1 twemproxy
```
[root@asm bin]# ./redis-benchmark -h 127.0.0.1 -p 22121 -c 100 -n 1000000 -r 1000000 -d 1024 -t get,set
====== SET ======
  1000000 requests completed in 5.66 seconds
  100 parallel clients
  1024 bytes payload
  keep alive: 1
99.66% <= 1 milliseconds
99.76% <= 2 milliseconds
99.93% <= 3 milliseconds
99.96% <= 4 milliseconds
99.98% <= 5 milliseconds
99.99% <= 6 milliseconds
100.00% <= 7 milliseconds
100.00% <= 7 milliseconds
176584.84 requests per second
====== GET ======
  1000000 requests completed in 5.72 seconds
  100 parallel clients
  1024 bytes payload
  keep alive: 1
99.98% <= 1 milliseconds
100.00% <= 1 milliseconds
174947.52 requests per second
```

# 4.2 redis
```
[root@asm bin]# ./redis-benchmark -h 127.0.0.1 -p 6380 -c 100 -n 1000000 -r 1000000 -d 1024 -t get,set
====== SET ======
  1000000 requests completed in 4.97 seconds
  100 parallel clients
  1024 bytes payload
  keep alive: 1
99.95% <= 1 milliseconds
99.95% <= 2 milliseconds
100.00% <= 3 milliseconds
100.00% <= 3 milliseconds
201166.77 requests per second
====== GET ======
  1000000 requests completed in 4.96 seconds
  100 parallel clients
  1024 bytes payload
  keep alive: 1
100.00% <= 0 milliseconds
201775.62 requests per second
```

原生GET、SET有20w的吞吐量，twemproxy是17w，有15%的差异。目前开发环境已经连到twemproxy，需要连开发机的可使用22121端口。同时，redis3.0版本已提供集群功能，后续可以试试redis本身的集群功能。
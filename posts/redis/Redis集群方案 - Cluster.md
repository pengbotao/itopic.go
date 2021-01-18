```
{
    "url": "redis-cluster",
    "time": "2019/11/21 05:53",
    "tag": "Redis"
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

配置文件做些适当调整：

```
bind 127.0.0.1
port 6380
daemonize yes
pidfile /Users/peng/data/redis/6380/6380.pid
logfile "/Users/peng/logs/redis/6380.log"
dir /Users/peng/data/redis/6380/data/
```

尝试启动服务：`$ /usr/local/server/redis6.0.6/bin/redis-server ~/data/redis/6380/6380.conf`

如果启动报错可到日志文件中查看原因。

# 二、主从模式

Redis的主从常用来做数据的备份，从节点只读，但一般也不提供服务，只有当主节点异常时才将从节点提升为主节点，实现业务的灾备。配置方式也比较简单：

## 2.1 主从配置

增加重节点6381，调整对应配置文件，主要为指定主节点的地址。如果有密码，则指定对应的密码。

```
replicaof 127.0.0.1 6380
```

我们启动2台重节点6381和6382，查看主库效果：

```
127.0.0.1:6380> info replication
# Replication
role:master
connected_slaves:2
slave0:ip=127.0.0.1,port=6381,state=online,offset=428,lag=1
slave1:ip=127.0.0.1,port=6382,state=online,offset=428,lag=1
master_replid:4d6ec07400aa5d7d4e1af9ecb2cb00aed4b71220
```

- 整个从节点上下线过程不需要重启主节点。
- 重节点上线后会自动同步主节点数据
- 一个主节点可以挂N个从节点

## 2.2 主从切换

停止主从，停止主从之后role从slave升级为master，可操作读写，原数据保留。

```
127.0.0.1:6382> slaveof no one
```

也可以重新添加到主从里，会重新做主从同步，差异数据会丢失。

```
127.0.0.1:6382> slaveof 127.0.0.1 6380
```

所以如果做数据迁移是可以考虑直接挂个从节点就可以实现数据的同步迁移。

## 2.3 哨兵（Sentinel）

当Master出现异常时需要人为参与进行主从的切换，而哨兵会监控Master实现自动提升从节点为主节点，同时通知其他从节点切换为新主节点，它可以实现对单份数据的高可用。配置方式同Redis-Server部分差不多，修改配置文件：

```
$ cat sentinel.conf  | grep -v ^# | grep -v ^$
port 26379
daemonize no
pidfile /var/run/redis-sentinel.pid
logfile ""
dir /tmp
sentinel monitor mymaster 127.0.0.1 6379 2
sentinel down-after-milliseconds mymaster 30000
sentinel parallel-syncs mymaster 1
sentinel failover-timeout mymaster 180000
sentinel deny-scripts-reconfig yes
```

启动Sentinel，启动之后配置文件被重写了。

```
$ /usr/local/server/redis6.0.6/bin/redis-sentinel ~/data/redis/sentinel/sentinel.conf

$ /usr/local/server/redis6.0.6/bin/redis-cli -p 26380 info sentinel
# Sentinel
sentinel_masters:1
sentinel_tilt:0
sentinel_running_scripts:0
sentinel_scripts_queue_length:0
sentinel_simulate_failure_flags:0
master0:name=mymaster,status=ok,address=127.0.0.1:6380,slaves=2,sentinels=1
```

当停掉6380的服务后，master自动切换到了6382，从Sentinel的日志文件可以看到切换步骤，同时Sentinel的配置文件也会被重写。也可以配置多个Sentinel，多个Sentinel之间会相互发现，提供可靠性，比如上面示例最后可配置成1主2从3哨兵的模式。当然从服务层面上主从已经进行了切换，但调用方如果不支持哨兵则也不会主动切换连接地址，依然连接到已经下线的节点上去。

# 三、集群模式

## 3.1 集群配置

需要开配置文件中开启cluster-enabled yes，设置下cluster-config-file路径。启动各个服务后可以通过redis-cli进行创建集群。

```
$ /usr/local/server/redis6.0.6/bin/redis-cli --cluster create \
> 127.0.0.1:6380 \
> 127.0.0.1:6381 \
> 127.0.0.1:6382 \
> 127.0.0.1:6383 \
> 127.0.0.1:6384 \
> 127.0.0.1:6385 \
> --cluster-replicas 1


>>> Performing hash slots allocation on 6 nodes...
Master[0] -> Slots 0 - 5460
Master[1] -> Slots 5461 - 10922
Master[2] -> Slots 10923 - 16383
Adding replica 127.0.0.1:6384 to 127.0.0.1:6380
Adding replica 127.0.0.1:6385 to 127.0.0.1:6381
Adding replica 127.0.0.1:6383 to 127.0.0.1:6382
>>> Trying to optimize slaves allocation for anti-affinity
[WARNING] Some slaves are in the same host as their master
M: 7a1da9ed8199d09f1276a6908467693405206f4b 127.0.0.1:6380
   slots:[0-5460] (5461 slots) master
M: b13aeca137a541cd6b8352de69b23b76be0c284e 127.0.0.1:6381
   slots:[5461-10922] (5462 slots) master
M: 9def2a7ff7045ad5afacd89ec1bedbe43466fb04 127.0.0.1:6382
   slots:[10923-16383] (5461 slots) master
S: 0e69804a4b3fb5450dbf089f64d615064b3341e1 127.0.0.1:6383
   replicates 7a1da9ed8199d09f1276a6908467693405206f4b
S: 47b80759dc1e6e1c088539c9bfb2e22bd3ccb730 127.0.0.1:6384
   replicates b13aeca137a541cd6b8352de69b23b76be0c284e
S: 6e48e04238625fc2a6cd55b060a00c7934ab4817 127.0.0.1:6385
   replicates 9def2a7ff7045ad5afacd89ec1bedbe43466fb04
Can I set the above configuration? (type 'yes' to accept): yes


>>> Nodes configuration updated
>>> Assign a different config epoch to each node
>>> Sending CLUSTER MEET messages to join the cluster
Waiting for the cluster to join
.
>>> Performing Cluster Check (using node 127.0.0.1:6380)
M: 7a1da9ed8199d09f1276a6908467693405206f4b 127.0.0.1:6380
   slots:[0-5460] (5461 slots) master
   1 additional replica(s)
M: b13aeca137a541cd6b8352de69b23b76be0c284e 127.0.0.1:6381
   slots:[5461-10922] (5462 slots) master
   1 additional replica(s)
M: 9def2a7ff7045ad5afacd89ec1bedbe43466fb04 127.0.0.1:6382
   slots:[10923-16383] (5461 slots) master
   1 additional replica(s)
S: 6e48e04238625fc2a6cd55b060a00c7934ab4817 127.0.0.1:6385
   slots: (0 slots) slave
   replicates 9def2a7ff7045ad5afacd89ec1bedbe43466fb04
S: 47b80759dc1e6e1c088539c9bfb2e22bd3ccb730 127.0.0.1:6384
   slots: (0 slots) slave
   replicates b13aeca137a541cd6b8352de69b23b76be0c284e
S: 0e69804a4b3fb5450dbf089f64d615064b3341e1 127.0.0.1:6383
   slots: (0 slots) slave
   replicates 7a1da9ed8199d09f1276a6908467693405206f4b
[OK] All nodes agree about slots configuration.
>>> Check for open slots...
>>> Check slots coverage...
[OK] All 16384 slots covered.
```

这里是配置了3主 + 3从，如果不需要从节点去掉副本数量或者设置为0：

```
>>> Performing hash slots allocation on 6 nodes...
Master[0] -> Slots 0 - 2730
Master[1] -> Slots 2731 - 5460
Master[2] -> Slots 5461 - 8191
Master[3] -> Slots 8192 - 10922
Master[4] -> Slots 10923 - 13652
Master[5] -> Slots 13653 - 16383
M: bc6d520f8362b68c7e04833e29f8838713acb108 127.0.0.1:6380
   slots:[0-2730],[3300],[7365],[11298],[15495] (2731 slots) master
M: ebd4c64e35eb85387729cec4eaa122f68be178eb 127.0.0.1:6381
   slots:[2731-5460],[7365],[11298],[15495] (2730 slots) master
M: f9114cdd71bced8c4597e96aa977eef82e92e602 127.0.0.1:6382
   slots:[3300],[5461-8191],[11298],[15495] (2731 slots) master
M: fdf2af81565f2f27f2cd7447ca139c390063006e 127.0.0.1:6383
   slots:[8192-10922] (2731 slots) master
M: f97013bd97fa8592c2664bef7eccd67e24460fd9 127.0.0.1:6384
   slots:[10923-13652] (2730 slots) master
M: 1b0b42b5f5c295e546b6e72fa7288e79c9de2900 127.0.0.1:6385
   slots:[13653-16383] (2731 slots) master
Can I set the above configuration? (type 'yes' to accept): no
```

查看集群信息

```
$ /usr/local/server/redis6.0.6/bin/redis-cli -p 6380 cluster nodes
b13aeca137a541cd6b8352de69b23b76be0c284e 127.0.0.1:6381@16381 master - 0 1610957372000 2 connected 5461-10922
9def2a7ff7045ad5afacd89ec1bedbe43466fb04 127.0.0.1:6382@16382 master - 0 1610957370997 3 connected 10923-16383
7a1da9ed8199d09f1276a6908467693405206f4b 127.0.0.1:6380@16380 myself,master - 0 1610957369000 1 connected 0-5460
6e48e04238625fc2a6cd55b060a00c7934ab4817 127.0.0.1:6385@16385 slave 9def2a7ff7045ad5afacd89ec1bedbe43466fb04 0 1610957370000 3 connected
47b80759dc1e6e1c088539c9bfb2e22bd3ccb730 127.0.0.1:6384@16384 slave b13aeca137a541cd6b8352de69b23b76be0c284e 0 1610957373074 2 connected
0e69804a4b3fb5450dbf089f64d615064b3341e1 127.0.0.1:6383@16383 slave 7a1da9ed8199d09f1276a6908467693405206f4b 0 1610957372042 1 connected
```

测试：

```
$ /usr/local/server/redis6.0.6/bin/redis-cli -c -p 6380
127.0.0.1:6380> set hi Redis
-> Redirected to slot [16140] located at 127.0.0.1:6382
OK
```

## 3.2 slots

```
Master[0] -> Slots 0 - 5460
Master[1] -> Slots 5461 - 10922
Master[2] -> Slots 10923 - 16383
```

Codis设置slots的数量是1024，Redis设置的是16384，通过crc16(k)%16384得出放在哪个slot中。上面示例是将slots平分到3个节点上，只有Master节点才会分配slots，另外三个节点作为备份。每个主节点存储的内容是不一致的，通过横向扩展机器调整slots分布实现扩容的目的，通过从节点增强可用性，但Redis的副本成本还是挺高的，一个副本机器的利用率就降低了一半。


```
{
    "url": "codis",
    "time": "2019/11/11 20:56",
    "tag": "Redis"
}
```

Codis是一个分部署的Redis解决方案，可以很方便的进行Redis的扩容、数据迁移。项目中大量的使用到Codis3，确实很方便。文档也有中文文档，用起来不难。遗憾的是项目已经很久没更新了，貌似凉了。

Github上给的架构图如下：



![](../../static/uploads/codis.png)

几个关键组件说明：

**Codis Server：**

它实际就是个Redis的服务端，基于redis-server做了二次开发，增加额外的数据结构以便支持slot的相关操作。

**Codis Proxy:**

它相当于是Redis的客户端，实现了Redis协议，调用方将请求发送给CodisProxy，CodisProxy再从CodisServer中获取。但在系统配置上他俩是解耦的。也就是可以单独增加Codis Server扩大内存容量，也可以单独提供更多的Proxy处理客户端请求。

**Codis Dashboard:**

集群管理工具，看到DashBoard可能理解为页面管理工具，而实际的页面工具是CodisFE，命令行工具是CodisAdmin。它支持CodisProxy、CodisServer的添加、删除以及数据的迁移等。同一时刻Dashboard只能有0个或者1个，所有集群的操作都必须通过CodisDashboard完成。好在它只是管理工具，数据都是存储在ZooKeeper中，访问是通过CodisProxy访问，这个节点的临时下线不会影响到业务。

**分布式存储：**

项目中使用的ZooKeeper，也支持Etcd、Fs。用来存储Codis元数据。

安装和使用文档上描述的比较清楚，基本组件装好之后就是在FE提供的可视化界面操作。

1、添加Proxy节点

2、添加Server节点。Server为一组，可以包含主从，如果不考虑备份，就是一对一的关系。

如果需要做数据迁移，将Server添加之后通过Migrate就可以将指定Slots迁移到其他机器上，迁移的速度和影响都挺好。要说遗憾的是目前使用的Redis3中不支持碎片的回收只能重启。

## 日常问题

**启动Codis Dashboard**

```
#! /bin/bash
num_1=`ps -ef |grep codis-dashboard |grep -v grep |wc -l`
num_2=`ps -ef |grep codis-fe |grep -v grep |wc -l`

if [ $num_1 -lt 1 ];then
	su - codis -c "/usr/local/codis/codis-dashboard --ncpu=2 --config=/data/codis/conf/dashboard.toml --log=/data/logs/dashboard/dashboard.log --log-level=WARN &"
else
	echo "Codis Dashboard is already running"
fi

if [ $num_2 -lt 1 ];then
        su - codis -c "/usr/local/codis/codis-fe --ncpu=2 --log=/data/logs/dashboard/fe.log --log-level=WARN --zookeeper=192.168.0.100:2181,192.168.0.101:2181,192.168.0.102:2181 --listen=192.168.0.100:8080 &"
else
        echo "Codis fe is already running"
fi
```

**启动Codis Proxy**

```
#! /bin/bash
num=`ps -ef |grep "proxy.toml" |grep -v grep |wc -l`

if [ $num -lt 1 ];then
	su - codis -c "/usr/local/codis/codis-proxy --ncpu=4 --config=/data/proxy/proxy.toml --log=/data/logs/proxy/proxy.log --log-level=WARN &"
else
	echo "Codis proxy is already running"
fi
```

**启动Codis Server**

```
/usr/local/codis/codis-server /data/codis/redis.conf
```

**停止Dashboard**

```
./codis-admin --dashboard=192.168.0.100:18080 --shutdown
```

**[ERROR] store: acquire lock of codistest failed**

重启或者迁移时容易碰到这个问题，需要移除锁。

```
./codis-admin --remove-lock --product=codistest --zookeeper=192.168.0.100:2181,192.168.0.101:2181,192.168.0.102:2181
```

**添加Proxy**

```
./codis-admin --dashboard=192.168.0.100:18080 --create-proxy --addr 192.168.0.100:11080
```

**移除Codis Proxy**

```
./codis-admin --dashboard=192.168.0.100:18080 --remove-proxy --addr=192.168.0.100:11080 --force 
```



---

- [1] https://github.com/CodisLabs/codis
- [2] [使用codis-admin搭建codis集群](https://www.cnblogs.com/zhoujinyi/p/9950105.html)
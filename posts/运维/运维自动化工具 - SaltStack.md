```
{
    "url": "saltstack",
    "time": "2019/07/26 23:50",
    "tag": "运维"
}
```

# 一、概述

## 1.1 关于Salt

`Salt` 是：

- 一个配置管理系统，能够维护预定义状态的远程节点(比如，确保指定的包被安装，指定的服务在运行)
- 一个分布式远程执行系统，用来在远程节点（可以是单个节点，也可以是任意规则挑选出来的节点）上执行命令和查询数据

## 1.2 Salt架构

`Salt`采用`Client/Server`架构。

## 1.3 Salt安装

`Salt`的安装比较简单，通过`yum`就可以直接安装。

### 1.3.1 CentOS7机器

| Hostname      | IP             | 说明                        |
| ------------- | -------------- | --------------------------- |
| peng-master-1 | 172.16.196.200 | 安装salt-master,salt-minion |
| peng-node-1   | 172.16.196.201 | 安装salt-minion             |
| peng-node-2   | 172.16.196.202 | 安装salt-minion             |

### 1.3.2 安装epel源

```
[root@peng-master-1 ~]# wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
[root@peng-master-1 ~]# wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo
[root@peng-master-1 ~]# yum clean all
[root@peng-master-1 ~]# yum makecache
```

### 1.3.3 Master节点安装

```
[root@peng-master-1 ~]# yum install -y salt-master

# 启动Master并设置开机启动(CentOS7)
[root@peng-master-1 ~]# systemctl start salt-master
[root@peng-master-1 ~]# systemctl enable salt-master
```

配置文件在`/etc/salt/master`

### 1.3.4 Node节点安装

```
[root@peng-node-1 ~]# yum install -y salt-minion

# 启动Minion并设置开机启动(CentOS7)
[root@peng-node-1 ~]# systemctl start salt-minion
[root@peng-node-1 ~]# systemctl enable salt-minion
```

配置文件在`/etc/salt/minion`

# 二、Salt配置

## 2.1 Master配置

```
grep -v ^# /etc/salt/master|grep -v ^$
```

## 2.2 Node配置

调整配置文件`/etc/salt/minion`，在`Node`节点中指定`Master`地址。

```
[root@peng-node-1 ~]# vi /etc/salt/minion

# Set the location of the salt master server. If the master server cannot be
# resolved, then the minion will fail to start.
#master: salt
master: 172.16.196.200

[root@peng-node-1 ~]# systemctl restart salt-minion
```



## 2.3 配置认证

`Node`节点配置`master`地址之后重启服务，在`master`通过`salt-key`就可以看到列表。

```
[root@peng-master-1 ~]# salt-key -L
Accepted Keys:
Denied Keys:
Unaccepted Keys:
peng-node-1
peng-node-2
Rejected Keys:
```

需要接受对应的`KEY`， 通过`salt-key -h`可以看到接受与拒绝，这里接受所有的`KEY`。

```
[root@peng-master-1 ~]# salt-key -A
The following keys are going to be accepted:
Unaccepted Keys:
peng-node-1
peng-node-2
Proceed? [n/Y] Y
Key for minion peng-node-1 accepted.
Key for minion peng-node-2 accepted.
```

## 2.4 Salt测试

执行`ping`操作

```
[root@peng-master-1 ~]# salt '*' test.ping
peng-node-1:
    True
peng-node-2:
    True
```

远程执行命令

```
[root@peng-master-1 ~]# salt '*' cmd.run 'free -m'
peng-node-2:
                  total        used        free      shared  buff/cache   available
    Mem:            972         434          62           8         475         383
    Swap:             0           0           0
peng-node-1:
                  total        used        free      shared  buff/cache   available
    Mem:            972         415          89           8         467         402
    Swap:             0           0           0
```

# 三、 常用命令

## 3.1 基础命令

指定所有机器

```
# salt '*' cmd.run 'cp /etc/hosts /etc/hosts.bak'
```

指定多个IP机器

```
# salt -L "127.0.0.1,128.0.0.1" cmd.run 'cat /etc/hosts | grep -E "127.0.0.1|128.0.0.1"'
```

指定设置好的组

```
# salt -N test cmd.run 'netstat -tlnp'
```

## 3.2 统计CPU

统计`CPU`核数

```
# salt '*' cmd.run 'cat /proc/cpuinfo| grep "processor"| wc -l'
```

统计`peng`分组下`CPU`空闲率

```
# 各台机器CPU空闲率，统计10次的平均CPU
$ salt -N peng cmd.run 'top -n 10 -b | grep "Cpu(s):" | awk -F "," "{print \$4}" | awk -F "%" "{a+=\$1}END{print a/10}"'

# 所有机器平均空闲率
$ salt -N peng cmd.run 'top -n 10 -b | grep "Cpu(s):" | awk -F "," "{print \$4}" | awk -F "%" "{a+=\$1}END{print a/10}"' | grep -v ":" | awk '{sum+=$1}END {print sum/NR}'
```

## 3.3 统计内存

统计内存大小

```
# salt '*' cmd.run 'cat /proc/meminfo | head -n 1'
```

统计`peng`分组下内存使用

```
# 各台机器的内存使用率
$ salt -N peng cmd.run 'free |head -n 2| tail -n 1' | grep Mem | awk '{print ($4+$6+$7)*100/$2}'

# 所有机器的内存使用率
$ salt -N peng cmd.run 'free |head -n 2| tail -n 1' | grep Mem | awk '{a+=$4;b+=$6;c+=$7;d+=$2}END{print (a+b+c)*100/d}'
```

## 3.4 Host进行IP替换

替换`test`分组下的`Host`配置，`sed`去掉`-i`不会执行，只打印替换后的信息。

```
# salt -N test cmd.run 'sed -i "s/.*test.salt.com/128.0.0.1 test.salt.com/" /etc/hosts'
```



![](../../static/uploads/saltstack.png)

- [1] [saltstack高效运维](https://www.cnblogs.com/xintiao-/p/10380656.html)


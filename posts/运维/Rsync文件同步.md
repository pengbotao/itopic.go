```
{
    "url": "rsync",
    "time": "2018/08/01 15:21",
    "tag": "运维"
}
```

# 一、简介

rsync是一个类似scp的同步工具，会检测变化的文件，只传输有变化的部分，同时可以方便的保留原文件的所有者、时间、软连等信息。可以直接当命令行工具使用，也可以在服务端以后台的形式使用。

# 二、安装

```
$ rpm -qa | grep rsync
rsync-3.0.6-9.el6_4.1.x86_64

$ yum -y install rsync
```

# 三、客户端工具

可以当命令使用，安装之后不需要做任何配置，直接使用即可，用法和scp差不多。

## 3.1 同步本机目录

```
$ rsync -av /demo/backup/ /demo/package/
```

将`backup`目录下的所有文件同步到`package`目录下。需要注意`backup`后面的斜杠，有斜杠的时候表示目录下的文件，没有斜杠的时候同步的是backup目录。

## 3.2 同步本地到远程

```
$ rsync -avzP /data root@172.16.0.1:/
```

## 3.3 同步远程到本地

```
$ rsync -avzP root@172.16.0.1:/data/ /data/
```

# 四、服务端运行

`rsync`也可以在服务端以后台形式运行，默认会启动873端口。该方式需要做一些简单的配置。

## 4.1 密码配置

创建文件并设置文件权限为600

```
$ vi /etc/rsyncd.pass
rsync:123456

$ chmod 600 /etc/rsyncd.pass
```

## 4.2 主配置文件

配置文件位于`/etc/rsyncd.conf`，配置示例：

```
$ vi cat /etc/rsyncd.conf
port = 873
max connections = 50
pid file = /var/run/rsyncd.pid
lock file = /var/run/rsyncd.lock
log file = /var/log/rsyncd.log
# hosts allow = 172.16.0.0/12

[demo]
uid = rsync
gid = rsync
path = /demo/
comment = this is rsync comment
use chroot = yes
read only = no
auth users = rsync
secrets file = /etc/rsyncd.pass
```

启动服务

```
$ systemctl start rsyncd
```

或者

```
$ rsync --daemon --config=/etc/rsyncd.conf
```

## 4.3 客户端

客户端需要安装`rsync`，配置秘钥文件，文件中只需要写密码即可，然后给与600权限

```
$ rsync -avzP --password-file=rsync.secrets rsync@172.16.196.200::demo/ .
```

# 五、示例

同步目录下指定类型的文件(只同步py文件)

```
$ rsync -avzP --include="*.py" --exclude="*.*"  root@172.16.0.1:/data/project/  /data/project/
```





---

- [1] [rsync 用法教程](http://www.ruanyifeng.com/blog/2020/08/rsync.html)
```
{
    "url": "nfs",
    "time": "2018/09/02 09:06",
    "tag": "运维"
}
```

# 一、NFS服务端安装

查询是否有安装`nfs`服务，没有输出则表示没有安装。

```
$ rpm -qa | grep nfs
```

安装`nfs`只需要安装： 

```
$ yum install -y nfs-utils
```

`nfs`依赖`rpcbind`，`yum`安装会自动解决依赖，安装之后再查询就可以查询到了

```
$ rpm -qa | grep nfs
nfs-utils-1.2.3-78.el6_10.2.x86_64
nfs-utils-lib-1.1.5-13.el6.x86_64
```

如果要设置开机启动，可以这么操作：

```
$ chkconfig --level 35 nfs on
```

到这里安装过程就完成了。

# 二、NFS服务端启动

查看`nfs`服务状态：

```
$ service nfs status
rpc.svcgssd is stopped
rpc.mountd is stopped
nfsd is stopped
rpc.rquotad is stopped
```

如果这时直接启动`nfs`，在`CentOS 6`下会报错：

```
$ service nfs start
Starting NFS services:                                     [  OK  ]
Starting NFS quotas: Cannot register service: RPC: Unable to receive; errno = Connection refused
rpc.rquotad: unable to register (RQUOTAPROG, RQUOTAVERS, udp).
                                                           [FAILED]
Starting NFS mountd: rpc.mountd: svc_tli_create: could not open connection for udp6
rpc.mountd: svc_tli_create: could not open connection for tcp6
rpc.mountd: svc_tli_create: could not open connection for udp6
rpc.mountd: svc_tli_create: could not open connection for tcp6
rpc.mountd: svc_tli_create: could not open connection for udp6
rpc.mountd: svc_tli_create: could not open connection for tcp6
                                                           [FAILED]
Starting NFS daemon: rpc.nfsd: writing fd to kernel failed: errno 111 (Connection refused)
rpc.nfsd: address family inet6 not supported by protocol TCP
rpc.nfsd: unable to set any sockets for nfsd
                                                           [FAILED]
```

需要先启动`rpcbind`：

```
$ service rpcbind status
rpcbind is stopped

$ service rpcbind start
Starting rpcbind:                                          [  OK  ]
```

然后在启动`nfs`服务：

```
$ service nfs start
Starting NFS services:                                     [  OK  ]
Starting NFS quotas:                                       [  OK  ]
Starting NFS mountd: rpc.mountd: svc_tli_create: could not open connection for udp6
rpc.mountd: svc_tli_create: could not open connection for tcp6
rpc.mountd: svc_tli_create: could not open connection for udp6
rpc.mountd: svc_tli_create: could not open connection for tcp6
rpc.mountd: svc_tli_create: could not open connection for udp6
rpc.mountd: svc_tli_create: could not open connection for tcp6
                                                           [  OK  ]
Starting NFS daemon: rpc.nfsd: address family inet6 not supported by protocol TCP
                                                           [  OK  ]
Starting RPC idmapd:                                       [  OK  ]
```

查看状态：

```
$ service nfs status
rpc.svcgssd is stopped
rpc.mountd (pid 27786) is running...
nfsd (pid 27806 27805 27804 27803 27802 27801 27800 27799) is running...
rpc.rquotad (pid 27769) is running...
```

# 三、NFS服务端配置

配置文件在`/etc/exports`，配置文件的格式是：

```
/PATH/TO/DIR HOST([OPTIONS])
```

如，

```
$ cat /etc/exports
/data/logs 172.16.0.20(rw,async,no_root_squash)
```

其中：`/data/logs`表示共享目录，`172.16.0.20`表示共享给指定`IP`，括号内为`nfs`共享的相关参数。

文件编辑后可以通过`exportfs -v`校验，通过`exportfs -r`重新加载。

## 3.1 Client设置

1、共享给所有主机

```
/data/logs *(sync)
```

2、共享给特定IP段，共享给`172.16.0.*`的机器。IP也可以写成：`172.16.0.0/24`

```
/data/logs 172.16.0.0/255.255.255.0(sync)
```

3、共享给多主机，共享给`172.16.0.*` 和 `172.16.1.*`

```
/data/logs 172.16.0.0/255.255.255.0(sync)
/data/logs 172.16.1.0/255.255.255.0(sync)
```

## 3.2 OPTIONS参数说明：

| 选项             | 说明                                                         |
| ---------------- | ------------------------------------------------------------ |
| `rw`             | 读写访问                                                     |
| `ro`             | 只读访问                                                     |
| `root_squash`    | 把客户端`root`账号当普通用户对待（默认）                     |
| `no_root_squash` | 客户端`root`具有超级权限                                     |
| `all_squash`     | 共享文件的`UID`和`GID`映射匿名用户`anonymous`，适合公用目录  |
| `no_all_squash`  | 保留共享文件的`UID`和`GID`（默认）                           |
| `sync`           | 同步写入，有修改时同步写入                                   |
| `async`          | 可以异步写入，通常可以提升性能，但数据没有实时落地，有异常时可能有丢失。 |

关于`root_squash`说明：

```
$ cat /etc/exports
/data/test 172.16.60.7(sync,rw,root_squash)
$ exportfs -v
/data/test    	172.16.0.30(rw,wdelay,root_squash,no_subtree_check,sec=sys,rw,root_squash,no_all_squash)
```

挂载之后，客户端`root`账号创建的文件属于`nfsnobody`，如果服务端`nfsnobody`没有`test`目录权限，那么客户端`root`用户也无法写入，即把客户端`root`用户当普通用户看待。

如果调整为`no_root_squash`，则`root`相当于有超级用户的权限，可以创建文件，同时创建的文件所有者为`root`。

**配置完成之后重新加载`exports`文件：**

```
$ exportfs -r
```

# 四、客户端挂载

客户端也需要安装`yum install -y nfs-utils`，挂载命令是：

```
$ mount -t nfs 172.16.0.10:/data/logs/ /data/logs/
```

查看挂载情况：

```
$ df -Th
Filesystem           Type   Size  Used Avail Use% Mounted on
/dev/vda1            ext4    40G   11G   27G  30% /
tmpfs                tmpfs  7.8G   16K  7.8G   1% /dev/shm
/dev/vdb1            ext4   985G  679G  257G  73% /data
172.16.0.10:/data/logs/
                     nfs    985G  349G  586G  38% /data/logs
```

如果要取消挂载：

```
 $ umount 172.16.0.10:/data/logs
```

也可以通过编辑`/etc/fstab`文件实现客户端开机自动挂载：

```
172.16.0.10:/data/logs /data/logs   nfs    defaults        0  0
```



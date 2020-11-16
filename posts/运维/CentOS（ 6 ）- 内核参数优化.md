```
{
    "url": "linux-sysctl",
    "time": "2017/12/24 19:30",
    "tag": "运维",
    "toc": "yes",
    "public": "no"
}
```



# 文件句柄数限制

`Linux`系统对单用户`单个进程`可打开的文件数量有限制，每一个`TCP`连接都需要创建一个`Socket`句柄，每个`Socket`句柄同时也是一个文件句柄，这个设置会直接影响到系统可支撑的最大连接数。

默认单个进程最大可打开`1024`个文件，通过`ulimit -n`可以查看打开的最大文件数限制：

```
$ ulimit -n
1024
```

限制分为软限制和硬限制，`ulimit -n`显示的是软限制的数量，用户可以设置自己的限制数量，但受限于最大硬限制设置的值。

```
[peng@peng-master-1 ~]$ ulimit -SHn 1000
[peng@peng-master-1 ~]$ ulimit -SHn 10000
-bash: ulimit: open files: 无法修改 limit 值: 不允许的操作
```

> 注：ulimit命令只影响当前Shell环境

如果要永久生效，可以通过调整文件`/etc/security/limits.conf`

```
$ cat /etc/security/limits.conf
* soft nofile 65536
* hard nofile 65536
```

其中类型的值可以为`hard`，`soft`或者`-`，其中`-`相当于同时配置了`hard`和`soft`两行。

软限制(`soft nofile`)是指`Linux`在当前系统能够承受的范围内进一步限制用户同时打开的文件数；硬限制(`hard nofile`)则是根据系统硬件资源状况(主要是系统内存)计算出来的系统最多可同时打开的文件数量。通常软限制小于或等于硬限制。

修改后重新登录即可生效。同时系统层面也有可打开的文件数量限制，根据系统的资源情况计算出来的。

**系统层面限制**

查看与临时修改系统打开文件总数限制：

```
$ cat /proc/sys/fs/file-max
95079
$ echo 1620826 > /proc/sys/fs/file-max
```

永久修改系统打开文件最大值限制：

```
$ sysctl -a | grep file-max
fs.file-max = 1620826

$ vi /etc/sysctl.conf
fs.file-max = 1000000

# 立即生效：
$ sysctl -p
```

**查看进程资源限制与打开的文件**

```
$ cat /proc/2309/limits
Limit                     Soft Limit           Hard Limit           Units
Max cpu time              unlimited            unlimited            seconds
Max file size             unlimited            unlimited            bytes
Max data size             unlimited            unlimited            bytes
Max stack size            10485760             unlimited            bytes
Max core file size        0                    unlimited            bytes
Max resident set          unlimited            unlimited            bytes
Max processes             63691                63691                processes
Max open files            1024                 4096                 files
Max locked memory         65536                65536                bytes
Max address space         unlimited            unlimited            bytes
Max file locks            unlimited            unlimited            locks
Max pending signals       63691                63691                signals
Max msgqueue size         819200               819200               bytes
Max nice priority         0                    0
Max realtime priority     0                    0
Max realtime timeout      unlimited            unlimited            us

$ lsof -p 2309
COMMAND    PID   USER   FD   TYPE             DEVICE SIZE/OFF    NODE NAME
zabbix_se 2309 zabbix  cwd    DIR              252,1     4096       2 /
zabbix_se 2309 zabbix  rtd    DIR              252,1     4096       2 /
zabbix_se 2309 zabbix  txt    REG              252,1  3506157 1050689 /usr/local/zabbix/sbin/zabbix_server
zabbix_se 2309 zabbix  mem    REG              252,1   803410  267760 /lib64/ld-2.18.so
zabbix_se 2309 zabbix  mem    REG              252,1 10114254  269365 /lib64/libc-2.18.so
zabbix_se 2309 zabbix  mem    REG              252,1   858324  269378 /lib64/libpthread-2.18.so
zabbix_se 2309 zabbix  mem    REG              252,1   102750  269377 /lib64/libdl-2.18.so
zabbix_se 2309 zabbix  mem    REG              252,1  2452986  269371 /lib64/libm-2.18.so
zabbix_se 2309 zabbix  mem    REG              252,1    91096  269369 /lib64/libz.so.1.2.3
zabbix_se 2309 zabbix  mem    REG              252,1   172659  269388 /lib64/librt-2.18.so
zabbix_se 2309 zabbix  mem    REG              252,1   340018  269509 /lib64/libresolv-2.18.so
zabbix_se 2309 zabbix  mem    REG              252,1   528075  269246 /lib64/libnsl-2.18.so
```



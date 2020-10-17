```
{
    "url": "linux-monitor",
    "time": "2017/12/10 21:30",
    "tag": "运维",
    "toc": "yes"
}
```





```
[root@peng ~]# top
top - 04:57:54 up 236 days, 14:03,  1 user,  load average: 0.00, 0.01, 0.00
Tasks: 377 total,   1 running, 375 sleeping,   0 stopped,   1 zombie
Cpu0  :  0.3%us,  0.3%sy,  0.0%ni, 99.0%id,  0.3%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu1  :  0.3%us,  0.7%sy,  0.0%ni, 99.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu2  :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu3  :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu4  :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu5  :  0.3%us,  0.7%sy,  0.0%ni, 99.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu6  :  0.3%us,  0.0%sy,  0.0%ni, 99.7%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu7  :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu8  :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu9  :  0.0%us,  0.3%sy,  0.0%ni, 99.7%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu10 :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu11 :  0.3%us,  0.3%sy,  0.0%ni, 99.3%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu12 :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu13 :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu14 :  0.0%us,  0.3%sy,  0.0%ni, 99.7%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Cpu15 :  0.0%us,  0.0%sy,  0.0%ni,100.0%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Mem:  16465312k total,  2866440k used, 13598872k free,   324020k buffers
Swap:        0k total,        0k used,        0k free,  2144800k cached
```



可以看到平均负载，按1可以看到Cpu数量和内存信息。



# 一、CPU

## 1.1 sar

```
[root@peng ~]# sar 1
Linux 2.6.32-754.25.1.el6.x86_64 (peng) 	2020年10月13日 	_x86_64_	(16 CPU)

05时02分19秒     CPU     %user     %nice   %system   %iowait    %steal     %idle
05时02分20秒     all      0.06      0.00      0.06      0.00      0.00     99.87
05时02分21秒     all      0.06      0.00      0.06      0.00      0.00     99.88
05时02分22秒     all      0.00      0.00      0.06      0.00      0.00     99.94
05时02分23秒     all      0.19      0.00      0.38      0.00      0.00     99.44
05时02分24秒     all      0.13      0.00      0.19      0.06      0.00     99.62
05时02分25秒     all      0.12      0.00      0.06      0.00      0.00     99.81
05时02分26秒     all      0.00      0.00      0.06      0.00      0.00     99.94
05时02分27秒     all      0.06      0.00      0.00      0.00      0.00     99.94
05时02分28秒     all      0.12      0.00      0.12      0.00      0.00     99.75
05时02分29秒     all      0.13      0.00      0.13      0.00      0.00     99.75
05时02分30秒     all      0.06      0.00      0.06      0.00      0.00     99.87
```

## 1.2 vmstat

```
[root@peng ~]# vmstat 1
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 1  0      0 13522672 344480 2144948    0    0     0     3    0    0  6  4 90  0  0
 0  0      0 13522608 344480 2144948    0    0     0     0 1530 2430  0  0 100  0  0
 0  0      0 13522944 344480 2144948    0    0     0     0 1504 2426  0  0 100  0  0
 0  0      0 13522864 344480 2144948    0    0     0     0 1563 2476  0  0 100  0  0
```

# 二、内存

## 2.1 free

```
[root@peng ~]# free -h
             total       used       free     shared    buffers     cached
Mem:           15G       2.7G        12G       1.1M       316M       2.0G
-/+ buffers/cache:       386M        15G
Swap:           0B         0B         0B
```

# 三、磁盘

(io)

## 3.1 df

查看磁盘整体情况

```
[root@peng ~]# df -h
Filesystem      Size  Used Avail Use% Mounted on
/dev/vda1        40G  2.9G   35G   8% /
tmpfs           7.8G   12K  7.8G   1% /dev/shm
/dev/vdb1       197G  1.6G  186G   1% /data
```



## 3.2 du

查看某个目录下文件占比

```
[root@peng ~]# du -h /home --max-depth=1
16K	/home/zabbix
104K	/home/peng
188K	/home
```

## 3.3 iostat



# 四、进程

## 4.1 ps



 ## 4.2 pstree



## 4.3 strace

# 五、网络连接

## 5.1 netstat

```
# netstat -tlnp
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address               Foreign Address             State       PID/Program name
tcp        0      0 0.0.0.0:22                  0.0.0.0:*                   LISTEN      2281/sshd
tcp        0      0 0.0.0.0:10050               0.0.0.0:*                   LISTEN      16856/zabbix_agentd
```

## 5.2 lsof



## 5.3 ss

```
[root@peng ~]# ss -s
Total: 65 (kernel 173)
TCP:   71 (estab 4, closed 59, orphaned 0, synrecv 0, timewait 59/0), ports 18

Transport Total     IP        IPv6
*	  173       -         -
RAW	  0         0         0
UDP	  3         3         0
TCP	  12        12        0
INET	  15        15        0
FRAG	  0         0         0
```



## 5.4 smokeping



# 六、网络流量

iftop







https://www.cnblogs.com/peida/archive/2013/03/11/2953420.html

https://www.jb51.net/article/50436.htm
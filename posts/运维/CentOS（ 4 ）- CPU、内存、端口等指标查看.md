```
{
    "url": "linux-monitor",
    "time": "2017/12/10 21:30",
    "tag": "运维,CentOS",
    "toc": "yes"
}
```

# 一、CPU

## 1.1 top

```
$ top
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

可以看到平均负载，按1可以看到Cpu数量和内存信息。CPU指标说明：

- us(%user)：用户程序使用CPU时间片占比
- sy(%system)：内核使用的CPU时间占比
- id（%idle）：CPU空闲时间比
- wa($iowait)：CPU处理IO等待的时间占比

## 1.2 sar

```
$ sar 1
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

## 1.3 vmstat

```
$ vmstat 1
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
$ free -h
             total       used       free     shared    buffers     cached
Mem:           15G       2.7G        12G       1.1M       316M       2.0G
-/+ buffers/cache:       386M        15G
Swap:           0B         0B         0B

$ free -h -s 3
             total       used       free     shared    buffers     cached
Mem:           30G        30G       692M       696K       446M       4.5G
-/+ buffers/cache:        25G       5.6G
Swap:           0B         0B         0B

             total       used       free     shared    buffers     cached
Mem:           30G        30G       659M       696K       446M       4.6G
-/+ buffers/cache:        25G       5.6G
Swap:           0B         0B         0B
```

指标说明：

- `- buffers/cache -> used`：实际使用内存，等于Mem行中used - buffers - cached
- `+ buffers/cache -> free`：实际空闲内存，等于Mem行中free + buffers + cached

## 2.2  /proc/meminfo

```
cat /proc/meminfo
MemTotal:       32478432 kB
MemFree:          469660 kB
Buffers:          457052 kB
Cached:          4982444 kB
SwapCached:            0 kB
Active:         27215592 kB
Inactive:        4155144 kB
Active(anon):   25931508 kB
Inactive(anon):      664 kB
Active(file):    1284084 kB
Inactive(file):  4154480 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:             0 kB
SwapFree:              0 kB
Dirty:             96576 kB
Writeback:             0 kB
AnonPages:      25931476 kB
Mapped:            53964 kB
Shmem:               696 kB
Slab:             335560 kB
SReclaimable:     292696 kB
SUnreclaim:        42864 kB
KernelStack:       23056 kB
PageTables:        58696 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    16239216 kB
Committed_AS:   30251252 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       65396 kB
VmallocChunk:   34359641124 kB
HardwareCorrupted:     0 kB
AnonHugePages:  24539136 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:        5944 kB
DirectMap2M:     2615296 kB
DirectMap1G:    30408704 kB
```

# 三、磁盘

## 3.1 df

查看磁盘整体情况

```
$ df -h
Filesystem      Size  Used Avail Use% Mounted on
/dev/vda1        40G  2.9G   35G   8% /
tmpfs           7.8G   12K  7.8G   1% /dev/shm
/dev/vdb1       197G  1.6G  186G   1% /data
```

## 3.2 du

查看某个目录下文件占比

```
$ du -h /home --max-depth=1
16K	/home/zabbix
104K	/home/peng
188K	/home

$ du -sh
208K	.
```

## 3.3 iostat

```
$ iostat -x 1
Linux 2.6.32-754.28.1.el6.x86_64 (peng-master-1) 	2021年01月14日 	_x86_64_	(8 CPU)

avg-cpu:  %user   %nice %system %iowait  %steal   %idle
           3.02    0.00    1.27    0.25    0.00   95.46

Device:         rrqm/s   wrqm/s     r/s     w/s   rsec/s   wsec/s avgrq-sz avgqu-sz   await r_await w_await  svctm  %util
vda               0.01     1.14    0.24    1.15    13.01    18.38    22.55     0.00    0.93    0.57    1.01   0.32   0.04
vdb               4.26  1459.26   65.23  115.17  5675.69 12595.45   101.28     1.61    8.94    1.91   12.92   0.19   3.44
```

指标说明：

- rrqm/s： 每秒进行merge的读操作数
- wrqm/s： 每秒进行merge的写操作数
- r/s： 每秒读次数
- w/s：每秒写次数
- rsec/s：每秒读扇区数
- wsec/s：每秒写扇区数
- rkB/s：每秒读字节数，单位kB
- wkB/s：每秒写字节数
- avgrq-sz：平均每次设备IO操作的数据大小（扇区）
- avgqui-sz：平均IO队列长度
- await：平均每次设备IO操作的等待时间（毫秒）
- svctm：平均每次设备IO操作的服务时间（毫秒）
- %util：一秒之中有多少时间用于IO操作，如果接近100%说明产生的IO请求太多，可能存在瓶颈。

# 四、进程

## 4.1 ps

```
$ ps -ef | grep api-test
peng  16118 14141  0 10:35 ?        00:00:01 gunicorn: master [api-test]
peng  16156 16118  0 10:36 ?        00:00:20 gunicorn: worker [api-test]
peng  16164 16118  0 10:36 ?        00:00:20 gunicorn: worker [api-test]
peng  16171 16118  0 10:36 ?        00:00:20 gunicorn: worker [api-test]
peng  16185 16118  0 10:36 ?        00:00:20 gunicorn: worker [api-test]
```

## 4.2 lsof

```
$ lsof -i:33040
COMMAND     PID    USER   FD   TYPE  DEVICE SIZE/OFF NODE NAME
gunicorn: 16118    peng    6u  IPv4 2542870      0t0  TCP *:33040 (LISTEN)
gunicorn: 16156    peng    6u  IPv4 2542870      0t0  TCP *:33040 (LISTEN)
gunicorn: 16164    peng    6u  IPv4 2542870      0t0  TCP *:33040 (LISTEN)
gunicorn: 16171    peng    6u  IPv4 2542870      0t0  TCP *:33040 (LISTEN)
gunicorn: 16185    peng    6u  IPv4 2542870      0t0  TCP *:33040 (LISTEN)
```

## 4.3 pstree

```
$ pstree -p 16118
gunicorn:\040maste(16118)─┬─gunicorn:\040worke(16156)─┬─{gunicorn:\040work}(16196)
                          │                           ├─{gunicorn:\040work}(16393)
                          │                           └─{gunicorn:\040work}(16394)
                          ├─gunicorn:\040worke(16164)─┬─{gunicorn:\040work}(16202)
                          │                           ├─{gunicorn:\040work}(16389)
                          │                           ├─{gunicorn:\040work}(16390)
                          │                           ├─{gunicorn:\040work}(16391)
                          │                           └─{gunicorn:\040work}(16392)
                          ├─gunicorn:\040worke(16171)─┬─{gunicorn:\040work}(16204)
                          │                           ├─{gunicorn:\040work}(16383)
                          │                           ├─{gunicorn:\040work}(16384)
                          │                           └─{gunicorn:\040work}(16385)
                          └─gunicorn:\040worke(16185)─┬─{gunicorn:\040work}(16215)
                                                      ├─{gunicorn:\040work}(16380)
                                                      ├─{gunicorn:\040work}(16381)
                                                      └─{gunicorn:\040work}(16382)
```

## 4.4 strace

```
$ strace -p 16118
Process 16118 attached - interrupt to quit
select(5, [4], [], [], {0, 209202})     = 0 (Timeout)
gettimeofday({1610618741, 795149}, NULL) = 0
fstat(9, {st_mode=S_IFREG, st_size=0, ...}) = 0
```

# 五、网络连接

## 5.1 netstat

```
$ netstat -tlnp
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address               Foreign Address             State       PID/Program name
tcp        0      0 0.0.0.0:22                  0.0.0.0:*                   LISTEN      2281/sshd
tcp        0      0 0.0.0.0:10050               0.0.0.0:*                   LISTEN      16856/zabbix_agentd
```

## 5.2 ss

```
$ ss -s
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

# 六、网络流量

## 6.1 iftop

```
$ iftop -i eth0
19.1Mb              38.1Mb                    57.2Mb                 76.3Mb
└─────────────────────┴───────────────────────┴────────────────────────┴─────────────────────────
peng-master-1                       => 172.16.0.1                   740Kb  1.11Mb   988Kb
                                    <=                              36.9Mb  47.4Mb  44.0Mb

─────────────────────────────────────────────────────────────────────────────────────────────────
TX:             cum:   85.3MB   peak:   20.0Mb                    rates:   5.33Mb  8.62Mb  8.07Mb
RX:                     563MB           66.0Mb                             42.6Mb  54.7Mb  50.8Mb
TOTAL:                  649MB           81.4Mb                             47.9Mb  63.3Mb  58.9Mb
```

指标说明：

- `TX`：发送流量（Outgoing network traffic on eth0）
- `RX`：接收流量（Incoming network traffic on eth0）
- `TOTAL`：总流量，cum：运行iftop累计统计到的流量，rates：平均流量
- `=>` 发送流量，`<=`接收流量
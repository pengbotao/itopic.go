```
{
    "url": "redis-memory-check",
    "time": "2019/08/28 15:29",
    "tag": "Redis"
}
```

# 存在大KEY

背景说明：
发现线上Redis存储占用异常，某些机器内存资源异常，查看内存信息`mem_fragmentation_ratio`为1.03

![](/static/uploads/codis-memory.jpg)

根据经验怀疑是某些较大的KEY占用了内存而没有释放。需要找出REDIS中占用内存资源较大的KEY。



- 工具：rdb
- 说明：rdbtools 是解析Redis rdb文件、分析其内存，导出其数据等比較好的工具，用python编写。
- 安装：`pip install rdbtools`
- Github: https://github.com/sripathikrishnan/redis-rdb-tools


需要将dump文件导出，生成CSV格式的内存报告，输出**内存**使用前10 （-l命令）：

`rdb -c memory -l 10 dump.rdb`

5G的rdb文件，大概执行1个小时左右（可考虑放后台执行），顺利找到大KEY，清除后正常。结果示例：

database|type|key|size_in_bytes|encoding|num_elements|len_largest_element|expiry
---|---|---|---|---|---|---|---
0|list|k1|2566030649|quicklist|58448402|43|
0|list|k2|25597691|quicklist|403218|63|
0|list|k3|30883254|quicklist|505723|61|

# 存在碎片

背景说明：Zabbix报警Codis内存不足，查看机器内存基本是Codis使用，而Dashboard上只有6G的存储，查看内存信息发现`mem_fragmentation_ratio=3.79`

```
# Memory
used_memory:7502225896
used_memory_human:6.99G
used_memory_rss:28436201472
used_memory_rss_human:26.48G
used_memory_peak:22335272544
used_memory_peak_human:20.80G
total_system_memory:33807069184
total_system_memory_human:31.49G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:3.79
mem_allocator:jemalloc-4.0.3
```

**字段说明：**

字段|说明|示例
---|---|---
**used_memory_human**|Redis分配的内存总量，即存储的所有数据占用的内存|6.99G
**used_memory_rss_human**|从系统角度,显示Redis进程占用的物理内存总量|26.48G
used_memory_peak_human|内存使用的最大值，表示used_memory峰值|20.80G
total_system_memory_human|系统总内存|31.49G
used_memory_lua_human|Lua进程使用内存|37.00K
**mem_fragmentation_ratio**|内存碎片率，等价于(used_memory_rss /used_memory)|3.79，表示有19.49G的碎片空间
mem_allocator|使用的内存分配器|jemalloc-4.0.3


**解决方法：**

Redis 4之前的方式就是重启， Redis 4支持了碎片清理功能：

1、自动清理：默认情况下自动清理碎片的参数是关闭的，可以按如下命令查看

```
127.0.0.1:6379> config get activedefrag 
1) "activedefrag"
2) "no"
```

启动自动清理内存碎片

```
127.0.0.1:6379> config set  activedefrag yes
OK
```

2、手动清理

```
127.0.0.1:6379> memory purge
OK
```
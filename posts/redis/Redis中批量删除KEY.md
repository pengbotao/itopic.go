```
{
    "url": "redis-delete-keys",
    "time": "2014/02/28 22:30",
    "tag": "Redis"
}
```

Redis server went away

查看系统日志文件时发现每天定时有该错误抛出：

> PHP Fatal error: Uncaught exception 'RedisException' with message 'Redis server went away'

抛出该问题的脚本为统计脚本，需要读取前一天数据并入库，最初以为是REDIS读取太频繁造成的，但将数据导到测试机后执行脚本发现不会出现该情况，仔细调试发现手动执行时有一行代码没有执行，若执行该行则十分缓慢。该行代码为：

```
$Redis->delete($Redis->keys($pre_key_del.'*'));
```

查看手册有相应提示：

> KEYS 的速度非常快，但在一个大的数据库中使用它仍然可能造成性能问题，如果你需要从一个数据集中查找特定的 key ，你最好还是用 Redis 的集合结构(set)来代替。

登录redis通过info查看，内存使用25G多，而KEY也有1.44亿了。。。REIDS中有大量无用而又未设置过期时间的KEY存在。设置个过期时间，举手之劳的事，还是有必要的。

```
used_memory_human:24.72G
db0:keys=144856453,expires=25357
```

通过测试机执行 keys prefix* 导致REDIS卡死，其他连接也连不上。所以定位到问题出现在keys命令上，也正如手册上说的造成性能问题。

**如何删除未用到的KEY？**

大部分KEY是有规律的，有特定前缀，需要拿到特定前缀的KEY然后删除，网上有这样的命令：

```
redis-cli -a redis-pwd -n 0 keys "preffix*" | xargs redis-cli -p 6379 -a redis-pwd -n 0 del
```

测试机执行keys "preffix-1*"时间大概40多s，这意味着redis要停40s+，而前缀是按天设置的，这样子需要操作多次，因为业务的原因，不允许这么操作，分分钟都是钱~最后想到的办法是先从测试机上把满足条件的key导到文本，前面的语句通过cat文本去拿。如：

```
redis-cli -p 6380 -a redis-pwd keys "preffix-1*" > /home/keys_redis/preffix-1
```

然后通过这些数据删掉生产环境上的key。

```
cat /home/keys_redis/preffix-1 | xargs redis-cli -a redis-pwd -n 0 del
```

删除的速度非常快，内存耗的也挺快，感觉像是有多少耗多少的。执行之后KEY的数量减少了95%+，内存也从25G降到了2G。不过有一个指数升高：`mem_fragmentation_ratio`，前后的memory对比：

```
# Memory 处理前
used_memory:26839186032
used_memory_human:25.00G
used_memory_rss:23518339072
used_memory_peak:26963439000
used_memory_peak_human:25.11G
used_memory_lua:31744
mem_fragmentation_ratio:0.88
mem_allocator:jemalloc-3.2.0

# Memory 处理后
used_memory:2399386704
used_memory_human:2.23G
used_memory_rss:4621533184
used_memory_peak:26963439000
used_memory_peak_human:25.11G
used_memory_lua:31744
mem_fragmentation_ratio:1.93
mem_allocator:jemalloc-3.2.0
```

mem_fragmentation_ratio的问题可能还需要优化下，从redis这个问题可以看到，设置cache的时候我们也需要考虑到cache的维护问题，是否该设置cache的过期时间，key的命名方式如何管理，不能只想着把数据塞进去就万事大吉了。

# 命令

通过`keys *`来删除

```
$ redis-cli --raw keys "test*" | xargs redis-cli del
```

通过`scan` 来删除

```
$ redis-cli --scan --pattern "test*" | xargs -L 1000 redis-cli del
```

示例是删除本机，如果非本机管道前后需要指定host和port。
```
{
    "url": "redis-data-migration",
    "time": "2019/12/28 13:49",
    "tag": "Redis"
}
```

# 同步工具 - RedisShake

Redis-shake是阿里云自研的开源Redis数据传输工具，基于Linux环境，支持对Redis数据进行解析（decode）、恢复（restore）、备份（dump）和同步（sync或rump），灵活高效。在不方便使用DTS的迁移场景，您可以尝试使用Redis-shake进行迁移。

```
$ ./redis-shake.linux -conf=redis-shake.conf -type=sync
```

**通过RedisShake迁移Codis：**

1、配置源信息(Codis Serer列表)

```
source.type = standalone
source.address = 127.0.0.1:6379;127.0.0.1:6380;127.0.0.1:6381
```

2、配置目的信息(Codis Proxy列表)

```
target.type = proxy
target.address = 127.0.0.2:10000;127.0.0.1:10000
```

3、设置`big_key_threshold = 1`，以及启用`filter.lua = true`。

```
big_key_threshold = 1
filter.lua = true
```

启动后会执行`bgsave`，可以看看各个节点是否有正确生成。如果内存不够bgsave失败可能会丢数据，可以尝试 `/etc/sysctl.conf`增加

```
vm.overcommit_memory = 1
```

[1] [RedisShake wiki](https://github.com/alibaba/RedisShake/wiki)


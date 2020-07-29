```
{
    "url": "redis-find-big-keys",
    "time": "2019/08/28 15:29",
    "tag": "Redis"
}
```

背景说明：发现线上REDIS存储占用异常，某些机器内存资源异常。怀疑是某些较大的KEY占用了内存而没有释放。需要找出REDIS中占用内存资源较大的KEY。

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

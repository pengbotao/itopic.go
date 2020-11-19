```
{
    "url": "mysql-data-free",
    "time": "2019/04/15 06:44",
    "tag": "Mysql"
}
```

源于线上阿里云RDS数据库报磁盘空间不足，1000G的空间，占用到800多G，需要查看是否存在较大的库或表，以及是否可以做清理操作。

# 一、获取库/表占用存储大小

`information_schema.tables`单表中有各个表的基本信息，可以看到存储大小，行数，可以通过汇总做库级和表级别的统计：

**统计数据库占用存储**

```
SELECT TABLE_SCHEMA AS `db`, 
"MB" AS `unit`,
round(sum(DATA_LENGTH)/1024/1024) AS `data`, 
round(sum(INDEX_LENGTH)/1024/1024) AS `index`, 
round(sum(DATA_FREE)/1024/1024) AS `free`,
sum(TABLE_ROWS) AS `rows`
FROM information_schema.tables
GROUP BY TABLE_SCHEMA
ORDER BY `data` DESC;
```



db|unit|data|index|free|rows
---|---|---|---|---|---
x|MB|335091|0|8|27313843
demo|MB|87044|2878|88|24664309

- data: 数据占用存储空间
- index: 索引占用存储空间
- free: 碎片空间
- rows: 行数

**统计表占用存储**

表级别的类似，只不过将分组改按表分组即可。

```
SELECT TABLE_NAME AS `table`, 
"MB" AS `unit`,
round(sum(DATA_LENGTH)/1024/1024) AS `data`, 
round(sum(INDEX_LENGTH)/1024/1024) AS `index`, 
round(sum(DATA_FREE)/1024/1024) AS `free`,
sum(TABLE_ROWS) AS `rows`
FROM information_schema.tables
WHERE TABLE_SCHEMA = 'room'
GROUP BY TABLE_NAME
ORDER BY `data` DESC;
```

# 二、释放空间

尝试用`delete`删除数据，发现存储空间并未释放。找到这些资料：

```
1、drop table table_name 立刻释放磁盘空间 ，不管是 InnoDB和MyISAM

2、truncate table table_name 立刻释放磁盘空间 ，不管是 Innodb和MyISAM 。

truncate table其实有点类似于drop table 然后create。只不过这个create table 的过程做了优化，比如表结构文件之前已经有了等等，就不需要重新再搞一把。所以速度上应该是接近drop table的速度。

3、对于delete from table_name  删除表的全部数据

对于MyISAM 会立刻释放磁盘空间 （应该是做了特别处理，也比较合理）；InnoDB 不会释放磁盘空间

4、对于delete from table_name where xxx带条件的删除

不管是innodb还是MyISAM都不会释放磁盘空间。

5、delete操作以后 使用optimize table table_name 会立刻释放磁盘空间。不管是InnoDB还是MyISAM 。

所以要想达到清理数据的目的，请delete以后执行optimize table 操作。

6、delete from表 以后虽然未释放磁盘空间，但是下次插入数据的时候，仍然可以使用这部分空间

From: https://blog.csdn.net/seven_3306/article/details/30254299
```


存储较大，担心锁库、锁表的情况，选择时间点执行，335G的存储执行`truncate`耗时4s，存储很快就释放出来了。

```
mysql> TRUNCATE TABLE x;
Query OK, 0 rows affected (4.03 sec)
```

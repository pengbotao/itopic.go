```
{
    "url": "mysql-suipian-youhua",
    "time": "2013/03/06 22:16",
    "tag": "Mysql",
    "toc" : "no"
}
```

起因：查看线上数据库中Table Information时发现有一个日志表数据大小和索引大小有915M，但实际行数只有92行。该表需要频繁插入并且会定时去删掉旧的记录。表类型为Myisam，已建立一个索引，所以应该是产生了大量碎片，使用 Optimize table 表名 优化后大小变为2.19M，少了很多， 同时可以看出该表上的索引建的多余，因为插入操作比查询操作要多很多，而且查询不多，查询的数据量也一般比较小。

借此延伸下MYSQL中Myisam、InnoDB碎片优化方式：

# Myisam清理碎片

OPTIMIZE TABLE table_name

# InnoDB碎片优化

if you frequently delete rows (or update rows with variable-length data types), you can end up with a lot of wasted space in your data file(s), similar to filesystem fragmentation.

If you’re not using the innodb_file_per_table option, the only thing you can do about it is export and import the database, a time-and-disk-intensive procedure.

But if you are using innodb_file_per_table, you can identify and reclaim this space!

Prior to 5.1.21, the free space counter is available from the table_comment column of information_schema.tables. Here is some SQL to identify tables with at least 100M (actually 97.65M) of free space:

SELECT table_schema, table_name, table_comment FROM information_schema.tables WHERE engine LIKE 'InnoDB' AND table_comment RLIKE 'InnoDB free: ([0-9]{6,}).*';

Starting with 5.1.21, this was moved to the data_free column (a much more appropriate place):

SELECT table_schema, table_name, data_free/1024/1024 AS data_free_MB FROM information_schema.tables WHERE engine LIKE 'InnoDB' AND data_free > 100*1024*1024;

You can reclaim the lost space by rebuilding the table. The best way to do this is using ‘alter table’ without actually changing anything:

ALTER TABLE foo ENGINE=InnoDB;

This is what MySQL does behind the scenes if you run ‘optimize table’ on an InnoDB table. It will result in a read lock, but not a full table lock. How long it takes is completely dependent on the amount of data in the table (but not the size of the data file). If you have a table with a high volume of deletes or updates, you may want to run this monthly, or even weekly.

# 什么是mysql碎片?怎样知道表的碎片有多大呢?

简单的说,删除数据必然会在数据文件中造成不连续的空白空间,而当插入数据时,这些空白空间则会被利用起来.于是造成了数据的存储位置不连续,以及物理存储顺序与理论上的排序顺序不同,这种是数据碎片.实际上数据碎片分为两种,一种是单行数据碎片,另一种是多行数据碎片.前者的意思就是一行数据,被分成N个片段,存储在N个位置.后者的就是多行数据并未按照逻辑上的顺序排列.当有大量的删除和插入操作时,必然会产生很多未使用的空白空间,这些空间就是多出来的额外空间.索引也是文件数据,所以也会产生索引碎片,理由同上,大概就是顺序紊乱的问题.Engine 不同,OPTIMIZE 的操作也不一样的,MyISAM 因为索引和数据是分开的,所以 OPTIMIZE 可以整理数据文件,并重排索引.

OPTIMIZE 操作会暂时锁住表,而且数据量越大,耗费的时间也越长,它毕竟不是简单查询操作.所以把 Optimize 命令放在程序中是不妥当的,不管设置的命中率多低,当访问量增大的时候,整体命中率也会上升,这样肯定会对程序的运行效率造成很大影响.比较好的方式就是做个 Script,定期检查mysql中 `information_schema`.`TABLES`字段,查看 DATA_FREE 字段,大于0话,就表示有碎片.脚本多长时间运行一次,可以根据实际情况来定,比如每周跑一次.
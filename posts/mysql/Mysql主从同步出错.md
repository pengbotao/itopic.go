```
{
    "url": "mysql-master-slave-error",
    "time": "2014/08/07 20:13",
    "tag": "Mysql",
    "toc" : "no"
}
```

Mysql主从运行有一段时间了，没有出过什么问题。但最近接着出了两次问题，记录下方便后面排查！

# Slave_IO_Running和Slave_SQL_Running均为YES，主从同步出错

首先还是确认下各服务器状态。查看主库状态正常，binlog position一直在变，进程状态也正常。
```
mysql> show master status;
+------------------+-----------+--------------+------------------+
| File             | Position  | Binlog_Do_DB | Binlog_Ignore_DB |
+------------------+-----------+--------------+------------------+
| mysql-bin.000364 | 232554068 |              |                  |
+------------------+-----------+--------------+------------------+
 
mysql> show processlist;
+-------------+----------+-----------------------------------------------------------------------------+
| Command     | Time     | State                                                                       |
+-------------+----------+-----------------------------------------------------------------------------+
| Connect     | 14536445 | Slave has read all relay log; waiting for the slave I/O thread to update it |
| Binlog Dump |    22459 | Master has sent all binlog to slave; waiting for binlog to be updated       |
+-------------+----------+-----------------------------------------------------------------------------+
```
查看重库状态，整体上看重库只是有延迟。
```
mysql> show slave status\G;
 
Master_Log_File: mysql-bin.000364
Read_Master_Log_Pos: 246924389
Relay_Log_File: mysql-relay-bin.3831269
Relay_Log_Pos: 244389572
Relay_Master_Log_File: mysql-bin.000363
Slave_IO_Running: Yes
Slave_SQL_Running: Yes
Seconds_Behind_Master: 23423
 
mysql> show processlist;
+---------+-------+-----------------------------------------------------------------------------+------------------+
| Command | Time  | State                                                                       | Info             |
+---------+-------+-----------------------------------------------------------------------------+------------------+
| Connect | 22800 | Waiting for master to send event                                            | NULL             |
| Connect |    99 | Slave has read all relay log; waiting for the slave I/O thread to update it | NULL             |
+---------+-------+-----------------------------------------------------------------------------+------------------+
```
但等一段时间查看重库却一直不更新，重启后Seconds_Behind_Master为0，Slave_IO_Running和Slave_SQL_Running状态均为YES。确认了Master_Host、Master_User等参数，也匹配了Master_Server_Id都是正常的。在网上也查到了SQL_SLAVE_SKIP_COUNTER来跳过一步操作，但因为对数据完整性要求比较高，担心产生数据异常而不敢操作。于是到此基本上就没辙了。

等一天还找不到就打算重做了，但重做也不是办法，总得找到问题，数据比较多也不可能每次去重做。之前查看过Binlog没有明显发现，于是还是得再去查看下Binlog看能不能发现什么？
```
mysqlbinlog mysql-relay-bin.3831269 --start-position=244389572 --stop-position=246924461 | more
mysqlbinlog mysql-relay-bin.3831269 --start-datetime="2014-08-07 21:30:00" --stop-datetime="2014-08-07 21:35:00" --base64-output=decode-rows -v | more
```
binlog基于行的复制带上了--base64-output=decode-rows -v参数。

慢慢的还真的发现了点东西，发现有执行很多的删除语句，当通过wc统计时发现竟然有70多万。在通过业务查看是有执行一条SQL，删除表中的所有记录，数据太多，此时查看主从这个表的记录，主库为空，重库记录全在，那可能就是这个原因导致的。该操作可以跳过，于是尝试跳过之：
```
mysql>slave stop;
mysql>SET GLOBAL SQL_SLAVE_SKIP_COUNTER = 1;
mysql>slave start;
```
跳过后Mysql恢复正常，最后手动清空重库中该表的数据。至于为什么这个大的删除导致重库停止，还有待深究。

# max_allowed_packet限制导致主从同步出错

产生的原因也是执行了一个较大的更新，往数据库中更新几十兆的数据（可见更新的不合理），导致主从同步出错，查看重库状态显示

> Last_IO_Error: Got fatal error 1236 from master when reading data from binary log:'log event entry exceeded max_allowed_packet; Increase max_allowed_packet on master' 

有明显的错误描述好查很多，描述上说增大主库的max_allowed_packet。**max_allowed_packet**

> mysql 服务是通过网络包来传输数据的(通信信息包是指发送至MySQL服务器的单个SQL语句或发送至客户端的单一行),mysql协议能够识别的数据包的大小是由max_allowed_packet控制的。当MySQL客户端或mysqld服务器收到大于max_allowed_packet字节的信息包时,将发出“log event entry exceeded max_allowed_packet;”错误,并关闭连接。就像此次主从复制遇到的，IO 进程从主库获取日志，但是单个日志中的sql 大小超过了max_allowed_packet的限制，于是报错,IO thread 进程停止，sql thread 显示为yes。 对于客户端,如果通信信息包过大,在执行查询期间,可能会遇到“丢失与MySQL服务器的连接”错误。

> 参考：http://blog.itpub.net/22664653/viewspace-752580/ 

停止重库，主从都调整下，然后启动重库即可！
```
stop slave;
set global max_allowed_packet=1035543552;
start slave;
```
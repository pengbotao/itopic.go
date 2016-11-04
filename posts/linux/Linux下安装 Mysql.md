url: linux-mysql-install
des: 
time: 2014/08/14 23:55
category: linux
++++++++

# 1、通过官网下载mysql源码包。

访问http://dev.mysql.com/downloads/ 点击MySQL Community Server，选择Source Code， 点击 Generic Linux (Architecture Independent), Compressed TAR Archive后的Download
```
# wget http://dev.mysql.com/get/Downloads/MySQL-5.6/mysql-5.6.20.tar.gz
# tar zxvf mysql-5.6.20.tar.gz 
# cd mysql-5.6.20
```
# 2、 安装cmake
mysql5.5以后源码安装都得通过cmake编译，并安装了ncurses ncurses-devel
```
# yum -y install cmake ncurses ncurses-devel
# groupadd mysql
# useradd -g mysql mysql
```
# 3、编译并安装
```
# cmake . -DCMAKE_INSTALL_PREFIX=/usr/local/webserver/mysql -DMYSQL_DATADIR=/usr/local/webserver/mysql -DSYSCONFDIR=/usr/local/webserver/mysql -DDEFAULT_CHARSET=utf8 -DDEFAULT_COLLATION=utf8_general_ci  -DEXTRA_CHARSETS=all -DENABLED_LOCAL_INFILE=1
# make && make install
```
**参数说明：**

- -DCMAKE_INSTALL_PREFIX=/usr/local/webserver/mysql //指定安装目录
- -DINSTALL_DATADIR=/usr/local/webserver/mysql //指定数据存放目录
- -DSYSCONFDIR=/usr/local/webserver/mysql //指定配置文件目录（本例的配置文件为/opt/mysql/my.cnf）
- -DDEFAULT_CHARSET=utf8 //指定字符集
- -DDEFAULT_COLLATION=utf8_general_ci //指定校验字符
- -DEXTRA_CHARSETS=all //安装所有扩展字符集
- -DENABLED_LOCAL_INFILE=1 //允许从本地导入数据

编译出错需删掉CMakeCache.txt
```
# rm CMakeCache.txt
```
# 4、配置my.cnf
拷贝mysql配置文件，并进行相应配置，这里是服务器是阿里云的最低配置，单核 512M内存。
```
# cd /usr/local/webserver/mysql
# chown -R mysql:mysql data/
# cp support-files/my-default.cnf  my.cnf
# vi my.cnf 
```
**编辑my.cnf**
```
[mysqld]
 
innodb_buffer_pool_size = 100M
 
basedir = /usr/local/webserver/mysql
datadir = /usr/local/webserver/mysql/data
port = 3306
server_id = 1
socket = /tmp/mysql.sock
 
join_buffer_size = 10M
sort_buffer_size = 10M
read_rnd_buffer_size = 12M 
 
query_cache_size = 32M
tmp_table_size = 32M
key_buffer_size = 32M
 
performance_schema_max_table_instances=1000
table_definition_cache=800
table_open_cache=512
 
long_query_time=1
slow_query_log=1
slow_query_log_file=/usr/loca/webserver/mysql/data/slow-queries.log
log_queries_not_using_indexes=1
 
sql_mode=NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES 
```
# 5、初始化Mysql数据库
```
/usr/loca/webserver/mysql/scripts/mysql_install_db --user=mysql
```
# 6、启动Mysql
```
# ./support-files/mysql.server start
```
报错
```
Starting MySQL. ERROR! The server quit without updating PID file (/usr/local/webserver/mysql/data/AY121218115148c506503.pid).
 
2014-08-14 11:29:38 1678 [Note] InnoDB: Using mutexes to ref count buffer pool pages
2014-08-14 11:29:38 1678 [Note] InnoDB: The InnoDB memory heap is disabled
2014-08-14 11:29:38 1678 [Note] InnoDB: Mutexes and rw_locks use InnoDB's own implementation
2014-08-14 11:29:38 1678 [Note] InnoDB: Memory barrier is not used
2014-08-14 11:29:38 1678 [Note] InnoDB: Compressed tables use zlib 1.2.3
2014-08-14 11:29:38 1678 [Note] InnoDB: Not using CPU crc32 instructions
2014-08-14 11:29:38 1678 [Note] InnoDB: Initializing buffer pool, size = 100.0M
InnoDB: mmap(106840064 bytes) failed; errno 12
2014-08-14 11:29:38 1678 [ERROR] InnoDB: Cannot allocate memory for the buffer pool
2014-08-14 11:29:38 1678 [ERROR] Plugin 'InnoDB' init function returned error.
2014-08-14 11:29:38 1678 [ERROR] Plugin 'InnoDB' registration as a STORAGE ENGINE failed.
2014-08-14 11:29:38 1678 [ERROR] Unknown/unsupported storage engine: InnoDB
2014-08-14 11:29:38 1678 [ERROR] Aborting
```
无法给innodb_buffer_pool_size分配100M内存，但启动Mysql之前实际上是有内存的。
Mysql5.6有几个默认值，按照这些值启动需要消耗几百兆内存，然后再分配给innodb_buffer_pool_size就不足了，服务器上可怜的512M内存。。。
```
performance_schema_max_table_instances = 12500
table_definition_cache = 1400
table_open_cache = 2000
```
调整一下
```
performance_schema_max_table_instances=600
table_definition_cache=400
table_open_cache=256
```
就只使用40---60M左右的内存了,重新启动mysql
```
# ./support-files/mysql.server start
Starting MySQL. SUCCESS! 
 
# cp ./support-files/mysql.server /etc/rc.d/init.d/mysqld
# chmod 755 /etc/init.d/mysqld 
# chkconfig mysqld on
```
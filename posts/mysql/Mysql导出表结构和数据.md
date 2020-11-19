```
{
    "url": "mysqldump",
    "time": "2013/05/15 18:15",
    "tag": "Mysql"
}
```

# 导出数据库

```
-- 导出dbname表结构
mysqldump -uroot -p123456 -d dbname > dbname.sql
 
-- 导出dbname表数据
mysqldump -uroot -p123456 -t dbname > dbname.sql
 
-- 导出dbname表结构和数据
mysqldump -uroot -p123456 dbname > dbname.sql
```



```
mysqldump --opt --default-character-set=utf8 --hex-blob test --skip-triggers --skip-lock-tables > /data/peng/test.sql
```



# 导出数据库中指定表

```
-- 导出dbname下的test表结构
mysqldump -uroot -p123456 -d dbname test > test.sql
 
-- 导出dbname下的test表数据
mysqldump -uroot -p123456 -t dbname test > test.sql
 
-- 导出dbname下的test表结构和数据
mysqldump -uroot -p123456 dbname test > test.sql

-- 导出dbname下的test1 test2表结构和数据
mysqldump -uroot -p123456 dbname test1 test2 > test.sql

-- 导出t_开头的表结构和数据
mysqldump -uroot -p 库名 $(mysq l -uroot -p 库名 -Bse "show tables like 't_%'") > "导出位置.sql"
```

# 还原

```
-- 创建数据库
CREATE DATABASE dbname DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
 
-- 还原数据库
mysql -uroot -p123456 dbname < /path/dbname.sql
 
-- 还原数据库
source dbname.sql
 
-- 查看表结构
desc wp_users;
show create table wp_users \G;
```

# SELECT INTO OUTFILE导出

```
SELECT * INTO OUTFILE 'temp.txt' 
FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\n' 
FROM table_name 
WHERE createtime < 1382716800;
```

# LOAD DATA 导入

```
LOAD DATA INFILE '/home/temp.txt' 
INTO TABLE table_name 
FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\n'
(product_id,uuid,mac,monitor,win_version,ip,createtime) ;
```

注： 从本地导入远程服务器需使用`LOAD DATA LOCAL INFILE`

确保local_infile打开且有权限的情况下LOAD DATA若报错：

```
message:The used command is not allowed with this MySQL version
```

1. mysql_connect指定第五个参数128测试，即`mysql_connect($host, $user, $pwd, false, CLIENT_LOCAL_FILES)`

2. PHP中mysql扩展版本（待确认）
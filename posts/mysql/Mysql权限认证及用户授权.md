```
{
    "url": "mysql-priv",
    "time": "2019/05/08 18:15",
    "tag": "Mysql",
    "toc": "yes"
}
```

# 一、概述

`Mysql`中权限控制相关的表有（Version：5.7.17）：`mysql.user`、`mysql.db`、`mysql.tables_priv`、`mysql.columns_priv`、`mysql.procs_priv`、`mysql.proxies_priv`

## 1.1 user

| 字段名      | 字段类型 | 说明            |
| ----------- | -------- | --------------- |
| Host        | char(60) | 主机            |
| User        | char(64) | 数据库名        |
| ...         |          |                 |
| PRIMARY KEY |          | (`Host`,`User`) |

- 1、账号唯一键和我们日常设计的有点不同，它是由Host+User来确定唯一主键，意味着相同的User名对不同的主机可以设置不同的权限。
- 2、User表存放的是全局权限，即对所有的库、表的通用权限。
- 3、Host可以设置指定的主机，如`192.168.0.1`；也可以设置为通配符，比如`192.168.%`代表`192.168`下所有的`IP`，如果直接用`%`则代表不限制`IP`。
- 4、因为存在通配符，用用户名去连接的时候可能存在多个主机的情况，精确匹配的会优先，可以通过查看登录后匹配到的用户进行查看：`select current_user()`

```
mysql> SELECT USER(),CURRENT_USER();
+----------------+----------------+
| USER()         | CURRENT_USER() |
+----------------+----------------+
| root@localhost | root@localhost |
+----------------+----------------+
1 row in set (0.00 sec)
```

## 1.2 db

控制数据库级别的权限。

| 字段名      | 字段类型 | 说明                 |
| ----------- | -------- | -------------------- |
| Host        | char(60) | 主机                 |
| Db          | char(64) | 数据库名             |
| User        | char(32) | 用户名               |
| ...         |          |                      |
| PRIMARY KEY |          | (`Host`,`Db`,`User`) |

## 1.3 tables_priv

控制表级别的权限。

| 字段名      | 字段类型                                                     | 说明                                                         |
| ----------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| Host        | char(60)                                                     | 主机                                                         |
| Db          | char(64)                                                     | 数据库名                                                     |
| User        | char(32)                                                     | 用户名                                                       |
| Table_name  | char(64)                                                     | 表名                                                         |
| Grantor     | char(93)                                                     | 修改该记录的用户                                             |
| Timestamp   | timestamp                                                    | 修改该记录的时间                                             |
| Table_priv  | set('Select','Insert','Update','Delete',' Create','Drop','Grant','References', 'Index','Alter','Create View','Show view','Trigger') | 表示对表的操作权限，包括 Select、Insert、Update、Delete、Create、Drop、Grant、References、Index 和 Alter 等 |
| Column_priv | set('Select','Insert','Update','References')                 | 表示对表中的列的操作权限，包括 Select、Insert、Update 和 References |
| PRIMARY KEY |                                                              | (`Host`,`Db`,`User`,`Table_name`)                            |

## 1.4 columns_priv

控制列级别的权限。

| 字段名      | 字段类型                                     | 说明                                                         |
| ----------- | -------------------------------------------- | ------------------------------------------------------------ |
| Host        | char(60)                                     | 主机                                                         |
| Db          | char(64)                                     | 数据库名                                                     |
| User        | char(32)                                     | 用户名                                                       |
| Table_name  | char(64)                                     | 表名                                                         |
| Column_name | char(64)                                     | 数据列名称，用来指定对哪些数据列具有操作权限                 |
| Timestamp   | timestamp                                    | 修改该记录的时间                                             |
| Column_priv | set('Select','Insert','Update','References') | 表示对表中的列的操作权限，包括 Select、Insert、Update 和 References |
| PRIMARY KEY |                                              | (`Host`,`Db`,`User`,`Table_name`,`Column_name`)              |

## 1.5 procs_priv

控制存储过程的权限。

| 字段名       | 字段类型                               | 说明                                                         |
| ------------ | -------------------------------------- | ------------------------------------------------------------ |
| Host         | char(60)                               | 主机名                                                       |
| Db           | char(64)                               | 数据库名                                                     |
| User         | char(32)                               | 用户名                                                       |
| Routine_name | char(64)                               | 表示存储过程或函数的名称                                     |
| Routine_type | enum('FUNCTION','PROCEDURE')           | 表示存储过程或函数的类型，Routine_type 字段有两个值，分别是 FUNCTION 和 PROCEDURE。FUNCTION 表示这是一个函数；PROCEDURE 表示这是一个 存储过程。 |
| Grantor      | char(93)                               | 插入或修改该记录的用户                                       |
| Proc_priv    | set('Execute','Alter Routine','Grant') | 表示拥有的权限，包括 Execute、Alter Routine、Grant 3种       |
| Timestamp    | timestamp                              | 表示记录更新时间                                             |
| PRIMARY KEY  |                                        | (`Host`,`Db`,`User`,`Routine_name`,`Routine_type`)           |

## 1.6 proxies_priv

| 字段名       | 字段类型   | 说明                                          |
| ------------ | ---------- | --------------------------------------------- |
| Host         | char(60)   | 主机名                                        |
| User         | char(32)   | 用户名                                        |
| Proxied_host | char(60)   |                                               |
| Proxied_user | char(32)   |                                               |
| With_grant   | tinyint(1) |                                               |
| Grantor      | char(93)   | 插入或修改该记录的用户                        |
| Timestamp    | timestamp  | 表示记录更新时间                              |
| PRIMARY KEY  |            | (`Host`,`User`,`Proxied_host`,`Proxied_user`) |

# 二、用户授权

## 2.1 创建用户

```
mysql> CREATE USER 'username'@'host' IDENTIFIED BY 'password';
```

- `username`：用户名
- `host`: 指定的主机，如果不限制主机可以使用通配符%
- `password`: 密码，可以为空

如：

```
mysql> create user 'demo'@'localhost';
mysql> create user 'test'@'localhost' identified by '123456';
```

此时，还没对用户授权，可以看到用户表里的权限都是`N`。

## 2.2 用户授权

```
mysql> GRANT privileges ON databasename.tablename TO 'username'@'host';
```

- `privileges`： 操作权限。比如，`SELECT`,`INSERT`,`UPDATE`等，如果要授予所有权限可使用`ALL`
- `databaseName`: 数据库名
- `tableName`：表名，可以用`*`代替所有表。

如，给`'test'@'localhost'`用户授予`test.test`表的查询权限，此时会在`tables_priv`表增加记录，`db`表不会插入记录。

```
mysql> grant select on test.test to 'test'@'localhost';
mysql> flush privileges;
```

如果说给授予test库所有表的查询权限，则可以设置为`test.*`，这时会在`db`表插入记录。

```
mysql> grant select on test.* to 'test'@'localhost';
```

**也可以对字段授权：**

```
mysql> grant select(id,tag_name,tag_val), update(tag_val) on test.test to 'test'@'localhost';
```

对`test.test`表中的`id,tag_name,tag_val`授予`select`权限，`tag_val`授予更新权限，如果此时查询`select *` 则会报：

> ERROR 1142 (42000): SELECT command denied to user 'test'@'localhost' for table 'test'

## 2.3 查看授权

```
mysql> show grants;
+---------------------------------------------------------------------+
| Grants for root@localhost                                           |
+---------------------------------------------------------------------+
| GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION |
| GRANT PROXY ON ''@'' TO 'root'@'localhost' WITH GRANT OPTION        |
+---------------------------------------------------------------------+
2 rows in set (0.00 sec)

mysql> show grants for current_user();
+---------------------------------------------------------------------+
| Grants for root@localhost                                           |
+---------------------------------------------------------------------+
| GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION |
| GRANT PROXY ON ''@'' TO 'root'@'localhost' WITH GRANT OPTION        |
+---------------------------------------------------------------------+
2 rows in set (0.00 sec)
```

查看指定用户：

```
mysql> show grants for 'test'@'localhost';
+-------------------------------------------------------------+
| Grants for test@localhost                                   |
+-------------------------------------------------------------+
| GRANT USAGE ON *.* TO 'test'@'localhost'                    |
| GRANT SELECT ON `test`.* TO 'test'@'localhost'              |
| GRANT SELECT, UPDATE ON `test`.`test` TO 'test'@'localhost' |
+-------------------------------------------------------------+
3 rows in set (0.00 sec)
```

## 2.4 更改密码

```
mysql> SET PASSWORD = PASSWORD("newpassword");
mysql> SET PASSWORD FOR 'username'@'host' = PASSWORD('newpassword');
```

但在`5.7`下会有个`Warning`：

> SET PASSWORD FOR <user> = PASSWORD('<plaintext_password>')' is deprecated and will be removed in a future release. Please use SET PASSWORD FOR <user> = '<plaintext_password>' instead

**所以，可以这么操作：**

```
mysql> set password for 'test'@'localhost' = 'abc213';
```

## 2.5 撤销权限

```
mysql> REVOKE privilege ON databasename.tablename FROM 'username'@'host';
```

比如：

```
mysql> revoke update on test.test from 'test'@'localhost';
```

## 2.6 删除用户

```
mysql> DROP USER 'username'@'host';
```

比如：

```
mysql> drop user 'test'@'localhost';
```

## 2.7 刷新权限

```
mysql> flush privileges;
```

# 三、用户管理

用户表（mysql.user）除了基本的权限之外，还可以对用户的资源进行控制。

## 3.1 资源限制

资源控制列的字段用来限制用户使用的资源，可设置的字段有：

| 字段名               | 字段类型         | 是否为空 | 默认值 | 说明                             |
| -------------------- | ---------------- | -------- | ------ | -------------------------------- |
| max_questions        | int(11) unsigned | NO       | 0      | 规定每小时允许执行查询的操作次数 |
| max_updates          | int(11) unsigned | NO       | 0      | 规定每小时允许执行更新的操作次数 |
| max_connections      | int(11) unsigned | NO       | 0      | 规定每小时允许执行的连接操作次数 |
| max_user_connections | int(11) unsigned | NO       | 0      | 规定允许同时建立的连接次数       |


以上字段的默认值为 0，表示没有限制。一个小时内用户查询或者连接数量超过资源控制限制，用户将被锁定，直到下一个小时才可以在此执行对应的操作。可以使用`GRANT`语句更新这些字段的值。

## 3.2 安全列

安全列主要用来判断用户是否能够登录成功，可设置的字段有：

| 字段名                | 字段类型                          | 是否为空 | 默认值                | 说明                                                         |
| --------------------- | --------------------------------- | -------- | --------------------- | ------------------------------------------------------------ |
| ssl_type              | enum('','ANY','X509','SPECIFIED') | NO       |                       | 支持ssl标准加密安全字段                                      |
| ssl_cipher            | blob                              | NO       |                       | 支持ssl标准加密安全字段                                      |
| x509_issuer           | blob                              | NO       |                       | 支持x509标准字段                                             |
| x509_subject          | blob                              | NO       |                       | 支持x509标准字段                                             |
| plugin                | char(64)                          | NO       | mysql_native_password | 引入plugins以进行用户连接时的密码验证，plugin创建外部/代理用户 |
| password_expired      | enum('N','Y')                     | NO       | N                     | 密码是否过期 (N 未过期，y 已过期)                            |
| password_last_changed | timestamp                         | YES      |                       | 记录密码最近修改的时间                                       |
| password_lifetime     | smallint(5) unsigned              | YES      |                       | 设置密码的有效时间，单位为天数                               |
| account_locked        | enum('N','Y')                     | NO       | N                     | 用户是否被锁定（Y 锁定，N 未锁定）                           |

> 注意：即使 password_expired 为“Y”，用户也可以使用密码登录 MySQL，但是不允许做任何操作。

通常标准的发行版不支持`ssl`，读者可以使用`SHOW VARIABLES LIKE "have_openssl"`语句来查看是否具有`ssl`功能。如果 `have_openssl`的值为`DISABLED`，那么则不支持`ssl`加密功能。

# 四、权限列表

可设置的权限参考：

```
mysql> show privileges;
+-------------------------+---------------------------------------+-------------------------------------------------------+
| Privilege               | Context                               | Comment                                               |
+-------------------------+---------------------------------------+-------------------------------------------------------+
| Alter                   | Tables                                | To alter the table                                    |
| Alter routine           | Functions,Procedures                  | To alter or drop stored functions/procedures          |
| Create                  | Databases,Tables,Indexes              | To create new databases and tables                    |
| Create routine          | Databases                             | To use CREATE FUNCTION/PROCEDURE                      |
| Create temporary tables | Databases                             | To use CREATE TEMPORARY TABLE                         |
| Create view             | Tables                                | To create new views                                   |
| Create user             | Server Admin                          | To create new users                                   |
| Delete                  | Tables                                | To delete existing rows                               |
| Drop                    | Databases,Tables                      | To drop databases, tables, and views                  |
| Event                   | Server Admin                          | To create, alter, drop and execute events             |
| Execute                 | Functions,Procedures                  | To execute stored routines                            |
| File                    | File access on server                 | To read and write files on the server                 |
| Grant option            | Databases,Tables,Functions,Procedures | To give to other users those privileges you possess   |
| Index                   | Tables                                | To create or drop indexes                             |
| Insert                  | Tables                                | To insert data into tables                            |
| Lock tables             | Databases                             | To use LOCK TABLES (together with SELECT privilege)   |
| Process                 | Server Admin                          | To view the plain text of currently executing queries |
| Proxy                   | Server Admin                          | To make proxy user possible                           |
| References              | Databases,Tables                      | To have references on tables                          |
| Reload                  | Server Admin                          | To reload or refresh tables, logs and privileges      |
| Replication client      | Server Admin                          | To ask where the slave or master servers are          |
| Replication slave       | Server Admin                          | To read binary log events from the master             |
| Select                  | Tables                                | To retrieve rows from table                           |
| Show databases          | Server Admin                          | To see all databases with SHOW DATABASES              |
| Show view               | Tables                                | To see views with SHOW CREATE VIEW                    |
| Shutdown                | Server Admin                          | To shut down the server                               |
| Super                   | Server Admin                          | To use KILL thread, SET GLOBAL, CHANGE MASTER, etc.   |
| Trigger                 | Tables                                | To use triggers                                       |
| Create tablespace       | Server Admin                          | To create/alter/drop tablespaces                      |
| Update                  | Tables                                | To update existing rows                               |
| Usage                   | Server Admin                          | No privileges - allow connect only                    |
+-------------------------+---------------------------------------+-------------------------------------------------------+
31 rows in set (0.00 sec)
```



---

- [1] [MySQL用户管理](http://c.biancheng.net/mysql/100/)
- [2] [Privileges Provided by MySQL](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html#priv_all)
- [3] [mysql-system-schema](https://github.com/xiaoboluo768/mysql-system-schema)
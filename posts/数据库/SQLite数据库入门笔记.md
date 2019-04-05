```
{
    "url": "sqlite",
    "time": "2017/03/03 22:22",
    "tag": "数据库,SQLite"
}
```

# 一、简介

`SQLite`是一个零配置的关系型数据库，支持大部分`SQL`语句。就像Linux系统下会自带Python一样，通常也会安装`Sqlite3`，可以通过`sqlite3`命令来确认是否已经安装：

```
$ sqlite3
SQLite version 3.24.0 2018-06-04 14:10:15
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite>
```

虽然`SQLite`支持的功能大部分`Mysql`都有，但对比`SQLite`就会发现，`SQLite`小巧、零配置、移植方便、不需要额外启动服务端进程、功能也相当完善，较擅长在一些独立项目上提供本地存储，比纯文本方式方便，比`Mysql`清爽。

安装上可直接从[官网](https://www.sqlite.org/download.html)上下载，相关文档可从[SQLite TuTorial](http://www.sqlitetutorial.net/)上查看。操作工具可以直接使用`命令行`或者[SQLite Studio](https://sqlitestudio.pl/index.rvt?act=download)或者`Navicat`。

# 二、建表操作

本章节围绕`CREATE TABLE`来进行展开说明。

```
CREATE TABLE [IF NOT EXISTS] [schema_name].table_name (
 column_1 data_type PRIMARY KEY,
 column_2 data_type NOT NULL,
 column_3 data_type DEFAULT 0,
 table_constraint
) [WITHOUT ROWID];
```

## 2.1 语句简介

- `CREATE TABLE [IF NOT EXISTS]`：可以通过 `IF NOT EXISTS`建表，表不存在是创建，存在时忽略。
- `schema_name`：指定数据库。
- `table_name`: 表前缀不可以为`sqlite_`，该前缀仅限内部使用。报错示例：`Error: object name reserved for internal use: sqlite_test`
- `column_1 data_type`指定字段名和字段类型。
- 约束：可以指定`PRIMARY KEY`，`UNIQUE`，`NOT NULL`和`CHECK`约束，可以指定字段上指定，也有一些可以在表上指定。
- `WITHOUT ROWID`：默认情况下，`SQLite`中的每一行都有一个特殊的列，通常称为`rowid`，它唯一地标识表中的那一行。 但是，如果在CREATE TABLE语句的末尾添加了短语`WITHOUT ROWID`，则省略特殊的`rowid`列。 省略rowid有时候有空间和性能优势。 WITHOUT ROWID表是使用聚簇索引作为主键的表。
- `SQLite`不支持`COMMENT`语句，建表时可以使用 `--` 来表示注释。

如：

```
CREATE TABLE IF NOT EXISTS article (
    article_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title text NOT NULL, -- 标题
    content text, -- 内容
    status INTEGER NOT NULL DEFAULT 1
);
```

## 2.2 数据库

```
$ sqlite3 demo.db
SQLite version 3.24.0 2018-06-04 14:10:15
Enter ".help" for usage hints.
sqlite>
```
这样子会在当前目录创建`demo.db`文件，后续在命令行里建表、插入等操作会记录到该文件，也可以先直接输入`sqlite3`操作会记录到内存中，然后调用`.save`方法保存到磁盘。

```
$ sqlite3
SQLite version 3.24.0 2018-06-04 14:10:15
Enter ".help" for usage hints.
Connected to a transient in-memory database.
sqlite> CREATE TABLE IF NOT EXISTS article (
   ...>     article_id INTEGER PRIMARY KEY AUTOINCREMENT,
   ...>     title text NOT NULL,
   ...>     content text,
   ...>     status INTEGER NOT NULL DEFAULT 1
   ...> );
sqlite> .save demo.db
```

通常情况下一个`database`一个文件，有时候也会碰到跨库查询的需求。我们创建一个新的数据库

```
$ sqlite3 test.db
SQLite version 3.24.0 2018-06-04 14:10:15
Enter ".help" for usage hints.
sqlite> create table comment(artice_id integer, content text);
sqlite> .exit
```

然后使用`attach`语句将其他数据库附加到当前数据库连接。将两个数据库加载到同一个程序中，这样子建表可以指定`schema_name`，也可以进行关联查询。

attach语法：`ATTACH DATABASE file_name AS database_name;`

```
$ sqlite3
SQLite version 3.24.0 2018-06-04 14:10:15
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite> attach "./test.db" as test;
sqlite> attach "./demo.db" as demo;
sqlite> .database
main:
test: /Users/peng/workspace/demo/./test.db
demo: /Users/peng/workspace/demo/./demo.db
```

## 2.3 数据类型

SQLite中数据库类型比较简单，只有以下几种：

类型|说明
---|---
NULL|值是一个 NULL 值。
INTEGER|值是一个带符号的整数，根据值的大小存储在 1、2、3、4、6 或 8 字节中。
REAL|值是一个浮点值，存储为 8 字节的 IEEE 浮点数字。
TEXT|值是一个文本字符串，使用数据库编码（UTF-8、UTF-16BE 或 UTF-16LE）存储。
BLOB|值是一个 blob 数据，完全根据它的输入存储。


根据值的格式，SQLite根据以下规则确定其数据类型：

- 如果文字没有封闭引号和小数点或指数，则SQLite会分配INTEGER存储类。
- 如果文字由单引号或双引号括起，则SQLite会分配TEXT存储类。
- 如果文字没有引号，小数点也没有指数，SQLite会分配REAL存储类。
- 如果文字是NULL而没有引号，则它分配了NULL存储类。
- 如果文字具有X'ABCD'或x'abcd'，则SQLite分配BLOB存储类。

关于数据类型的获取可根据`typeof()`函数来获取。

```
sqlite> insert into demo.article (title, content) values ("Title", 123);
sqlite> select typeof(article_id), typeof(title), typeof(content), typeof(status) from article;
integer|text|text|integer
sqlite> create table demo.category(category_id, category_name);
sqlite> insert into demo.category values (1, 1.0);
sqlite> insert into demo.category values ('A', NULL);
sqlite> select typeof(category_id), typeof(category_name) from demo.category;
integer|real
text|null
```

## 2.3 日期类型及函数

接上一节可以看到没有日期相关类型。`SQLite`不支持内置的日期和时间存储类。 但是可以使用`TEXT`，`INT`或`REAL`来存储日期和时间值。

**1. 使用TEXT存储SQLite日期和时间**

如果使用TEXT存储类来存储日期和时间值，则需要使用ISO8601字符串格式，如：`YYYY-MM-DD HH:MM:SS.SSS`。然后使用`DATETIME`函数来获取当前时间。

```
DATETIME('now');
DATETIME('now','localtime');
```

**2. 使用INTEGER存储SQLite日期和时间**

我们通常使用INTEGER来存储UNIX时间，从1970-01-01 00:00:00到当前的秒数。可以使用`strftime('%s', 'now')`,然后读取格式化使用：`datetime(d1, 'unixepoch')`


**3. 使用REAL存储SQLite日期和时间**

也可以使用REAL存储类将日期和/或时间值存储为Julian日数，这是自公元前4714年11月24日格林威治中午以来的天数。 基于公历`Gregorian`利历。

存储：`julianday('now')`，读取`date(d1)`, `time(d1)`

```
$ sqlite3
SQLite version 3.24.0 2018-06-04 14:10:15
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite> create table time_test(text_time text, int_time integer, real_time real);
sqlite> insert into time_test values (datetime('now', 'localtime'), strftime('%s', 'now'), julianday('now'));
sqlite> .header on
sqlite> .mode column
sqlite> select * from time_test;
text_time            int_time    real_time
-------------------  ----------  ----------------
2019-04-04 23:32:48  1554391968  2458578.14778737
sqlite> select text_time, datetime(int_time, 'unixepoch'),date(real_time), time(real_time) from time_test;
text_time            datetime(int_time, 'unixepoch')  date(real_time)  time(real_time)
-------------------  -------------------------------  ---------------  ---------------
2019-04-04 23:32:48  2019-04-04 15:32:48              2019-04-04       15:32:48
```

`SQLite`中关于时间的函数主要有5个：

函数|定义|说明
---|---|---
DATE|date(timestring, modifier, modifier,...)|以 YYYY-MM-DD 格式返回日期。
TIME|time(timestring, modifier, modifier, ...)|以 HH:MM:SS 格式返回时间。
DATETIME|datetime(timestring, modifier, modifier, ...)|以 YYYY-MM-DD HH:MM:SS 格式返回。
JULIANDAY|julianday(timestring, modifier, modifier, ...)|这将返回从格林尼治时间的公元前 4714 年 11 月 24 日正午算起的天数。
STRFTIME|strftime(format_string, timestring, modifier, modifier, ...)|根据指定的格式字符串格式化日期值。

**DATE函数**

在此语法中，每个修饰符用于将变换应用于其左侧的时间值。 修改器从左到右应用，顺序很重要。例如，以下语句返回该月的最后一天：

```
SELECT DATE('now', 'start of month', '+1 month', '-1 day');
```

在这个例子中：`now`是一个指定当前日期的时间字符串。`start of month`,`+1 month`, `-1 day`是修饰符。执行步骤如下：

- 首先，将月份开始应用于由现在时间字符串指定的当前日期，因此结果是当前月份的第一天。
- 其次，+1个月适用于当月的第一天，导致下个月的第一天。
- 第三，-1天应用于下个月的第一天，这导致前一个月的最后一天。

`timestring`支持常用时间格式（用`now`表示当前时间），`modifier`格式支持：

序号|修饰符|描述
---|---|---
1|N days|`± N days` 加减N天
2|N hours|加减N小时
3|N minutes|加减N分钟
4|N.N seconds|加减N秒
5|N months|加减N月
6|N years|加减N年
7|start of month|月初
8|start of year|年初
9|start of day|当天0点
10|weekday N|将日期提前到工作日编号为N的下一个日期
11|unixepoch|Unix时间
12|localtime|本地时间
13|utc|UTC时间

如：

```
sqlite> select date("2020-01-01", "-1 day") as day;
day
----------
2019-12-31
```

理解DATE函数后，后面的函数就比较好理解了。

**TIME函数示例**

```
sqlite> select time("12:00:00", '-2 hours');
time("12:00:00", '-2 hours')
----------------------------
10:00:00
```

**DATETIME函数示例**

```
sqlite> SELECT datetime('now','localtime');
datetime('now','localtime')
---------------------------
2019-04-05 00:07:51
```

**STRFTIME函数**

```
strftime(format_string, timestring, modifier, modifier, ...)
```

格式化时间，和前面几个函数可以相互转换：

函数|等价于 strftime()|示例
---|---|---
date(...)|strftime('%Y-%m-%d', ...)|date('now') = strftime('%Y-%m-%d', 'now')
time(...)|strftime('%H:%M:%S', ...)|
datetime(...)|strftime('%Y-%m-%d %H:%M:%S', ...)|select strftime('%Y-%m-%d %H:%M:%S', 'now', "start of year");
julianday(...)|strftime('%J', ...)|

`format_string`的支持格式如下：

格式|描述
---|---
%d|一月中的第几天，01-31
%f|带小数部分的秒，SS.SSS
%H|小时，00-23
%j|一年中的第几天，001-366
%J|儒略日数，DDDD.DDDD
%m|月，00-12
%M|分，00-59
%s|从 1970-01-01 算起的秒数
%S|秒，00-59
%w|一周中的第几天，0-6 (0 is Sunday)
%W|一年中的第几周，01-53
%Y|年，YYYY
%%|% symbol


## 2.4 数据约束

前面`PRIMARY KEY`是定义在字段后面在，如果有多个主键会有问题，所以`PRIMARY KEY`也可以定义在`table_constraint`中。如：

```
create table article_control(
    article_id integer, 
    category_id integer, 
    stat integer,
    primary key(article_id, category_id)
);
```

**唯一建**

和主键定义方式一样：

```
sqlite> create table user_info(
    id integer PRIMARY KEY AUTOINCREMENT, 
    mobile text not null unique
);
sqlite> drop table user_info;
sqlite> create table user_info(
    id integer PRIMARY KEY AUTOINCREMENT, 
    mobile text not null,
    email text not null,
    unique(mobile, email)
);
```

和`MYSQL`一样，可以`DEFAULT NULL UNIQUE`.

**外键**

```
CREATE TABLE IF NOT EXISTS supplier_groups (
 group_id integer PRIMARY KEY,
 group_name text NOT NULL
);
 
CREATE TABLE suppliers (
 supplier_id integer PRIMARY KEY,
 supplier_name text NOT NULL,
 group_id integer NOT NULL,
        FOREIGN KEY (group_id) REFERENCES supplier_groups(group_id)
);

```

## 2.5 创建索引

```
-- 创建索引
CREATE [UNIQUE] INDEX index_name ON table_name(indexed_column);

-- 删除索引
DROP INDEX [IF EXISTS] index_name;
```

同样可以通过`explain`字段来查看索引使用情况。

```
sqlite> explain query plan select * from user_info where email = 'xxx';
QUERY PLAN
`--SCAN TABLE user_info
sqlite> explain query plan select * from user_info where mobile = 'xxx';
QUERY PLAN
`--SEARCH TABLE user_info USING INDEX idx_mobile (mobile=?)
```

## 2.6 修改操作

```
-- 修改表名
ALTER TABLE existing_table RENAME TO new_table;

-- 增加字段(不支持删除字段)
ALTER TABLE table ADD COLUMN column_definition [after column_name];

-- 删除表
DROP TABLE [IF EXISTS] [schema_name.]table_name;
```

## 2.7 视图

```
CREATE [TEMP] VIEW [IF NOT EXISTS] view_name(column-name-list)
AS 
   select-statement;

DROP VIEW [IF EXISTS] view_name;
```

SQLite也支持触发器、`CHECK Constraint`，不常用就不一一介绍了。到这里SQLite的基本建表相关的操作就差不多了。表建好后数据上的增删查改和Mysql基本一致，所以后续会介绍的相对精简一点。

# 三、数据读取

![](/static/uploads/sqlite-select-stmt.gif)

# 四、命令行

SQLite项目提供了一个名为sqlite3（或Windows上的sqlite3.exe）的简单命令行工具，允许使用SQL语句和命令与SQLite数据库进行交互。默认情况下，SQLite会话使用内存数据库，因此会话结束时所有更改都将消失。

要打开或者创建一个数据库文件，可以使用`.open FILENAME` 或者使用`sqlite3 FILENAME`(不存在时会自动创建FILENAME)。要显示所有可用命令及其用途，请使用.help命令，如下所示：

命令|描述
---|---
.backup ?DB? FILE|备份 DB 数据库（默认是 "main"）到 FILE 文件。
.bail ON\|OFF|发生错误后停止。默认为 OFF。
.databases|列出数据库的名称及其所依附的文件。
.dump ?TABLE?|以 SQL 文本格式转储数据库。如果指定了 TABLE 表，则只转储匹配 LIKE 模式的 TABLE 表。
.echo ON\|OFF|开启或关闭 echo 命令。
.exit|退出 SQLite 提示符。
.explain ON\|OFF|开启或关闭适合于 EXPLAIN 的输出模式。如果没有带参数，则为 EXPLAIN on，及开启 EXPLAIN。
.header(s) ON\|OFF|开启或关闭头部显示。
.help|显示消息。
.import FILE TABLE|导入来自 FILE 文件的数据到 TABLE 表中。
.indices ?TABLE?|显示所有索引的名称。如果指定了 TABLE 表，则只显示匹配 LIKE 模式的 TABLE 表的索引。
.load FILE ?ENTRY?|加载一个扩展库。
.log FILE\|off|开启或关闭日志。FILE 文件可以是 stderr（标准错误）/stdout（标准输出）。
.mode MODE|设置输出模式，MODE 可以是下列之一<BR>- csv 逗号分隔的值<BR>- column 左对齐的列<BR>- html HTML 的 `<table>` 代码<BR>- insert TABLE 表的 SQL 插入（insert）语句<BR>- line 每行一个值<BR>- list 由 .separator 字符串分隔的值<BR>- tabs 由 Tab 分隔的值<BR>- tcl TCL 列表元素
.nullvalue STRING|在 NULL 值的地方输出 STRING 字符串。
.output FILENAME|发送输出到 FILENAME 文件。
.output stdout|发送输出到屏幕。
.print STRING...|逐字地输出 STRING 字符串。
.prompt MAIN CONTINUE|替换标准提示符。
.quit|退出 SQLite 提示符。
.read FILENAME|执行 FILENAME 文件中的 SQL。
.schema ?TABLE?|显示 CREATE 语句。如果指定了 TABLE 表，则只显示匹配 LIKE 模式的 TABLE 表。
.separator STRING|改变输出模式和 .import 所使用的分隔符。
.show|显示各种设置的当前值。
.stats ON\|OFF|开启或关闭统计。
.tables ?PATTERN?|列出匹配 LIKE 模式的表的名称。
.timeout MS|尝试打开锁定的表 MS 毫秒。
.width NUM NUM|为 "column" 模式设置列宽度。
.timer ON\|OFF|开启或关闭 CPU 定时器。

## 4.1 调整显示

按照不同的模式进行SQL查询后的结果展示。

```
sqlite> .header on
sqlite> .mode column
sqlite> select * from test;
id
----------
1
2
sqlite> .mode insert
sqlite> select * from test;
INSERT INTO "table"(id) VALUES(1);
INSERT INTO "table"(id) VALUES(2);
```

## 4.2 查看数据库及表

```
-- 显示数据库
sqlite> .database
main: /Users/peng/workspace/demo/test.db

-- 显示表名
sqlite> .table
comment

-- 显示DDL
sqlite> .schema
CREATE TABLE comment(article_id integer, content text);

-- 显示索引
sqlite> create index idx_article_id on comment(article_id);
sqlite> .index
idx_article_id
```

## 4.3 数据导出

**通过.output将输出信息写到文件中。**

```
-- 查询结果写到test.txt中
sqlite> .output "./test.txt"
sqlite> select * from comment;

-- 回复在标准输出打印
sqlite> .output stdout
sqlite> select * from comment;
1|Test
```

要将数据库转储到文件中，可以使用`.dump`命令。 `.dump`命令将`SQLite`数据库的整个结构和数据转换为单个文本文件。默认情况下直接在标准输出显示，需要配合上面的`.output`指定输出到文件。

## 4.4 导入CSV格式

先将模式`.mode`设置为`csv`以指示命令行shell程序将输入文件解释为CSV文件，然后通过`.import`导入。

```
sqlite> .mode csv
-- 将csv文件导入到test表
sqlite> .import "./test.csv" test
```

## 4.5 导出CSV格式

```
sqlite> .headers on
sqlite> .mode csv
sqlite> .output data.csv
sqlite> select * from test;
sqlite> .exit
```
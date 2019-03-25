```
{
    "url": "python-database",
    "time": "2016/03/21 00:12",
    "tag": "Python,Python常用库"
}
```


# 一、概述

Python常用连接Mysql的库有三种：

- MySQLdb: MySQLdb is an thread-compatible interface to the popular MySQL database server that provides the Python database API.
    - mysqlclient: This is a fork of MySQLdb1.
- PyMySQL: This package contains a pure-Python MySQL client library, based on PEP 249.
- MySQL-Connector/Python: MySQL Connector/Python is a standardized database driver for Python platforms and development. 

## MySQLdb, mysqlclient and MySQL connector/Python的区别?

There are three MySQL adapters for Python that are currently maintained:

- mysqlclient - By far the fastest MySQL connector for CPython. Requires the mysql-connector-c C library to work.
- PyMySQL - Pure Python MySQL client. According to the maintainer of both mysqlclient and MyPySQL, you should use PyMySQL if:
    - You can't use libmysqlclient for some reason.
    - You want to use monkeypatched socket of gevent or eventlet.
    - You wan't to hack mysql protocol.
- mysql-connector-python - MySQL connector developed by the MySQL group at Oracle, also written entirely in Python. It's performance appears to be the worst out of the three. Also, due to some licensing issues, you can't download it from PyPI (but it's now available through conda).

**Benchmarks**

According to the following benchmarks, mysqlclient is faster (sometimes > 10x faster) than the pure Python clients.

From [`stackoverflow`](https://stackoverflow.com/questions/43102442/whats-the-difference-between-mysqldb-mysqlclient-and-mysql-connector-python).

# 二、pymysql

文档地址：https://pymysql.readthedocs.io/en/latest/

安装方法: `pip install PyMySQL`

## 2.1 执行查询
```
import pymysql

db = pymysql.connect(host="localhost", user="root", password="123456", db="test", port=3306)

cur = db.cursor()
sql = "select * from test where tag_name = %s and tag_val = %s"

try:
    cur.execute(sql, ("test", "test"))
    results = cur.fetchall()
    for row in results:
        print(row)
except Exception as e:
    raise e
finally:
    cur.close()
    db.close()
```

上面为按照索引返回，也可以按字段来返回

```
cur = db.cursor(cursor=pymysql.cursors.DictCursor)
sql = "select * from test where tag_name = %s"

cur.execute(sql, ("test",))
results = cur.fetchone()
print(results)
cur.close
# Output:
{u'created_ts': datetime.datetime(2019, 3, 24, 22, 12, 52), u'tag_name': 'test', u'tag_val': 'test1', u'id': 6}
```

代码方式也可以通过with语句来实现，结束后自动关闭游标：

```
with db.cursor(cursor=pymysql.cursors.DictCursor) as cursor:
    cursor.execute(sql, ("test", "test1"))
    results = cursor.fetchall()
    for row in results:
        print(row)
```

## 2.2 执行插入
```
sql = "insert into test (tag_name, tag_val) values (%s, %s)"

cur.execute(sql, ("test", "test"))
db.commit()
print(cur.lastrowid)
cur.close
db.close()
```

批量插入

```
sql = "insert into test (tag_name, tag_val) values (%s, %s)"

data = [
    ("test1", "test1"),
    ("test2", "test2"),
]

# 返回影响函数
affected_rows = cur.executemany(sql, data)
db.commit()
cur.close
db.close()
```

## 2.3 更新/删除操作

```
sql = "update test set tag_val = %s where tag_name = %s"

# 返回影响函数
affected_rows = cur.execute(sql, ("test", 'test'))
db.commit()
cur.close
db.close()
```

# 三、MySQLdb

- [MySQLdb User's Guide](https://mysqlclient.readthedocs.io/user_guide.html)

用法与`PyMySQL`一致，官网上最后一次发布还是在2014年1月，目前还不支持`Python3`，网上查询`Python3`中可用`PyMySQL`代理。兼容方法（未验证）：

```
import pymysql

pymysql.install_as_MySQLdb()
```

同时`MySQLdb`还Fork出来的一个分支：`mysqlclient`，增加了对`Python3`的支持。

# 四、MySQL-Connector/Python
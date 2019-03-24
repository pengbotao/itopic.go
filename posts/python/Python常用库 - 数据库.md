```
{
    "url": "python-database",
    "time": "2016/03/21 00:12",
    "tag": "Python,Python常用库"
}
```


# 一、pymysql

文档地址：https://pymysql.readthedocs.io/en/latest/

安装方法: `pip install PyMySQL`

## 1.1 执行查询
```
import pymysql

db = pymysql.connect(host="localhost", user="root", password="123456", db="test", port=3306)

cur = db.cursor()
sql = "select * from test where tag_name = %s and tag_val = %s"

try:
    cur.execute(sql, ("test", "test"))
    #results = cur.fetchone()
    results = cur.fetchall()
    for row in results:
        print(row)

except Exception as e:
    raise e
finally:
    db.close()
```

## 1.2 执行插入
```
sql = "insert into test (tag_name, tag_val) values (%s, %s)"

cur.execute(sql, ("test", "test"))
db.commit()
print(cur.lastrowid)
db.close()
```

批量插入

```
sql = "insert into test (tag_name, tag_val) values (%s, %s)"

data = [
    ("test1", "test1"),
    ("test2", "test2"),
]

cur.executemany(sql, data)
db.commit()
db.close()
```

## 1.3 更新/删除操作

```
sql = "update test set tag_val = %s where channel = %s"

cur.execute(sql, ("test", 'test'))
db.commit()
db.close()
```
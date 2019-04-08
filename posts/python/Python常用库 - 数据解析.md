```
{
    "url": "python-data-parse",
    "time": "2016/04/12 13:14",
    "tag": "Python,爬虫"
}
```

# 一、XML解析

pass

# 二、JSON编解码

主要包含4个方法，没有`s`结尾的需要传入文件对象。

- json.load()
- json.dump()
- json.loads()
- json.dumps()

## 2.1 对象编码 - 字符串解码

```
import json

data = {
    "from": "itopic",
    "name": "JSON数据",
}

json_str = json.dumps(data)
json_data = json.loads(json_str)

print(json_str, json_data)
```

## 2.2 JSON与文件

```
with open('data.json', 'w') as f:
    json.dump(data, f)

with open('data.json', 'r') as f:
    data_obj = json.load(f)
    print(data_obj)
```

# 三、CSV读取与解析

- csv.reader()
- csv.DictReader
- csv.writer()
- csv.DictWriter()

## 3.1 CSV写入

```
import csv

headers = ["Time", "From", "Title"]
rows = [
    ("2019-04-06", "web", "Title1"),
    ("2019-04-07", "android", "Title2"),
]

with open("test.csv", "w") as f:
    f_csv = csv.writer(f)
    f_csv.writerow(headers)
    f_csv.writerows(rows)

# 按字典的方式进行写入
rows = [
    {"Time": "2019-04-08", "From": "web", "Title": "Title1"},
    {"Time": "2019-04-09", "From": "android", "Title": "Title2"},
]
with open("test.csv", "w") as f:
    f_csv = csv.DictWriter(f, headers)
    f_csv.writeheader()
    f_csv.writerows(rows)
```

## 3.2 CSV读取
```
import csv
from collections import namedtuple

# 读取CSV文件，并按索引方式读取
with open("test.csv", "r") as f:
    f_csv = csv.reader(f)
    headers = next(f_csv)
    for row in f_csv:
        print(row[0], row)

# 读取CSV文件，并按字段方式读取
with open("test.csv", "r") as f:
    f_csv = csv.DictReader(f)
    for row in f_csv:
        print(row["Time"], row)

# 读取CSV文件，按对方的方式读取
with open("test.csv", "r") as f:
    f_csv = csv.reader(f)
    headers = next(f_csv)
    Row = namedtuple('Row', headers)
    for r in f_csv:
        row = Row(*r)
        print(row.Time, row)
```

# 四、HTML解析

通过`requests`获取内容，然后在用`lxml`库通过`xpath`来解析节点。下面为抓取本博客所有文章。

```
import requests
from lxml import html

URL = "http://itopic.org"

for topic in html.fromstring(requests.get(URL).text).xpath('//div[@id="left-sider" or @id="right-sider"]/ul/li'):
    print("%s %s %s%s" % (topic.text, topic.xpath('./a/text()')[0], URL, topic.xpath('./a/@href')[0]))
```

requests+lxml可参考：`http://x-wei.github.io/python_crawler_requests_lxml.html`

# 五、BASE64编解码

- base64.b64encode()
- base64.b64decode()

```
import base64

s = base64.b64encode("Base64数据")
print(s)
print(base64.b64decode(s))
```

# 六、MD5编码

```
import hashlib

s = "Hello Python."

m = hashlib.md5()
m.update(s)
print(m.hexdigest())
```
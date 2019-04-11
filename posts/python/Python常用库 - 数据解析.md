```
{
    "url": "python-data-parse",
    "time": "2016/04/12 13:14",
    "tag": "Python,爬虫"
}
```

# 一、XML解析

这里使用`ElementTree`方式解析XML

```
#! /usr/bin/env python
# coding=utf-8

import xml.etree.ElementTree as ET
xml_str = '''<SearchRequest from="tools">
    <SearchDetails Nationality = "CN">
        <ArrivalDate>2019-05-02</ArrivalDate>
        <PropertyReferenceID>135155</PropertyReferenceID>
        <Duration>1</Duration>
        <RoomRequests>
            <RoomRequest>
                <Adults>2</Adults>
                <Children>0</Children>
                <Infants>0</Infants>
            </RoomRequest>
            <RoomRequest>
                <Adults>2</Adults>
                <Children>1</Children>
                <Infants>1</Infants>
            </RoomRequest>
        </RoomRequests>
    </SearchDetails>tail
</SearchRequest>'''
```

加载数据的方式有两种，从XML文件获取:

```
#tree = ET.parse("filename.xml")
#root = tree.getroot()
```

从XML字符串获取

```
root = ET.fromstring(xml_str)
```

第一种方式中tree为`ElementTree`对象，第二种方式中root为`Element`对象。每个Element有以下属性：

- root.tag 标签名称
- root.attrib 获取属性字典，也可以通过root.get('x')来获取指定的属性。
- root.text 获取xml文本内容
- root.tail 如上面xml中的`</SearchDetails>tail`，相对用的比较少。

示例如下：

```
print("Tag: {}, Attributes: {}, Text: {}, Tail: {}".format(root.tag, root.attrib, root.text.strip(), root.tail))

for child in root:
    print("Tag: {}, Attributes: {}, Text: {}, Tail: {}".format(child.tag, child.get("Nationality"), child.text.strip(), child.tail))
    print(child.find("ArrivalDate").text)
    print(child.find("PropertyReferenceID").text)
```

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
```
{
    "url": "python-html-xml-parse",
    "time": "2016/04/12 13:14",
    "tag": "Python,Python常用库,爬虫",
    "public": "yes"
}
```

# lxml

通过`requests`获取内容，然后在用`lxml`库通过`xpath`来解析节点。下面为抓取本博客所有文章。

```
import requests
from lxml import html

URL = "http://itopic.org"

for topic in html.fromstring(requests.get(URL).text).xpath('//div[@id="left-sider" or @id="right-sider"]/ul/li'):
    print("%s %s %s%s" % (topic.text, topic.xpath('./a/text()')[0], URL, topic.xpath('./a/@href')[0]))
```

http://x-wei.github.io/python_crawler_requests_lxml.html
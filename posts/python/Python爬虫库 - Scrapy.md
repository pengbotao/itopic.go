```
{
    "url": "scrapy",
    "time": "2017/08/24 23:12",
    "tag": "Python,爬虫"
}
```

# 一、关于Scrapy

## 1.1 文档地址
```
https://doc.scrapy.org/en/latest/
```

## 1.2 安装方式
```
pip install scrapy
```

# 二、抓取示例

## 2.1 初始化

```
$ scrapy startproject demo
New Scrapy project 'demo', using template directory '/Users/peng/.pyenv/versions/3.7.2/lib/python3.7/site-packages/scrapy/templates/project', created in:
    /Users/peng/workspace/demo

You can start your first spider with:
    cd demo
    scrapy genspider example example.com
$ cd demo/
$ scrapy genspider itopic itopic.org
Created spider 'itopic' using template 'basic' in module:
  demo.spiders.itopic
```

![](/static/uploads/scrapy-project-init.png)

## 2.2 启动抓取
```
scrapy crawl itopic
```

当然在这一步执行抓取除了可以看到一些日志之外，并不会得到其他的东西。目前抓取的需求还没有明确，还需要对`itopic.org`页面返回的数据进行解析。

## 2.3 解析列表

保存首页内容，并解析文章和超链地址。在这一步之后若在次执行抓取，顺利的话就可以看到应有的内容了。

```
# -*- coding: utf-8 -*-
import scrapy
from lxml import html


class ItopicSpider(scrapy.Spider):
    name = 'itopic'
    allowed_domains = ['itopic.org']
    start_urls = ['https://itopic.org/']

    def parse(self, response):
        with open("index.html", "wb") as f:
            f.write(response.body)
        for x in html.fromstring(response.body).xpath('//div[@id="left-sider" or @id="right-sider"]/ul/li/a'):
            print(x.text, x.attrib.get('href'))
```


## 2.4 解析详情


```
# -*- coding: utf-8 -*-
import scrapy
from lxml import html


class ItopicSpider(scrapy.Spider):
    name = 'itopic'
    allowed_domains = ['itopic.org']
    start_urls = ['https://itopic.org/']

    def parse(self, response):
        for x in html.fromstring(response.body).xpath('//div[@id="left-sider" or @id="right-sider"]/ul/li/a'):
            url = "https://itopic.org" + x.attrib.get('href')
            yield scrapy.Request(url=url, callback=self.parse_itopic_detail)

    def parse_itopic_detail(self, response):
        filename = response.url.split("/")[-1]
        with open(filename, "wb") as f:
            f.write(response.body)
        print(response.css(".title::text").extract()[0])
```

## 2.5 保存结果

### 2.5.1 编写Item
```
# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class ItopicDetailItem(scrapy.Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    title = scrapy.Field()
    url = scrapy.Field()
    pass
```

### 2.5.2 提取数据到Item
```
# -*- coding: utf-8 -*-
import scrapy
from demo.items import ItopicDetailItem
from lxml import html


class ItopicSpider(scrapy.Spider):
    name = 'itopic'
    allowed_domains = ['itopic.org']
    start_urls = ['https://itopic.org/']

    def parse(self, response):
        for x in html.fromstring(response.body).xpath('//div[@id="left-sider" or @id="right-sider"]/ul/li/a'):
            url = "https://itopic.org" + x.attrib.get('href')
            yield scrapy.Request(url=url, callback=self.parse_itopic_detail)

    def parse_itopic_detail(self, response):
        filename = response.url.split("/")[-1]
        with open(filename, "wb") as f:
            f.write(response.body)
        item = ItopicDetailItem()
        item['title'] = response.css(".title::text").extract()[0]
        item['url'] = response.url
        yield item
```

### 2.5.3 保存结果

执行抓取并保存结果到`itopic.json`文件。执行完成后查看`itopic.json`就可以看到抓取到的数据了。

```
$ scrapy crawl itopic -o itopic.json

[
    {
        "title": "Python常用库 - HTTP请求",
        "url": "https://itopic.org/python-requests.html"
    },
    {
        "title": "Python常用库 - 数据库",
        "url": "https://itopic.org/python-database.html"
    },
    {
        "title": "Python常用库 - 数据解析",
        "url": "https://itopic.org/python-data-parse.html"
    }
    ...
]
```
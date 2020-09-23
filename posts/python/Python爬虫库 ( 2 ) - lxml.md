```
{
    "url": "lxml",
    "time": "2017/03/12 20:14",
    "tag": "Python,爬虫"
}
```

lxml 是一个HTML/XML的解析器，主要的功能是如何解析和提取 HTML/XML 数据。

# 一、安装

`pip install lxml`

安装失败可指定镜像：

`pip install --index https://pypi.mirrors.ustc.edu.cn/simple/ lxml`


# 二、lxml结构


- Element
    - Property
      - .attrib
      - .base
      - .sourceline
      - .tag
      - .text
      - .prefix
      - .nsmap
    - Method
      - xpath()
      - getparent()
      - getprevious()
      - getnext()
      - getchildren()
      - getroottree()
      - find()
      - findall()
      - findtext()
      - clear()
      - get()
      - items()
      - keys()
      - values()
      - set()
- etree
    - etree.fromstring()
    - etree.tostring()
    - etree.parse()
- html
    - html.fromstring()
    - html.tostring()


# 三、生成xml

```
#! /usr/local/env python
# coding: utf-8

from lxml import etree


# 创建一个节点
root = etree.Element("Service_SearchHotel", nsmap={
    'SOAP-ENV': 'http://schemas.xmlsoap.org/soap/envelope/',
})


# 往节点里添加元素
xml_request = etree.SubElement(root, "SearchHotel_Request", Language="en", Citizenship="zh")

# 添加元素并设置内容
etree.SubElement(xml_request, "HotelId").text = "10086"
etree.SubElement(xml_request, "Currency").text = "CNY"

# 也可以通过append方式设置
avail = etree.Element("Avail")
avail.text = "1"
xml_request.append(avail)

# 添加元素并设置属性
p = etree.SubElement(xml_request, "Period", Checkin="2018-08-08")
# 另一种方式设置属性
p.set("Checkout", "2018-08-09")

xml_room = etree.SubElement(xml_request, "RoomInfo")
# 还有一种添加属性的方法：attrib
attribs = {"Adult": "2", "Child": "1"}
etree.SubElement(xml_room, "Room", attrib=attribs, RoomType="Twin", RoomNum="1")
etree.SubElement(xml_room, "Room", attrib=attribs, RoomType="Double", RoomNum="2")

# 可查看生成的XML文件
x = etree.tostring(root, pretty_print=True)
print(x)

# 存储到文件
tree = etree.ElementTree(root)
tree.write('search.xml', pretty_print=True, xml_declaration=True, encoding='utf-8')
```

```
<Service_SearchHotel xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SearchHotel_Request Citizenship="zh" Language="en">
    <HotelId>10086</HotelId>
    <Currency>CNY</Currency>
    <Avail>1</Avail>
    <Period Checkin="2018-08-08" Checkout="2018-08-09"/>
    <RoomInfo>
      <Room RoomNum="1" RoomType="Twin" Adult="2" Child="1"/>
      <Room RoomNum="2" RoomType="Double" Adult="2" Child="1"/>
    </RoomInfo>
  </SearchHotel_Request>
</Service_SearchHotel>
```

# 四、解析xml

```
#! /usr/local/env python
# coding: utf-8

from lxml import etree

# 从文件加载xml文件
xml = etree.parse('search.xml')

# 从字符串加载xml
# x = etree.fromstring(etree.tostring(xml))

root = xml.getroot()
print(type(root), root.tag, root.getchildren()[0].get('Language'))
print("\n")

# 按列表方式访问root下第一个节点（SearchHotel_Request）
for x in root[0]:
    print(x.tag, x.text, x.keys(), x.items())

print("\n")

# 通过xpath遍历Room节点
for x in root.xpath('//SearchHotel_Request/RoomInfo/Room'):
    print(x.tag, x.text, x.keys(), x.items())
```

```
(<type 'lxml.etree._Element'>, 'Service_SearchHotel', 'en')


('HotelId', '10086', [], [])
('Currency', 'CNY', [], [])
('Avail', '1', [], [])
('Period', None, ['Checkin', 'Checkout'], [('Checkin', '2018-08-08'), ('Checkout', '2018-08-09')])
('RoomInfo', '\n      ', [], [])


('Room', None, ['RoomNum', 'RoomType', 'Adult', 'Child'], [('RoomNum', '1'), ('RoomType', 'Twin'), ('Adult', '2'), ('Child', '1')])
('Room', None, ['RoomNum', 'RoomType', 'Adult', 'Child'], [('RoomNum', '2'), ('RoomType', 'Double'), ('Adult', '2'), ('Child', '1')])
```

# 五、解析html

```
#! /usr/local/env python
# coding: utf-8

from lxml import html

s = '''
<div id="article">
    <ul id="article-list">
        <li class="selected">item-1</li>
        <li><a href="#url2">item-2</a></li>
        <li><a href="#url3"><img src="#img3"/>item-3</a></li>
    </ul>
</div>
'''

x = html.fromstring(s)
print(html.tostring(x))

# 打印text，过滤掉A标签内的其他非文本标签
for x in html.fromstring(s).xpath('//ul[@id="article-list"]/li'):
    print(x.text_content())


# 获取a标签的标题和超链接地址，text可调为text_content()
for x in html.fromstring(s).xpath('//a'):
    print(x.text, x.attrib.get('href'))

# 获取class为selected的li的文本内容，结果为数组
print(html.fromstring(s).xpath('//li[@class="selected"]/text()'))


# 移除元素后重新组成html
c = ''
for x in html.fromstring(s).xpath('//ul'):
    for sv in x.xpath('./li[@class="selected"]'):
        x.remove(sv)
    c += html.tostring(x, pretty_print=True, encoding='unicode')
print(c)
```

# 参考文档

- Python读写XML文档(lxml方式)(http://yshblog.com/blog/151)
- lxml简明教程(https://www.cnblogs.com/ospider/p/5911339.html)
- lxml - 用Python解析XML和HTML(https://www.jianshu.com/p/282fbf3d0a0c)

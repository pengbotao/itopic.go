```
{
    "url": "lxml",
    "time": "2017/04/12 20:14",
    "tag": "Python,爬虫"
}
```

lxml 是一个HTML/XML的解析器，主要的功能是如何解析和提取 HTML/XML 数据。

# 一、安装

`pip install lxml`

安装失败可指定镜像：

`pip install --index https://pypi.mirrors.ustc.edu.cn/simple/ lxml`


# 二、lxml结构
- etree
    - etree.fromstring
    - etree.tostring
    - etree.parse

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

# 五、解析html


# 参考文档

- Python读写XML文档(lxml方式)(http://yshblog.com/blog/151)
- lxml简明教程(https://www.cnblogs.com/ospider/p/5911339.html)
- lxml - 用Python解析XML和HTML(https://www.jianshu.com/p/282fbf3d0a0c)

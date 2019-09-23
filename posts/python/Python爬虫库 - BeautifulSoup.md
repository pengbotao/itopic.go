```
{
    "url": "beautifulsoup",
    "time": "2017/08/01 19:42",
    "tag": "Python,爬虫"
}
```

Beautiful Soup 是一个可以从HTML或XML文件中提取数据的Python库.它能够通过你喜欢的转换器实现惯用的文档导航,查找,修改文档的方式.Beautiful Soup会帮你节省数小时甚至数天的工作时间.

# 一、安装

```
pip install bs4
```

## 1.1 加载（从字符串加载）

```
html_doc = """
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; img-src data:; connect-src 'self'">
    <title>Page not found &middot; GitHub Pages</title>
  </head>
  <body>

    <div class="container">

      <h1>404</h1>
      <p><strong>File not found</strong></p>

      <p>
        The site configured at this address does not
        contain the requested file.
      </p>

      <p>
        If this is your site, make sure that the filename case matches the URL.<br>
        For root URLs (like <code>http://example.com/</code>) you must provide an
        <code>index.html</code> file.
      </p>

      <p>
        <a href="https://help.github.com/pages/">Read the full documentation</a>
        for more information about using <strong>GitHub Pages</strong>.
      </p>

      <div id="suggestions">
        <a href="https://githubstatus.com">GitHub Status</a> &mdash;
        <a href="https://twitter.com/githubstatus">@githubstatus</a>
      </div>

      <a href="/" class="logo logo-img-1x">
        <img width="32" height="32" title="" alt="" src="http://1.png">
      </a>

      <a href="/" class="logo logo-img-2x">
        <img width="32" height="32" title="" alt="" src="http://2.png">
      </a>
    </div>
  </body>
</html>
"""
```

以上`Github Pages`的404页面，使用`BeautifulSou`p解析这段代码,能够得到一个`BeautifulSoup`的对象,并能按照标准的缩进格式的结构输出:

```
from bs4 import BeautifulSoup

# Load Html字符串，并使用标准库来解析
soup = BeautifulSoup(html_doc, 'html.parser')

print(soup.prettify())
```

## 1.2 加载（从文件加载）

```
soup = BeautifulSoup(open("index.html"))
```

## 1.3 解析器列表

解析器|使用方法|优势|劣势
---|---|---|---
**Python标准库（默认）**|BeautifulSoup(markup, "html.parser")|Python的内置标准库<BR>执行速度适中<BR>文档容错能力强|Python 2.7.3 or 3.2.2前的版本中文档容错能力差
**lxml HTML 解析器**|BeautifulSoup(markup, "lxml")|速度快<BR>文档容错能力强|需要安装C语言库
**lxml XML 解析器**|BeautifulSoup(markup, ["lxml-xml"])<BR>BeautifulSoup(markup, "xml")|速度快<BR>唯一支持XML的解析器|需要安装C语言库
**html5lib**|BeautifulSoup(markup, "html5lib")|最好的容错性<BR>以浏览器的方式解析文档<BR>生成HTML5格式的文档|速度慢<BR>不依赖外部扩展

## 1.4 解析示例

**1. 获取标题**

```
print(soup.title)
# <title>Page not found · GitHub Pages</title>
print(soup.title.name)
# title
print(soup.title.string)
# Page not found · GitHub Pages
print(soup.title.parent.name)
# head
```

**2. 获取特定标签**

```
print(soup.p)
# 有多个标签时打印第一个
# <p><strong>File not found</strong></p>

print(soup.h1)
# <h1>404</h1>

print(soup.div)
print(soup.img)
```

**3. 获取属性**

```
print(soup.div['class'])
print(soup.div.attrs['class'])
# ['container']

print(soup.div.attrs)
# {'class': ['container']}
```

# 二、解析说明

`Beautiful Soup`将复杂`HTML`文档转换成一个复杂的树形结构,每个节点都是`Python`对象,所有对象可以归纳为4种: `Tag` , `NavigableString`, `BeautifulSoup`, `Comment`.

## 2.1 Tag

`Tag`对象与`XML`或`HTML`原生文档中的`tag`相同:

```
soup = BeautifulSoup('<b class="boldest">Extremely bold</b>')
tag = soup.b
type(tag)
# <class 'bs4.element.Tag'>
```

### 2.1.1. Name

每个`tag`都有自己的名字,通过`.name`来获取:

```
tag.name
# u'b'
```

如果改变了tag的name,那将影响所有通过当前Beautiful Soup对象生成的HTML文档:

```
tag.name = "blockquote"
tag
# <blockquote class="boldest">Extremely bold</blockquote>
```

### 2.1.2 Attributes

一个tag可能有很多个属性. tag `<b class="boldest">`有一个 “class” 的属性,值为 “boldest” . tag的属性的操作方法与字典相同:

```
tag['class']
# u'boldest'
```

也可以直接”点”取属性, 比如: .attrs :

```
tag.attrs
# {u'class': u'boldest'}
```

`tag`的属性可以被添加,删除或修改. 再说一次, `tag`的属性操作方法与字典一样

```
tag['class'] = 'verybold'
tag['id'] = 1
tag
# <blockquote class="verybold" id="1">Extremely bold</blockquote>

del tag['class']
del tag['id']
tag
# <blockquote>Extremely bold</blockquote>

tag['class']
# KeyError: 'class'
print(tag.get('class'))
# None
```

**1.多值属性**

`HTML4`定义了一系列可以包含多个值的属性.在`HTML5中`移除了一些,却增加更多.最常见的多值的属性是`class`(一个`tag`可以有多个CSS的class). 还有一些属性`rel`, `rev` , `accept-charset` , `headers` , `accesskey` . 在`Beautiful Soup`中多值属性的返回类型是`list`:

```
css_soup = BeautifulSoup('<p class="body strikeout"></p>')
css_soup.p['class']
# ["body", "strikeout"]

css_soup = BeautifulSoup('<p class="body"></p>')
css_soup.p['class']
# ["body"]
```

如果某个属性看起来好像有多个值,但在任何版本的HTML定义中都没有被定义为多值属性,那么`Beautiful Soup`会将这个属性作为字符串返回

```
id_soup = BeautifulSoup('<p id="my id"></p>')
id_soup.p['id']
# 'my id'
```

将tag转换成字符串时,多值属性会合并为一个值

```
rel_soup = BeautifulSoup('<p>Back to the <a rel="index">homepage</a></p>')
rel_soup.a['rel']
# ['index']
rel_soup.a['rel'] = ['index', 'contents']
print(rel_soup.p)
# <p>Back to the <a rel="index contents">homepage</a></p>
```

如果转换的文档是XML格式,那么tag中不包含多值属性

```
xml_soup = BeautifulSoup('<p class="body strikeout"></p>', 'xml')
xml_soup.p['class']
# u'body strikeout'
```

## 2.2 NavigableString
字符串常被包含在tag内.Beautiful Soup用 NavigableString 类来包装tag中的字符串:

```
tag.string
# u'Extremely bold'
type(tag.string)
# <class 'bs4.element.NavigableString'>
```

一个`NavigableString`字符串与`Python`中的`Unicode`字符串相同,并且还支持包含在`遍历文档树`和`搜索文档树`中的一些特性. 通过`unicode()`方法可以直接将 `NavigableString`对象转换成`Unicode`字符串:

```
unicode_string = unicode(tag.string)
unicode_string
# u'Extremely bold'
type(unicode_string)
# <type 'unicode'>
```

`tag`中包含的字符串不能编辑,但是可以被替换成其它的字符串,用`replace_with()`方法:

```
tag.string.replace_with("No longer bold")
tag
# <blockquote>No longer bold</blockquote>
```

`NavigableString`对象支持`遍历文档树`和`搜索文档树`中定义的大部分属性, 并非全部.尤其是,一个字符串不能包含其它内容(tag能够包含字符串或是其它tag),字符串不支持`.contents`或`.string`属性或`find()`方法.

如果想在`Beautiful Soup`之外使用`NavigableString`对象,需要调用`unicode()`方法,将该对象转换成普通的Unicode字符串,否则就算Beautiful Soup已方法已经执行结束,该对象的输出也会带有对象的引用地址.这样会浪费内存.

## 2.3 BeautifulSoup

BeautifulSoup 对象表示的是一个文档的全部内容.大部分时候,可以把它当作 Tag 对象,它支持 遍历文档树 和 搜索文档树 中描述的大部分的方法.

因为 BeautifulSoup 对象并不是真正的HTML或XML的tag,所以它没有name和attribute属性.但有时查看它的 .name 属性是很方便的,所以 BeautifulSoup 对象包含了一个值为 “[document]” 的特殊属性 .name

```
soup.name
# u'[document]'
```

## 2.4 Comment

Tag , NavigableString , BeautifulSoup 几乎覆盖了html和xml中的所有内容,但是还有一些特殊对象.容易让人担心的内容是文档的注释部分:

```
markup = "<b><!--Hey, buddy. Want to buy a used parser?--></b>"
soup = BeautifulSoup(markup)
comment = soup.b.string
type(comment)
# <class 'bs4.element.Comment'>
```

Comment 对象是一个特殊类型的 NavigableString 对象:

```
comment
# u'Hey, buddy. Want to buy a used parser'
```

但是当它出现在HTML文档中时, Comment 对象会使用特殊的格式输出:

```
print(soup.b.prettify())
# <b>
#  <!--Hey, buddy. Want to buy a used parser?-->
# </b>
```

`Beautiful Soup`中定义的其它类型都可能会出现在XML的文档中: CData , ProcessingInstruction , Declaration , Doctype .与 Comment 对象类似,这些类都是 NavigableString 的子类,只是添加了一些额外的方法的字符串独享.下面是用CDATA来替代注释的例子:

```
from bs4 import CData
cdata = CData("A CDATA block")
comment.replace_with(cdata)

print(soup.b.prettify())
# <b>
#  <![CDATA[A CDATA block]]>
# </b>
```

# 三、查找

文档地址：https://beautifulsoup.readthedocs.io/zh_CN/latest/
```
{
    "url": "markdown",
    "time": "2016/11/08 19:45",
    "tag": "Markdown"
}
```

# H1标题
Markdown 是一种轻量级标记语言，创始人为约翰·格鲁伯（`John Gruber`）。它允许人们“使用易读易写的纯文本格式编写文档，然后转换成有效的`XHTML`(或者`HTML`)文档”。这种语言吸收了很多在电子邮件中已有的纯文本标记的特性。

## H2粗体、斜体、删除线
`Go` is an open source *programming language* that ~~makes it easy~~ to build **simple**, **reliable**, and **efficient** software.

### H3引用
> `John Gruber` 在 2004 年创造了 `Markdown` 语言，在语法上有很大一部分是跟 `Aaron Swartz `共同合作的。这个语言的目的是希望大家使用“易于阅读、易于撰写的纯文字格式，并选择性的转换成有效的 `XHTML` (或是`HTML`)”。

#### H4列表

- 高效率，自动排版。
- 语法简洁优雅，轻量级，记忆负担小。
- 更专注于内容，标签对内容的侵入性低。
- 纯文本，不受编辑工具限制
- `Markdown` 转 `HTML` 非常方便。`HTML `是整个万维网（`web`）的标记语言，但更重要的是，它也是目前主流电子书格式所用的标记语言。若采用 `Markdown`，对于日后的文件转换工作也大有裨益。

##### H5超链接和图片
http://www.baidu.com  [百度](http://www.baidu.com)  [腾讯](http://www.qq.com "腾讯官网")  ![](/static/favicon.ico)

###### H6分割线
---

# 代码块
```
package main

import "fmt"

func main() {
    fmt.Println("Hello, 世界")
}
```

# 表格
站点   | URL
---    | ---
百度   | https://www.baidu.com
腾讯   | http://www.qq.com
谷歌   | https://www.google.com.hk

---

# Markdown源文件
<pre>
# H1标题
Markdown 是一种轻量级标记语言，创始人为约翰·格鲁伯（`John Gruber`）。它允许人们“使用易读易写的纯文本格式编写文档，然后转换成有效的`XHTML`(或者`HTML`)文档”。这种语言吸收了很多在电子邮件中已有的纯文本标记的特性。

## H2粗体和斜体
`Go` is an open source *programming language* that ~~makes it easy~~ to build **simple**, **reliable**, and **efficient** software.

### H3引用
> `John Gruber` 在 2004 年创造了 `Markdown` 语言，在语法上有很大一部分是跟 `Aaron Swartz `共同合作的。这个语言的目的是希望大家使用“易于阅读、易于撰写的纯文字格式，并选择性的转换成有效的 `XHTML` (或是`HTML`)”。

#### H4列表

- 高效率，自动排版。
- 语法简洁优雅，轻量级，记忆负担小。
- 更专注于内容，标签对内容的侵入性低。
- 纯文本，不受编辑工具限制
- `Markdown` 转 `HTML` 非常方便。`HTML `是整个万维网（`web`）的标记语言，但更重要的是，它也是目前主流电子书格式所用的标记语言。若采用 `Markdown`，对于日后的文件转换工作也大有裨益。

##### H5超链接和图片
http://www.baidu.com  [百度](http://www.baidu.com)  [腾讯](http://www.qq.com "腾讯官网")  ![](/static/favicon.ico)

###### H6分割线
---

# 代码块
```
package main

import "fmt"

func main() {
    fmt.Println("Hello, 世界")
}
```

# 表格
站点   | URL
---    | ---
百度   | https://www.baidu.com
腾讯   | http://www.qq.com
谷歌   | https://www.google.com.hk
</pre>
采用`Go`语言和`Markdown`实现的一个简易博客系统，主要包括以下功能：

- 按日期、按标签展现文章列表
- 首页、文章详情页
- 可生成静态页面，配合Github/Coding Pages可实现简单的博客，http://itopic.org 采用此方案。

功能比较简单，不需要依赖数据库，不需要管理后台，使用者只需要关注文章内容的书写即可，同时写好的文章可直接在`Github`上查看。

# 安装Golang
从golang.org下载并设置环境变量。

```
export GOROOT=/usr/local/server/go1.8
export GOPATH=/Users/peng/workspace/gopath:/Users/peng/workspace/golang
export PATH=$PATH:$GOROOT/bin
```

# 如何写文章？
文章采用`Markdown`的写法，需要先了解`Markdown`的写法，基本用法可查看[Markdown基本用法](/posts/Markdown基本用法.md)。除此之外有几点需要注意：

1、文章放在`posts`目录中，文件夹可多层嵌套（无影响），文件需以md为后缀，文件名即文章标题。

2、`md`文件头部需写入文章头部，文章头部和文章正文以换行区分，示例如下：

> ```
> {
>     "url": "markdown",
>     "time": "2016/11/01 19:45",
>     "tag": "Markdown"
> }
> ```
>
> 文章正文

文章头部采用`json`来描述文章信息，字段定义如下：

字段   | 必选 | 说明
---    | --- | ---
url    | 是  | 文章URL
time   | 是  |  文章发表时间
tag    | 是  | 标签，多个标签用英文逗号分隔
public | 否  | 为no的时候表示文章不可被浏览器访问到

# 安装
```
go get github.com/pengbotao/itopic.go
```

# 启动
```
go run main.go
```
说明：

- 不需要生成静态页面时将`main.go`中`isCreateHTML`设置为false
- `domain`用来定义模版页链接前缀，可设置为空
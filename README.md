采用`Go`语言和`Markdown`实现的一个简易博客系统[\[点此查看文章\]](/posts/)，主要包括以下功能：

- 按日期、按标签展现文章列表
- 首页、文章详情页
- 可生成静态页面，配合Github/Coding Pages可实现简单的博客，https://itopic.org 采用此方案。

功能比较简单，不需要依赖数据库，不需要管理后台，使用者只需要关注文章内容的书写即可，同时写好的文章可直接在`Github`上查看。

# 安装Golang
已安装环境可忽略。从golang.org下载并根据实际路径设置环境变量。

```
export GOROOT=/usr/local/server/go1.8
export GOPATH=/user/local/server/gopath:/Users/peng/workspace/golang
export PATH=$PATH:$GOROOT/bin
```

# 如何写文章？
文章采用`Markdown`的写法，需要先了解`Markdown`的写法，基本用法与文档示例写法可查看[Markdown基本用法](/posts/Markdown基本用法.md)。除此之外有几点需要注意：

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

```
  -debug
        debug mode
  -host string
        host (default "127.0.0.1:8001")
  -html
        is create html
  -prefix string
        html folder (default "../itopic.org")
```

- `-debug`: 文档调整后访问浏览器实时看到效果
- `-html`: 往`prefix`目录写静态数据

# 发布

这里以为发布到`github pages`上为例。由于只支持静态页面，所以可以通过`go`的程序指定下生成静态页面的目录，然后将静态页面提交到`github`上来。参考`https://github.com/pengbotao/pengbotao.github.io`

1. 创建`username.github.io`，将生成的静态页面提交到master分支。
2. 访问`username.github.io/index.html`是否正常。
3. 访问`username.github.io`，若出现404，
    - 需要添加README文件
    - 在仓库的`Settings` -> `Change Theme` 选择一个主题然后保存，此时会在项目中生成一个`_config.yml`文件，重新访问即可。
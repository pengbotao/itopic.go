```
{
    "url": "go-module",
    "time": "2020/04/06 16:00",
    "tag": "Golang"
}
```


# 一、GOPATH

## 1.1 关于GOPATH

使用`GO`语言开发的项目中有一个必须弄清楚的环境变量`GOPATH`，该变量指向一个目录，用来存储`GO`开发的各种项目。`GOPATH`目录下约定有3个子目录：

目录|说明
---|---
src|源文件目录
pkg|编译时生成的中间文件
bin|编译后生成的可执行文件

按这种思路，多个项目会是这种组织形：

![](/static/uploads/go-project.png)


## 1.2 GOPATH配置

这种约定多少会带来点使用起来的别扭，第三方包和自己的包混在了一起，混在一起了也无法区分引用的不同版本。好在`GOPATH`还可以设置多个项目，当有多个`GOPATH`时，`go get`安装的包默认安装在第一个目录下，包的查询会按顺序往后查找。这样子可以第一个目录默认存储依赖包，后续的目录存储自定义项目，配置如下：

```
$ go version
go version go1.10.3 darwin/amd64

vi ~/.bash_profile
export GOROOT=/usr/local/server/go1.10.3/
export GOPATH=/usr/local/server/gopath:/Users/peng/workspace/golang
export GOPROXY=https://goproxy.cn
```

由于语言及工具链强依赖`GOPATH`，所以下面示例可以看到当项目不放在`$GOPATH/src`下时引入会报错（`helper`为同目录下的包，实现了一个`ToUpper`方法）

```
package main

import (
	"fmt"
	"helper"
)

func main() {
	fmt.Println(helper.ToUpper("HelloWorld"))

}
```

从报错也可以看到查找顺序，先从`$GOROOT/src`去找`helper`包，然后在`$GOPATH/src`中依次查找，最后找不到就报错了。

```
$ go run main.go 
main.go:4:2: cannot find package "helper" in any of:
        /usr/local/server/go1.10.3/src/helper (from $GOROOT)
        /usr/local/server/gopath/src/helper (from $GOPATH)
        /Users/peng/workspace/golang/src/helper
```

按照上面方式配置一般的都没问题，但可能面临同一个包的不同版本依赖问题。比如`project1`依赖v1的版本，`project2`也用到了这个包，同时这个包升级到v2版本了，可能不兼容v1的用法，此时会发现`GOPATH`里解决不了版本号的问题。

# 二、Vendor

## 2.1 vendor简介

在`GO 1.5`版本引入了`vendor`进制，也就是如果项目目录下存在一个`vendor`目录则优先从`vendor`目录查找依赖，该进制默认关闭，需要设置环境变量

```
GO15VENDOREXPERIMENT=1
```

在`Go 1.6`之后`vendor`进制默认开启。当然前提还是项目需要存在`GOPATH`中。我们在前面的示例目录里创建`vendor`目录，再次执行可以看到查找的先后顺序发生变化，会优先从`vendor`目录查找。

```
$ go run main.go 
main.go:5:2: cannot find package "helper" in any of:
        /Users/peng/workspace/golang/src/gotest/vendor/helper (vendor tree)
        /usr/local/server/go1.10.3/src/helper (from $GOROOT)
        /usr/local/server/gopath/src/helper (from $GOPATH)
        /Users/peng/workspace/golang/src/helper
```

## 2.2 第三方包管理工具

可以看到这个目录基本可以解决我们日常开发上的问题，每个项目可以有每个项目的依赖，甚至提交`vendor`目录后可以省去各种包被墙的问题，可以适当缓解版本依赖问题。由于还没有官方的包管理工具，各种第三方依赖管理工具也百花齐放：`govendor`、`godep`、`glide`，官方也在`GO 1.9`版本推出了官方管理工具：`dep`。

随着2018年8月24号`GO1.11`的发布，`Go modules`呈现出来，进一步规范依赖包管理。目前现状来看这些工具基本都完成了历史使命，推荐使用到`Go modules`：

**godep**

```
Please use dep or another tool instead.
```

**govendor**

```
Go modules Work well now. Go modules is considered production with go1.14, but work well on 1.13 and 1.12.

Please don't use this tool anymore.

Please use Go modules.
```

**glide**

```
Go Modules

The Go community is now using Go Modules to handle dependencies. 
Please consider using that instead of Glide. Glide is now mostly unmaintained.
```

# 三、Go Modules

## 3.1 GO111MODULE环境变量

带版本号的环境变量通常是特定版本下某些实验性功能的开关，`GO1.11`中模块并没有默认打开，需要设置`GO111MODULE`环境变量，该变量在不同版本下变现有不同。

**GO111MODULE 与 Go 1.11 和 1.12**

```
即使项目在您的GOPATH中，GO111MODULE = on 仍将强制使用 Go 模块。需要 go.mod 正常工作。
GO111MODULE = off 强制 Go 表现出 GOPATH 方式，即使在 GOPATH 之外。
GO111MODULE = auto 是默认模式。在这种模式下，Go 会表现
    - 当您在 GOPATH 外部时， 设置为 GO111MODULE = on，
    - 当您位于 GOPATH 内部时，即使存在 go.mod, 设置为 GO111MODULE = off。

每当您进入 GOPATH 中，并且您希望执行的操作需要 Go 模块 (例如，Go get 特定版本的二进制文件)，您需要执行以下操作：

GO111MODULE=on go get github.com/golang/mock/tree/master/mockgen@v1.3.1
```

**Go 1.13 下的 GO111MODULE**

```
在Go 1.13下， GO111MODULE的默认行为 (auto) 改变了：

当存在 go.mod 文件时或处于 GOPATH 外， 其行为均会等同于于 GO111MODULE=on。
这意味着在 Go 1.13 下你可以将所有的代码仓库均存储在 GOPATH 下。
当处于 GOPATH 内且没有 go.mod 文件存在时其行为会等同于 GO111MODULE=off。
```

## 3.2 GO模块基本使用

使用模块之后就不强制要求项目必须在`GOPATH/src`中了，可以将前面的测试项目移到`GOPATH`之外，并执行模块初始化。

```
$ go version
go version go1.14.1 darwin/amd64

$ go mod init gotest
go: creating new go.mod: module gotest
```

init操作会创建`go.mod`文件，定义了模块名称，相当于整个项目目录为一个模块，项目内的引用需要基于该模块。`import "helper"`改为`import "gotest/helper"`后前面的示例就可以运行成功了。

```
package main

import (
	"fmt"
	"gotest/helper"
)

func main() {
	fmt.Println(helper.ToUpper("HelloWorld"))

}
```

也可以正常引入外部包，引入之后执行自动下载依赖包，下载的安装包存储在`$GOPATH/pkg/mod`下，同时`go.mod`文件也会相应更新引用及对应的版本。


```
package main

import (
	"fmt"
	"gotest/helper"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println(helper.ToUpper("HelloWorld"))
	log.Info("A walrus appears")
}
```

```
$ go run main.go 
go: finding module for package github.com/sirupsen/logrus
go: downloading github.com/sirupsen/logrus v1.5.0
go: found github.com/sirupsen/logrus in github.com/sirupsen/logrus v1.5.0
go: downloading golang.org/x/sys v0.0.0-20190422165155-953cdadca894
HELLOWORLD
INFO[0000] A walrus appears
```

## 3.3 go.mod文件


**replace**

可以通过replace指令替换一些无法下载的包；或者做一些调试等功能。

```
module gotest

go 1.14

require github.com/sirupsen/logrus v1.5.0 // indirect

replace golang.org/x/crypto => github.com/golang/crypto latest
```

## 3.4 go mod使用

命令|说明
---|---
go mod download|download modules to local cache
go mod edit|edit go.mod from tools or scripts
go mod graph|print module requirement graph
go mod init|initialize new module in current directory
go mod tidy|add missing and remove unused modules
go mod vendor|可以将依赖包缓存到vendor目录
go mod verify|verify dependencies have expected content
go mod why|explain why packages or modules are needed


# 四、GOPROXY

国内下载GO的包可能会碰到被墙的问题，可以通过设置代理解决。


名称|地址
---|---
Goproxy China|`https://goproxy.cn`
官方|`https://goproxy.io`
阿里|`https://mirrors.aliyun.com/goproxy/`

**设置方式**

```
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

---

- [1] [Go 模块解惑：到处都是 GO111MODULE ，这到底什么？](https://learnku.com/go/t/39086)
- [2] [Go 包管理的前世今生](https://www.infoq.cn/article/history-go-package-management)
- [3] [Go 模块存在的意义与解决的问题](https://zhuanlan.zhihu.com/p/86631181)
- [4] [GO 依赖管理工具go Modules（官方推荐）](https://blog.csdn.net/guyan0319/article/details/101783164)
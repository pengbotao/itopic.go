```
{
    "url": "gin",
    "time": "2022/12/01 09:00",
    "tag": "Golang",
    "toc": "yes"
}
```

# 一、概述

Gin 是一个用 Go (Golang) 编写的 HTTP Web 框架。 它具有类似 Martini 的 API，但性能比 Martini 快 40 倍。如果你需要极好的性能，使用 Gin 吧。

> From : https://gin-gonic.com/zh-cn/docs/

基本的代码片段如下：

```
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.Run("127.0.0.1:8080")
}
```

代码说明：

> 1、`gin.Default()`为默认启动方式

> 2、`r.GET()`：提交方式，可以有GET、POST、PUT、DELETE、ANY等，其中ANY会把九种提交方式都注册一下。

> 3、`r.Run()`：启动服务，默认监听0.0.0.0:8080端口。

# 二、输入处理

处理请求的方法定义为：`type HandlerFunc func(*Context)`，在包外面看到的就是`func(ctx *gin.Context)`，只要实现了`gin.HandlerFunc`就可以作为请求处理的入口函数。

## 2.1 参数获取

```
r.POST("/user/:name", func(ctx *gin.Context) {
	//获取路径中的参数:name
	name := ctx.Param("name")

	//获取GET参数 ?age=18&height=100
	age := ctx.Query("age")
	height := ctx.DefaultQuery("height", "18")

	//获取POST参数
	intro := ctx.DefaultPostForm("intro",)
	job := ctx.PostForm("job", "Nothing.")

	ctx.JSON(http.StatusOK, gin.H{
		"name":   name,
		"age":    age,
		"height": height,
		"intro":  intro,
		"job":    job,
	})
})
```

## 2.2 模型绑定

如果提交的参数比较多使用前面获取的方式则不太方便，可以通过定义一个结构体将请求参数直接映射到结构体上。比如这里定义了`tag`为`json:"name"`，则可以将提交的json内容中的`name`字段映射到`req`结构体上。其他字段类似。其他的`xml:"age"`字段用来从xml报文中解析字段，`form:"intro"`则为从表单内容中解析字段。

```
r.POST("/user", func(ctx *gin.Context) {
	var req struct {
		Name  string `form:"name" json:"name" xml:"name"`
		Age   int    `form:"age" json:"age" xml:"age"`
		Intro string `form:"intro" json:"intro" xml:"intro"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Code":    http.StatusBadRequest,
			"Message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, req)
})
```

类似的方法有：`ShouldBind`，`ShouldBindJSON`，`ShouldBindQuery`，`ShouldBindXML`，`ShouldBindYAML`等，与之对应的是一组`Bind`，`BindJSON`，`BindQuery`等，底层调用的是`MustBindWith`方法，如果绑定失败则会调用`c.Abort()`结束后面的请求。

## 2.3 模型验证

Gin使用`https://github.com/go-playground/validator`进行验证，可以通过在`tag`上加上`binding:""`来指定验证规则。

```
var req struct {
	Name  string `form:"name" json:"name" xml:"name" binding:"required,min=2,max=30"`
	Age   int    `form:"age" json:"age" xml:"age" binding:"omitempty,gte=18"`
	Intro string `form:"intro" json:"intro" xml:"intro" binding:"required"`
}
```

如果提交json内容为：`{"name": "Test", "age": 10, "intro": ""}`则会提示年龄需要大于等于18，介绍必须有值。

```
{
    "Code": 400,
    "Message": "
    	Key: 'Age' Error:Field validation for 'Age' failed on the 'gte' tag
    	Key: 'Intro' Error:Field validation for 'Intro' failed on the 'required' tag
    "
}
```

# 三、输出处理

可以输出为`JSON`、`XML`、`HTML`等格式。

## 3.1 输出为JSON

```
ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
```

## 3.2 输出为HTML

```
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "ping.html", gin.H{
			"name": "Lion",
		})
	})
	r.Run("127.0.0.1:8080")
}
```

使用HTML模板时需要先指定好模板文件的路径。

# 四、路由和中间件

## 4.1 路由

在示例中`r.Get("/ping")`后面参数是`...gin.HandlerFunc`，属于可变长参数，可传入多个`gin.HandlerFunc`，程序会按顺序往后面执行，也就是针对一个路由可以有多个`gin.HandlerFunc`处理函数。

```
func start(ctx *gin.Context) {
	fmt.Println("Start HandlerFunc")
}

func end(ctx *gin.Context) {
	fmt.Println("End HandlerFunc")
}

func main() {
	r := gin.Default()
	r.GET("/ping", start, func(ctx *gin.Context) {
		fmt.Println("ping HandlerFunc")
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	}, end)
	r.Run("127.0.0.1:8080")
}
```

## 4.2 路由组

有些时候同一类业务有相同的路由前缀，比如用户相关的接口都类似`/user/xx`，类似功能则可以通过路由组来实现。通过`router.Group("/user")`定义一个组，然后往组里添加路由即可。下面示例中的花括号`{}`只为标识一个代码块，类似`if`语句后的代码块，变量作用域受代码块影响，它与前一句`userGroup`并不直接关联，可以去掉花括号或者在中间插入代码都是可以。实际访问地址为：`/user/add`与`/order/:order_id`

```
func placeholder(ctx *gin.Context) {}

func main() {
	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/add", placeholder)
	}
	orderGroup := r.Group("/order")
	{
		orderGroup.GET("/:order_id", placeholder)
	}

	r.Run("127.0.0.1:8080")
}
```

## 4.3 中间件

我们可以在路由处理函数前后增加一些钩子函数从而做一些通用处理逻辑，这个钩子函数就叫中间件。比如：登录认证、耗时统计等；可以通过`r.Use(middleware ...gin.HandlerFunc)`来使用中间件。中间件的处理函数定义也是`gin.HandlerFunc`。默认启动方式`gin.Default()`其实是包含了两个中间件，效果如下：

```
r := gin.New()
r.Use(gin.Logger())
r.Use(gin.Recovery())
```

这里对所有路由添加一个耗时打印为例，

```
func middleware(ctx *gin.Context) {
	fmt.Println("Middleware Start")
	start := time.Now()
	ctx.Next()
	t := time.Since(start)
	fmt.Println("Middleware End", t)
}

func start(ctx *gin.Context) {
	fmt.Println("Start HandlerFunc")
}

func end(ctx *gin.Context) {
	fmt.Println("End HandlerFunc")
}

func main() {
	r := gin.Default()
	r.Use(middleware)
	r.GET("/ping", start, func(ctx *gin.Context) {
		fmt.Println("ping HandlerFunc")
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	}, end)
	r.Run("127.0.0.1:8080")
}
```

执行后程序控制台输出如下：

```
Middleware Start         -> 1. 进入中间件
Start HandlerFunc        -> 2. 执行第一个HandlerFunc
ping HandlerFunc         -> 3. 执行第二个HandlerFunc
End HandlerFunc          -> 4. 执行第三个HandlerFunc
Middleware End 117.391µs -> 5. 回到中间件处理函数
```

可以看到中间件函数和普通的处理函数没有区别，通过`ctx.Next()`将处理逻辑执行完后再回到了中间件函数，从而实现计算整个函数耗时的目的。当然并不是一定需要使用`r.Use`来使用中间件，直接加在路由处理函数`...gin.HandlerFunc`里也是可以的。

```
r.GET("/ping", middleware, func(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
})
```

同样针对路由组也可以使用中间件。

```
userGroup := r.Group("/user", middleware)

或

userGroup := r.Group("/user"）
userGroup.Use(middleware)
```

多个处理函数间的交互有几个重要函数可以了解：

- `ctx.Abort()`：会停止后续的`gin.HanderFunc`执行
- `ctx.Next()`：会往后执行`gin.HandlerFunc`，执行完毕后再回到此处
- `ctx.Set(key string, value any)`：设置后可以在后续的`gin.HandlerFunc`里通过`ctx.Get()`来获取到值
- `ctx.Get(key string) (value any, exists bool)`：跨`gin.HanderFunc`取值

所以，当一个路由可以支持多个路由处理函数后，路由处理函数和中间件其实讲的是同一个东西，可以针对所有的路由执行一段逻辑，也可以在特定路由或者路由组前后执行一段逻辑。

> 注意：在中间件或 handler 中启动新的 Goroutine 时，**不能**使用原始的上下文，必须使用只读副本（ctx.Copy()）。



---

[1] [Gin Web Framework 示例](https://gin-gonic.com/zh-cn/docs/examples/)

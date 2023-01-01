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

> 1、`gin.Default()`为默认启动方式，包含Logger、Recovery中间件。也可以使用`r := gin.New() `然后通过`r.Use(gin.Recovery())`使用中间件。

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



## 3.2 输出为HTML



# 四、路由

## 4.1 路由



## 4.2 路由组

可以通过`router.Group("/user")`来定义一个组，相当于在路由上加了一个前缀。至于后面的花括号`{}`只为标识一个代码块，类似`if`语句后的代码块，变量作用域受代码块影响，它与前一句`userGroup`并不直接关联，可以去掉花括号或者在中间插入代码都是可以。实际访问地址为：`/user/add`与`/order/:order_id`

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

# 五、中间件

```
{
    "url": "go-types",
    "time": "2020/05/12 20:00",
    "tag": "Golang",
    "toc": "yes"
}
```

# 一、概述

`Go`属于强类型静态语言，变量和类型是始终相互关联的一对。程序运行过程中变量的值可以变，但类型一经定义就不再可变了，所以变量传递过程中需确保和接受的类型一致，另外类型的运算和比较上也需要确保相同的类型，不同的类型运算和比较也需要保持一致，如：


1、变量`i`传给函数时类型非`int64`: `./main.go:12:17: cannot use i (type int) as type int64 in argument to sum`

```
package main

import "fmt"

func sum(x, y int64) int64 {
	return x + y
}

func main() {
	var i = 1
	var j int64 = 2
	fmt.Println(sum(i, j))
}
```

2、`int`和`int64`属于不同的类型，无法进行运算操作：`./main.go:10:16: invalid operation: i + j (mismatched types int and int64)` 

```
package main

import (
	"fmt"
)

func main() {
	var i int = 1
	var j int64 = 2
	fmt.Println(i + j)
}
```

3、 变量`t`和整数`1`属于不同的类型，无法进行比较操作：`./main.go:7:7: cannot use 1 (type untyped int) as type bool`


```
package main

import "fmt"

func main() {
	var t bool = true
	if t == 1 {
		fmt.Println("t is true")
	} else {
		fmt.Println("t is false")
	}
}
```

Golang中有一些基础类型：`byte` 、 `int` 、 `uint` 、 `string` 、`bool` 、 `float64`、`chan` 等，也有一些组合的高级类型：`struct` 、 `[3]int` 、 `[]int` 、`map[int]string`等。我们可以从以下几方面对类型做一下总结：

- 如何自定义类型，查看当前变量的类型
- 不同类型之间的类型转换
- 各类型是传值还是传引用


# 二、类型初始化

## 2.1 零值

当一个声明一个变量不赋初值时，`Go`语言会自动初始化值为此类型对应的零值。

- 布尔类型的零值：`false`
- 数值类型的零值：`0` 或`0.0`
- 字符串类型的零值为空字符串
- 数组类型的零值为对应多个元素的零值，如`[3]int`的零值为[0, 0, 0]
- 指针、`slice`、`map`、`channel`、`func`和`interface{}`的零值：`nil`

关于`nil`:

- nil不是关键字或保留字
- nil不可比较，比如比较两个slice会提示：`invalid operation: x == y (slice can only be compared to nil)`


## 2.2 make && new

主要通过make和new做内存分配：

1、make是用来分配并且初始化`slice`,`map`,`channel`类型的对象，返回类型为`Type`，这三种类型都是引用类型，零值为nil。

```
s := make([]int, 5, 10)
m := make(map[int]string)
c := make(chan int, 5)
```

2、new也是用来分配内存的，初始化类型的零值，返回指向这片内存地址的指针，类型为`*Type`，也可以通过 `&`来获取变量的指针。

```
x := new(int)
*x = 1
```

## 2.3 type关键字

这里存在两种用法，一种是类型定义，比如通常定义一个新的结构体。

```
type TypeName Type
```

还有一种是类型别名，考虑系统兼容重新起一个名字。

```
type TypeAlias = Type
```

如果是类型别名，则说明他们是等效的，它俩之间可直接替换，比如下面输出为：`int - int 0`

```
func main(){
	type X = int
	var x X
	var y int
	fmt.Printf("%T - %T %d", x, y, x+y)
}
```

如果没有等号则是类型定义，会产生一个新的类型，并拥有该类型已定义的操作。既然是新类型，它俩之间就没法直接替代，但可以做强制类型转换，去掉上面的等号和加法输出为：`main.X - int`，如果做加法则会报：`invalid operation: x + y (mismatched types X and int)`

type定义除了上面可以新定义类型、别名外，还可以定义结构体、接口：

**1、定义结构体：**

```
type person struct {
	Name string
	Age int
}
```

**2、定义接口**

```
type skill interface {
	Read()
	Write()
}
```

做类型断言时也会用到type：

```
func main(){
	var x interface{} = "abc"
	switch x.(type) {
	case string:
		fmt.Println("x type is string")
	case int:
		fmt.Println("x type is int")
	default:
		fmt.Println("unknown type")
	}
}
```

# 三、类型转换

Go属于强类型语言，类型之间的转换需要显式的转换。

## 3.1 强制转换

```
表达式 T(v) 将值 v 转换为类型 T
```

一些关于数值的转换：

```
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)
```

或者，更加简单的形式：

```
i := 42
f := float64(i)
u := uint(f)
```

适用类型：

- 整型、浮点型的互转



## 3.2 字符串转其他

1、字符串转整型

```
i, err := strconv.Atoi(s)
```

2、字符串转Bool类型

```
b, err := strconv.ParseBool(s)
```

3、字符串转浮点

```
f, err := strconv.ParseFloat(s, 32)
```

4、字符串转整型（int64）

```
i64, err := strconv.ParseInt(s, 10, 64)
```

5、字符串转无符号整型（int64）

```
ui64, err := strconv.ParseUint(s, 10, 64)
```

6、字符串转`[]byte`

```
b := []byte(s)
```

## 3.3 其他转字符串

1、整型转换为字符串：

```
strconv.Itoa(x)
```

2、浮点型转换为字符串：

```
str := strconv.FormatFloat(f float64, fmt byte, prec, bitSize int)
```

| 参数      | 描述                                                         |
| --------- | ------------------------------------------------------------ |
| *f*       | 需要转换的 float64 类型的变量。                              |
| *fmt*     | 使用 f 表示不使用指数的形式。                                |
| *prec*    | 保留几位小数。                                               |
| *bitSize* | 如果为 32，表示是 float32 类型，如果是 64，表示是 float64 类型。 |

3、`[]byte`转字符串：

```
s := string(b)
```

# 四、传值与传引用





---

- [1] [Go语言type关键字（类型别名）](http://c.biancheng.net/view/25.html)
- [2] [Go语言nil：空值/零值](https://www.cnblogs.com/lurenq/p/12013168.html)
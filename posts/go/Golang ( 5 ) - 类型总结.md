```
{
    "url": "go-types",
    "time": "2020/05/24 13:00",
    "tag": "Golang",
	"public": "no",
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

Golang中有一些基础类型：`byte` 、 `int` 、 `uint` 、 `string` 、`bool` 、 `float64` 等，也有一些组合的高级类型：`struct` 、 `[3]int` 、 `[]int` 、`map[int]string`等。我们可以从以下几方面对类型做一下总结：

- 如何自定义类型，查看当前变量的类型
- 不同类型之间的类型转换
- 各类型是传值还是传引用


# 二、类型定义

## 2.1 零值

当一个声明一个变量不赋初值时，`Go`语言会自动初始化值为此类型对应的零值。

- 布尔类型的零值：`false`
- 数值类型的零值：`0`
- 字符串类型的零值为空字符串
- 数组类型的零值为对应多个元素的零值，如`[3]int`的零值为[0, 0, 0]
- 指针、`slice`、`map`、`channel`、函数和`interface{}`的零值：`nil`

关于`nil`:


## 2.2 make && new

主要通过make和new做内存分配：

- make是用来分配并且初始化slice,map,chan等类型的对象
- new也是用来分配内存的,返回对应内向的0值的指针,但并不初始化对象

## 2.3 类型初始化



## 2.4 type关键字




## 2.5 其他

fmt包用法
打印类型与获取类型

print打印
reflectof



# 三、类型转换


## 3.1 字符串转其他

```
package main

import (
	"fmt"
	"strconv"
)

func main() {
	var s = "1.23"

	//字符串转整型（int)
	i, err := strconv.Atoi(s)
	fmt.Printf("%d %T, %v\n", i, i, err)

	//字符串转Bool类型
	b, err := strconv.ParseBool(s)
	fmt.Printf("%t %T, %v\n", b, b, err)

	//字符串转浮点型
	f, err := strconv.ParseFloat(s, 32)
	fmt.Printf("%.2f %T, %v\n", f, f, err)

	//字符串转整型（int64）
	i64, err := strconv.ParseInt(s, 10, 64)
	fmt.Printf("%d %T, %v\n", i64, i64, err)

	//字符串转无符号整型（int64）
	ui64, err := strconv.ParseUint(s, 10, 64)
	fmt.Printf("%d %T, %v\n", ui64, ui64, err)
}
```

## 3.2 其他转字符串


整型
浮点型

string与[]byte之间的转换：

## 3.3 数字类型互转



# 四、传值与传引用

## 指针

- new

- type
- make


# 七、 接口类型

断言

---

- [1] [Go语言type关键字（类型别名）](http://c.biancheng.net/view/25.html)
- [2] [Go语言nil：空值/零值](https://www.cnblogs.com/lurenq/p/12013168.html)
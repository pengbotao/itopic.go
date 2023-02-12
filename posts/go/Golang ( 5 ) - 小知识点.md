```
{
    "url": "go-tips",
    "time": "2020/05/24 13:00",
    "tag": "Golang",
    "toc": "yes"
}
```



### 1. 零值可用

当声明一个变量不赋初值时，`Go`语言会自动初始化值为此类型对应的零值。

- 布尔类型的零值：`false`
- 数值类型的零值：`0` 或`0.0`
- 字符串类型的零值为空字符串
- 数组类型的零值为对应多个元素的零值，如`[3]int`的零值为[0, 0, 0]
- 指针、`slice`、`map`、`channel`、`func`和`interface{}`的零值：`nil`

零值可用是指变量为零值时依然可以使用，这是一种设计理念。比如：

```
var s []int
s = append(s, 1)
fmt.Println(s)
```

append函数做了零值的判断，从而传入零值时可用的，直接使用`s[0] = 1`则是不可以的。类似的还有：

```
var mu sync.Mutex
mu.Lock()
mu.Unlock()

var b bytes.Buffer
b.Write([]byte("Hi"))
fmt.Println(b.String())
```

另外零值可用的类型要注意尽量避免值复制。

### 2. type关键字

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

### 3.包导入的是路径名还是包名？

1、import 后面部分是一个路径，只是大部分时候路径名与包名一致

2、调用地方使用的是包名，两者可以不一致

```
//main.go
package main

import v2 "demo/test"

func main() {
	v2.Show()
}

//test/v2.go
package v2

import "fmt"

func Show() {
	fmt.Println("v2")
}
```

当包名与到包导入路径中的最后一个目录名不同时，最后通过上面的方式将包名显示放入到导入语句中。

### 4. Go语言表达式求值顺序

在Go包中，包级别变量的初始化按照变量声明的先后顺序执行。如果某个变量的初始化直接或间接依赖其他变量，那么该变量在被依赖的变量之后。如：

```
var (
	a = b + 1
	b = 1
)
```

此处会先初始化b，在初始化a。

**代码块与作用域**

```
if c := f(); c > 1 {

}

等价于：
{
	c := f()
	if c > 1 {
		
	}
}
```

空的代码块`{}`也受其作用域控制，比如上面的`c`变量在其作用域范围之外则无法使用。

### 5. itoa实现枚举

```
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

Sunday = 0, Monday = 1 ...

const (
	_ = iota
	Weixin
	Alipay
)

Weixin = 1, Alipay = 2

const (
	StepOne = 1 << iota
	StepTwo
	StepThree
)
StepOne = 1, StepTwo = 2, StepThree = 4
```

### 6. defer关键字

1、多个defer按后进先出的规则执行，return之前执行

2、defer后面只能跟函数，如果是通过函数返回则会先对defer求值

```
func deferStart(name string) func() {
	fmt.Println("Defer Start", name)
	return func() {
		fmt.Println("Defer End", name)
	}
}

func demo() {
	defer deferStart("A")()
	defer deferStart("B")()
	fmt.Println("Run f()")
}

-------------------------------------------------
Defer Start A
Defer Start B
Run f()
Defer End B
Defer End A
```

### 易错知识点

#### 1. for ... range

```
package main

import "fmt"

func main() {
	m := []int{1, 2, 3}
	n := make(map[int]*int)
	for k, v := range m {
		fmt.Println(k, v)
		n[k] = &v
	}
	fmt.Printf("%+v", n)
}
```

> Output:

```
0, 1, 0xc000014168, 0xc000014180
1, 2, 0xc000014168, 0xc000014180
2, 3, 0xc000014168, 0xc000014180
map[0:0xc000014180 1:0xc000014180 2:0xc000014180]
```

可以看到Range里k,v的变量地址未发生改变，多次循环使用了同一块内存地址接收。

#### 2. 切片陷阱

```
package main

import "fmt"

func main() {
	a := [...]int{1, 2, 3, 4, 5}
	s1 := a[1:3]
	fmt.Printf("%#v, len:%d, cap: %d\n", s1, len(s1), cap(s1))
	s1[0] = 6
	fmt.Println(a, s1)
}
```

切片的底层结构是数组，上面切片s1指向底层数组a，起始点为a1[1], 长度为2, 容量为4。当s1的值变化时，底层数组也变了，所以会看到a数组的值也变化了。
如果把`s1[0] = 6` 改为 `s1 = append(s1, 6, 7, 8)`，超过了切片的容量，切片开辟一块新的空间扩容，与原数组脱离关系，从而不会被改变。

#### 3. 随机到相同的值

```
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(rand.Intn(100))
}
```

发现每次都随机到相同的值，需要指定下随机种子：`rand.Seed(time.Now().UnixNano())`

#### 4. Map元素不可寻址

```
type Person struct {
	Name string
	Age  int
}

func main() {
	x := make(map[string]Person)
	x["Lion"] = Person{"Lion", 3}

	x["Lion"].Name = "Test"
}
```

> ./main.go:17:17: cannot assign to struct field x["Lion"].Name in map

#### 5. map并发读写

```
func main() {
	m := make(map[int]struct{})
	go func() {
		for {
			m[0] = struct{}{}
		}
	}()
	go func() {
		for {
			fmt.Println(m[0])
		}
	}()
	select {}
}
```

> fatal error: concurrent map read and map write


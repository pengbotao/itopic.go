```
{
    "url": "golang",
    "time": "2020/03/01 19:00",
    "tag": "Golang",
    "toc": "yes"
}
```

# 一、基本数据类型

Go的数据类型有两种：一种是`语言内置的数据类型`，另外一种是通过语言提供的自定义数据类型方法自己定义的`自定义数据类型`。先看看语言内置的基础数据类型

## 1.1 数值型

数值型有三种，一种是`整数类型`，另外一种是`带小数的类型`(一般计算机里面叫做浮点数类型)，还有一种`虚数类型`。`整数类型`和数学里面不同的地方在于计算机里面正整数和零统称为无符号整型，而负整数则称为有符号整型。

### 1.1.1 整形

有符号整数|无符号整数|说明
---|---|---
|int8|uint8|占1个字节
|int16|uint16|占2个字节
|int32|uint32|占4个字节
|int64|uint64|占8个字节
|int|uint|32 位操作系统上64 位,64 位操作系统64 位
| |uintptr|32 位操作系统上为32位的指针,64 位操作系统为64位的指针

**取值范围：**

- 有符号：`[-2^(N-1) ~ 2^(N-1)-1]`
- 无符号：`[0 ~ (2^N - 1)]`

另外，还有一些别名类型：

- `byte类型`：这个类型和uint8是一样的，表示字节类型
- `rune类型`：这个类型和int32是一样的，用来表示unicode的代码点，就是unicode字符所对应的整数。

### 1.1.2 浮点型

Go的浮点数类型有两种，`float32`大约可以提供小数点后6位的精度，而`float64`可以提供小数点后15位的精度。

- `float32`：单精度浮点型
- `float64`：双精度浮点型

### 1.1.3 虚数类型

另外Go还有两个其他语言所没有的类型，虚数类型。

- `complex64`
- `complex128`

### 1.1.4 数的计算

对于数值类型，其所共有的操作为加法(＋)，减法(－)，乘法(＊)和除法(/)。另外对于整数类型，还定义了求余运算(%)。求余运算为整型所独有。如果对浮点数使用求余，比如这样

```
package main

import (
    "fmt"
)

func main() {
    var a float64 = 12
    var b float64 = 3

    fmt.Println(a % b)
}
```
编译时候会报错
```
invalid operation: a % b (operator % not defined on float64)
```

所以，这里我们可以知道所谓的数据类型有`两层意思`，一个是定义了`该类型所能表示的数`，另一个是定义了`该类型所能进行的操作`。

## 1.2 字符串类型
字符串就是一串固定长度的字符连接起来的字符序列。Go的字符串是由单个字节连接起来的。（对于汉字，通常由多个字节组成）。这就是说，传统的字符串是由字符组成的，而Go的字符串不同，是由字节组成的。这一点需要注意。

字符串的表示很简单。用(双引号"")或者(``号)来描述。

```
"hello world"
```
或者
```
`hello world`
```

唯一的区别是，**双引号之间的转义字符会被转义，而``号之间的转义字符保持原样不变**。

## 1.3 布尔类型

布尔型是表示真值和假值的类型。可选值为`true`和`false`。任何空值(`nil`)或者零值(0, 0.0, "")都不能作为布尔型来直接判断。所能进行的操作如下：

操作|说明
---|---
&&|与
\|\||或
!|非


# 二、常量和变量

## 2.1 变量

所谓的变量就是一个拥有指定`名称`和`类型`的数据存储位置。变量之所以称为变量，就是因为它们的值在程序运行过程中可以发生变化，但是它们的变量类型是无法改变的。Go语言是静态语言，并不支持程序运行过程中变量类型发生变化。如果强行将一个字符串值赋值给定义为`int`的变量，那么会发生编译错误。

```
package main

import (
	"fmt"
)

var x string = "Hello World"

func main() {
	var y string
	y = "Hello World"

	var z = "Hello World"

	m := "Hello World"

	fmt.Println(x, y, z, m)
}
```
通过`var`关键字定义变量，变量的定义包含以下四种方式：

- x: 声明变量`x`，并指定类型为`string`，同时赋初始值
- y: 先声明变量`y`, 之后再进行赋值
- z: 相当于变量`x`简化了指定类型，让Go语言推断变量的类型
- m: 把`var`关键字也省略了，需要知道变量的初始值，只能用于函数内部。


## 2.2 常量

所谓`常量`就是在程序运行过程中保持值`不变`的变量定义。常量的定义和变量类似，只是用`const`关键字替换了`var`关键字，另外常量在定义的时候必须有初始值。

```
package main

import (
	"fmt"
)

func main() {
	const x string = "Hello World"
	const y = "Hello World"
	fmt.Println(x, y)
}
```

## 2.3 多变量/常量定义

Go还提供了一种同时定义多个变量或者常量的快捷方式。

```
package main

import (
	"fmt"
)

func main() {
	var (
		a int     = 10
		b float64 = 32.45
		c bool    = true
	)
	const (
		Pi   float64 = 3.14
		True bool    = true
	)

	fmt.Println(a, b, c)
	fmt.Println(Pi, True)
}
```

# 三、控制结构

## 3.1 if

```
score := 80
if score >= 90 {
    fmt.Println("A")
} else if score >= 80 && score < 90 {
    fmt.Println("B")
} else if score >= 60 {
    fmt.Println("C")
} else {
    fmt.Println("D")
}
```

## 3.2 switch

```
score := 80

switch score / 10 {
case 10:
case 9:
    fmt.Println("A")
case 8:
    fmt.Println("B")
case 7:
case 6:
    fmt.Println("C")
default:
    fmt.Println("D")
}
```

**说明：**

- `switch`的判断条件可以为任何数据类型。
- 每个`case`后面跟的是一个完整的程序块，该程序块不需要`{}`，也不需要`break`结尾，因为每个`case`都是独立的。
- 可以为`switch`提供一个默认选项`default`，在上面所有的`case`都没有满足的情况下，默认执行`default`后面的语句。

## 3.3 for
```
for ...; ...; ...{
	...
}

for ...{
	...
}

for{
	...
}
```

# 四、高级数据类型

## 4.1 数组
数组是一个具有`相同数据类型`的元素组成的`固定长度`的`有序集合`。数组的定义:

```
var x [5]int
var y = [5]int{1, 2, 3, 4}
var z = [...]int{1, 2, 3, 4, 5}
```

- 1. 定义数组`x`，指定长度为5，没有赋初值，默认为零值。比如对于整数，零值就是0，浮点数，零值就是0.0，字符串，零值就是""，对象零值就是nil，所以x的5个元素都是0.
- 2. 定义数组`x`，同时初始化前4个元素。也可以通过`x[4] = 5`来赋值。
- 3. 长度用`...`代替，Go会自动计算出数组的长度。这种方式定义的数组需要赋初值。


## 4.2 切片

数组的长度是固定的，数组一旦定义后将无法增加新的元素，只能修改已有元素的值。所以切片诞生了，可以支持元素个数不确定的场景。切片有两种定义方式：

- 1. 先声明一个变量是切片，然后使用内置函数make去初始化这个切片
- 2. 通过取数组切片来赋值。

### 4.2.1 切片定义

**方法一：通过make初始化**

```
func main() {
	//定义并初始化切片
	var x = make([]int, 5, 10)

	fmt.Println(x)
	fmt.Println("Length:", len(x), "Capcity:", cap(x))
	//赋值
	for i := 0; i < len(x); i++ {
		x[i] = i
	}
	fmt.Println(x)

	//追加数据
	x = append(x, 5, 6, 7, 8, 9, 10)
	fmt.Println("Length:", len(x), "Capcity:", cap(x))
	fmt.Println(x)
}


//output:
[0 0 0 0 0]
Length: 5 Capcity: 10
[0 1 2 3 4]
Length: 11 Capcity: 20
[0 1 2 3 4 5 6 7 8 9 10]
```

- 切片通过make函数初始化，可以指定长度(`length`)和容量(`capacity`)两个属性，当不指定容量时，与长度相同。
- 如果追加的元素超过了`容量`大小，Go会自动地重新为切片分配容量，容量大小为原来的两倍。
- 虽然切片的`容量`可以大于`长度`，但是赋值的时候要注意最大的索引仍然是`len(x)－1`。否则会报索引超出边界错误。
- 虽然切片会自动扩容，但设置合理的容量大小可以减少内存的重新分配。

**方法二：通过数组切片**

```
func main() {
	var x = [5]int{1, 2, 3, 4, 5}
	var y = x[1:3]
	fmt.Println(y)
	fmt.Println("Length:", len(y), "Capcity:", cap(y))
	y = append(y, 4, 5, 6)
	fmt.Println(y)
	fmt.Println("Length:", len(y), "Capcity:", cap(y))
}

//output:
[2 3]
Length: 2 Capcity: 4
[2 3 4 5 6]
Length: 5 Capcity: 8
```

### 4.2.2 访问切片

```
func main() {
	var x = make([]int, 5, 10)
	for k := range x {
		fmt.Println(k)
	}

	for k, v := range x {
		fmt.Println(k, v)
	}
}
```

## 4.3 字典

字典是一组`无序`的，`键值对`的集合。数组通过`索引`来查找元素，而字典通过`键`来查找元素。

### 4.3.1 字典定义

**方法一：直接赋值**

```
func main() {
	var x = map[string]string{
		"A": "Apple",
		"B": "Banana",
		"O": "Orange",
		"P": "Pear",
	}

	for key, val := range x {
		fmt.Println("Key:", key, "Value:", val)
	}
}
```

**方法二：通过make函数初始化**

```
func main() {
	var x = make(map[string]string)

	x["A"] = "Apple"
	x["B"] = "Banana"
	x["O"] = "Orange"
	x["P"] = "Pear"

	for key, val := range x {
		fmt.Println("Key:", key, "Value:", val)
	}
}
```

### 4.3.2 访问字典

如果访问的元素所对应的键不存在于字典中，返回零值。对于字符串零值就是""，对于整数零值就是0。可以通过下面的方式来进行遍历和判断：

```
func main() {
	var x = make(map[string]string)

	x["A"] = "Apple"
	x["B"] = "Banana"
	x["O"] = "Orange"
	x["P"] = "Pear"

	//字典遍历
	for key, val := range x {
		fmt.Println("Key:", key, "Value:", val)
	}

	//判断键是否存在
	if val, ok := x["C"]; ok {
		fmt.Println(val)
	}

	//删除键
	delete(x, "P")
}
```

# 五、函数

函数，简单来讲就是一段将`输入数据`转换为`输出数据`的公用代码块。

## 5.1 函数定义

**方法一：基本定义**

```
func func1(x int, y int) int {
	return x + y
}
```

**方法2：可变长参数**

```
func func2(nums ...int) int {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return sum
}
```

**方法3：多返回值**

```
func func3(x, y int, z string) (int, string) {
	return x + y, z + " World"
}
```

**方法4：命名返回值**

- 可以为返回值预先定义一个名称，在函数结束的时候，直接return就可以返回所有的预定义返回值
- 如果定义了命名返回值，那么在函数内部将不能再重复定义一个同样名称的变量。

```
func func4(x int, y int) (sum int, m, n int) {
	m = x
	n = y
	sum = x + y
	return
}
```

## 5.2 闭包函数

所谓闭包函数就是将整个函数的定义一气呵成写好并赋值给一个变量。然后用这个变量名作为函数名去调用函数体。

```
func main() {
	sum := func(nums ...int) int {
		s := 0
		for _, v := range nums {
			s += v
		}
		return s
	}
	fmt.Println(sum(1, 2, 3))
}
```

## 5.3 递归函数

所谓递归，就是在函数的内部重复调用一个函数的过程。需要注意的是这个函数必须能够一层一层分解，并且有出口。

```
func fibonacci(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	for i := 1; i < 10; i++ {
		fmt.Println(fibonacci(i))
	}
}
```


# 六、错误和异常处理

## 6.1 defer关键字

Go语言提供了关键字`defer`来在函数运行结束的时候运行一段代码或调用一个清理函数。

```
func main() {
	defer func() {
		fmt.Println("Print From Defer.")
	}()
	fmt.Println("Hello World.")
}
```

上面示例虽然defer操作写在打印的前面，但实际会在main函数结束前调用，运行示例就可以看到Defer的打印在Hello之后。所以常用`defer`来处理一些清理和释放资源的操作，比如常见的：

```
mu.Lock()
defer mu.Unlock()
```


`defer`有一些特性：

- 如果有多个`defer`操作，按照`FILO`（先进后出）的方式执行。
- `defer`还可以修改函数中的命名返回值。

## 6.2 panic & recover

- `panic`: 抛出一条异常信息。程序执行中调用了panic则正常执行流程终止，但是该函数中panic之前定义的defer语句将被依次执行。
- `recover`: 用于将panic的信息捕捉。recover必须定义在panic之前的defer语句中。

```
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s", r)
		}
	}()
	panic("panic info.")
}
```

## 6.3 错误判断

```
func division(x, y int) (int, error) {
	if y == 0 {
		return 0, errors.New("integer divide by zero")
	}
	return x / y, nil
}

func main() {
	x, err := division(1, 0)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(x)
	}
}
```

也可以通过下划线`_`忽略掉error信息，但潜在的错误就被忽略掉了，所以经常看到一堆的`if err != nil`的判断。

```
func main() {
	x, _ := division(1, 0)
	fmt.Println(x)
}
```

---

- [1] [Go轻松学](https://www.kancloud.cn/itfanr/go-quick-learn)
- [2] [golang--数据类型](https://studygolang.com/articles/16011)

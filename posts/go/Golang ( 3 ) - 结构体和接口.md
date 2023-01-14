```
{
    "url": "go-struct",
    "time": "2020/04/20 11:00",
    "tag": "Golang",
    "toc": "yes"
}
```

# 一、关于结构体

## 1.1 基本定义

Go提供的`结构体`就是把使用`各种数据类型定义`的`不同变量组合起来`的`高级数据类型`。示例：

```
type Person struct {
	Name string
	Age  int
}
```

通过`type`关键字指定`Person`为一个结构体`struct`，并且包含2个属性`Name`和`Age`，类型分别是`string`和`int`。所以说结构体可以把多种数据类型组合起来形成新的数据类型。 可以将`Person`结构体可以理解成定义了一个`Person`类，有两个成员属性`Name`和`Age`，只是暂时还没有成员方法，接下来我们给结构体定义一个方法。

## 1.2 定义方法

给结构体定义了一个`Info`方法，返回一个`string`。

```
type Person struct {
	Name string
	Age  int
}

func (p Person) Info() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}
```

- 首先可以看到`Info`方法并非定义在结构体对应的花括号内部，而是独立于结构体之外，这和一些语言不同，倒和Lua的table有点相似。
- 结构体方法和普通函数定义唯一的区别是在`func`关键字和函数名`Info`中间增加了一个结构体类型限定：`(p Person)`，表示这是一个`Person`结构体的方法。
- 结构体类型限定`(p Person)`，`p`表示形参名可以用任意变量名;`Person`是`p`的数据类型。
- 可通过`.`形式来获取到结构体的成员属性和方法，如`p.Name`


## 1.2 初始化

```
func main() {
	//方式一：先定义然后进行赋值
	var p Person
	p.Name = "iTopic"
	p.Age = 18

	//方式二：定义并初始化
	var q = Person{Name: "iTopic", Age: 18}

	//方法三：方法二相同，只不过这种写法最后一行需要保留逗号
	m := Person{
		Name: "iTopic",
		Age:  18,
	}

	//方式四：按字段顺序进行赋值就无需指定字段了,需要按顺序对每一个字段进行初始化
	n := Person{"iTopic", 18}

	fmt.Println(p, q, m, n)

	//调用结构体的Info方法
	fmt.Println(p.Info())
}
```

- 属性和方法都通过`.`来访问。
- 调用过程就像是初始化了一个类，并调用了类的方法。`Info()`方法也可以访问当前结构体实例`p`的成员属性。

# 二、结构体指针

## 2.1 指针接收器

前面结构体方法上也可以通过指针的方式来接收。

```
type Person struct {
	Name string
	Age  int
}

func (p *Person) Info() string {
	p.Age = p.Age + 1
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func main() {
	p := &Person{"iTopic", 18}
	fmt.Println(p)
	fmt.Println(p.Info())
	fmt.Println(p)
}

//output:

&{iTopic 18}
Name: iTopic, Age: 19
&{iTopic 19}
```

和前面的区别可以看到

- `Info`方法的结构体类型改为了`(p *Person)`，多了一个星号（`*`）。
- `Info`方法内部使用无区别，还是用`p.Age`访问，由于是接收的指针，可以直接修改结构的属性。
- 如果把`Info`方法的`(p *Person)`定义改回`(p Person)`，重新执行可以看到执行完成之后`p.Age`还会是18，即值拷贝。
- 此处为指针接收器，则不可以通过`Person{}.Info()`来调用，函数传递上需要是`*Person`类型，可以写成：`(&Person{}).Info()`

## 2.2 结构体调用

这里一种有4种调用方式：

```
func (p *Person) Info() {
	fmt.Printf("Pointer: Name: %s, Age: %d\n", p.Name, p.Age)
}

func (p Person) Job() {
	fmt.Printf("Value: Name: %s, Age: %d\n", p.Name, p.Age)
}

func main() {
	//第一种情况：不通，报错 cannot call pointer method on Person{}
	//Person{}.Info()
	Person{}.Job()     // 输出：Value: Name: , Age: 0
	(&Person{}).Info() //输出：Pointer: Name: , Age: 0
	(&Person{}).Job()  //输出：Value: Name: , Age: 0
}
```

第一种方式报错，提示不能通过结构体值类型来调用指针方法。如果我们将初始化赋值给一个新的变量，则Go语言会帮我们做自动转换，都可以调用成功：

```
func main() {
	p1 := Person{}
	p2 := &Person{}
	p1.Info()
	p1.Job()
	p2.Info()
	p2.Job()
}
```


# 三、结构体内嵌

## 3.1 基本嵌套

结构体是由多种数据类型组合起来的高级类型，自然结构体也可以嵌套结构体。按理所当然的去设置和访问就好了，比如：`p.Job.CompanyName`

```
type Person struct {
	Name string
	Age  int
	Job  Job
}

type Job struct {
	CompanyName string
}
```

## 3.2 内嵌匿名字段

我们去掉了`Person`中前面的`Job`变量名，同时给`Job`结构体增加了`Show`方法，可以看到：

- 可以通过`p.CompanyName`直接访问`Job`结构体中的字段
- 通过通过`p.Show()`直接访问方法

`Person`获得了`Job`的属性和方法，这就像从`Person`类继承了`Job`类的成员属性和成员方法。

```
type Person struct {
	Name string
	Age  int
	Job
}

type Job struct {
	CompanyName string
}

func (p *Person) Info() string {
	p.Age = p.Age + 1
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func (j *Job) Show() string {
	return fmt.Sprintf("Job: %s", j.CompanyName)
}

func main() {
	var p Person
	p.Name = "iTopic"
	p.Age = 18
	p.CompanyName = "iTopic.org"
	fmt.Println(p, p.Info(), p.Show())
}

//output:
{iTopic 18 {iTopic.org}} Name: iTopic, Age: 19 Job: iTopic.org
```

如果存在同名的字段和方法则会优先获取到当前结构体的属性和方法，不存在才会去父结构体找。比如将`Job.CompanyName`改成`Name`，方法改成`Info()`，则可以这样子访问：

```
func main() {
	var p Person
	p.Name = "iTopic"
	p.Age = 18
	p.Job.Name = "iTopic.org"
	fmt.Println(p, p.Info(), p.Job.Info())
}
```

如果有多个这种继承字段也存在同名的方法，则就需要上面这种显示的指定了，否则会报：`ambiguous selector p.Info`

## 3.2 结构体锁

这里为前一章节的一个示例，比如对`Age`进行一千次调用累加操作，如果要确保数据准确就需要上锁，而`sync.Mutex`也是一个结构体，所以可以直接用内嵌的方式，`Person`就继承了`Lock()`和`Unlock()`方法了。

```
type Person struct {
	Name string
	Age  int
	C    chan bool
	sync.Mutex
}

func (p *Person) Add() {
	p.Lock()
	defer p.Unlock()
	p.Age++
}

func main() {
	p := new(Person)
	p.Age = 18
	p.C = make(chan bool)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			p.Add()
			p.C <- true
		}(i)
	}
	for i := 0; i < 1000; i++ {
		<-p.C
	}
	fmt.Println(p.Age)
}
```

# 四、关于接口

## 4.1 接口定义

接口是一组方法的集合。接口的定义，首先是关键字`type`，然后是接口名称，最后是关键字`interface`表示这个类型是接口类型。示例，

- 定义了一个接口`ExportHandler`，包含一个方法：`Export()`。
- `CsvExporter` 和 `ExcelExporter` 都实现了该方法，但仅仅是实现了改方法，代码层面并无显示关联。
- 任何类型，只要实现了该接口中方法集，那么就属于这个类型，所以`CsvExporter`和`ExcelExporter`都可以定义在`ExportHandler`数组中。

```
type ExportHandler interface {
	Export()
}

type CsvExporter struct {
}

func (e CsvExporter) Export() {
	fmt.Println("Export CVS File.")
}

type ExcelExporter struct {
}

func (e ExcelExporter) Export() {
	fmt.Println("Export Excel File.")
}

func main() {
	var handlerList = []ExportHandler{CsvExporter{}, ExcelExporter{}}

	for _, handle := range handlerList {
		handle.Export()
	}
}
```

注：不能声明一个接口的指针，比如函数使用`*ExportHandler`类型， 则会提示：

```
type *ExportHandler is pointer to interface, not interface
```

## 4.2 空接口

空接口就是不含有方法的接口：`interface{}`，相当于所有的类型都实现了空接口，从示例可以看到类型都可以传给函数`show`而不会报类型错误。

```
func show(v interface{}) {
	fmt.Println(v)
}

func main() {
	show(123)
	show("123456")
	show([]int{1, 2, 3})
}
```

但上面的类型只是表示可以传给函数`show`， 而实际上无法做运算：

```
func show(v interface{}) {
	fmt.Println(v + 1)
}

func main() {
	show(123)
}
//invalid operation: v + 1 (mismatched types interface {} and int)
```

所以要做运算需要做断言：

```
func show(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Println(val + 1)
	case string:
		fmt.Println(val + "123456")
	default:
		fmt.Println("unknown data type")
	}
}

func main() {
	show(1)
}
```

---

- [1] [Go结构体和接口](https://www.kancloud.cn/itfanr/go-quick-learn/81641)
- [2] [Golang 入门 : 结构体(struct)](https://www.cnblogs.com/sparkdev/p/10761825.html)
- [3] [Golang关键字--type 类型定义](https://www.jianshu.com/p/a02cf41c0520)
- [4] [Golang 之 interface接口全面理解](https://www.cnblogs.com/echojson/p/10746807.html)
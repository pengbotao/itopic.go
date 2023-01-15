```
{
    "url": "goroutine",
    "time": "2020/05/24 13:00",
    "tag": "Golang",
    "toc": "yes"
}
```

# 一、Go协程

只需要在方法前加一个`go`关键字就可以让一个普通方法协程化。以下面的代码为例，一般同步阻塞的编码方式下会按顺序打印012然后再输出Finish. 示例中启动了3个协程之后主进程会继续往下执行，不会等待函数返回，大概率会先看到Finish输出，然后看到012或者210。

即加上`go`关键字之后程序不会同步阻塞主进程，协程的执行速度跟程序复杂性关系，无法保证先启动的协程先执行完毕。

```
func main() {
	for i := 0; i < 3; i++ {
		go func(v int) {
			fmt.Println(v)
		}(i)
	}
	fmt.Println("Finish.")
	time.Sleep(time.Second)
}
```

通常情况下同步的逻辑方式书写是最方便的，无须考虑程序逻辑的先后关系，上面示例最后一行休眠1秒，确保协程可以执行完毕，但正常逻辑下无法确保1秒内协程能执行完毕，也不会只执行一个打印这么简单，通常还需要能获取到函数的返回，所以这里面存在两个主要问题：

- 协程同步：多个协程之间同步等待、超时退出等
- 协程通信：多个协程之间的数据通信

Channel正是为协程间通信而产生，接下来分开看看同步等待、超时、通信的问题。

# 二、 Channel简介

通道（channel）提供了协程之间的`通信方式`以及`运行同步机制`。

## 2.1 通道定义

Channel是Go中的一个核心类型，Goroutine通过它发送或者接收数据。就像`map`和`slice`数据类型一样,`channel`必须先创建再使用:

```
ch := make(chan int)
```

操作符是箭头 `<-`(箭头的指向就是数据的流向)。

```
ch <- v    // 发送值v到Channel ch中
v := <-ch  // 从Channel ch中接收数据，并将数据赋值给v
```

chan按读写方向可以分为双向Channel和单向Channel（只读Channel和只写Channel），根据是否缓冲可以分为带缓冲的Channel和无缓冲Channel。

## 2.2 select语句

select 是`Go`中的一个控制结构，类似于用于通信的`switch`语句。每个`case`必须是一个通信操作，要么是发送要么是接收。`select`随机执行一个可运行的`case`。如果没有`case`可运行，它将阻塞，直到有 `case`可运行。一个默认的子句应该总是可运行的。

**基本用法**

```
select {
case <- chan1:
// 如果chan1成功读到数据，则进行该case处理语句
case chan2 <- 1:
// 如果成功向chan2写入数据，则进行该case处理语句
default:
// 如果上面都没有成功，则进入default处理流程
}
```
以下描述了 select 语句的语法：

- 每个`case`都必须是一个通信
- 所有`channel`表达式都会被求值
- 所有被发送的表达式都会被求值
- 如果任意某个通信可以进行，它就执行，其他被忽略。
- 如果有多个`case`都可以运行，`Select`会随机公平地选出一个执行。其他不会执行。
- 否则：
	- 如果有`default`子句，则执行该语句。
	- 如果没有`default`子句，`select`将阻塞，直到某个通信可以运行；`Go`不会重新对`channel`或值进行求值。

**示例**

```
func main() {
	ch1 := make(chan int)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Hello"
	}()

	fmt.Println("Start")
	select {
	case <-ch1:
		fmt.Println("Read From ch1")
	case <-ch2:
		fmt.Println("Read From ch2")
	default:
		fmt.Println("Read From Default")
	}
	fmt.Println("END")
}
```

- 因为两个通道都要等1秒，存在default则直接执行了default语句，打印`Read From Default`，退出`swith`
- 如果去掉default，则阻塞等待可能打印出`Read From ch1`

如果`select`里啥都没有 `select{}`，则会等待，达到阻塞Goroutine的目的。如果当前运行环境没有新的协程而使用该语句则会抛错。


# 三、协程同步

## 3.1 通过Channel同步

创建了一个存储10个bool类型的通道，函数执行成功向通道里写入true，执行失败向通道里写入false。启动一个循环从通道读取数据，读取10次之后程序在打印最后的结果：

`true true true true true true true true true true [0 1 2 3 4 5 6 7 8 9]`

```
func main() {
	ch := make(chan bool, 10)

	data := make([]int, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			defer func() {
				if err := recover(); err != nil {
					ch <- false
				}
			}()
			data[idx] = idx
			ch <- true
		}(i)
	}
	for j := 0; j < 10; j++ {
		fmt.Print(<-ch, " ")
	}
	fmt.Println(data)
}
```


## 3.2 通过sync同步

控制流程同步等待也可以通过`sync.WaitGroup`来实现，`WaitGroup`对象内部有一个计数器，最初从0开始，它有三个方法：

- `Add()`: 计数器增加N
- `Done()`: 完成一个任务，计数器减少1
- `Wait()`: 同步阻塞，计数器为0之后才继续向下执行

```
func main() {
	var wg sync.WaitGroup

	data := make([]int, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int, w *sync.WaitGroup) {
			defer func() {
				w.Done()
			}()
			data[idx] = idx
		}(i, &wg)
	}
	wg.Wait()
	fmt.Println(data)
}
```

---

通过这两种方法可以实现多协程的同步等待，但这会涉及到另一个超时问题，如果多个协程时间较长，主线程就一直等待了，所以还得考虑同步超时问题。

# 四、同步超时

这个是前面的补充，正常单个协程可以通过`select + time.After`实现超时控制。

```
func main() {
	ch := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- true
	}()

	select {
	case <-ch:
		fmt.Println("Read From CH")
	case <-time.After(2 * time.Second):
		fmt.Println("timeout")
	}
}
```

## 4.1 Channel同步超时

来看看`3.1`中多个goroutine的超时问题，如果主线程想最多只等待2s：

```
func main() {
	ch := make(chan bool, 10)

	data := make([]int, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			defer func() {
				if err := recover(); err != nil {
					ch <- false
				}
			}()
			time.Sleep(time.Duration(idx) * time.Second)//Sleep idx sec
			data[idx] = idx
			ch <- true
		}(i)
	}
	timeout := make(chan bool)
	go func() {
		time.Sleep(2 * time.Second)
		timeout <- true
	}()

	for{
		select {
		case t := <- ch:
			fmt.Println(t)
		case <- timeout:
			fmt.Println("Timeout")
			goto EXIT
		}
	}
EXIT:
	fmt.Println(data)
}
```

## 4.2 sync同步超时

`3.2`的超时也借助了Channel来实现。

```
func main() {
	var wg sync.WaitGroup

	data := make([]int, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer func() {
				wg.Done()
			}()
			time.Sleep(time.Duration(idx) * time.Second)
			data[idx] = idx
			fmt.Println(idx)
		}(i)
	}
	wc := make(chan bool)
	go func() {
		wg.Wait()
		wc <- true
	}()
	select {
	case <- time.After(3 * time.Second):
		fmt.Println("Timeout")
	case <- wc:
		fmt.Println("Wg Wait")
	}
	fmt.Println(data)
}
```

# 五、协程通信

协程之间数据交互上主要有两种方式，一种为全局变量然后通过锁来控制原子性，另一种则是通过channel来进行通信。

## 5.1 全局变量

启动10个协程来执行1加到10的操作，s变量为协程共享，所以需要加锁才会正确输出55，若去掉锁的三行代码，则会出现非55的情况。

```
func sum() int {
	s := 0
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(10)
	for i := 1; i <= 10; i++ {
		go func(i int) {
			mutex.Lock()
			s += i
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	return s
}
func main() {
	//执行100次结果都是55
	for i := 0; i < 100; i++ {
		fmt.Print(sum(), " ")
	}
}
```

## 5.2 通过Channel来通信

将结果写到通道中，从通道中读取结果进行累加。

```
func sum() int {
	ch := make(chan int)
	for i := 1; i <= 10; i++ {
		go func(i int) {
			ch <- i
		}(i)
	}
	s := 0
	for i := 0; i < 10; i++ {
		s += <-ch
	}
	return s
}
func main() {
	for i := 0; i < 100; i++ {
		fmt.Print(sum(), " ")
	}
}
```

# 六、并发控制

## 6.1 先处理再消费

这里为等待所有任务处理完成后再对结果进行处理，通过`sync.WaitGroup`来确认同步阻塞，通过有缓冲Channel的阻塞特性来控制最大执行的并发数，在启动Goroutine之前往通道里写值，处理完成之后再读取。处理成功后将结果写到`result`通道中，这里模拟了一个延时与错误操作，也就是`result`通道可能写不满，所以在任务执行完毕后做了`close(result)`操作，避免后续遍历通道时阻塞。

```
var (
	taskNum     = 10
	maxParalNum = 3
)

func main() {
	wg := sync.WaitGroup{}
	limit := make(chan struct{}, maxParalNum)
	result := make(chan int, taskNum)

	task := func() (int, error) {
		time.Sleep(time.Second * 2)
		if rand.Intn(3) == 1 {
			return 0, fmt.Errorf("task error")
		}
		return 1, nil
	}

	for i := 1; i <= taskNum; i++ {
		wg.Add(1)
		limit <- struct{}{}
		go func(i int) {
			fmt.Println("start task", i, time.Now())
			defer func() {
				<-limit
				wg.Done()
			}()
			if t, err := task(); err == nil {
				result <- t
			}
		}(i)
	}
	wg.Wait()
	close(result)
	for res := range result {
		fmt.Printf("%+v", res)
	}
}
```

## 6.2 边处理边消费

如果想限制并发的同时，一边处理一边消费。则可以调整为发送时起新的协程去发送，通过`limit`来控制发送的数量，根据任务次数来读取，不过这就要求执行任务时需往通道里写入对应数量的结果。

```
var (
	taskNum     = 10
	maxParalNum = 3
)

type data struct {
	code int
	res  int
}

func main() {
	limit := make(chan struct{}, maxParalNum)
	result := make(chan data, taskNum)

	task := func() (data, error) {
		time.Sleep(time.Second * 2)
		if rand.Intn(3) == 1 {
			return data{}, fmt.Errorf("task error")
		}
		return data{code: 1}, nil
	}
	go func() {
		for i := 1; i <= taskNum; i++ {
			limit <- struct{}{}
			go func(i int) {
				fmt.Println("start task", i, time.Now())
				defer func() {
					<-limit
				}()
				t, _ := task()
				result <- t
			}(i)
		}
	}()

	for i := 0; i < taskNum; i++ {
		res := <-result
		fmt.Printf("%+v\n", res)
	}
}
```

协程之间的通信进制大体如上，下一篇将介绍怎么通过上下文`Context`来控制协程之间的交互。

---

- [1] [go 语言之行--golang 核武器 goroutine 调度原理、channel 详解](https://learnku.com/articles/41668)
- [2] [一文读懂什么是进程、线程、协程](https://www.jianshu.com/p/80bde972196d)
- [3] [七周七并发模型](http://yuedu.163.com/source/a4b77ff9abaf4109acd11c38e5c8babc_4)
- [4] [Go语言通道（chan）——goroutine之间通信的管道](http://c.biancheng.net/view/97.html)
- [5] [golang sync包互斥锁和读写锁的使用](http://www.361way.com/rwmutex/5984.html)
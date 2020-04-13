```
{
    "url": "goroutine",
    "time": "2020/04/20 11:00",
    "tag": "Golang"
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

通常情况下同步的逻辑方式书写是最方便的，无须考虑程序逻辑的先后关系，上面示例最后一行休眠1秒，确保协程可以执行完毕，但正常逻辑下无法确保1秒内协程能执行完毕，也不会只执行一个打印这么简单，通常还需要能获取到函数的返回。所以需要有一种通信方式能解决此类问题，而通道正是为协程间通信而产生。

# 二、 通道

通道（channel）提供了协程之间的`通信方式`以及`运行同步机制`。

## 2.1 通道定义

Channel是Go中的一个核心类型，你可以把它看成一个管道，通过它并发核心单元就可以发送或者接收数据进行通讯，它的操作符是箭头 `<-`(箭头的指向就是数据的流向)。

```
ch <- v    // 发送值v到Channel ch中
v := <-ch  // 从Channel ch中接收数据，并将数据赋值给v
```

就像`map`和`slice`数据类型一样,`channel`必须先创建再使用:

```
ch := make(chan int)
```

# 三、协程同步与通信示例

## 3.1 通过Channel同步

创建了一个存储10个bool类型的通道，函数执行成功向通道里写入true，执行失败向通道里写入false。启动一个循环从通道读取数据，读取10次之后程序在打印最后的结果：

`true true true true false true false false false false [0 1 2 3 4]`

```
func main() {
	ch := make(chan bool, 10)

	data := make([]int, 5)
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


## 3.2 通过sync.WaitGroup同步

控制流程同步等待也可以通过`sync.WaitGroup`来实现，`WaitGroup`对象内部有一个计数器，最初从0开始，它有三个方法：

- `Add()`: 计数器增加N
- `Done()`: 完成一个任务，计数器减少1
- `Wait()`: 同步阻塞，计数器为0之后才继续向下执行

```
func main() {
	var wg sync.WaitGroup

	data := make([]int, 5)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}()
			data[idx] = idx
		}(i)
	}
	wg.Wait()
	fmt.Println(data)
}
```

## 3.3 模拟生产者与消费者

```
func main() {
	ch := make(chan string)

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			ch <- time.Now().Format("2006-01-02 15:04:05")
		}
	}()

	go func() {
		for {
			t, ok := <-ch
			if ok {
				fmt.Println(t)
			}
		}
	}()

	select {}
}
```

这个示例启动了2个协程，一个用来每一秒往通道里写一个时间，另一个用来从通道里读取，模拟生产者和消费者的情况。当然就示例本身实现起来只需要上面生产者并打印即可。

```
func main() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}
```


- [1] [go 语言之行--golang 核武器 goroutine 调度原理、channel 详解](https://learnku.com/articles/41668)
- [2] [一文读懂什么是进程、线程、协程](https://www.jianshu.com/p/80bde972196d)
- [3] [七周七并发模型](http://yuedu.163.com/source/a4b77ff9abaf4109acd11c38e5c8babc_4)
- [4] [Go语言通道（chan）——goroutine之间通信的管道](http://c.biancheng.net/view/97.html)
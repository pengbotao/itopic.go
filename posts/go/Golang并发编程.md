```
{
    "url": "goroutine",
    "time": "2020/04/20 11:00",
    "tag": "Golang",
    "public": "no"
}
```


# Goroutine的同步与通信

## 通过sync包同步goroutine

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


## 通过Channel通信

```
func main() {
	ch := make(chan bool, 10)

	data := make([]int, 5)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					ch <- false
				}
			}()
			data[idx] = idx
			ch <- true
		}(i)
	}
	for j := 0; j < 10; j++ {
		<-ch
	}
	fmt.Println(data)
}
```

## 生产者与消费者

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

**精简**

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
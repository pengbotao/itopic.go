```
{
    "url": "go-context",
    "time": "2021/08/28 19:00",
    "tag": "Golang",
    "toc": "yes"
}
```

# 一、概述

用于协程之间传递上下文信息，包含取消信号、超时信号、传值。

# 二、基本结构

`context`包中定义了一个`context.Context`结构体，可通过两种方式获取：

```
context.Background()
context.TODO()
```

支持四种使用方式，使用时需要传入`context.Context`对象并返回新的Context对象，其中取消、超时会返回一个取消函数。

## 2.1 WithCancel

取消信号，返回`context.Context`和取消函数`CancelFunc`，调用`CancelFunc`则会终止调用树上协程的执行。

```
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

## 2.2 WithDeadline

超时截止时间，设置具体的截止时间点。

```
WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

## 2.3 WithTimeout

超时信号，和`WithDeadline`的相似，时间点从`WithDeadline`的具体时间点变为从当前时间开始的相对时间。

```
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

## 2.4 WithValue

键值，可以在协程之间进行简单的数据传递。

```
func WithValue(parent Context, key, val interface{}) Context
```

# 三、基础理解





---

- [1] [Go语言设计与实现 - 6.1 上下文 Context](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/)
- [2] [用 10 分鐘了解 Go 語言 context package 使用場景及介紹](https://blog.wu-boy.com/2020/05/understant-golang-context-in-10-minutes/)




```
{
    "url": "go-coding",
    "time": "2024/08/10 11:35",
    "tag": "Golang",
    "toc": "no"
}
```

# 一、Option Pattern

在 Go 语言中，“选项”模式（Option Pattern），用于灵活配置结构体或者函数的行为，通过可选的配置参数，允许用户在使用一个函数或初始化一个结构体时，选择性地传递参数，而不需要显式地提供所有的配置项。

### Option 模式的基本原理

1. **定义一个结构体**: 包含所有可配置的选项字段。
2. **定义 Option 类型**: Option 通常是一个函数类型，用于接收并修改配置结构体。
3. **提供 Option 的函数**: 创建一组用于设置配置项的函数，每个函数返回一个 `Option`。
4. **构造函数**: 在构造函数中应用这些选项。

假设我们要创建一个 `Server` 结构体，允许用户配置 `Host`、`Port` 和 `Timeout`。

```
package main

import (
    "fmt"
    "time"
)

// Server 配置结构体
type Server struct {
    Host    string
    Port    int
    Timeout time.Duration
}

// Option 是一个函数类型，用于修改 Server 的配置
type Option func(*Server)

// NewServer 创建并初始化一个 Server
func NewServer(opts ...Option) *Server {
    // 默认配置
    server := &Server{
        Host:    "localhost",
        Port:    8080,
        Timeout: 30 * time.Second,
    }

    // 应用所有的 Option
    for _, opt := range opts {
        opt(server)
    }

    return server
}

// WithHost 返回一个设置 Host 的 Option
func WithHost(host string) Option {
    return func(s *Server) {
        s.Host = host
    }
}

// WithPort 返回一个设置 Port 的 Option
func WithPort(port int) Option {
    return func(s *Server) {
        s.Port = port
    }
}

// WithTimeout 返回一个设置 Timeout 的 Option
func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.Timeout = timeout
    }
}

func main() {
    // 使用默认配置创建 Server
    server1 := NewServer()
    fmt.Printf("Server1: %+v\n", server1)

    // 使用自定义 Host 和 Port 创建 Server
    server2 := NewServer(WithHost("127.0.0.1"), WithPort(9090))
    fmt.Printf("Server2: %+v\n", server2)

    // 使用自定义 Timeout 创建 Server
    server3 := NewServer(WithTimeout(60 * time.Second))
    fmt.Printf("Server3: %+v\n", server3)
}

```

还有一种使用结构体的 setter 方法并返回指针也可以实现链式调用。这种方法在某些场景下与 Option 模式类似，可以灵活地设置结构体的属性。

```
type Server struct {
    Host    string
    Port    int
    Timeout time.Duration
}

// 设置 Host 并返回 Server 指针
func (s *Server) SetHost(host string) *Server {
    s.Host = host
    return s
}

// 设置 Port 并返回 Server 指针
func (s *Server) SetPort(port int) *Server {
    s.Port = port
    return s
}

// 设置 Timeout 并返回 Server 指针
func (s *Server) SetTimeout(timeout time.Duration) *Server {
    s.Timeout = timeout
    return s
}

func main() {
    server := &Server{}
    server.
        SetHost("127.0.0.1").
        SetPort(9090).
        SetTimeout(60 * time.Second)

    fmt.Printf("Server: %+v\n", server)
}
```

Option 模式可以通过不同的 Option 组合来实现灵活的配置，而不会影响构造函数的签名，添加新配置项时，只需增加一个新的 Option 函数，不需要改动结构体的内部逻辑。

# 二、New函数

`New` 函数是一种常见的设计方式，用于初始化和返回结构体的实例。`New` 函数通常用于封装结构体的初始化逻辑，特别是当结构体具有复杂的默认值设置或需要执行额外的初始化操作时。

### `New` 函数的基本形式

一个典型的 `New` 函数会创建并返回结构体的指针。它的命名通常以 `New` 开头，后跟结构体的名称，例如 `NewServer`、`NewClient` 等。

通过结合可选参数的 Option 模式，可以使 `New` 函数更加灵活，允许用户传递额外的配置。

```
// NewClient 创建并返回带有可选参数的 Client 结构体
func NewClient(host string, port int, opts ...Option) *Client {
    client := &Client{
        Host:    host,
        Port:    port,
        Timeout: 30, // 默认 30 秒超时
    }

    for _, opt := range opts {
        opt(client)
    }

    return client
}
```

New也可以直接返回接口，用于初始化实现某个接口的结构体实例。

```
package main

import "fmt"

type Database interface {
    Connect() string
}

type MySQL struct {
    DSN string
}

// NewMySQL 创建并返回实现了 Database 接口的 MySQL 实例
func NewMySQL(dsn string) Database {
    return &MySQL{DSN: dsn}
}

func (db *MySQL) Connect() string {
    return "Connected to MySQL with DSN: " + db.DSN
}

func main() {
    db := NewMySQL("user:pass@tcp(localhost:3306)/dbname")
    fmt.Println(db.Connect())
}
```

# 三、面向接口编程

Go 的接口提供了极大的灵活性，因为它们是隐式实现的。这意味着一个类型只要实现了接口中所有的方法，就自动满足该接口，不需要显式地声明实现了哪个接口。这种特性使得代码更容易扩展和组合，如对支付系统中不同支付方式的处理。

```
type PaymentProcessor interface {
    Pay(amount float64) string
}
```

接着，实现几种不同的支付方式，例如信用卡支付和 PayPal 支付。

```
// CreditCard 支付方式
type CreditCard struct {
    CardNumber string
    CardHolder string
}

// Pay 实现 PaymentProcessor 接口的 Pay 方法
func (cc CreditCard) Pay(amount float64) string {
    return fmt.Sprintf("Paid %.2f using Credit Card: %s", amount, cc.CardNumber)
}

// PayPal 支付方式
type PayPal struct {
    Email string
}

// Pay 实现 PaymentProcessor 接口的 Pay 方法
func (pp PayPal) Pay(amount float64) string {
    return fmt.Sprintf("Paid %.2f using PayPal: %s", amount, pp.Email)
}
```

创建一个处理支付的函数 `ProcessPayment`，该函数接受一个 `PaymentProcessor` 接口类型和支付金额。

```
// ProcessPayment 处理支付
func ProcessPayment(p PaymentProcessor, amount float64) {
    result := p.Pay(amount)
    fmt.Println(result)
}
```

面向接口编程对扩展性支持较好，要添加新的支付方式（例如，Apple Pay），只需实现 `PaymentProcessor` 接口，而无需修改现有代码。

# 四、传统继承

在传统的面向对象语言中，子类可以继承父类并重写父类的方法。当子类的实例调用方法时，如果子类有自己的实现，则优先调用子类的实现；如果没有，则调用父类的实现。

在 Go 语言中，虽然没有传统的继承机制，但可以通过组合和接口实现类似的功能。这种方法主要依赖于以下几点：

1. **接口 (Interface)**：用于定义一组方法，任何实现这些方法的类型都可以被视为该接口的实现。
2. **组合 (Composition)**：通过在结构体中嵌入另一个结构体，来复用已有的逻辑和方法。
3. **方法转发**：父类的方法调用子类的方法，通过接口的方式实现动态的行为选择。

```
// 定义一个 Processor 接口，要求实现 Before 和 Run 方法
type Processor interface {
	Before()
	Run()
}

// 定义一个基础结构体 base，包含一个私有的 run 方法
type base struct {
	name string
}

// base 的私有方法 run，接收一个 Processor 接口
func (b *base) run(p Processor) {
	// 调用传入的 Processor 的 Before 方法
	p.Before()

	// 执行 run 逻辑
	fmt.Println("Base run logic")
}

// base 的 Before 方法，供子类重写
func (b *base) Before() {
	fmt.Println("Base before logic")
}

// 定义一个派生结构体 derived，嵌入 base 并实现 Processor 接口
type derived struct {
	base
}

// derived 实现 Processor 接口的 Run 方法
func (d *derived) Run() {
	// 调用 base 的 run 方法，将 derived 自己传递进去
	d.base.run(d)
}

// 如果 derived 实现了 Before 方法，它将优先于 base 的 Before 方法被调用
func (d *derived) Before() {
	fmt.Println("Derived before logic")
}
```

`base.run()` 方法接收一个实现了 `Processor` 接口的类型，通过调用 `p.Before()`，它可以在运行时选择是调用子类的 `Before()` 还是自己的 `Before()` 方法。

通过这种用法，Go 可以模拟出一种类似于传统继承的行为：

- **方法重写**：子类可以通过实现接口中的方法来重写父类的方法。
- **方法转发**：父类方法可以通过接口调用子类的方法，实现动态行为选择。
- **组合与继承**：组合允许子类复用父类的逻辑，同时通过接口实现动态的多态行为。

# 五、函数转接口

接口定义了一个或多个方法的集合，而在某些情况下，你可能希望将一个函数转换为一个实现了某个接口的类型。这可以通过定义一个带有接口方法签名的函数类型来实现，然后将该函数类型的实例传递给接口。

### 示例：将函数转换为接口

假设你有一个 `Handler` 接口，定义了一个 `Serve` 方法。你可以将一个函数转换为实现 `Handler` 接口的类型。

```
package main

import "fmt"

// 定义 Handler 接口
type Handler interface {
    Serve(data string)
}

// 定义一个函数类型，签名与接口方法一致
type HandlerFunc func(string)

// 实现接口的方法，将该函数类型转换为接口
func (f HandlerFunc) Serve(data string) {
    f(data)
}

// 一个具体的函数
func myHandler(data string) {
    fmt.Println("Handling data:", data)
}

func main() {
    // 将函数 myHandler 转换为 Handler 接口类型
    var h Handler = HandlerFunc(myHandler)

    // 调用接口方法
    h.Serve("Example data")
}

```

可以将现有的函数转换为接口类型，而不需要定义新的结构体，增加了代码的灵活性。这种模式可以简化代码，减少为每个函数定义新类型的需要。例如，Go 的标准库中，`http.HandlerFunc` 就是使用这种模式来将函数转换为 `http.Handler` 接口的。

# 六、回调/闭包函数

回调函数模式是一种常用的编程模式，允许将某些行为或操作以函数的形式传递给另一个函数，以便在特定事件或操作发生时调用。这种模式在许多编程语言中都有广泛的应用，包括 Go 语言。回调函数模式可以提高代码的灵活性和可重用性，因为它将操作与控制流程分离。

下面是一个简单的示例，展示如何使用回调函数模式来处理一组数据。我们假设有一组整数数组，我们希望遍历数组，并对每个元素执行一些操作。

```
// 定义回调函数类型
type ProcessFunc func(int) int

// 遍历并处理数组的函数
func ProcessArray(arr []int, process ProcessFunc) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = process(v)
	}
	return result
}

func main() {
	// 原始数组
	arr := []int{1, 2, 3, 4, 5}
	// 处理数组并获得结果
	result := ProcessArray(arr, func(n int) int {
		return n * 2
	})
}
```

闭包函数允许函数访问其外部作用域中的变量。通过闭包，我们可以创建一个函数，它不仅仅执行自己的任务，还可以捕获并利用在其创建时的上下文中的变量

```
// counter 是一个返回函数的函数，返回的函数会保存对 count 的引用
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	// 创建一个新的计数器函数
	c := counter()

	// 调用计数器函数
	fmt.Println(c()) // 输出 1
	fmt.Println(c()) // 输出 2
	fmt.Println(c()) // 输出 3
}
```

# 七、依赖注入

依赖注入用于将依赖关系通过构造函数或 setter 方法传递给对象，而不是在对象内部创建这些依赖关系。这种模式有助于提高代码的可测试性和灵活性。

```
package main

import "fmt"

type Service struct {
	Repository Repository
}

type Repository interface {
	Find(id int) string
}

type UserRepository struct{}

func (r UserRepository) Find(id int) string {
	return fmt.Sprintf("User %d", id)
}

func NewService(repo Repository) *Service {
	return &Service{Repository: repo}
}

func main() {
	repo := UserRepository{}
	service := NewService(repo)
	fmt.Println(service.Repository.Find(1))
}
```

依赖注入模式能够有效地解耦模块，使模块之间的依赖关系更加清晰，并且更容易进行单元测试。golang中可以借助wire来进行依赖注入。

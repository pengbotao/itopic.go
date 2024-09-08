```
{
    "url": "go-stdlib",
    "time": "2020/06/04 19:00",
    "tag": "Golang",
    "toc": "yes"
}
```

Golang 的标准库非常丰富，涵盖了从基础数据类型到高级并发控制的各个方面。以下是一些常用的标准库，按照使用频率从高到低介绍，并给出详细介绍。

### 1. `fmt`

- **用途**: 格式化 I/O 操作，包括打印和格式化字符串。

- 常用函数:

  - `Println`, `Printf`, `Sprintf`: 打印格式化输出。
  - `Errorf`: 创建一个格式化的错误。

- 示例:

  ```
  fmt.Println("Hello, World!")
  fmt.Printf("Number: %d\n", 42)
  ```

### 2. `os`

- **用途**: 提供与操作系统交互的接口，如文件操作、环境变量、命令行参数等。

- 常用函数:

  - `Open`, `Create`, `ReadFile`, `WriteFile`: 文件操作。
  - `Getenv`, `Setenv`: 获取和设置环境变量。
  - `Exit`: 终止程序。

- 示例:

  ```
  file, err := os.Open("file.txt")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()
  ```

### 3. `io` 和 `io/ioutil`

- **用途**: 提供基本的 I/O 接口和常用 I/O 操作，如文件读写、流复制等。

- 常用函数:

  - `Copy`, `CopyN`: 复制数据流。
  - `ReadAll`, `ReadFile`, `WriteFile`: 文件和流的读写操作。

- 示例:

  ```
  data, err := ioutil.ReadFile("file.txt")
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println(string(data))
  ```

### 4. `log`

- **用途**: 提供简单的日志记录功能。

- 常用函数:

  - `Println`, `Printf`, `Fatal`, `Panic`: 日志输出。
  - `SetOutput`: 设置日志输出的目标。

- 示例:

  ```
  log.Println("This is a log message")
  log.Fatalf("Error: %v", err)
  ```

### 5. `strings`

- **用途**: 提供对字符串的操作函数，如分割、拼接、替换等。

- 常用函数:

  - `Split`, `Join`, `Replace`, `Contains`, `HasPrefix`, `HasSuffix`: 字符串处理。

- 示例:

  ```
  parts := strings.Split("a,b,c", ",")
  joined := strings.Join(parts, "-")
  fmt.Println(joined) // "a-b-c"
  ```

### 6. `strconv`

- **用途**: 提供字符串与基本数据类型之间的转换功能。

- 常用函数:

  - `Atoi`, `Itoa`: 字符串与整数的相互转换。
  - `ParseBool`, `ParseFloat`, `ParseInt`: 字符串解析为布尔值、浮点数、整数。
  - `FormatBool`, `FormatFloat`, `FormatInt`: 格式化布尔值、浮点数、整数为字符串。

- 示例:

  ```
  num, err := strconv.Atoi("123")
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println(num)
  ```

### 7. `time`

- **用途**: 提供时间处理功能，如时间获取、格式化、解析、计时等。

- 常用函数:

  - `Now`, `Sleep`, `Parse`, `Format`: 获取当前时间、暂停执行、时间解析与格式化。
  - `After`, `NewTicker`, `NewTimer`: 定时器和计时器。

- 示例:

  ```
  now := time.Now()
  fmt.Println(now.Format("2006-01-02 15:04:05"))
  ```

### 8. `math/rand`

- **用途**: 提供伪随机数生成功能。

- 常用函数:

  - `Int`, `Intn`, `Float64`: 生成随机整数、浮点数。
  - `Seed`: 设置随机数种子。

- 示例:

  ```
  rand.Seed(time.Now().UnixNano())
  fmt.Println(rand.Intn(100)) // 生成0-99之间的随机数
  ```

### 9. `sync`

- **用途**: 提供基础的并发控制原语，如互斥锁、等待组等。

- 常用类型:

  - `Mutex`, `RWMutex`: 互斥锁与读写锁。
  - `WaitGroup`: 等待一组 goroutine 完成。
  - `Once`: 确保某个操作只执行一次。

- 示例:

  ```
  var mu sync.Mutex
  mu.Lock()
  // critical section
  mu.Unlock()
  ```

### 10. `net/http`

- **用途**: 提供 HTTP 客户端和服务器的实现，是构建 Web 服务的基础库。

- 常用类型和函数:

  - `ListenAndServe`: 启动一个 HTTP 服务器。
  - `Get`, `Post`: 发起 HTTP 请求。
  - `HandleFunc`, `ServeMux`: 注册路由和处理函数。

- 示例:

  ```
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "Hello, World!")
  })
  log.Fatal(http.ListenAndServe(":8080", nil))
  ```

### 11. `encoding/json`

- **用途**: 提供 JSON 数据的编码与解码功能。

- 常用函数:

  - `Marshal`, `Unmarshal`: 结构体与 JSON 字符串的相互转换。
  - `NewEncoder`, `NewDecoder`: 处理 JSON 数据流。

- 示例:

  ```
  type Person struct {
      Name string `json:"name"`
      Age  int    `json:"age"`
  }
  
  person := Person{Name: "Alice", Age: 30}
  data, _ := json.Marshal(person)
  fmt.Println(string(data))
  ```

### 12. `flag`

- **用途**: 提供命令行参数解析功能。

- 常用函数:

  - `String`, `Int`, `Bool`: 定义命令行标志。
  - `Parse`: 解析命令行参数。

- 示例:

  ```
  name := flag.String("name", "World", "a name to say hello to")
  flag.Parse()
  fmt.Printf("Hello, %s!\n", *name)
  ```

### 13. `regexp`

- **用途**: 提供正则表达式匹配和替换功能。

- 常用函数:

  - `MatchString`, `FindString`, `ReplaceAllString`: 字符串的正则匹配、查找、替换。
  - `Compile`: 编译正则表达式。

- 示例:

  ```
  re := regexp.MustCompile(`\d+`)
  fmt.Println(re.FindString("abc123def")) // "123"
  ```

### 14. `path/filepath`

- **用途**: 提供对文件路径的操作功能，如路径解析、拼接、查找等。

- 常用函数:

  - `Join`, `Split`, `Abs`, `Rel`: 路径操作。
  - `Walk`: 遍历文件目录树。

- 示例:

  ```
  path := filepath.Join("dir", "file.txt")
  fmt.Println(path)
  ```

### 15. `bytes`

- **用途**: 提供对字节切片的操作功能，如拼接、查找、替换等。

- 常用函数:

  - `Buffer`, `NewBuffer`, `Read`, `Write`: 字节缓冲区操作。
  - `Join`, `Split`, `Contains`: 字节切片操作。

- 示例:

  ```
  buffer := bytes.NewBufferString("Hello")
  buffer.WriteString(", World!")
  fmt.Println(buffer.String())
  ```

这些标准库涵盖了 Go 语言日常开发中的大多数需求。掌握这些库的使用，可以大大提高开发效率，并写出简洁、优雅的 Go 代码。

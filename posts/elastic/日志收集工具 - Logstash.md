```
{
    "url": "logstash",
    "time": "2021/12/01 22:15",
    "tag": "logstash",
    "toc": "yes"
}
```

# 一、概述

## 1.1 关于Logstash

![](../../static/uploads/logstash-pipeline.png)

> Logstash 是免费且开放的服务器端数据处理管道，能够从多个来源采集数据，转换数据，然后将数据发送到您最喜欢的“存储库”中。

从上图中可以看到Logstash的主要作用就是进行数据规整转换。主要模块有：

- 数据输入：支持
- 数据过滤：对数据进行结构化处理、过滤等。
- 数据输出：将规整之后的数据写到Elasticsearch、文件等

## 1.2 操作流程

- 安装Logstash：[Download Logstash Free](https://www.elastic.co/cn/downloads/logstash)
- 定义配置文件：通过该文件配置数据源、可以进行数据处理并输出到对应存储中。
- 启动收集程序：`bin/logstash -f logstash.conf`

## 1.3 示例 - 基础

创建配置文件: `logstash.conf`，表示从标准输出接受数据，并输出到标准输出，中间没有做数据过滤。

```
input {
    stdin { }
}

filter {}

output {
    stdout { }
}
```

执行：`logstash -f logstash.conf，待程序启动后即可在命令行输入并回车，可以看到类似如下输出：

```
Hello World

{
    "message" => "Hello World\r", 
    "host" => "peng",
    "@timestamp" => 2021-12-01T22:30:50.599Z,
    "@version" => "1"
}
```

## 1.4 示例 - 解析JSON

前面例子太简单，假设源数据为JSON格式，则解析JSON数据源只需要增加一条：

```
input {
    stdin { 
        codec => json
    }
}
```

重新输入json数据可以看到JSON数据已经被解析过了。

```
{"Name": "Peng", "Age": 18} 

{
          "host" => "peng",
           "Age" => 18,
          "Name" => "Peng",
      "@version" => "1",
    "@timestamp" => 2021-12-01T22:32:17.666Z
} 
```

## 1.5 示例 - 解析Nginx

来个相对复杂点的，解析Nginx日志，默认日志格式如下：

```
log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
```

还是从标准输入获取数据源，只记录GET请求，同时输出到标准输出和文件。

```
input {
    stdin { 
        
    }
}

filter {
    grok {
        match => { "message" => "%{IPORHOST:remote_addr} - %{HTTPDUSER:remote_user} \[%{HTTPDATE:time_local}\] \"(?:%{WORD:method} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion})?|%{DATA:rawrequest})\" %{NUMBER:status} (?:%{NUMBER:body_bytes}|-) %{QS:referrer} %{QS:user_agent} %{QS:x_forward_for}" }
        remove_field => "message"
    }
    if [method] != "GET" {
        drop{}
    }
}

output {
    stdout { }
    file {
        path => "./demo.log"
    }
}
```

测试结果如下：

```
172.19.0.1 - - [01/Dec/2021:22:35:27 +0000] "GET /test.html HTTP/1.1" 200 12 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36" "-"

{
       "time_local" => "01/Dec/2021:22:35:27 +0000",
           "method" => "GET",
    "x_forward_for" => "\"-\"",
       "body_bytes" => "12",
       "user_agent" => "\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36\"",
         "referrer" => "\"-\"",
      "remote_user" => "-",
             "host" => "peng",
      "remote_addr" => "172.19.0.1",
           "status" => "200",
       "@timestamp" => 2021-12-02T02:28:30.911Z,
      "httpversion" => "1.1",
          "request" => "/test.html",
         "@version" => "1"
}
```

到这里对Logstash的基本使用应该有个大致的了解：基本流程只有`输入-处理-输出`。但这里把这个流程玩出花了，可以从多个数据源获取数据，针对不同的数据源做不同的格式处理、过滤，最后也可以输出到多个地方。详细文档可查看：[Logstash Reference](https://www.elastic.co/guide/en/logstash/current/index.html)

# 二、数据源

官网文档上列了不下50种数据源的计入方式，可见Logstash支持的数据源有多强大。

## 2.1 文件输入

通过File Input 插件进行收集，示例：

```
input {
  file {
    path => "/data/logs/nginx.log"
    start_position => "beginning"
  }
}
```

## 2.2 Kafka输入





# 三、过滤器

## 3.1 Json

```
filter {
    json {
        source => "message"
    }
}
```

## 3.2 grok





# 四、输出

## 4.1 输出到文件



## 4.2 输出到Elasticsearch







---

- [1] [Input Plugins](https://www.elastic.co/guide/en/logstash/current/input-plugins.html)
- [2] [Filter Plugins](https://www.elastic.co/guide/en/logstash/current/filter-plugins.html)
- [3] [Output Plugins](https://www.elastic.co/guide/en/logstash/current/output-plugins.html)
- [4] [Transforming and sending Nginx log data to Elasticsearch using Filebeat and Logstash - Part 1](https://krakensystems.co/blog/2018/logstash-nginx-logs-part-1)
- [5] [Logstash 最佳实践](https://doc.yonyoucloud.com/doc/logstash-best-practice-cn/index.html)

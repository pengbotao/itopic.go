```
{
    "url": "logstash",
    "time": "2021/12/01 22:15",
    "tag": "logstash,ELK",
    "toc": "yes"
}
```

# 一、概述

## 1.1 关于Logstash

![](../../static/uploads/logstash-pipeline.png)

> Logstash 是免费且开放的服务器端数据处理管道，能够从多个来源采集数据，转换数据，然后将数据发送到您最喜欢的“存储库”中。

从上图中可以看到Logstash的主要作用就是进行数据规整转换。主要模块有：

- 数据输入：支持不下50种数据接入，常见的如文件、Redis、Kafka等。
- 数据过滤：对数据进行结构化处理、过滤等。
- 数据输出：将规整之后的数据写到Elasticsearch、文件等。

## 1.2 操作流程

- 安装Logstash：[Download Logstash Free](https://www.elastic.co/cn/downloads/logstash)
- 定义配置文件：通过该文件配置数据源、可以进行数据处理并输出到对应存储中。
- 启动收集程序：`bin/logstash -f logstash.conf`

# 二、常见示例

## 2.1 示例 - 基础

创建配置文件: `logstash.conf`，表示从标准输出接受数据，并输出到标准输出，中间没有做数据过滤。

```json
input {
    stdin { }
}

filter {}

output {
    stdout { }
}
```

执行：`logstash -f logstash.conf，待程序启动后即可在命令行输入并回车，可以看到类似如下输出：

```json
Hello World

{
    "message" => "Hello World\r", 
    "host" => "peng",
    "@timestamp" => 2021-12-01T22:30:50.599Z,
    "@version" => "1"
}
```

## 2.2 示例 - 解析JSON

前面例子太简单，假设源数据为JSON格式，则解析JSON数据源只需要增加一条：

```json
input {
    stdin { 
        codec => json
    }
}
```

重新输入json数据可以看到JSON数据已经被解析过了。

```json
{"Name": "Peng", "Age": 18} 

{
          "host" => "peng",
           "Age" => 18,
          "Name" => "Peng",
      "@version" => "1",
    "@timestamp" => 2021-12-01T22:32:17.666Z
} 
```

## 2.3 示例 - Nginx日志

来个相对复杂点的，解析Nginx日志，默认日志格式如下：

```
log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
```

为方便测试还是从标准输入获取数据源，只记录GET请求，同时输出到标准输出和文件。

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

到这里对Logstash的基本使用应该有个大致的了解：基本流程只有`输入-处理-输出`。但这里把这个流程玩出花了，可以从多个数据源获取数据，针对不同的数据源做不同的格式处理、过滤，最后也可以输出到多个地方，详细文档可查看：[Logstash Reference](https://www.elastic.co/guide/en/logstash/current/index.html)

# 三、输入

## 3.1 文件

通过File Input 插件进行收集，示例：

```json
input {
  file {
    path => "/data/logs/nginx.log"
    start_position => "beginning"
  }
}
```

## 3.2 Kafka

可以通过启动多个Logstash增加吞吐量，前提条件：

- 多个Logstash属于同一个group_id，即group_id要相同
- 有多个分区，让更多的示例可以消费

```json
input {
    kafka {
        bootstrap_servers => "192.168.0.100:9092"
        group_id => "http-nginx-log"
        client_id => "logstash-nginx-01"
        auto_offset_reset => "latest"
        topics => ["http-nginx-log"]
        codec => json {charset => "UTF-8"}
        type => "nginx"
    }

    kafka {
        bootstrap_servers => "192.168.0.101:9092"
        group_id => "http-apache-log"
        client_id => "logstash-apache-01"
        auto_offset_reset => "latest"
        topics => ["http-nginx-log"]
        codec => json {charset => "UTF-8"}
        type => "apache"
    }
}

```

这里增加了一个type参数，用于演示`5.1`的分开输出，可以去掉。

# 四、过滤器

## 4.1 Json

```json
filter {
    json {
        source => "message"
        remove_field => "message"
    }
}
```

## 4.2 grok

grok支持标准的正则，对于不规整的数据格式可以用它来规整。同时，它内置了120多种可以直接用的匹配规则。比如他也定义了HTTP日志格式，用此就不需要写前面那一长串了，当然前提是和他定义的格式得匹配。

```
filter {
    grok {
        match => { "message" => "%{HTTPD_COMMONLOG}" }
        remove_field => "message"
    }
}
```

也有一些grok在线测试匹配的站点，更多匹配规则可以参考：https://github.com/logstash-plugins/logstash-patterns-core/tree/main/patterns

## 4.3 mutate

它提供了丰富的基础类型数据处理能力。包括类型转换，字符串处理和字段处理等。比如下面示例：按 `.` 切割，增加字段，重命名。

```json
filter {
    mutate {
        split => { "hostname" => "." }
        add_field => { "shortHostname" => "%{[hostname][0]}" }
    }

    mutate {
        rename => {"shortHostname" => "hostname"}
    }
}
```

# 五、输出

## 5.1 文件

```json
output{
    if [type]  == "nginx" {
        file {
            path => "/data/logs/nginx-%{+YYYY-MM-dd}-%{host}.log"
            gzip => true
        }
    }

    if [type]  == "apache" {
        file {
            path => "/data/logs/apache-%{+YYYY-MM-dd}-%{host}.log"
            gzip => true
        }
    }
}
```

## 5.2 Elasticsearch

```json
output {
　　elasticsearch{
   　　hosts => ["192.168.0.100:9200"]
      index => "http-log-%{+YYYY-MM-dd}"
   }
}
```

---

- [1] [Input Plugins](https://www.elastic.co/guide/en/logstash/current/input-plugins.html)
- [2] [Filter Plugins](https://www.elastic.co/guide/en/logstash/current/filter-plugins.html)
- [3] [Output Plugins](https://www.elastic.co/guide/en/logstash/current/output-plugins.html)
- [4] [Transforming and sending Nginx log data to Elasticsearch using Filebeat and Logstash - Part 1](https://krakensystems.co/blog/2018/logstash-nginx-logs-part-1)
- [5] [Logstash 最佳实践](https://doc.yonyoucloud.com/doc/logstash-best-practice-cn/index.html)

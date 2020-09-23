```
{
    "url": "httpie",
    "time": "2016/02/04 21:28",
    "tag": "常用软件",
    "toc" : "no"
}
```

# 一、简介
HTTPie是一个命令行下的类似CURL的HTTP请求工具，使用Python开发，相比CURL而言参数更精简、人性化，通过http命令很容易发送一些日常的http请求。

HTTPie的安装可参考：`https://github.com/jakubroztocil/httpie`，安装完成后即可使用`http`命令

# 二、http命令

```
Usage: http [flags] [METHOD] URL [REQUEST_ITEM [REQUEST_ITEM ...]]
```

唯一不可缺少的是`URL`参数，其他参数都有一些默认行为，请求示例：

```
$ http -v itopic.org
GET / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: itopic.org
User-Agent: HTTPie/1.0.2

HTTP/1.1 301 Moved Permanently
Connection: keep-alive
Content-Length: 185
Content-Type: text/html
Date: Tue, 12 Mar 2019 13:52:14 GMT
Location: https://itopic.org/
Server: nginx/1.15.3
```

## 2.1 常用flags定义

## 2.2 METHOD定义

可选的`METHOD`参数，支持`GET`、`POST`、`PUT`、`DELETE`等。
当没有指定该参数时，如果有数据传递则会设置为`POST`，否则会设置为`GET`，如：

```
http -v itopic.org from=httpie
POST / HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 18
Content-Type: application/json
Host: itopic.org
User-Agent: HTTPie/1.0.2

{
    "from": "httpie"
}
```

## 2.3 URL定义

支持`http://`、`https://`，默认不指定时为`http://`，如果是localhost也可以省略为`http :/demo` 等同于 `http http://localhost/demo`

## 2.4 REQUEST_ITEM 定义

像`curl`里通过`-H Content-Type:application/json`来指定HTTP头部信息，`-d`来指定`POST`数据信息。而`http`命令里主要根据一些简单的符号来区分如何传递数据，配合示例就一清二楚了。

序号|符号|说明
---|---|---
2.4.1|:|HTTP头部信息
2.4.2|==|URL参数
2.4.3|=|键值对，用来传输`json`数据或者表单数据。默认会当json处理，指定`-f`时会当表单数据提交。
2.4.4|:=|非JSON字符串类型的数据。如：`awesome:=true  amount:=42  colors:='["red", "green", "blue"]'`。只在`json`传输起作用。
2.4.5|@|上传文件


### 2.4.1 指定头部示例

```
$ http -v itopic.org X-Powered-By:itopic
GET / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: itopic.org
User-Agent: HTTPie/1.0.2
X-Powered-By: itopic
```

### 2.4.2 设置Query参数

```
$ http -v itopic.org from==tools
GET /?from=tools HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: itopic.org
User-Agent: HTTPie/1.0.2
```

### 2.4.3 设置请求数据

不指定类型时默认为`json`数据类型。因为有提交数据，默认请求会用`POST`方式请求，也可以显示指定。

```
$ http -v itopic.org from=tools name=httpie
POST / HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 35
Content-Type: application/json
Host: itopic.org
User-Agent: HTTPie/1.0.2

{
    "from": "tools",
    "name": "httpie"
}
```

**指定`-f`(`--form`)时以表单类型提交。**

```
$ http -f -v itopic.org from=tools name=httpie
POST / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 22
Content-Type: application/x-www-form-urlencoded; charset=utf-8
Host: itopic.org
User-Agent: HTTPie/1.0.2

from=tools&name=httpie
```

### 2.4.4 设置JSON数组

因为`http`默认以JSON格式传递，所以省去了显示`-j`(`--json`)指定。通过下面的方式可以指定一些简单的json格式。

```
$ http -v itopic.org from:='["app", "web"]' name=itopic show:=true
POST / HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 56
Content-Type: application/json
Host: itopic.org
User-Agent: HTTPie/1.0.2

{
    "from": [
        "app",
        "web"
    ],
    "name": "itopic",
    "show": true
}
```

# 三、Curl请求

```
$ curl -w "@curl-format.txt" -o /dev/null -s http://www.baidu.com
              http: 200
     time_namelookup: 0.001757s
       time_redirect: 0.000000s
        time_connect: 0.024540s
     time_appconnect: 0.000000s
    time_pretransfer: 0.024580s
  time_starttransfer: 0.044848s
          time_total: 0.045182s

     size_download: 2381bytes
      speed_download: 52911.000B/s
```
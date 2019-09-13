```
{
    "url": "requests",
    "time": "2017/04/20 23:11",
    "tag": "Python,爬虫"
}
```

# 一、发送请求

使用`Requests`发送网络请求非常简单。一开始要导入`Requests`模块：

```
>>> import requests
```

然后，尝试获取某个网页。本例子中，我们来获取`Github`的公共时间线：

## 1.1 GET请求

```
>>> r = requests.get('https://api.github.com/events', params={})
```

现在，我们有一个名为`r`的`Response`对象。我们可以从这个对象中获取所有我们想要的信息。`Requests`简便的`API`意味着所有`HTTP`请求类型都是显而易见的。例如，你可以这样发送一个 **HTTP POST** 请求：

## 1.2 POST请求
```
>>> r = requests.post('http://httpbin.org/post', data = {'key':'value'})
```

漂亮，对吧？那么其他`HTTP`请求类型：`PUT`，`DELETE`，`HEAD`以及`OPTIONS`又是如何的呢？都是一样的简单：

```
>>> r = requests.put('http://httpbin.org/put', data = {'key':'value'})
>>> r = requests.delete('http://httpbin.org/delete')
>>> r = requests.head('http://httpbin.org/get')
>>> r = requests.options('http://httpbin.org/get')
```

# 二、传递URL参数

你也许经常想为`URL`的查询字符串(`query string`)传递某种数据。如果你是手工构建 `URL`，那么数据会以键/值对的形式置于`URL`中，跟在一个问号的后面。例如， `httpbin.org/get?key=val`。 `Requests`允许你使用`params`关键字参数，以一个字符串字典来提供这些参数。举例来说，如果你想传递`key1=value1`和`key2=value2`到 `httpbin.org/get`，那么你可以使用如下代码：

```
>>> payload = {'key1': 'value1', 'key2': 'value2'}
>>> r = requests.get("http://httpbin.org/get", params=payload)
```

通过打印输出该`URL`，你能看到`URL`已被正确编码：

```
>>> print(r.url)
http://httpbin.org/get?key2=value2&key1=value1
```

注意字典里值为`None`的键都不会被添加到`URL`的查询字符串里。你还可以将一个列表作为值传入：

```
>>> payload = {'key1': 'value1', 'key2': ['value2', 'value3']}

>>> r = requests.get('http://httpbin.org/get', params=payload)
>>> print(r.url)
http://httpbin.org/get?key1=value1&key2=value2&key2=value3
```

# 三、响应内容

我们能读取服务器响应的内容。再次以 GitHub 时间线为例：

## 3.1 响应内容
```
>>> import requests
>>> r = requests.get('https://api.github.com/events')
>>> r.text
u'[{"repository":{"open_issues":0,"url":"https://github.com/...
```

`Requests`会自动解码来自服务器的内容。大多数`unicode`字符集都能被无缝地解码。

请求发出后，`Requests`会基于`HTTP`头部对响应的编码作出有根据的推测。当你访问 `r.text`之时，`Requests`会使用其推测的文本编码。你可以找出`Requests`使用了什么编码，并且能够使用`r.encoding`属性来改变它：

```
>>> r.encoding
'utf-8'
>>> r.encoding = 'ISO-8859-1'
```

如果你改变了编码，每当你访问`r.text`，`Request`都将会使用`r.encoding`的新值。你可能希望在使用特殊逻辑计算出文本的编码的情况下来修改编码。比如`HTTP`和`XML`自身可以指定编码。这样的话，你应该使用`r.content`来找到编码，然后设置`r.encoding`为相应的编码。这样就能使用正确的编码解析`r.text`了。

在你需要的情况下，`Requests`也可以使用定制的编码。如果你创建了自己的编码，并使用 `codecs`模块进行注册，你就可以轻松地使用这个解码器名称作为`r.encoding`的值， 然后由`Requests`来为你处理编码。

## 3.2 响应状态码
```
>>> r.status_code
200
```


# 四、请求头与响应头

## 4.1 查看请求头

```
>>> r = requests.post('http://httpbin.org/post', data = {'key':'value'})
>>> r.request.headers
{
	'Content-Length': '9',
	'Accept-Encoding': 'gzip, deflate',
	'Accept': '*/*',
	'User-Agent': 'python-requests/2.22.0',
	'Connection': 'keep-alive',
	'Content-Type': 'application/x-www-form-urlencoded'
}
>>> r.request.body
'key=value'
```

## 4.2 查看返回头

```
>>> r.headers
{
	'Content-Length': '266',
	'X-XSS-Protection': '1; mode=block',
	'X-Content-Type-Options': 'nosniff',
	'Content-Encoding': 'gzip',
	'Server': 'nginx',
	'Connection': 'keep-alive',
	'Access-Control-Allow-Credentials': 'true',
	'Date': 'Fri, 13 Sep 2019 12:44:56 GMT',
	'Access-Control-Allow-Origin': '*',
	'Referrer-Policy': 'no-referrer-when-downgrade',
	'Content-Type': 'application/json',
	'X-Frame-Options': 'DENY'
}
```

## 4.3 设置请求头

如果你想为请求添加 HTTP 头部，只要简单地传递一个 dict 给 headers 参数就可以了。

```
>>> r = requests.get('http://httpbin.org', headers={"User-Agent": "my-app/0.0.1"})
>>> r.request.headers
{'Connection': 'keep-alive', 'Accept-Encoding': 'gzip, deflate', 'Accept': '*/*', 'User-Agent': 'my-app/0.0.1'}
```

# 五、常用参数设置

## 5.1 超时设置

为防止服务器不能及时响应，大部分发至外部服务器的请求都应该带着`timeout`参数。在默认情况下，除非显式指定了`timeout`值，`requests`是不会自动进行超时处理的。如果没有`timeout`，你的代码可能会挂起若干分钟甚至更长时间。

连接超时指的是在你的客户端实现到远端机器端口的连接时（对应的是`connect()`_），`Request` 会等待的秒数。一个很好的实践方法是把连接超时设为比 3 的倍数略大的一个数值，因为 TCP 数据包重传窗口 (TCP packet retransmission window) 的默认大小是 3。

一旦你的客户端连接到了服务器并且发送了`HTTP`请求，读取超时指的就是客户端等待服务器发送请求的时间。（特定地，它指的是客户端要等待服务器发送字节之间的时间。在`99.9%`的情况下这指的是服务器发送第一个字节之前的时间）。

如果你制订了一个单一的值作为`timeout`，如下所示：

```
r = requests.get('https://github.com', timeout=5)
```
这一`timeout`值将会用作`connect`和`read`二者的`timeout`。如果要分别制定，就传入一个元组：

```
r = requests.get('https://github.com', timeout=(3.05, 27))
```
如果远端服务器很慢，你可以让`Request`永远等待，传入一个`None`作为`timeout`值，然后就冲咖啡去吧。
```
r = requests.get('https://github.com', timeout=None)
```

## 5.2 代理设置

如果需要使用代理，你可以通过为任意请求方法提供`proxies`参数来配置单个请求:

```
import requests

proxies = {
  "http": "http://10.10.1.10:3128",
  "https": "http://10.10.1.10:1080",
}

requests.get("http://example.org", proxies=proxies)
```

你也可以通过环境变量`HTTP_PROXY`和`HTTPS_PROXY`来配置代理。

```
$ export HTTP_PROXY="http://10.10.1.10:3128"
$ export HTTPS_PROXY="http://10.10.1.10:1080"

$ python
>>> import requests
>>> requests.get("http://example.org")
```

若你的代理需要使用`HTTP Basic Auth`，可以使用 `http://user:password@host/` 语法：

```
proxies = {
    "http": "http://user:pass@10.10.1.10:3128/",
}
```

要为某个特定的连接方式或者主机设置代理，使用`scheme://hostname`作为`key`， 它会针对指定的主机和连接方式进行匹配。
```
proxies = {'http://10.20.1.128': 'http://10.10.1.10:5323'}
```
注意，代理`URL`必须包含连接方式。





文档：https://2.python-requests.org//zh_CN/latest/user/quickstart.html
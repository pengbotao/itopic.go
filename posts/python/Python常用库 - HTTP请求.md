```
{
    "url": "python-requests",
    "time": "2016/04/07 20:00",
    "tag": "Python,爬虫,HTTP",
    "public": "yes"
}
```

# HTTP请求 - requests库

requests库的文档还比较齐全，可以参考

- 文档地址：http://docs.python-requests.org/zh_CN/latest/user/quickstart.html
- 高级用法：http://docs.python-requests.org/zh_CN/latest/user/advanced.html

## GET请求

```
import requests
```

设置Query参数。会自动拼接，形式可参考返回值。也可以传递完成的`GET`请求地址。

```
payload = {
    'from': 'python',
}
```

设置HTTP请求头

```
headers = {
    'X-Powered-By': 'iTopic',

}
```

设置代理

```
proxies = {
    'http': '',
    'https': '',
}
```

设置Cookie，也可以放在headers中

```
cookies = {
    'test': '1',
}
```

发送请求

```
r = requests.get("https://itopic.org/index.html", params=payload, 
                 headers=headers, 
                 proxies=proxies, 
                 cookies=cookies,
                 timeout=10)
```


打印请求URL

```
print(r.url)

https://itopic.org/index.html?from=python
```

打印请求头

```
print(r.request.headers)

{
    'Accept-Encoding': 'gzip, deflate', 
    'X-Powered-By': 'iTopic', 
    'Accept': '*/*', 
    'User-Agent': 'python-requests/2.18.4', 
    'Connection': 'keep-alive', 
    'Cookie': 'test=1'
}
```

打印HTTP状态码

```
print(r.status_code)

200
```

打印返回头

```
print(r.headers)

{
    'Date': 'Thu, 21 Mar 2019 13:50:41 GMT', 
    'Transfer-Encoding': 'chunked', 
    'Connection': 'keep-alive', 
    'Content-Type': 'text/html; charset=UTF-8', 
    'Server': 'nginx/1.15.3'
}
```

打印内容

```
# print(r.cookies)
# print(r.text)
```


## POST请求

post实际用法和上面post一致。

```
>>> import requests
>>> payload = {'key1': 'value1', 'key2': 'value2'}
>>> r = requests.post("http://httpbin.org/post", data=payload)
>>> print(r.text)
{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "key1": "value1",
    "key2": "value2"
  },
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip, deflate",
    "Content-Length": "23",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "httpbin.org",
    "User-Agent": "python-requests/2.18.4"
  },
  "json": null,
  "origin": "27.18.253.121, 27.18.253.121",
  "url": "https://httpbin.org/post"
}
```

# HTTP请求 - urllib/urllib2

## URL编码/解码

```
import urllib

data = {
    "from": "web",
    "remark": "中文",
}

# 字典编码
print(urllib.urlencode(data))
# 解码
print(urllib.unquote(urllib.urlencode(data)))

# 字符串编码
print(urllib.quote(":"))
print(urllib.unquote("%3A"))
```
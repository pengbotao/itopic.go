```
{
    "url": "http-get-length",
    "time": "2014/08/05 16:13",
    "tag": "PHP,HTTP"
}
```

曾做过一个plist的接口，需要将参数编码后放在URL上传递，编码后的URL很长，长到让人担心这么传有没有问题？要弄清这个问题得先弄明白HTTP报文请求格式，借用网上一张图片：

![](/static/uploads/http-proxy.png)

一个HTTP请求由四部分组成：Request Line、Headers、空行和Request Body。见下面GET示例：

![](/static/uploads/http_get.png)

可以看到GET请求时，数据放在Request Line中Request-URI传递；若发送POST请求时，数据则放在Request Body中传递，不同的地方决定了其享受不同的待遇。

HTTP手册上是这么说的：

HTTP协议不会对URI做任何限制，服务端必须能处理不限长度的URI，如果不能处理则返回414

> The HTTP protocol does not place any a priori limit on the length of a URI. Servers MUST be able to handle the URI of any resource they serve, and SHOULD be able to handle URIs of unbounded length if they provide GET-based forms that could generate such URIs. A server SHOULD return 414 (Request-URI Too Long) status if a URI is longer than the server can handle (see section 10.4.15).

> http://www.w3.org/Protocols/rfc2616/rfc2616-sec3.html#sec3.2.1

虽然HTTP协议中并未限制Request-URI的大小，但不同浏览器和WEB服务器则有不同的限制。

**浏览器限制**

不同的浏览器限制不同，若通过浏览器访问，2000个字符长度以内基本上都可以兼容。

**WEB服务器**

这里在SAE上写了个脚本，访问地址 http://ipbt.vipsinaapp.com/demo/Http/request.php?len=8151 ，看到可以正常访问，若len参数改为8152则提示`414 Request-URI Too Large`

可见这是当前NGINX环境下的极限值，此时Request Line的总长度为8190。同样的请求本地Apache配置的服务器，当Request Line长度超过8190时，Aapche也会提示

> Request-URI Too Large

> The requested URL's length exceeds the capacity limit for this server.

可见Nginx和Apache在请求行的长度控制上是一致的，默认都是8190个字符长度。但同时该值也是可以调整的，Nginx中可以通过large_client_header_buffers来控制请求行的最大长度；Apache中可设置LimitRequestLine的值来控制。

**LimitRequestLine指令**

- 语法：LimitRequestLine bytes
- 默认：LimitRequestLine 8190

这个指令用来设置客户端发送的HTTP请求行的最大字节数。请求行包括HTTP方法、URL、协议版本等。因此LimitRequestLine指令能够限制URL的长度，服务器会需要这个值足够大以装载它所有的资源名，包括可能在GET请求中所传递的查询部分的所有信息。一般情况下，不要更改这个值，使用默认即可。 

当把本地Apache LimitRequestLine 调整为 8190000后，加3个0在请求也没问题。。。

```
GET /ipbt/1/demo/Http/response.php?s=11111111.....
Host: localhost:80
Connection: Close
 
HTTP/1.1 200 OK
Date: Tue, 05 Aug 2014 23:19:17 GMT
Server: Apache/2.2.11 (Win32) PHP/5.4.7
X-Powered-By: PHP/5.4.7
Content-Length: 14
Connection: close
Content-Type: text/html
 
Length:8000000
```

可见GET请求实际上也是没有长度限制的，只是浏览器和WEB服务器做了控制。。。
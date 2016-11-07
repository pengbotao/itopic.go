```
{
    "url": "session-application",
    "time": "2013/12/22 21:04",
    "tag": "PHP,SESSION"
}
```

为加强对SESSION的理解，这里通过实现OSChina登录并发表一条动态来说明。

首先需要去oschina上查看登录的提交地址和参数，然后通过程序向该地址传递参数，并把COOKIE记录下来，下次读取其他页面使用该COOKIE即可。具体登录流程见下面代码：

```
<?php
$tmp_folder = dirname(__FILE__).'/tmp';
if(! file_exists($tmp_folder)) {
    mkdir($tmp_folder, 0777, true);
}
$cookie_file = tempnam($tmp_folder, "sess");
$post_data = array(
    'email' => 'youremail@163.com',
    'pwd' => hash('sha1', 123456),
    'save_login' => 1,
);
$ch = curl_init("https://www.oschina.net/action/user/hash_login");
curl_setopt($ch, CURLOPT_HEADER, true);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($post_data));
curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, false);
//模拟浏览器头
curl_setopt($ch, CURLOPT_USERAGENT, 'Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25');
//将COOKIE保存在$cookie_file中
curl_setopt($ch, CURLOPT_COOKIEJAR,  $cookie_file);
//curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true);
$os_login_rtn = curl_exec($ch);
if(curl_errno($ch)) {
    throw new Exception('Curl Error:'.curl_error($ch));
}
curl_close($ch);
echo $os_login_rtn;
```

通过CURLOPT_COOKIEJAR将COOKIE保存在当前tmp目录下，并打印显示HTTP头和返回值，获取到的结果是：

```
HTTP/1.1 200 OK
Server: Tengine/1.4.6
Date: Sun, 22 Dec 2013 01:00:51 GMT
Content-Type: text/html;charset=utf-8
Content-Length: 0
Connection: keep-alive
Set-Cookie: oscid=""; Domain=.oschina.net; Expires=Thu, 01-Jan-1970 00:00:10 GMT; Path=/; HttpOnly
Set-Cookie: oscid=LklqaSHEgZygp%2BnwFOGLO7qmtVSOCMGv1Ea4aiyl0o2ybYt2B20bCyafAHYO9dEICMXhndAeBeg%2BwxAuRJIcbTtmAfGA503D1jVWgm1pbXUsR2BEmFMoLOsIy9RU4DtH; Domain=.oschina.net; Expires=Mon, 22-Dec-2014 00:58:05 GMT; Path=/; HttpOnly
```

可以看到HTTP中Content-Length实体长度为0，只是告诉浏览器设置oscid的COOKIE，oscid存在两个Set-Cookie的动作，前面一个的过期时间是1970年，应该是为了先删除该COOKIE再创建。执行成功后可以看到tmp目录下有一个tmp文件，记录了该COOKIE值。

接下来就是使用该COOKIE，要想发送动态需要先知道空间的地址，可以先通过首页获取到空间地址，程序接上面：

```
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "http://www.oschina.net");
curl_setopt($ch, CURLOPT_HEADER, 0);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
//读取cookie file
curl_setopt($ch, CURLOPT_COOKIEFILE, $cookie_file); 
//CURL设置AJAX请求
//curl_setopt($ch, CURLOPT_HTTPHEADER, array("X-Requested-With: XMLHttpRequest"));
$os_index_rtn = curl_exec($ch);
curl_close($ch);
preg_match('#http://my.oschina.net/u/(\d+)#', $os_index_rtn, $matches);
print_r($matches);exit;
```

直接读取oschin首页，根据首页源文件匹配出空间地址。打印如下：

```
Array
(
    [0] => http://my.oschina.net/u/1425813
    [1] => 1425813
)
```

这就表示登录状态已经获取到了，可以进行接下来的发帖、评论等其他事情了，之后的操作跟这里都是一样的，就不给出演示代码了。

知道了程序可以模拟很多操作，这是很危险的。所以在很多地方我们要有这个意识，防止程序批量去注册、登录后做某些恶意操作，最简单的解决方法还是加验证码，验证码用掉后就删掉。

之前听说OSC的SESSION是存在COOKIE中，这里看到的不是SESSION ID，而是一串加密串，应该是存放在COOKIE中。OSC的COOKIE很干净，访问过很多次，但是只有一个oscid的cookie，值得我们学习。OSC的COOKIE设置了HttpOnly属性，一般情况下该属性可以阻止JS读取COOKIE，达到防止xss的目的。

虽然说OSC用COOKIE存放会话状态，但还是不建议这么做。主要还是开发的程序员能力有好有差，很多开发人员开发过程为了方便把很多数据扔在COOKIE中，把所有的COOKIE都设置主域，难以避免乱用的情况。
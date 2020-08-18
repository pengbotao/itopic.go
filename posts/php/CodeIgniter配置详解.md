```
{
    "url": "codeigniter-config",
    "time": "2013/10/28 23:12",
    "tag": "PHP,CodeIgniter"
}
```

# 一、CodeIgniter配置之autoload
CodeIgniter带了自动加载的功能，可以全局加载类库、模型、配置、语言包等，对于需要全局使用的功能相当方便。例如：有个全局函数写在app_helper.php中，需要全局加载这个函数，只需设置autoload.php：
```
$autoload['helper'] = array('app');
```
接下来，所有的地方都可以使用了，配置、模型等配置相似。但方便的同时也需要考虑下该种加载方式有何弊端。如果一个项目中分了两块，如前台、后台，那这个功能是否为前后台都必须？ 如果前后台还有不同的业务模块区分， 是否是每个模块都要用到？如果都需要， 那写在这里就很好， 如果不需要， 就不建议写在这里。

**对于相关的类库、函数调用应该按需加载**

实现加载的方式有很多，可以在指定的页面load， 可以在公用的控制器、函数里面load， 一经load即可全局使用。个人的常用做法是忽略该文件，手动加载全局函数等。说到这里，顺便说下CI的加载机制。下面为类库、函数等的加载方式：
```
$this->load->library('session');

$this->load->model('hello_model');

$this->load->helper(array('url', 'array'));

$this->load->language(array('user_menu', 'user_tips'));
```
加载方式统一，使用起来比较简单，但load类库时传参有点不方便。再次load类库时不会再去加载，而是从保存好的静态数组中拿出来，需要注意下成员属性的状态，防止因为值还存在而导致程序异常。

# 二、CodeIgniter配置之URL
应该有很多项目中会有这样的情况，通过 http://pc.local 可以访问，若通过 http://localhost/pc/public 则会出现一些图片、样式显示不到，超链接出错的情况，项目的移植性不太好，主要原因就是创建的URL不够灵活，接下来看看CI中是怎么处理。配置文件中有几个有关URL的配置，影响到路由、参数的获取和URL的创建，它们是：
```
$config['base_url'] = '';
  
$config['index_page'] = 'index.php';
  
$config['uri_protocol'] = 'AUTO';
  
$config['url_suffix'] = '';
  
$config['allow_get_array']      = TRUE;
  
$config['enable_query_strings'] = FALSE;
  
$config['controller_trigger']   = 'c';
  
$config['function_trigger']     = 'm';
  
$config['directory_trigger']    = 'd';
```
**$config['uri_protocol']**

uri_protocol可选项有 AUTO、PATH_INFO、QUERY_STRING、REQUEST_URI、ORIG_PATH_INFO。含义分别如下：

- QUERY_STRING:查询字符串；
- PATH_INFO：客户端提供的路径信息，即在实际执行脚本后面尾随的内容，会去掉Query String；
- REQUEST_URI：包含HTTP协议中定义的URI内容。

访问： http://pc.local/index.php/product/pc/summary?a=1 时PATH_INFO为/product/pc/summary；REQUEST_URI为`/index.php/product/pc/summary?a=1`；QUERY_STRING为a=1。实际的配置跟服务器配置也会有点关系，有的服务器会配置ORIG_PATH_INFO，大部分没有。

`uri_protocol`的值决定了CI路由和参数的获取方式，CI会通过这些值来判断解析到哪一个控制器，所以需要确保服务器配置了正确的值。大部分情况下设置AUTO即可，AUTO会依次检测`REQUEST_URI` `PATH_INFO` `QUERY_STRING` `$_GET`的值，直到读到内容。

影响路由解析还有`enable_query_strings`参数,当该参数为TRUE时， 并且传入了controller_trigger参数，则会以查询字符串的方式来获取参数，如index.php?c=products&m=view&id=345则解析到product控制器中的view方法。

**$config['allow_get_array']**

CI中封装了$this->input->get()方法来获取get参数，这里设置为true则表示也允许通过$_GET方式来获取参数，否则会清空$_GET中的值。

**$config['base_url']**

该参数的设置会影响到base_url() site_url()函数创建的URL，为空时程序会自动获取当前地址，否则会根据设置的地址进行创建URL。

**$config['index_page']**

默认主页，使用site_url()创建时会带上该参数，若需要去掉index.php设置为空即可。

**$config['url_suffix']**

后缀名， 使用site_url()创建时也会带上该参数。

以上两个参数的设置对base_url()无效，即不会追加 index_page 和 url_suffix的值。
```
//pc.local/bootstrap/js/bootstrap.min.js
echo base_url('bootstrap/js/bootstrap.min.js');
  
//pc.local/login.htm（设置了index_page='' url_suffix='.htm'）
echo site_url('login');
```
所以可以发现，base_url() 用来创建静态资源的地址，如JS、CSS等地址需要用该函数来生成。site_url() 用来创建跟控制器地址相关的URL。使用base_url()和site_url()函数创建的URL地址具有更好的移植性，建议系统URL相关的地址统一使用该函数创建。

CI提供的site_url对于参数的处理还不够灵活，可以自己扩展一个create_url函数用来创建跟控制器相关的地址。下面为一种实现方式。扩展url_helper，在helpers中创建MY_url_helper.php，内容如下：
```
function create_url($route = NULL, $params = array(), $ampersand = '&')
{
    $route = site_url($route);
    if(!empty($params)) {
        return $route.'?'.http_build_str($params, NULL, $ampersand);
    }
    return $route;
}
if(! function_exists('http_build_str'))
{
    function http_build_str($query, $prefix='', $arg_separator='')
    {
        if (!is_array($query)) {
            return null;
        }
        if ($arg_separator == '') {
            $arg_separator = ini_get('arg_separator.output');
        }
        $args = array();
        foreach ($query as $key => $val) {
            $name = $prefix.$key;
            if (!is_numeric($name)) {
                if(is_array($val)){
                    http_build_str_inner($val, $name, $arg_separator, $args);
                }else{
                    $args[] = rawurlencode($name).'='.urlencode($val);
                }
            }
        }
        return implode($arg_separator, $args);
    }
}
  
if(! function_exists('http_build_str_inner'))
{
    function http_build_str_inner($query, $prefix, $arg_separator, &$args)
    {
        if (!is_array($query)) {
            return null;
        }
        foreach ($query as $key => $val) {
            $name = $prefix."[".$key."]";
            if (!is_numeric($name)) {
                if(is_array($val)){
                    http_build_str_inner($val, $name, $arg_separator, $args);
                }else{
                    $args[] = rawurlencode($name).'='.urlencode($val);
                }
            }
        }
    }
}
```

# 三、CodeIgniter配置之XSS和CSRF
XSS：跨站脚本攻击，通常情况下为将一段脚本入库，显示时未做处理，导致脚本执行。大多数情况下获取到COOKIE中保存的SESSION ID就表示攻击成功。比如，在评论时填写
```
<script>alert(document.cookie)</script>
```
如果入库、显示均为处理，则其他人看到该条评论时就会弹出他的COOKIE信息，这里只是弹出可能没有攻击性，但如果将COOKIE信息以参数传回给服务器呢？服务器则可以拿到用户的SESSION信息，从而就获取到了用户状态，具体机制查看《CodeIgniter配置之SESSION》 和 《SESSION机制应用》。

**CSRF**：跨站请求伪造，也就是常说的钓鱼网站，当用户登录A网站之后，会话状态内，用户点了B网站上某个请求地址，该地址发送一个请求给A，告诉A将当前用户余额转到我的账户中，A网站服务器判断是在会话状态内，认为是用户提交的请求，所以成功伪造了一个请求。

XSS和CSRF涉及到用户会话状态和用户操作权限，若收到攻击是比较危险的。接下来看看CI中对XSS和XSRF处理。

![](../../static/uploads/CodeIgniter-crsf-error.png)

先看一下上面的错误提示，你能很快找到问题所在吗？ 

该问题就是因为csrf引起的！当csrf_protection为true时，会对每个POST提交页面进行csrf保护校验，验证不通过就会报上面的错误了。

csrf就好比在每个提交页面设置了一个隐藏的验证码，验证码的名字叫csrf_test_name，值为cookie名csrf_cookie_name的值。当提交时，会把cookie中的值和表单中提交的值进行比较，如果表单没有设置该值或者由于cookie过期，则提示上面错误。所以可以看到全局打开csrf_protection后会导致每个POST提交页面都需要验证并且校验依赖cookie。

创建普通表单时可通过CI提供的form_open函数来创建，但ajax提交呢？就需要手动传入csrf_cookie_name到服务端了。以及cookie会在一定时间后过期，虽然时间默认设置了7200s，但还是存在过期的可能性。难道你愿意因为某些原因让用户看到上面提示？还有个问题可能会想到，可否全局打开，在指定控制器关闭呢？csrf是在初始化Input类之前就调用了，在控制器里更改是无效的，如果要更改可以生效的可能就是前面的钩子类了，但操作起来会相对麻烦，不适合使用该方式。

所以，不建议开启csrf_protection，可以在需要设置的地方使用验证码或自己设置。接下来，说一下xss

当global_xss_filtering设置为true，会对GPC进行过滤，如会将`<script>1</script>`过滤成`[removed]1[removed]`。也就是对用户输入的原始数据进行了修改，这是很严重的。也许用户想要的就是过滤前的数据呢？

比如，用户设置密码为a，而经过过滤处理后变成了b，这将会导致用户的密码在数据库中实际为b，用户每次通过系统登录时，都会先做一次转换工作再登录，或者输入密码b也能登录成功。那如果其他系统也要做登录功能，用户是共享的，那该用户通过密码a登录就失败了。或者也可能系统升级更换框架， 那这批用户也就登录不成功了。这就是潜在的问题，表面上看起来没问题，但可能为以后造成重大影响。所以**只对需要过滤的参数进行过滤，慎重考虑全局参数控制**

# 四、CodeIgniter配置之session
刚使用Codeigniter时也被其中的SESSION迷惑过，后来就再也没用过CI自带的SESSION，想必还是有必要整理一下SESSION。为弄清CI中的SESSION，先来说一下PHP中SESSION是如何工作的。由于HTTP协议本身是无状态的，所以当保留某个用户的访问状态信息时，需要客户端有一个唯一标识传给服务端，这个唯一标识就是SESSION ID，存放在客户端的COOKIE中，然后服务端根据该标识读取存放的用户状态信息，达到保存会话状态的目的。PHP中启动一个会话需要执行：`session_start();`

1、客户端每次请求时会有一些信息存放中HTTP头中发送给服务端，以用户第一次访问为例：
```
Request Headers
Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Encoding:gzip,deflate,sdch
Accept-Language:zh-CN,zh;q=0.8
Cache-Control:max-age=0
Connection:keep-alive
Host:s.local
User-Agent:Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36
```
2、服务端接到请求处理后并返回给客户端，并在HTTP Response中加上添加COOKIE的请求，告诉浏览器需要设置一个COOKIE，COOKIE名为PHPSESSID，值为r887k5n4scg32d4ba34huuhmq7，如：
```
Response Headers
Cache-Control:no-store, no-cache, must-revalidate, post-check=0, pre-check=0
Connection:Keep-Alive
Content-Length:0
Content-Type:text/html
Date:Sun, 08 Dec 2013 12:56:56 GMT
Expires:Thu, 19 Nov 1981 08:52:00 GMT
Keep-Alive:timeout=5, max=100
Pragma:no-cache
Server:Apache/2.2.11 (Win32) PHP/5.4.7
Set-Cookie:PHPSESSID=r887k5n4scg32d4ba34huuhmq7; path=/
X-Powered-By:PHP/5.4.7
```
3、当客户端再次访问该网站的页面时，浏览器会将该COOKIE发送给服务端，服务端根据COOKIE的值去读取服务器上存放SESSION的文件，拿到到会话信息，如：
```
Request Headers
Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Encoding:gzip,deflate,sdch
Accept-Language:zh-CN,zh;q=0.8
Cache-Control:max-age=0
Connection:keep-alive
Cookie:PHPSESSID=r887k5n4scg32d4ba34huuhmq7
Host:s.local
User-Agent:Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63
```
从而达到保存会话状态的目的。但也需要注意，如果获取到用户A登录的SESSION ID会怎么样？根据上面的逻辑，如果在请求过程中把获取到的SESSION ID一并发送给服务端，服务端根据SESSION ID读取文件，发现文件内容存在，从而判定用户为A用户，也就是获取到了A用户的用户状态，从而可能可以进行一些敏感操作。所以在会话有效期内，获取到了SESSION ID即获取到了用户的授权，这是比较危险的，以本地的一个管理系统为例，通过chrome登录后查看到客户端COOKIE如下图：

![](../../static/uploads/CodeIgniter-Cookie.png)

假如如果通过某种手段获取到了SESSION ID， 可以模拟发送一个相同的COOKIE过去即可实现登录。FireFox中可添加COOKIE，打开Firebug后Cookies中新建cookie，确定之后刷新页面即可登录到管理系统，如下图：

![](../../static/uploads/CoIgniter-Add-Cookie.png)

通常情况下可通过js获取到cookie，所以需要注意转义，防止数据展示时被执行了。接下来看看CI中的SESSION。在配置文件中有几个跟Session配置相关的参数，影响到Session的使用，它们是：
```
//session保存在cookie中的名称
$config['sess_cookie_name'] = 'ci_session';
//session的有效时间
$config['sess_expiration']  = 7200;
//是否关闭浏览器session失效
$config['sess_expire_on_close'] = FALSE;
//SESSION是否加密存放在COOKIE中
$config['sess_encrypt_cookie']  = FALSE;
//是否保存在数据库中
$config['sess_use_database']    = FALSE;
//存在数据库中，则数据库表名
$config['sess_table_name']  = 'ci_sessions';
//是否匹配IP
$config['sess_match_ip']    = FALSE;
//是否匹配UserAgent
$config['sess_match_useragent'] = TRUE;
//更新时间时间
$config['sess_time_to_update']  = 300;
```
CI自带的SESSION没有服务端文件存储，所有的信息都存放在客户端COOKIE中，当调用$this->load->library('session');时会启动一个会话，即设置一个COOKIE，COOKIE的内容如下：
```

Array
(
    [session_id] => f05138a9513e4928cb0a57672cfe3b53
    [ip_address] => 127.0.0.1
    [user_agent] => Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36
    [last_activity] => 1386569398
    [user_data] =>
)
```
当客户端请求时会将这些信息在HTTP头中传输给服务端，服务端从HTTP头中读取到SESSION信息。同样的可以实现会话，但该方式有很多的不确定因素，根据源码说几点吧：

- 1、 如果日志文件中出现：The session cookie data did not match what was expected. This could be a possible hacking attempt.说明两个问题：a.sess_encrypt_cookie为false，SESSION在COOKIE中未加密存放 b.读取到COOKIE后，校验失败。涉及到加解密、参数处理的情况，容易出现匹配不通过的情况，若不通过则清空SESSION。
- 2、如果sess_match_ip为true，当客户端IP变化时，SESSION将校验不通过，从而清空SESSION。
- 3、sess_match_useragent默认为true，当客户端UserAgent变化时，校验不通过，清空SESION。简单的例子，通过IE浏览器访问，若切换到不同的IE模式，Agent不同，所以校验不通过，清空SESSION。

可以看到，当出现上面任何一种情况时，SESSION都会清空，出现登录不成功或者跳转到登录页面的情况。如果说不加密、不校验IP、UserAgent呢？因为COOKIE是存放在客户端，需要伴随HTTP请求发给服务端，一来过多的COOKIE会影响速度，对一些图片等资源来说完全时浪费带宽；二来COOKIE只能存储4K的数据，加密处理后能存放的更小。种种的不确定因素将产生各种奇怪的问题，避免过多的纠结，果断改用其他方式吧。

# 五、CodeIgniter配置之config

## 配置说明

$config['language']：指定项目语言包。需要注意的时Codeigniter自带的类库错误提示语言包位于/system/language/english/目录下，当这里配置非english时， 如果需要用到这些类库，则需要拷贝语言包到指定的目录中，否则会出现load出错。

$config['charset']：设置系统使用的编码，在某些需要指定编码的函数中会用到，系统、数据库统一编码即可。

$config['enable_hooks']：钩子开关控制，设置为true表示允许使用钩子，否则不允许。

$config['subclass_prefix']：设置自定义类库、函数的前缀，默认为MY_，比如需要重写language helper中的lang方法时，只需要在helper目录下创建MY_language_herper.php，并实现lang函数即可实现“重载”。这里MY_即为subclass_prefix中定义的值。

$config['permitted_uri_chars']：设置URL中允许的字符。

$config['log_threshold']：设置日志记录等级，为0则关闭日志记录，为4则记录所有信息，一般情况设置为1即可。设置之后需要确认下logs目录是否有写入权限。

$config['proxy_ips']：当服务器使用了代理时，REMOTER_ADDR获取的就是代理服务器的IP了，需要从HTTP_X_FORWARDED_FOR、HTTP_CLIENT_IP、HTTP_X_CLIENT_IP、HTTP_X_CLUSTER_CLIENT_IP或其他设定的值中获取。这里设定的就是代理服务器的IP，逗号分隔。

$config['encryption_key']：加密值，如果要用到CI自带的SESION则必须要设置该值。CI的自带SESSION存储与Cookie中，为安全起见，作加密处理。

## 配置读取
CI初始化开始过程中会通过get_config函数加载config.php文件，同时也提供了config_item来获取config的值，如：`echo config_item('charset');`

CI也提供了一个配置类用来维护配置文件。也可以通过下面方式来获取和设置config的值，当设置之后调用get_config的结果同样会变化，所以可以在某些逻辑前修改config的值。
```
//获取config中配置的charset值
echo $this->config->item('charset');
//重新设置config中charset的值
$this->config->set_item('charset', 'gbk')
```

# 六、CodeIgniter配置之router
application/config/routes.php中定义了一个名为$route的数组，用来设置默认路由和404页面以及可以设置一些匹配方式。默认的配置如下：
```
$route['default_controller'] = "welcome";
$route['404_override'] = '';
```
default_controller指定默认的控制器名称，404_override指定当出现404时调用的控制器名称。有时候可能出现解析不成功，或者一直在默认页面，我们可以调用$this->router打印一下当前解析的控制器和Acion名称。比如可以在MY_Controller中如下打印：
```
var_dump($this->router->fetch_directory());
var_dump($this->router->fetch_class());
var_dump($this->router->fetch_method());
```
确定下解析到哪个控制器了， 然后在看看URL的配置、服务器配置，以及可以在Router.php 和URI.php中调试下。$route数组也可以通过通配符（:num, :any）、正则来设置重写规则，下面是一些简单的例子：

1、将 http://pc.local/admin/detail_1.htm 请求解析到 http://pc.local/admin/detail.htm?user_id=1 处理。Codeigniter并不支持包含查询字符串的重写规则，这个规则看起来应当这么写：
```
$route['admin/detail_(:num)'] = 'admin/detail?user_id=$1';
```
但实际上并未生效，程序匹配到admin/detail?user_id=1后用"/"分隔，索引为0的为控制器名，索引为1的为方法名，也就是会将上面的 detail?user_id=1赋值给方法名，结果可想而知就404了。搞清分隔原理后可以在detail后面增加一个斜杠，确保类名和方法名的正确，如：
```
$route['admin/detail_(:num)'] = 'admin/detail/?user_id=$1';
```
但此时又存在参数的获取问题了，会将第三个参数传递给方法，如果需要使用$_GET或者$this->input->get获取还需要对参数进行处理，如：
```
parse_str(ltrim($query_string, '?'), $_GET);
```

2、对PATH_INFO的URL形式重写规则还是比较支持的。如要实现http://pc.local/admin/1这种格式：
```
$route['admin/(:num)'] = 'admin/detail/$1';
```
参数的获取就只能通过段落的方式来获取了。需要注意的是路由将会按照定义的顺序来运行.高层的路由总是优先于低层的路由.最后，能使用CI来设置的路由还是建议使用CI来设置，不依赖服务器配置。

# 七、CodeIgniter配置之database
CodeIgniter的数据库配置文件位于application/config/database.php， 该文件中定义了$db的二维数组，参考文件如下：
```
$active_group = 'default';
$active_record = TRUE;

$db['default']['hostname'] = 'localhost';
$db['default']['username'] = 'root';
$db['default']['password'] = '123456';
$db['default']['database'] = 'test';
$db['default']['dbdriver'] = 'mysql';
$db['default']['dbprefix'] = '';
$db['default']['pconnect'] = FALSE;
$db['default']['db_debug'] = TRUE;
$db['default']['cache_on'] = FALSE;
$db['default']['cachedir'] = '';
$db['default']['char_set'] = 'utf8';
$db['default']['dbcollat'] = 'utf8_general_ci';
$db['default']['swap_pre'] = '';
$db['default']['autoinit'] = TRUE;
$db['default']['stricton'] = FALSE;
```
**配置说明**

$active_group 为$db中的一维键名，表示默认使用的数据库配置，即$this->load->database()不传入参数时，将默认使用$db[$active_group]来连接数据库。

$active_record 是否开启AR模式，开启后将可以使用AR类中的方法,该值可通过$this->load->database()的第三个参数传入。

$db数组需要注意的地方

- 1、port 默认只列出了主机、帐号、密码等，未配置端口号，如果需要特别指定端口号则需要配置该值。
- 2、pconnect 长连接的问题，值默认为TRUE表示默认使用长连接。长连接的使用需要特别小心，数据库可能会出现大量的sleep的进程而导致更多的请求执行不成功，这里不建议开启长连接。
- 3、db_debug 为TRUE时SQL执行出错则会直接在错误页面打印，开发环境可以打开，生产环境需关闭。
- 4、autoinit 是否自动初始化数据库，为true时则$this->load->database()就会连接数据库，否则在查询时连接数据库。CI的类都做了单例，所以不用担心多次链接。
- 5、stricton 当该值为TRUE时，初始化时会执行这样一条语句，会对不规范的数据，比如字符超过长度、自增主键传入‘’等将会直接抛错。

```
SET SESSION sql_mode="STRICT_ALL_TABLES"
```
**如何连接数据库？**

可通过Loader中的database方式调用，即$this->load->database(); 函数的定义如下：
```
/**
 * Database Loader
 *
 * @param    string  数据库连接值，数组或DSN字符串传递。
 * @param    bool    是否返回数据库对象，否则将数据库对象赋值给控制器的db属性
 * @param    bool    是否使用AR，这里的设置会覆盖database.php中设置
 * @return   object
 */
function database($params = '', $return = FALSE, $active_record = NULL){}
```
$params的值有3种情况，分别是：

1、字符串，传入$db数组一维键名，如 default test等，为空则默认$active_group定义的值

2、数组，可以直接传入类似$db的一维数组，如：
```
$this->load->database(array(
    'hostname' => 'localhost',
    'username' => 'root',
    'password' => '123456',
    'database' => 'test',
    'dbdriver' => 'mysql',
    'pconnect' => FALSE,
    'db_debug' => TRUE,
    'char_set' => 'utf8',
    'dbcollat' => 'utf8_general_ci',
));
```
3、DSN 字符串，如：
```
$dsn = 'mysql://root:123456@localhost/test?charset=utf8&dbcollat=utf8_general_ci';
$this->load->database($dsn);
```
PDO的初始化需要使用DSN字符串，那么在CI中该如何配置呢，可参考如下配置：
```
//当前版本2.x.x
$active_group = 'default';
$active_record = TRUE;
  
$db['default']['hostname'] = 'mysql:host=localhost;dbname=test';
$db['default']['username'] = 'root';
$db['default']['password'] = '123456';
$db['default']['database'] = 'test';
$db['default']['dbdriver'] = 'pdo';
$db['default']['dbprefix'] = '';
$db['default']['pconnect'] = FALSE;
$db['default']['db_debug'] = TRUE;
$db['default']['cache_on'] = FALSE;
$db['default']['cachedir'] = '';
$db['default']['char_set'] = 'utf8';
$db['default']['dbcollat'] = 'utf8_general_ci';
$db['default']['swap_pre'] = '';
$db['default']['autoinit'] = TRUE;
$db['default']['stricton'] = FALSE;
```

**如何连接多个数据库？**

$this->load->database()时会将数据库对象赋值给CI_Controller的db属性，如果已经存在了db则不会重新连接。也就是执行$this->load->database()之后再次$this->load->database('test')时则第二次load不会执行。但load的第二个参数允许返回，所以可以返回并赋值给变量，达到连不同库的目的。
```
$DB1 = $this->load->database('default', TRUE);
$DB2 = $this->load->database('test', TRUE);
```
但这种方式需要使用的时候主动去load，使用不太方便，我们可以在MY_Model的构造函数中实现，将返回的$DB1重新赋值给CI_Controller的一个属性，并将该属性赋值或者clone给$this->db，例如：
```
public function __construct($group_name = '')
{
    parent::__construct();
    if($group_name == '') {
        $db_conn_name = 'db';
    } else {
        $db_conn_name = 'db_'.$group_name;
    }
    $CI = & get_instance();
    if(isset($CI->{$db_conn_name}) && is_object($CI->{$db_conn_name})) {
        $this->db = $CI->{$db_conn_name};
    } else {
        $CI->{$db_conn_name} = $this->db = $this->load->database($group_name, TRUE);
    }
}
```
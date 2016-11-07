```
{
    "url": "win-web-server",
    "time": "2016/06/04 18:00",
    "tag": "PHP,Linux,环境部署"
}
```

# 一、Nginx安装
> Nginx是一款轻量级的Web 服务器/反向代理服务器及电子邮件（IMAP/POP3）代理服务器，并在一个BSD-like 协议下发行。

可以直接安装Nginx，也可以安装OpenResty。OpenResty介绍：

> OpenResty 是一个基于 Nginx 与 Lua 的高性能 Web 平台，其内部集成了大量精良的 Lua 库、第三方模块以及大多数的依赖项。用于方便地搭建能够处理超高并发、扩展性极高的动态 Web 应用、Web 服务和动态网关。

> OpenResty 的目标是让你的Web服务直接跑在 Nginx 服务内部，充分利用Nginx 的非阻塞 I/O 模型，不仅仅对 HTTP 客户端请求,甚至于对远程后端诸如 MySQL、PostgreSQL、Memcached 以及 Redis 等都进行一致的高性能响应。

这里安装OpenResty，安装方式同Nginx。安装过程十分简单，只需要下载解压缩即可。解压后根目录下有nginx.exe，双击即可运行Nginx，访问 http://localhost 看到如下页面则表示安装成功。

![](/static/uploads/nginx-start.png)

## 1.1 配置文件
nginx配置文件位于./conf/nginx.conf。这里采用vhost的方式来配置nginx，即一个站点配置在一个文件中，主配置文件引入该配置文件即可。nginx.conf文件引入vhost配置文件方式为在http区域增加以下部分，并创建vhost目录。
```
include vhost/*.conf;
```
## 1.2 示例站点
假设现在需要配置一个后台站点，我们可以设定访问域名为 http://asm_admin.local ，然后创建vhost/asm_admin.conf，写如以下配置信息：
```
server {
    listen       80;
    server_name  asm_admin.local;
    index index.html index.htm index.php;
    root C:/www/asm/asm_admin/public;
    location / {
        if (!-e $request_filename) {
            rewrite . /index.php last;
        }
    }
    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   html;
    }
    location ~ \.php$ {
        fastcgi_pass   127.0.0.1:9000;
        fastcgi_index  index.php;
        include        fastcgi_params;
    }
}
```
这里为PHP项目的示例配置，将请求转发到index.php，其中fastcgi_params需要添加以下行指定文件路径，否则可能出现找不到文件的提示。
```
fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
```
# 二、PHP安装
## 2.1 下载及安装
可直接通过PHP官网下载PHP包，下载地址：http://php.net/downloads.php ，Window用户点击Windows downloads即可跳转到Windows下的下载页面。可发现有多种包可选，VC11、VC14；X86、X64；Non Thread Safe与Thread Safe

**VC**

VC9意思就是该版本PHP是用VisualStudio2008编译的，而VC11则是用VisualStudio2012编译的。依此类推

如果你下载的是VC9版本的，就需要先安装VisualC++RedistributableforVisualStudio2008SP1，

如果你下载的是VC11版本，就需要先安装VisualC++RedistributableforVisualStudio2012.

如果你下载的是VC14版本，就需要先安装Visual C++ Redistributable for Visual Studio 2015

**X86 && X64**

64位系统选择X64，32位系统选择X86

**Non Thread Safe与Thread Safe**

None-Thread Safe就是非线程安全，在执行时不进行线程（thread）安全检查；Nginx、IIS选择此种类型。

Thread Safe就是线程安全，执行时会进行线程（thread）安全检查，以防止有新要求就启动新线程的 CGI 执行方式耗尽系统资源。Apache选择次类型。

## 2.2 解压及配置
下载对应的包后解压，并将PHP路径设置在环境变量PATH中，开发环境将根目录中的php.ini-development拷贝并重命名为php.ini
### 2.2.1 修改PHP扩展路径
```
extension_dir = "C:/soft/php-7.0.4-nts-Win32-VC14-x64/ext"
```
### 2.2.2 打开常用扩展
```
extension=php_curl.dll
extension=php_gd2.dll
extension=php_mbstring.dll
extension=php_mysql.dll
extension=php_mysqli.dll
extension=php_pdo_mysql.dll
extension=php_xmlrpc.dll
extension=php_sockets.dll
```
### 2.2.3 打开错误显示及设置错误级别
```
error_reporting = E_ALL
display_errors = On
```
### 2.2.4 设置时区
```
date.timezone = PRC
```
## 2.3 命令行运行
输入php -v，查看PHP版本信息，看到下面输出则表示成功。
```
php -v
PHP 7.0.4 (cli) (built: Mar  2 2016 14:42:35) ( NTS )
Copyright (c) 1997-2016 The PHP Group
Zend Engine v3.0.0, Copyright (c) 1998-2016 Zend Technologies
```

## 2.4 启动PHP
本地PHP路径为C:/soft/php-7.0.4-nts-Win32-VC14-x64/则，命令行执行下列命令：
```
C:/soft/php-7.0.4-nts-Win32-VC14-x64/php-cgi.exe -b 127.0.0.1:9000 -c C:/soft/php-7.0.4-nts-Win32-VC14-x64/php.ini
```
执行后无任何输出则表示启动成功，配合前面的Nginx，创建文件输出phpinfo并访问，如果能正确访问则说明PHP和Nginx安装成功。

## 2.5 PHP扩展
项目需要用到REDIS，PHP需安装REDIS扩展，下载地址： http://windows.php.net/downloads/pecl/snaps/redis/ 
选择对应的版本下载后将php_redis.dll放到extension_dir指定的ext目录中。php.ini中新增
```
extension = php_redis.dll
```
重启PHP，查看phpinfo中是否出现redis扩展，若出现则表示安装成功。

# 三、Mysql安装
## 3.1 下载
下载地址： http://dev.mysql.com/downloads/mysql/ ，选择对应版本下载后直接解压缩即可。
## 3.2 配置
将安装路径C:/soft/mysql-5.7.11-winx64/bin添加到环境变量；并将mysql-default.ini重命名为my.ini，设置相关路径
```
basedir = C:/soft/mysql-5.7.11-winx64
datadir = C:/soft/mysql-5.7.11-winx64/data
character_set_server = utf8
port = 3306
```
## 3.3 安装Mysql为系统服务
命令行执行
```
mysqld install mysql --defaults-file="C:\soft\mysql-5.7.11-winx64\my.ini"
```
#移除服务命令为：mysqld -remove

## 3.4 Mysql初始化
```
mysqld --initialize
```
初始化完成后root的密码为记录在./data/xx.err文件中。如：

A temporary password is generated for root@localhost: 2r-Pdtkihlu5

## 3.5 启动Mysql
```
net start/stop mysql
```

## 3.6 登录Mysql
```
mysql -u root -p
```
修改密码：
```
mysql> use mysql;
mysql> UPDATE user SET Password = PASSWORD('newpass') WHERE user = 'root';
mysql> FLUSH PRIVILEGES;
```
# 四、Redis服务端安装
这里为Redis服务端，本地需要用到Redis的场景可直接连接本机的Redis，Win下直接下载即可使用。64位可直接[点此下载](/static/attachments/Redis-x64-2.8.2400.zip)。
# 五、PHP、Nginx、启动
前面PHP的启动命令输入后命令行窗口不可关掉，关掉后PHP就退出了， 所以需要把PHP放后台执行，系统需要先安装一个[RunHiddenConsole.exe](/static/attachments/2010-5-RunHiddenConsole.zip)（解压后放在system32目录），然后可以直接将启动命令写在一个批处理脚本，如：
```
@echo off
echo Stopping nginx...
taskkill /F /IM nginx.exe > nul
echo Stopping PHP FastCGI...
taskkill /F /IM php-cgi.exe > nul
echo Stopping Redis...
taskkill /F /IM redis-server.exe > nul
 
set path_php=C:/soft/php-7.0.4-nts-Win32-VC14-x64/
set path_nginx=C:/soft/openresty-1.9.7.3-win32/
set path_redis=C:/soft/Redis-x64-2.8.2400/
 
echo =============================================================================
echo Start PHP FastCGI...
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9000 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9001 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9002 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9003 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9004 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9005 -c %path_php%php.ini
echo Start nginx...
cd /d %path_nginx%
%path_nginx%nginx.exe -t
start  %path_nginx%nginx.exe -c %path_nginx%conf/nginx.conf
echo Start Redis
cd /d %path_redis%
RunHiddenConsole.exe %path_redis%redis-server.exe %path_redis%redis.windows.conf
echo =============================================================================
tasklist /nh^ | findstr /i /s /c:"nginx.exe"
tasklist /nh^ | findstr /i /s /c:"php-cgi.exe"
tasklist /nh^ | findstr /i /s /c:"redis-server.exe"
pause
```
脚本的执行逻辑：

- kill掉Nginx、PHP、Redis进程
- 启动PHP
- 检测Nginx配置文件并启动Nginx
- 启动Redis
- 查看进程中是否有Nginx、PHP、Redis进程

需要注意的的，脚本里开了6个fastcgi进程用来处理PHP脚本。如果只用一个进程（9000端口）来处理，当页面需要执行多个PHP脚本时，比如本地curl则可能会出现超时的情况。主要原因是页面上的PHP交给后端fastcgi进程处理，若还没返回又进来一个请求，此时9000端口已经被占用，导致curl一致处理等待状态从而超时。解决办法是通过upstream模块做负载均衡。
可以在nginx.conf的http模块里增加upstream配置
```
upstream php_fastcgi_pass{
    server 127.0.0.1:9000;
    server 127.0.0.1:9001;
    server 127.0.0.1:9002;
    server 127.0.0.1:9003;
    server 127.0.0.1:9004;
    server 127.0.0.1:9005;
}
```
然后server中的 `fastcgi_pass 127.0.0.1:9000;` 改成 `fastcgi_pass php_fastcgi_pass;`
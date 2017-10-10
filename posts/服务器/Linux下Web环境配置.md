```
{
    "url": "linux-web-server",
    "time": "2016/05/23 19:00",
    "tag": "PHP,Linux,环境部署"
}
```

# 一、OpenResty
OpenResty 是一个基于 Nginx 与 Lua 的高性能 Web 平台，其内部集成了大量精良的 Lua 库、第三方模块以及大多数的依赖项。用于方便地搭建能够处理超高并发、扩展性极高的动态 Web 应用、Web 服务和动态网关。

## 1.1 编译
```
# yum -y install readline-devel pcre-devel openssl-devel perl gcc
# wget https://openresty.org/download/openresty-1.9.7.4.tar.gz
# tar zxvf openresty-1.9.7.4.tar.gz
# cd openresty-1.9.7.4
# ./configure --prefix=/usr/local/server/openresty1.9.7 \
--with-luajit \
--without-http_redis2_module \
--with-http_iconv_module
# gmake
# gmake install
```

### mac下编译方式
```
brew update
brew install pcre openssl curl
./configure --prefix=/usr/local/server/openresty1.11.2.2 \
--with-luajit \
--without-http_redis2_module \
--with-http_iconv_module \
--with-cc-opt="-I/usr/local/opt/openssl/include/ -I/usr/local/opt/pcre/include/" \
--with-ld-opt="-L/usr/local/opt/openssl/lib/ -L/usr/local/opt/pcre/lib/" 
sudo make
sudo make install
```

## 1.2 Nginx参考配置
配置nginx.conf，通过vhost目录将所有配置文件放在此目录下，特定站定只需要关注该站点的配置文件即可。

```
#user  nobody;
worker_processes 2;
error_log /data/log/nginx/nginx_error.log crit;
#Specifies the value for maximum file descriptors that can be opened by this process.
worker_rlimit_nofile 65536;
events
    {
        use epoll;
        worker_connections 65536;
    }
http
    {
        include       mime.types;
        default_type  application/octet-stream;
        server_names_hash_bucket_size 128;
        client_header_buffer_size 32k;
        large_client_header_buffers 4 32k;
        client_max_body_size 50m;
        sendfile on;
        tcp_nopush     on;
        keepalive_timeout 60;
        tcp_nodelay on;
        fastcgi_connect_timeout 300;
        fastcgi_send_timeout 300;
        fastcgi_read_timeout 300;
        fastcgi_buffer_size 64k;
        fastcgi_buffers 4 64k;
        fastcgi_busy_buffers_size 128k;
        fastcgi_temp_file_write_size 256k;
        gzip on;
        gzip_min_length  1k;
        gzip_buffers     4 16k;
        gzip_http_version 1.0;
        gzip_comp_level 2;
        gzip_types       text/plain application/x-javascript text/css application/xml;
        gzip_vary on;
        log_format  access  '$remote_addr - $remote_user [$time_local] "$request" '
             '$status $body_bytes_sent "$http_referer" '
             '"$http_user_agent" $http_x_forwarded_for';
        #limit_zone  crawler  $binary_remote_addr  10m;
        server {
                listen 80 default;
                return 500;
        }
        include vhost/*.conf;
}
```

配置server部分，此示例有判断文件不存在时解析到index.php

```
server
{
    listen       80;
    server_name localhost;
    index index.html index.php;
    root  /data/www/public;
    location /
    {
        if (!-e $request_filename)
        {
            rewrite . /index.php last;
        }
    }
    location ~ .*\.(php|php5)?$
    {
        fastcgi_pass  127.0.0.1:9000;
        fastcgi_index index.php;
        include fastcgi.conf;
    }
    location ~ .*\.(gif|jpg|jpeg|png|bmp|swf)$
    {
        expires      30d;
    }
    location ~ .*\.(js|css)?$
    {
        expires      12h;
    }
    access_log  /data/log/nginx/localhost_access.log  access;
}
```

## 1.3 Nginx启动
```
-- 检测配置文件配置是否正确，如错误会有报错信息
/usr/local/server/openresty1.9.7/nginx/sbin/nginx -t
-- 启动nginx
/usr/local/server/openresty1.9.7/nginx/sbin/nginx
-- 重启动nginx
/usr/local/server/openresty1.9.7/nginx/sbin/nginx reload
```

# 二、Mysql配置
## 2.1 安装
**配置mysql用户**

```
# groupadd mysql
# useradd -s /sbin/nologin -g mysql mysql
# mkdir -p /data/mysql/3306
# chown -R mysql:mysql 3306/

```

**安装依赖包**

```
# yum -y install gcc gcc-c++ ncurses ncurses-devel cmake
```

**下载并解压相应源码包，从MySQL 5.7.5开始Boost库是必需的**

```
# wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-boost-5.7.19.tar.gz
# wget http://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-5.7.12.tar.gz
```

**预编译**

```
# tar zxvf mysql-5.7.12.tar.gz
# cd mysql5.7.2
# cmake . -DCMAKE_INSTALL_PREFIX=/usr/local/server/mysql5.7.12 \
-DMYSQL_DATADIR=/data/mysql/3306 \
-DSYSCONFDIR=/data/mysql/3306 \
-DDEFAULT_CHARSET=utf8mb4 \
-DDEFAULT_COLLATION=utf8mb4_unicode_ci  \
-DEXTRA_CHARSETS=all \
-DENABLED_LOCAL_INFILE=1 \
-DDOWNLOAD_BOOST=1 \
-DWITH_BOOST=/root/package
# make
# make install
```

**初始化数据库**，之前版本mysql_install_db是在mysql_basedir/script下，5.7放在了mysql_install_db/bin目录下,且已被废弃

```
2016-04-15 17:25:52 [WARNING] mysql_install_db is deprecated. Please consider switching to mysqld --initialize
2016-04-15 17:25:52 [ERROR] The data directory needs to be specified.
```

`–initialize`会生成一个随机密码(~/.mysql_secret)，而`–initialize-insecure`不会生成密码。`–datadir`目标目录下不能有数据文件

```
/usr/local/server/mysql5.7.12/bin/mysqld \
--initialize-insecure \
--user=mysql \
--basedir=/usr/local/server/mysql5.7.12 \
--datadir=/data/mysql/3306/data
```

**拷贝配置文件**

```
# cp support-files/my-default.cnf /data/mysql/3306/my.cnf
```

**设置密码**

```
set password = '123456';
```

## 2.2 启动
```
-- 启动
/usr/local/server/mysql5.7.12/bin/mysqld_safe --defaults-file=/data/mysql/3306/my.cnf > /dev/null &
-- 关闭
/usr/local/server/mysql5.7.12/bin/mysqladmin -uroot -p -h localhost -P3306 shutdown
```

## 2.3 设置配置文件
```
# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/5.7/en/server-configuration-defaults.html
# *** DO NOT EDIT THIS FILE. It's a template which will be copied to the
# *** default location during install, and will be replaced if you
# *** upgrade to a newer version of MySQL.
[mysqld]
user = mysql
# Remove leading # and set to the amount of RAM for the most important data
# cache in MySQL. Start at 70% of total RAM for dedicated server, else 10%.
innodb_buffer_pool_size = 128M
# Remove leading # to turn on a very important data integrity option: logging
# changes to the binary log between backups.
# log_bin
# These are commonly set, remove the # and set as required.
basedir = /usr/local/webserver/mysql
datadir = /data/mysql/3306/data
port = 3306
server_id = 1
socket = /tmp/mysql.sock
# Remove leading # to set options mainly useful for reporting servers.
# The server defaults are faster for transactions and fast SELECTs.
# Adjust sizes as needed, experiment to find the optimal values.
join_buffer_size = 128M
sort_buffer_size = 2M
read_rnd_buffer_size = 2M
sql_mode=NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES
auto_increment_increment = 1
slow_query_log
long_query_time = 3
innodb_file_per_table = 1
```

# 三、REDIS安装
```
# yum -y install tcl
# wget http://download.redis.io/releases/redis-3.0.7.tar.gz
# tar zxvf redis-3.0.7.tar.gz
# cd redis-3.0.7
# mkdir -p /usr/local/server/redis3.0.7
# make PREFIX=/usr/local/server/redis3.0.7 install
# cp redis.conf /etc/
 
-- 启动REDIS
# /usr/local/server/redis3.0.7/bin/redis-server /etc/redis.conf&
```

# 四、PHP安装
## 4.1 安装依赖包
安装libmcrypt，未安装时会报如下错误

> configure: error: mcrypt.h not found. Please reinstall libmcrypt.

```
# wget https://sourceforge.net/projects/mcrypt/files/Libmcrypt/2.5.8/libmcrypt-2.5.8.tar.gz/download
# mv download libmcrypt-2.5.8.tar.gz
# tar zxvf libmcrypt-2.5.8.tar.gz
# cd libmcrypt-2.5.8
# ./configure --prefix=/usr/local
# make
# make install
```

> configure: error: Don't know how to define struct flock on this system, set --enable-opcache=no

```
# echo "/usr/local/lib" >> /etc/ld.so.conf.d/local.conf
# ldconfig -v
```

# 4.2 编译PHP
```
# wget http://cn2.php.net/get/php-7.0.5.tar.gz/from/this/mirror
# mv mirror php-7.0.5.tar.gz
# tar zxvf php-7.0.5.tar.gz
# cd php-7.0.5
# yum -y install gcc gcc-c++ libxml2 libxml2-devel autoconf libjpeg libjpeg-devel libpng libpng-devel freetype freetype-devel  zlib zlib-devel glibc glibc-devel glib2 glib2-devel libcurl libcurl--devel curl-devel curl libmcrypt libmcrypt-devel   bzip2 bzip2-devel libmcrypt libmcrypt-devel
# ./configure --prefix=/usr/local/server/php7.0.5 \
--enable-fpm \
--enable-opcache \
--with-mcrypt \
--with-zlib \
--enable-mbstring \
--with-curl \
--disable-debug  \
--disable-rpath \
--enable-inline-optimization \
--with-bz2  \
--with-zlib \
--enable-sockets \
--enable-sysvsem \
--enable-sysvshm \
--enable-pcntl \
--enable-mbregex \
--with-mhash \
--enable-shmop  \
--enable-zip \
--with-pcre-regex  \
--with-gd  \
--with-gettext  \
--enable-bcmath  \
--with-png-dir  \
--with-freetype-dir \
--with-jpeg-dir \
--with-openssl \
--enable-pdo \
--with-pdo-mysql \
--enable-mysqlnd \
--with-mysqli=mysqlnd \
--with-pdo-mysql=mysqlnd
# make
# make install
# cp php.ini-production /usr/local/server/php7.0.5/lib/php.ini
# cd /usr/local/server/php7.0.5/etc
# cp php-fpm.conf.default php-fpm.conf
# vi php-fpm.conf
# cp php-fpm.d/www.conf.default php-fpm.d/www.conf
```

### 内置扩展安装

```
cd php7.0.5/ext/soap
/usr/local/server/php7.0.5/bin/phpize
./configure --with-php-config=/usr/local/server/php7.0.5/bin/php-config --enable-soap
make
make install
```

### mac下编译

```
# brew install openssl libjpeg libpng freetype gettext libmcrypt
# ./configure --prefix=/usr/local/server/php7.1.2 \
--enable-fpm \
--enable-opcache \
--with-mcrypt \
--with-zlib \
--enable-mbstring \
--with-curl \
--disable-debug \
--disable-rpath \
--enable-inline-optimization \
--with-bz2 \
--with-zlib \
--enable-soap \
--enable-sockets \
--enable-sysvsem \
--enable-sysvshm \
--enable-pcntl \
--enable-mbregex \
--with-mhash \
--enable-shmop \
--enable-zip \
--with-pcre-regex \
--with-gd \
--with-gettext=/usr/local/opt/gettext \
--enable-bcmath \
--with-png-dir \
--with-freetype-dir \
--with-jpeg-dir \
--with-openssl=/usr/local/opt/openssl \
--enable-pdo \
--with-pdo-mysql \
--enable-mysqlnd \
--with-mysqli=mysqlnd \
--with-pdo-mysql=mysqlnd

# make
# sudo make install
```


## 4.3 php-fpm启动、重启、终止
php 5.3.3 源码中已经内嵌了 php-fpm，不用象以前的php版本一样专门打补丁了，只需要在configure的时候添加编译参数即可。不再支持 php-fpm 以前具有的 /usr/local/php/sbin/php-fpm (start|stop|reload)等命令，需要使用信号控制：

- INT, TERM 立刻终止
- QUIT 平滑终止
- USR1 重新打开日志文件
- USR2 平滑重载所有worker进程并重新载入配置和二进制模块

```
-- 启动PHP-Fpm
# /usr/local/server/php7.0.5/sbin/php-fpm
 
-- php-fpm 关闭：
# kill -INT `cat /usr/local/server/php7.0.5/var/run/php-fpm.pid`
 
-- php-fpm 重启：
# kill -USR2 `cat /usr/local/server/php/var/run/php-fpm.pid`
```

## 4.4 设置自启动
```
-- 设置自启动
# echo '/usr/local/server/php7.0.5/sbin/php-fpm' >>/etc/rc.local
 
-- 设置软链接
# ln -s /usr/local/webserver/php7.0.5/bin/php /usr/bin/php
```

## 4.5 PHP配置文件设置
```
-- 设置时区
date.timezone = PRC
-- 去掉头PHP返回
expose_php = Off
```

## 4.6 扩展安装
### 4.6.1 REDIS扩展
安装完成后会显示so路径，需修改php.ini添加到so到配置文件中。

```
# wget https://github.com/phpredis/phpredis/archive/php7.zip
# unzip php7.zip
# cd phpredis-php7
# /usr/local/server/php7.0.5/bin/phpize
# ./configure --with-php-config=/usr/local/server/php7.0.5/bin/php-config
# make
# make install
```
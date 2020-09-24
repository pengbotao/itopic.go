```
{
    "url": "linux-php-install",
    "time": "2014/08/14 23:16",
    "tag": "Linux,PHP",
    "toc": "no"
}
```

# PHP安装
```
# yum -y install gcc gcc-c++ libxml2 libxml2-devel autoconf libjpeg libjpeg-devel libpng libpng-devel freetype freetype-devel  zlib zlib-devel glibc glibc-devel glib2 glib2-devel libcurl libcurl--devel curl-devel curl libmcrypt libmcrypt-devel   bzip2 bzip2-devel libmcrypt libmcrypt-devel
 
# wget http://cn2.php.net/get/php-5.5.15.tar.gz/from/this/mirror
# tar zxvf php-5.5.15.tar.gz
# cd php-5.5.15
# ./configure --prefix=/usr/local/webserver/php --enable-fpm --with-mcrypt --with-zlib --enable-mbstring --enable-pdo --with-curl --disable-debug  --disable-rpath --enable-inline-optimization --with-bz2  --with-zlib --enable-sockets --enable-sysvsem --enable-sysvshm --enable-pcntl --enable-mbregex --with-mhash --enable-shmop  --enable-zip --with-pcre-regex  --with-gd  --with-gettext  --enable-bcmath  --with-png-dir  --with-freetype-dir --with-jpeg-dir --with-pdo-mysql
```
未安装opensll和mysql，此时mysql还未安装，之后通过phpize添加openssl和mysql扩展

- --with-openssl=/usr/local/openssl
- --with-mysql=/usr/local/webserver/mysql/

```
# make && make install
 
# cp php.ini-production /usr/local/webserver/php/lib/php.ini
# cd /usr/local/webserver/php/etc
# cp php-fpm.conf.default php-fpm.conf
# vi php-fpm.conf
# /usr/local/webser/php/sbin/php-fpm
```
添加到自启动
```
# echo '/usr/local/webserver/php/sbin/php-fpm' >>/etc/rc.local
```
# Nginx添加PHP解析
```
server {
    listen 80;
    server_name localhost;
    index index.htm index.html index.php;
 
    root /data/www;
 
    location ~ .*\.(php|php5)?$ {
        fastcgi_pass 127.0.0.1:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
    access_log /data/logs/nginx/default.access.log;
}
```
# php-fpm重启、终止

php 5.3.3 源码中已经内嵌了 php-fpm，不用象以前的php版本一样专门打补丁了，只需要在configure的时候添加编译参数即可。不再支持 php-fpm 以前具有的 /usr/local/php/sbin/php-fpm (start|stop|reload)等命令，需要使用信号控制：

master 进程可以理解以下信号：

- INT, TERM 立刻终止
- QUIT 平滑终止
- USR1 重新打开日志文件
- USR2 平滑重载所有worker进程并重新载入配置和二进制模块 

**示例：**
```
php-fpm 关闭：
# kill -INT `cat /usr/local/webserver/php/var/run/php-fpm.pid`
php-fpm 重启：
# kill -USR2 `cat /usr/local/webserver/php/var/run/php-fpm.pid`
 
# ps aux | grep php
root      1422  0.0  0.5  21908  2660 ?        Ss   Aug13   0:00 php-fpm: master process (/usr/local/webserver/php/etc/php-fpm.conf)
nobody    1423  0.0  0.8  22228  4316 ?        S    Aug13   0:00 php-fpm: pool www                    
nobody    1424  0.0  0.6  21908  3292 ?        S    Aug13   0:00 php-fpm: pool www                    
root     19713  0.0  0
```
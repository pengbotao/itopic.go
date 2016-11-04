url: linux-nginx-install
des: 
time: 2014/08/14 19:21
tag: linux
++++++++

# Nginx安装
```
# yum -y install gcc* pcre glib2-devel openssl-devel pcre-devel bzip2-devel gzip-devel lrzsz 
 
# groupadd www && useradd www -g www
 
# wget http://nginx.org/download/nginx-1.6.1.tar.gz
# tar zxvf nginx-1.6.1.tar.gz
# cd nginx-1.6.1
 
# ./configure --user=www --group=www --prefix=/usr/local/webserver/nginx --with-http_stub_status_module --with-http_ssl_module --error-log-path=/data/logs/nginx/error.log --http-log-path=/data/logs/nginx/access.log
# make && make install
```
机器为阿里云512MCentOS，刚初始化的机器发现没有make命令，通过yum安装即可。
```
# yum -y install make
```
# 常用操作
```
-- root软连接到nginx.conf
# ln -s /usr/local/webserver/nginx/conf/nginx.conf /root/nginx.conf
 
-- root目录下直接重启脚本
# echo -e  '#!/bin/bash \n /usr/local/webserver/nginx/sbin/nginx -s reload ' >> /root/nginx_restart.sh
 
-- 添加执行权限
# chmod +x /root/nginx_restart.sh
 
-- 添加到自启动
# echo '/usr/local/webserver/nginx/sbin/nginx' >>/etc/rc.local
```
# 启动Nginx
```
# /usr/local/webserver/nginx/sbin/nginx
 
-- 检测是否配置文件是否正确
# /usr/local/webserver/nginx/sbin/nginx -t
 
-- 重启nginx
# /usr/local/webserver/nginx/sbin/nginx -s reload
```

# nginx.conf配置
配置示例，通过vhost来配置新站点，避免nginx.conf文件过长，不方便管理。
```
# vi nginx.conf 
user  www www;
worker_processes  1;
 
error_log  logs/error.log;
pid        logs/nginx.pid;
 
events {
    use epoll;
    worker_connections  1024;
}
 
http {
    include       mime.types;
    default_type  application/octet-stream;
 
    server_names_hash_bucket_size 128;
    client_header_buffer_size 32k;
    large_client_header_buffers 4 32k;
    client_max_body_size 50m;
 
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';
    sendfile        on;
    tcp_nopush     on;
 
    keepalive_timeout  60;
    tcp_nodelay on;
 
    fastcgi_connect_timeout 300;
    fastcgi_send_timeout 300;
    fastcgi_read_timeout 300;
    fastcgi_buffer_size 64k;
    fastcgi_buffers 4 64k;
    fastcgi_busy_buffers_size 128k;
    fastcgi_temp_file_write_size 128k;
 
    gzip  on;
    gzip_min_length  1k;
    gzip_buffers     4 16k;
    gzip_http_version 1.0;
    gzip_comp_level 2;
    gzip_types       text/plain application/x-javascript text/css application/xml;
    gzip_vary on;
 
    include vhost/*.conf;
}
```
**添加80端口**
```
server {
    listen 80;
    server_name localhost;
    index index.htm index.html index.php;
 
    root /data/www;
    access_log /data/logs/nginx/default.access.log;
}
```

## upstream模块
nginx可以通过upstream实现反向代理，将请求分发到后端的不同机器上。
### 1、轮询
轮询是upstream的默认分配方式，即每个请求按照时间顺序轮流分配到不同的后端服务器，如果某个后端服务器down掉后，能自动剔除。
```
upstream backend {
    server 192.168.1.101:88;
    server 192.168.1.102:88;
    server 192.168.1.103:88;
}
```

### 2、weight
轮询的加强版，即可以指定轮询比率，weight和访问几率成正比，主要应用于后端服务器异质的场景下。
```
upstream backend {
    server 192.168.1.101 weight=1;
    server 192.168.1.102 weight=2;
    server 192.168.1.103 weight=3;
}
```
### 3、ip_hash
每个请求按照访问ip（即Nginx的前置服务器或者客户端IP）的hash结果分配，这样每个访客会固定访问一个后端服务器，可以解决session一致问题。
```
upstream backend {
    ip_hash;
    server 192.168.1.101:88;
    server 192.168.1.102:88;
    server 192.168.1.103:88;
}
```

使用上可在具体的配置文件里设置
```
proxy_pass http://backend; 
```
url: linux-memcache-install
des: 
time: 2014/08/16 16:33
tag: linux
++++++++

# 安装Memcached服务端

memcached依赖于libevent，需要先安装libevent。

**安装 libevent**
```
# tar zxvf libevent-2.0.21-stable.tar.gz
# cd libevent-2.0.21-stable
# ./configure --prefix=/usr/local/libevent
# make && make install
```
**安装memcached服务端**
```
# wget http://www.memcached.org/files/memcached-1.4.20.tar.gz
# tar zxvf memcached-1.4.20.tar.gz 
# cd memcached-1.4.20
 
# ./configure --prefix=/usr/local/webserver/memcached --with-libevent=/usr/local/libevent
# make && make install
```
# 启动memcached
```
# cd /usr/local/webserver/memcached/
# ./bin/memcached -d -m 32 -u root -p 11211 -c 128 –P /tmp/memcached.pid
```
启动参数说明：

- -d 选项是启动一个守护进程。
- -u root 表示启动memcached的用户为root。
- -m 是分配给Memcache使用的内存数量，单位是MB，默认64MB。
- -M return error on memory exhausted (rather than removing items)。
- -u 是运行Memcache的用户，如果当前为root 的话，需要使用此参数指定用户。
- -p 是设置Memcache的TCP监听的端口，最好是1024以上的端口。
- -c 选项是最大运行的并发连接数，默认是1024。
- -P 是设置保存Memcache的pid文件。 

# PHP安装memcache扩展
安装地址
- **memcached** : http://pecl.php.net/package/memcached
- **memcache** : http://pecl.php.net/package/memcache
```
# wget http://pecl.php.net/get/memcache-2.2.7.tgz
# tar zxvf memcache-2.2.7.tgz
# cd memcache-2.2.7
# ./configure --enable-memcache --with-php-config=/usr/local/webserver/php/bin/php-config --with-zlib-dir
# make && make install
```
编译完成后修改php.ini添加extension=memcache.so ，重启php-fpm

# PHP安装memcached扩展
```
# wget http://launchpadlibrarian.net/66527034/libmemcached-0.48.tar.gz
# cd libmemcached-0.48
# ./configure --prefix=/usr/local/libmemcached  --with-memcached
# make && make install
```
安装要注意的问题：

1， 安装过程中不要忘了，--with-memcached，不然会提示你
```
checking for memcached... no
configure: error: "could not find memcached binary"
```
2，你的memcached是不是1.2.4以上的，如果不是会提示你
```
clients/ms_thread.o: In function `ms_setup_thread':
/home/zhangy/libmemcached-0.42/clients/ms_thread.c:225: undefined reference to `__sync_fetch_and_add_4'
clients/ms_thread.o:/home/zhangy/libmemcached-0.42/clients/ms_thread.c:196: more undefined references to `__sync_fetch_and_add_4' follow
collect2: ld returned 1 exit status
make[2]: *** [clients/memslap] Error 1
make[2]: Leaving directory `/home/zhangy/libmemcached-0.42'
```
解决办法是`--disable-64bit CFLAGS="-O3-march=i686"`。
```
# wget http://pecl.php.net/get/memcached-2.2.0.tgz
# cd memcached-2.2.0
# /usr/local/webserver/php/bin/phpize
# ./configure --enable-memcached --with-php-config=/usr/local/webserver/php/bin/php-config --with-libmemcached-dir=/usr/local/libmemcached --disable-memcached-sasl

checking for sasl/sasl.h... no
configure: error: no, sasl.h is not available. Run configure with --disable-memcached-sasl to disable this check
 
# make && make install
 ```
编译完成后修改php.ini添加extension=memcached.so ，重启php-fpm
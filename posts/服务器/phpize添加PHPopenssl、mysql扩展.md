```
{
    "url": "php-phpize-ext",
    "time": "2014/08/14 19:14",
    "tag": "PHP,Linux,phpize"
}
```

# phpize
phpize 命令是用来准备 PHP 扩展库的编译环境的。下面例子中，扩展库的源程序位于 extname 目录中：
```
$ cd extname
$ phpize
$ ./configure
$ make
# make install
```
成功的安装将创建 extname.so 并放置于 PHP 的扩展库目录中。需要调整 php.ini，加入 extension=extname.so 这一行之后才能使用此扩展库。

如果系统中没有 phpize 命令并且使用了预编译的包（例如 RPM），那要安装 PHP 包相应的开发版本，此版本通常包含了 phpize 命令以及相应的用于编译 PHP 及其扩展库的头文件。

使用 phpize --help 命令可以显示此命令用法。

原文：http://php.net/manual/zh/install.pecl.phpize.php

# 添加openssl扩展
```
# cd /data/pkg/php-5.5.15/ext/openssl
# /usr/local/webserver/php/bin/phpize
Cannot find config.m4. 
Make sure that you run '/usr/local/webserver/php/bin/phpize' in the top level source directory of the module
 
# mv config0.m4 config.m4
 
# ./configure --with-openssl --with-php-config=/usr/local/webserver/php/bin/php-config
# make && make install
```
编译完成后将生成的openssl.so添加到php.ini中，并重启php-fpm
```
# vi /usr/local/webserver/php/lib/php.ini
 
# ps aux | grep php
root     19714  0.0  0.5  21912  2652 ?        Ss   09:20   0:00 php-fpm: master process (/usr/local/webserver/php/etc/php-fpm.conf)
nobody   19715  0.0  0.8  21912  4144 ?        S    09:20   0:00 php-fpm: pool www                    
nobody   19716  0.0  0.8  22232  4308 ?        S    09:20   0:00 php-fpm: pool www                    
root     22464  0.0  0.1   4028   676 pts/0    R+   09:37   0:00 grep php
# kill -USR2 19714
```
重启完成后即可调用上篇rsa中的create方法生成公钥、私钥对。
```
print_r(RsaUtil::create());
```
# 添加mysql扩展
```
# cd /data/pkg/php-5.5.15/ext/mysql
# /usr/local/webserver/php/bin/phpize 
Configuring for:
PHP Api Version:         20121113
Zend Module Api No:      20121212
Zend Extension Api No:   220121212
 
# ./configure --with-php-config=/usr/local/webserver/php/bin/php-config --with-mysql=/usr/local/webserver/mysql/
# make && make install
Build complete.
Don't forget to run 'make test'.
 
Installing shared extensions:     /usr/local/webserver/php/lib/php/extensions/no-debug-non-zts-20121212/
 
# vi /usr/local/webserver/php/lib/php.ini
```
mysql.so会自动添加extension_dir所指目录中，修改php.ini添加 extension=mysql.so 重启php
```
# ps aux | grep php
root      1422  0.0  0.4  22020  2404 ?        Ss   11:28   0:00 php-fpm: master process (/usr/local/webserver/php/etc/php-fpm.conf)
nobody    1423  0.0  0.4  22020  2144 ?        S    11:28   0:00 php-fpm: pool www                    
nobody    1424  0.0  0.4  22020  2144 ?        S    11:28   0:00 php-fpm: pool www                    
root      4767  0.0  0.1   4028   680 pts/0    R+   11:39   0:00 grep php
# kill -USR2 1422
```
查看PHPINFO是否已添加成功。
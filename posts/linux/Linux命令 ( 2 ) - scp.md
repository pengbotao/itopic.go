```
{
    "url": "linux-scp",
    "time": "2013/11/15 11:16",
    "tag": "Linux",
    "toc": "no"
}
```

拷贝远程文件到本地
```
scp root@192.168.0.1:/home/database.php /var/www/database.php
```
拷贝本地文件到其他服务器
```
scp /var/www/database.php root@192.168.0.1:/home/database.php
```
拷贝目录带上 -r 参数， ssh非22端口带上-P 参数。
```
scp -P 2133 -r root@192.168.0.1:/home/wwwroot/ /var/www/
scp -P 2133 /var/www/ root@192.168.0.1:/home/wwwroot/
```
**几个可能有用的参数 :**

- -v 和大多数 linux 命令中的 -v 意思一样 , 用来显示进度 . 可以用来查看连接 , 认证 , 或是配置错误 .
- -C 使能压缩选项 .
- -r 传递目录下的所有内容 .
- -P 选择端口 . 注意 -p 已经被 rcp 使用 .
- -4 强行使用 IPV4 地址 .
- -6 强行使用 IPV6 地址 .
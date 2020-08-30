```
{
    "url": "supervisor",
    "time": "2016/02/15 07:14",
    "tag": "Python,常用软件"
}
```

-- 安装supervisor
$ pip install supervisor

-- 生成配置文件
$ echo_supervisord_conf > /etc/supervisord.conf
$ vim /etc/supervisord.conf

-- 启动supervisor
$ supervisord -c /etc/supervisord.conf

-- 查看supervisor是否运行
$ ps aux | grep supervisord


**program 配置**

一份配置文件至少需要一个 `[program:x]` 部分的配置，来告诉`supervisord`需要管理那个进程。`[program:x]` 语法中的 x 表示 program name，会在客户端（supervisorctl 或 web 界面）显示，在 supervisorctl 中通过这个值来对程序进行 start、restart、stop 等操作。

program可以直接修改`supervisord.conf`也可以通过 include 的方式把不同的程序（组）写到不同的配置文件里。这里参考后面这种。

```
[include]
files = /etc/supervisor/*.conf
```

```
#/etc/supervisor/mgr.conf
[program:mgr]
directory=/home/www/mgr
command=/usr/bin/python /home/www/mgr/mgr_server.py --port=88%(process_num)02d --log-file-prefix=/home/www/mgr/log/tornado/88%(process_num)02d.log
numprocs=3
process_name=%(program_name)s_88%(process_num)02d;
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile = /home/www/mgr/log/supervisor.log
```


**supervisorctl**

`Supervisorctl` 是 `supervisord` 的一个命令行客户端工具，启动时需要指定与 `supervisord` 使用同一份配置文件，否则与 `supervisord` 一样按照顺序查找配置文件。

```
$ supervisorctl -c /etc/supervisord.conf

> status         # 查看程序状态
> stop mgr       # 关闭 usercenter 程序
> start mgr      # 启动 usercenter 程序
> restart mgr    # 重启 usercenter 程序
> reread         # 读取有更新（增加）的配置文件，不会启动新添加的程序
> update         # 重启配置文件修改过的程序

> reload         # 重启supervisor程序
```

---

- [1] [使用 supervisor 管理进程](http://liyangliang.me/posts/2015/06/using-supervisor/)
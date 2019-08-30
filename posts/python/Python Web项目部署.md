```
{
    "url": "python-project-deploy",
    "time": "2016/06/01 00:05",
    "tag": "Python"
}
```

在部署之前我们先配置一个HelloWorld的项目，这里以Flask框架为例：

```
# test.py

import time
from flask import Flask

app = Flask(__name__)


@app.route('/')
def hello_world():
    return 'Hello World!'


@app.route('/sleep')
def sleep():
    time.sleep(1)
    return "Sleep 1 sec.."
```

`.flaskenv`（需安装：pip3 `install python-dotenv`）

```
FLASK_APP=test
FLASK_ENV=development
FLASK_RUN_PORT=5000
FLASK_RUN_HOST=127.0.0.1
```

然后命令行执行`flask run`即可启动该项目了。但如果环境配置为生产环境会提示：`WARNING: This is a development server. Do not use it in a production deployment.`。Flask在开发环境下默认开启了多线程支持，但生产环境不建议这样子直接启动。大多数框架都是同步阻塞的模式，也不方便使用到多核。所以生产环境下一般会建议配合使用uwsgi、gunicorn、nginx等。


# gunicorn

```
$ pip3 install gunicorn
```

**通过gunicorn启动：**

`gunicorn -w 4 -b 127.0.0.1:5000 test:app`

- `-w 4`：是指预定义的工作进程数为4
- `-b 127.0.0.1:5000`：服务绑定的IP和端口
- `test:app`：test启动文件，app是应用实例

**指定使用gevent**

`gunicorn -k gevent -w 4 -b 127.0.0.1:5000 test:app`

也可以通过指定**配置文件**的方式来进行启动：

```
# gunicorn.conf

# 并行工作进程数
workers = 4
# 指定每个工作者的线程数
threads = 2
# 监听内网端口5000
bind = '127.0.0.1:5000'
# 设置守护进程
daemon = 'false'
# 工作模式协程
worker_class = 'gevent'
# 设置最大并发量
worker_connections = 2000
# 设置进程文件目录
pidfile = './log/gunicorn.pid'
# 设置访问日志和错误信息日志路径
accesslog = './log/gunicorn_acess.log'
errorlog = './log/gunicorn_error.log'
# 设置日志记录水平
loglevel = 'warning'
```

`gunicorn -c gunicorn.conf test:app`

然后也可以在前面在搭一个`Nginx`：

```
server {
    listen 80;

    server_name flask.local;

    access_log  /var/log/nginx/access.log;
    error_log  /var/log/nginx/error.log;

    location / {
        proxy_pass         http://127.0.0.1:5000/;

        proxy_set_header   Host             $host;
        proxy_set_header   X-Real-IP        $remote_addr;
        proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
    }
}
```

默认是阻塞的模式，通过`gevent`协程化之后效率差别比较大。上面后两种方式使用到了`gevent`。

```
ab -c 100 -n 1000 http://127.0.0.1:5000/sleep

# 启用gevent
Requests per second:    90.03 [#/sec] (mean)
Time per request:       1110.689 [ms] (mean)
Time per request:       11.107 [ms] (mean, across all concurrent requests)
Transfer rate:          15.21 [Kbytes/sec] received

# 未启用
Requests per second:    7.59 [#/sec] (mean)
Time per request:       13170.934 [ms] (mean)
Time per request:       131.709 [ms] (mean, across all concurrent requests)
Transfer rate:          1.28 [Kbytes/sec] received
```

# uWSGI

```
$ pip3 install uwsgi
$ pip3 install uwsgitop
```

**通过uwsgi启动：**

`uwsgi --socket 127.0.0.1:5000 --protocol=http -w test:app`

同样，也可以通过**配置文件**的方式启动：

```
# uwsgi.ini

[uwsgi]
#socket = 127.0.0.1:5000
http = 127.0.0.1:5000
wsgi-file = test.py
callable = app
master = true
processes = 4
threads = 2

chdir = /Users/peng/workspace/python/demo
#daemonize = %(chdir)/log/uwsgi.log
pidfile = %(chdir)/log/uwsgi.pid
stats = %(chdir)/log/uwsgi.status
max-requests = 100000
# 当服务器退出时自动删除unix socket文件和pid文件
vacuum = false
```

**启动：**

启动：`uwsgi --ini uwsgi.ini`

重启：`uwsgi --reload log/uwsgi.pid`

停止：`uwsgi --stop log/uwsgi.pid`

查看状态：`uwsgi --connect-and-read log/uwsgi.status`

和`gunicorn`一样，也可以通过`Nginx`将请求转到后端5000端口(若配置Nginx，配置文件中请使用socket)。

```
server {
    listen  80;
    server_name flask.local;

    location / {
        include      uwsgi_params;
        uwsgi_pass   127.0.0.1:5000;
    }
}
```
```
{
    "url": "python-flask",
    "time": "2016/05/10 17:12",
    "tag": "Python"
}
```

# 一、概述
Flask是一个使用 Python 编写的轻量级 Web 应用框架。其 WSGI 工具箱采用 Werkzeug ，模板引擎则使用 Jinja2 。Flask使用 BSD 授权。[查看文档](https://dormousehole.readthedocs.io/en/latest/)

## Hello Python

```
# /usr/bin/env python
# encoding: utf-8

from flask import Flask

app = Flask(__name__)


@app.route('/')
def hello():
    return "Hello Python"


if __name__ == "__main__":
    app.run(host="127.0.0.1", port=5000, debug=True)
```
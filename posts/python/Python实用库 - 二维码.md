```
{
    "url": "python-qrcode",
    "time": "2017/02/25 19:45",
    "tag": "Python"
}
```

# 一、简介

- https://github.com/lincolnloop/python-qrcode

# 二、安装

```
pip install qrcode[pil]
```

# 三、示例

## 3.1 命令行

```
qr "Some text" > test.png
```

## 3.2 保存文件

```
import qrcode

img = qrcode.make("some data here")
img.save("test.png")
```
```
{
    "url": "python-cli-parse",
    "time": "2016/03/02 07:53",
    "tag": "Python,Python常用库",
    "public": "yes"
}
```

# optparser


```
from optparse import OptionParser

# Simple usage example:
parser = OptionParser()
parser.add_option("-f", "--from", action="store", dest="from", type="string", default="web", help=u"帮助信息")

# 函数定义为def add_option(self, *args, **kwargs):。可指定多个参数名
parser.add_option("-c", "--call", "-x", "-y",  action="store_true", dest="call", default=False, help=u"帮助信息")

(options, arg) = parser.parse_args()

print(options, arg)

# 转换为字典
print(vars(options))
```

## add_option

字段|说明
---|---
action|默认为store。 store_true 和 store_false，store_const 、append 、count 
dest|将接收的数据存储到`dest`变量上
type|数据类型："string", "int", "long", "float", "complex", "choice"
default|默认值
help|帮助信息

## 使用说明

```
$ python x.py -f sina
{'call': False, 'from': 'sina'}

$ python x.py --from=sina
{'call': False, 'from': 'sina'}
```

# argparse

- http://blog.xiayf.cn/2013/03/30/argparse/
- https://docs.python.org/zh-cn/3.7/library/optparse.html?highlight=optparse#module-optparse
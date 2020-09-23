```
{
    "url": "python-config",
    "time": "2016/04/01 23:03",
    "tag": "Python"
}
```

# 一、简介

# 二、ConfigParser模块

## 2.1 ConfigParser 简介

`ConfigParser`是用来读取配置文件的包。配置文件的格式如下：中括号"[ ]"内包含的为`section`。`section`下面为类似于`key-value`的配置内容。

```
[db]
db_host = 127.0.0.1
db_port = 69
db_user = root
db_pass = root
host_port = 69

[concurrent]
thread = 10
processor = 20
```

## 2.2 ConfigParser 初始化对象

使用ConfigParser 首选需要初始化实例，并读取配置文件：

```
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
```

## 2.4 获取所用的section节点

```
# 获取所用的section节点
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
print(config.sections())
#运行结果
# ['db', 'concurrent']
```

## 2.5 获取指定section的options。

即将配置文件某个section 内key 读取到列表中：

```
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
r = config.options("db")
print(r)
#运行结果
# ['db_host', 'db_port', 'db_user', 'db_pass', 'host_port']
```

## 2.6 获取指点section下指点option的值

```
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
r = config.get("db", "db_host")
# r1 = config.getint("db", "k1") #将获取到值转换为int型
# r2 = config.getboolean("db", "k2" ) #将获取到值转换为bool型
# r3 = config.getfloat("db", "k3" ) #将获取到值转换为浮点型
print(r)
#运行结果
# 127.0.0.1
```

## 2.7 获取指点section的所用配置信息

```
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
r = config.items("db")
print(r)
#运行结果
#[('db_host', '127.0.0.1'), ('db_port', '69'), ('db_user', 'root'), ('db_pass', 'root'), ('host_port', '69')]
```

## 2.8 修改某个option的值，如果不存在则会出创建

```
# 修改某个option的值，如果不存在该option 则会创建
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
config.set("db", "db_port", "69")  #修改db_port的值为69
config.write(open("ini", "w"))
```

## 2.9 检查section或option是否存在，bool值

```
import configparser
config = configparser.ConfigParser()
config.has_section("section") #是否存在该section
config.has_option("section", "option")  #是否存在该option
```

## 2.10 添加section 和 option

```
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
if not config.has_section("default"):  # 检查是否存在section
    config.add_section("default")
if not config.has_option("default", "db_host"):  # 检查是否存在该option
    config.set("default", "db_host", "1.1.1.1")
config.write(open("ini", "w"))
```

## 2.11 删除section 和 option

```
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
config.remove_section("default") #整个section下的所有内容都将删除
config.write(open("ini", "w"))
```

## 2.12 写入文件

以下的几行代码只是将文件内容读取到内存中，进过一系列操作之后必须写回文件，才能生效。

```
import configparser
config = configparser.ConfigParser()
config.read("ini", encoding="utf-8")
```

写回文件的方式如下：（使用configparser的write方法）

```
config.write(open("ini", "w"))
```

From: https://www.cnblogs.com/ming5218/p/7965973.html
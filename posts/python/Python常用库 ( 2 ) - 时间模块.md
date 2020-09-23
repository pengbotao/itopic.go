```
{
    "url": "python-time-module",
    "time": "2016/03/06 19:16",
    "tag": "Python"
}
```

# 一、time模块
```
>>> print(time.__doc__)

Functions:

time() -- return current time in seconds since the Epoch as a float
clock() -- return CPU time since process start as a float
sleep() -- delay for a number of seconds given as a float
gmtime() -- convert seconds since Epoch to UTC tuple
localtime() -- convert seconds since Epoch to local time tuple
asctime() -- convert time tuple to string
ctime() -- convert time in seconds to string
mktime() -- convert local time tuple to seconds since Epoch
strftime() -- convert time tuple to string according to format specification
strptime() -- parse string to time tuple according to format specification
tzset() -- change the local timezone
```

## 1.1 获取时间戳

```
>>> time.time()
1552827527.447389
>>> int(time.time())
1552827538
```

## 1.2 停留x秒

支持浮点数。

```
>>> print time.sleep.__doc__
sleep(seconds)

Delay execution for a given number of seconds.  The argument may be
a floating point number for subsecond precision.
```

## 1.3 time转字符串

```
>>> print time.strftime.__doc__
strftime(format[, tuple]) -> string

Convert a time tuple to a string according to a format specification.
See the library reference manual for formatting codes. When the time tuple
is not present, current time as returned by localtime() is used.
```

输出当前时间：`time.strftime("%Y-%m-%d %H:%M:%S")`; `time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())`

输出指定时间：`time.strftime("%Y-%m-%d %H:%M:%S", time.strptime("2006-01-02 15:04:05", "%Y-%m-%d %H:%M:%S"))`

## 1.4 字符串转time

```
>>> print time.strptime.__doc__
strptime(string, format) -> struct_time

Parse a string to a time tuple according to a format specification.
See the library reference manual for formatting codes (same as strftime()).
```

# 二、datetime.datetime模块

## 2.1 获取当前时间/时间戳

```
>>> from datetime import datetime

# 初始化一个时间
>>> x = datetime(year=2006, month=1, day=2, hour=15, minute=54, second=5)

# 获取当前时间
>>> t = datetime.now()
datetime.datetime(2019, 3, 17, 20, 55, 8, 556405)

# 单项打印
>>> print(t.year, t.month, t.day, t.hour, t.minute, t.second, t.microsecond)

# 获取时间戳。应该是3.x版本才有
>>> datetime.now().timestamp()
1552827338.052106

>>> datetime.today()
datetime.datetime(2019, 3, 17, 21, 10, 55, 897950)
```

## 2.2 datetime转字符串

```
# 将当前时间打印为字符串
>>> datetime.now().strftime("%Y-%m-%d %H:%M:%S")
'2019-03-17 20:56:48'
```

## 2.3 字符串转datetime

```
>>> s = '2006-01-02 15:04:05'
>>> print(datetime.strptime(s, "%Y-%m-%d %H:%M:%S"))
2006-01-02 15:04:05
```

## 2.4 时间戳转datetime

```
>>> t = time.time()
>>> print(datetime.fromtimestamp(t))
2019-03-17 20:58:15.470147
>>> print(datetime.utcfromtimestamp(t))
2019-03-17 12:58:15.470147
```

## 2.5 时间移动

通过`timedelta`对象来计算过去或者未来的时间。

```
from datetime import datetime, timedelta

print(datetime.today() - timedelta(weeks=1,
                                   days=30,
                                   hours=1,
                                   minutes=30,
                                   seconds=20,
                                   microseconds=10,
                                   milliseconds=5))
```
## 2.6 时间间隔

支持两个`datetime`类型的减法，两个`datetime`相减得到的是`timedelta`对象，不支持加法。

```
from datetime import datetime, timedelta

s1 = '2006-01-02 15:04:05'
d1 = datetime.strptime(s1, "%Y-%m-%d %H:%M:%S")

d2 = datetime.now()
print(d1, d2)
# (datetime.datetime(2006, 1, 2, 15, 4, 5), datetime.datetime(2019, 3, 17, 21, 25, 53, 648149))

t = d2 - d1
# 两个datetime相减得到的是`timedelta`类型
print(type(t))
# <type 'datetime.timedelta'>

# 可查看天数和秒数
print(t.days, t.seconds)
# (4822, 22908)
```

# 三、常用时间操作

```
from datetime import datetime, timedelta

t1 = datetime.now()

# 获取本月第一天
t2 = datetime(t1.year, t1.month, 1)
print(t2.strftime("%Y-%m-%d"))

# 本月最后一天
t3 = datetime(year=t1.year, month=t1.month+1, day=1) - timedelta(days=1)
print(t3.strftime("%Y-%m-%d"))
```

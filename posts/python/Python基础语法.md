```
{
    "url": "python-language",
    "time": "2015/12/25 00:22",
    "tag": "Python"
}
```

# 一、数据类型

## 1.1 数值类型 - Number

- int（有符号整型）
- long（长整型，Python3去掉）
- float（浮点数）
- complex（复数）

## 1.2 字符串 - str

可以使用单引号或者双引号来创建字符串，如：

```
>>> a = "Hello World"
>>> b = 'iTopic.org'
```

**Python中文档查看方法：**

命令行中可以通过`dir(str)`查看字符串支持的方法：
```
>>> dir(str)
['...', 'capitalize', 'center', 'count', 'decode', 'encode', 'endswith', 'expandtabs', 'find',
 'format', 'index', 'isalnum', 'isalpha', 'isdigit', 'islower', 'isspace', 'istitle', 'isupper',
 'join', 'ljust', 'lower', 'lstrip', 'partition', 'replace', 'rfind', 'rindex', 'rjust',
 'rpartition', 'rsplit', 'rstrip', 'split', 'splitlines', 'startswith', 'strip', 'swapcase',
 'title', 'translate', 'upper', 'zfill']
```

查看对应函数的说明文档：

```
>>> print(str.find.__doc__)
S.find(sub [,start [,end]]) -> int

Return the lowest index in S where substring sub is found,
such that sub is contained within S[start:end].  Optional
arguments start and end are interpreted as in slice notation.

Return -1 on failure.

>>> help(str.join)
```

**字符串常用操作:**

**1. 字符串格式化**

```
# 字符串拼接
>>> "Hello world " + "2019"
'Hello world 2019'
# 字符串包含变量
>>> "Hello %s %d" % ("world", 2019)
'Hello world 2019'
# 字符串占位符
>>> "Hello {0} {1}".format("world", 2019)
'Hello world 2019'
# 字符串拼接
>>> ' '.join(["Hello world", "2019"])
'Hello world 2019'
```

**2. 字符串去掉两侧空白字符**

```
>>> s = "Hello world 2019 "
# 去掉两侧字符
>>> s.strip()
'Hello world 2019'
# 去掉左侧字符
>>> s.lstrip("H")
'ello world 2019 '
# 去掉右侧字符
>>> s.rstrip(" 2019 ")
'Hello world'
```

**3. 字符串分隔**

```
# 字符串分隔
>>> fruits = "apple, orange, banana"
>>> fruits.split(", ")
['apple', 'orange', 'banana']

# 字符串拼接
>>> "-".join(["Hello", "World"])
'Hello-World'
```

**4. 字符串替换**

```
>>> s = "Hello 2019."
>>> s.replace("2019", "2020")
'Hello 2020.'
```

**5. 判断前缀、后缀**

```
# 判断前缀
>>> s = "Hello 2019."
>>> s.startswith("Hello")
True

# 判断后缀
>>> s.endswith(".")
True
>>> s.endswith("2019")
False
```

**6. 字符串比较**

```
# 内置函数cmp比较
>>> a = "Hello"
>>> b = "hello"
>>> cmp(a, b)
-1
```


# 二、控制流

## 2.1 if语句

```
fruits = "apple"

if fruits == "apple":
    print("it's an apple")
elif fruits == "orange":
    print("it's an orange")
else:
    print("other")
```

`Python`不支持`switch`语句。

## 2.2 while语句

```
i = 0
while i < 10:
    print(i)
    i += 2
    # break
else:
    print("else i = %d " % i)
```

Python中while语句可以带一个else语句，当循环正常结束时执行，需要注意通过break结束的循环不会执行else语句。

## 2.3 for语句

```
for i in range(0, 10, 1):
    print(i)
```

示例中range可以理解成`(i = 0; i <10; i++)`，同时可以精简写成`rang(10)`, 等同于区间`[0, 10)`的数组.

`for`语句后同样可以带`else`语句，使用和`while`里一致。

---

可以看到控制语句后面都可以携带`else`语句。有一种查找的用法，

```
import random

m = random.randint(5, 15)

for i in range(10):
    if i == m:
        print("find %d" % m)
        break
else:
    print("can't find %d" % m)
```


# 三、数据结构

## 3.1 列表 - list

- 列表有序，是一种序列
- 列表可变长度，可通过`list`的相关操作对列表进行增删等操作
- 列表数据类型可不同，即并不要求所有列表元素都是同一数据类型
- 列表用中括号表示

**列表提供的方法列表**
(通过`dir(list)`可查看到)

方法|说明
---|---
`list.append(obj)`|在列表末尾添加新的对象
`list.count(obj)`|统计某个元素在列表中出现的次数
`list.extend(seq)`|在列表末尾一次性追加另一个序列中的多个值（用新列表扩展原来的列表）
`list.index(obj)`|从列表中找出某个值第一个匹配项的索引位置
`list.insert(index, obj)`|将对象插入列表
`list.pop(obj=list[-1])`|移除列表中的一个元素（默认最后一个元素），并且返回该元素的值
`list.remove(obj)`|移除列表中某个值的第一个匹配项
`list.reverse()`|反向列表中元素
`list.sort([func])`|对原列表进行排序

**1. 列表初始化及遍历**

```
# 初始化列表
fruits = ["apple", "orange", "banana"]

# 默认遍历列表。默认只有一个参数接收遍历的值。用于按`{VAL}`遍历。该参数标识列表的值。
for val in fruits:
    print(val)

# 根据长度来遍历。用于按`{KEY}`遍历
for idx in range(len(fruits)):
    print(fruits[idx])

# 按索引 - 值遍历，也可用于字典。用于按`{KEY - VAL}`遍历
for idx, val in enumerate(fruits):
    print(idx, val)
```

**2. 合并数据**

```
# 直接在原列表后面追加数据。可以是任何类型。
fruits.append("pear")
print(fruits) # output: ['apple', 'orange', 'banana', 'grape', 'pear']

# 在list后扩展数据，接收一个list，展平的方式追加到原列表后面
fruits.extend(["peach", "strawberry"])
print(fruits) # output: ['apple', 'orange', 'banana', 'grape', 'pear', 'peach', 'strawberry']

# 也可用加号来合并两个列表
fruits + ["mango"] # mango追加到fruits数据后面

# 删除数据的3中方法
fruits.pop(2)
furits.remove("banana")
del fruits[2]
```

**3. 健壮性判断**

```
# 判断值是否在列表中(in 与 not in)
if "strawberry" in fruits:
    print("strawberry is in frutis") # output: strawberry is in frutis
else:
    print("strawberry is not in frutis")

# 类型判断
if isinstance(fruits, list):
    print("fruits is a list")

# 强制类型转换
>>> x = "Hello Python"
>>> list(x)
['H', 'e', 'l', 'l', 'o', ' ', 'P', 'y', 't', 'h', 'o', 'n']

# 判断是否为空
if not fruits:
    print("fruits is empty")

if len(fruits) == 0:
    print(fruits is empty)
```

## 3.2 元祖 - tuple

- 元祖有序，也是一种序列
- 元祖一旦初始化不可修改
- 元祖数据类型也可不同
- 元祖用圆括号表示
- 元祖的访问和遍历同列表

**元祖提供的方法列表**

方法|说明
---|---
`tuple.count(obj)`|统计某个元素在元祖中出现的次数
`tuple.index(obj)`|从元祖中找出某个值第一个匹配项的索引位置

**1. 元祖初始化**

```
# 初始化元祖。只有1个元素时最后携带逗号
fruits = ("apple", "orange", "banana")

# 和列表一致，按索引的方式进行访问
print(fruits[0])

# 遍历元祖
for val in fruits:
    print(val)
```

## 3.3 字典 - dict

- 按`键-值`对的方式初始化，相同的键会覆盖，用大括号表示
- 字典无序


**字典提供的方法列表**

方法|说明
---|---
`dict.clear()`|清空字典
`dict.get(k[,d])`|D[k] if k in D, else d.  d defaults to None.
`dict.has_key(k)`|True if D has a key k, else False
`dict.items()`|list of D's (key, value) pairs, as 2-tuples
`dict.keys()`|list of D's keys
`dict.pop(k, d)`|remove specified key and return the corresponding value.If key is not found, d is returned if given, otherwise KeyError is raised
`dict.update(obj)`|合并两个字典

**1. 字典初始化及遍历**

```
# 初始化字典方式一
fruits = {"apple": 1, "orange": 2, "banana": 3}
# 初始化字典方式二
fruits = dict(apple = 1, orange = 2, banana = 3)

# 按KEY遍历
for idx in fruits:
    print(fruits[idx])

# 判断键是否在字典中
if "grape" in fruits:
    print(fruits["grape"])
else:
    fruits["grape"] = 4
# 判断键是否在字段中
if fruits.has_key("apple"):
    print(fruits["apple"])

# 按KEY - VALUE遍历
for idx, val in fruits.items():
    print(idx, val)
```

**2. 字典操作**

批量更新字典里的内容。
```
fruits = dict(apple = 1, orange = 2, banana = 3)
fruits.update(cherry = 4)
# {'orange': 2, 'cherry': 4, 'banana': 3, 'apple': 1}

fruits.update({"mango": 5})
# {'orange': 2, 'cherry': 4, 'mango': 5, 'banana': 3, 'apple': 1}

fruits.keys() # ['orange', 'cherry', 'mango', 'banana', 'apple']

# 删除数据的3中方法
del fruits["apple"]
fruits.pop("banana")

# 清空字典
fruits.clear()
```

## 3.4 集合 - set

- 集合里的VALUE可以是不同的类型，相同数据会`去重`。
- 集合`无序`

**集合提供的方法列表**

方法|说明
---|---
`set.add(obj)`|往集合里添加元素
`set.update(obj)`|更新合并集合
`set.discard(obj)`|丢弃一个元素
`set.remove(obj)`|移除一个元素
`set.pop()`|从集合中弹出一个元素
`set.clear()`|清空一个集合
-|更多交叉并补方法

**1. 集合初始化**

```
# 可变集合set初始化 - 初始化之后在进行赋值
fruits = set()
fruits = {"apple", "orange"}

# 直接创建
fruits = {"apple", "orange"}

# 通过列表转换
fruits = set(["apple", "orange", "banana"])
print(fruits)

# 通过字典转换
fruits = set({"apple":1, "orange":2, "banana":3}) # set(['orange', 'apple', 'banana'])


# 不可变集合用frozenset表示
websites = frozenset(["qq.com", "weibo.com"])
```

集合分为可变集合(`set`)和不可变集合(`frozenset`)。针对可变集合可以往集合里**添加元素**。

```
fruits = {"apple", "orange"}
fruits.add("banana") # set(['orange', 'apple', 'banana'])

fruits.update({"cherry"}) # set(['orange', 'cherry', 'apple', 'banana'])

fruits.update("grape") # set(['a', 'e', 'apple', 'g', 'cherry', 'p', 'r', 'orange', 'banana'])

```

**移除元素**

```
fruits = {"apple", "orange"}

# 移除不存在的元素不会报错
fruits.discard("test") 

# 移除不存在的元素会报错
fruits.remove("apple") 

# 弹出一个元素，集合为空时会报错
fruits.pop()

# 清空集合
fruits.clear()
```

**2. 集合的交叉并补等操作**

pass

# 四、类型总结

## 4.1 类型对比

类型|定义|序列|可变类型|传引用|获取|
---|---|---|---|---|---
`数值类型`|`i=1`|-|否|是
`字符串`|`str="Hello"`|是|否|是
`list`|`[]`|是|是|是|索引
`dict`|`{}`|否|是|是|键
`set`|`set()`|否|是|是|
`tuple`|`()`|是|否|-|索引

列表、元组和字符串都是序列。字符串是字符的序列，列表和元祖是任意类型的序列。

序列的两个主要特点是索引操作符和切片操作符。索引操作符让我们可以从序列中抓取一个特定项目。切片操作符让我们能够获取序列的一个切片，即一部分序列。

```
fruits = "apple"

# 索引操作
print(fruits[1])

# 切片操作
print(fruits[1:-1])
print(fruits[1:])
print(fruits[:-1])
print(fruits[:])
```

## 4.2 类型转换

方法|说明
---|---
bool(x)|将x转bool类型
int(x)|将x转换为一个整数
float(x)|将x转换到一个浮点数
str(x)|将对象 x 转换为字符串
tuple(s)|将序列 s 转换为一个元组
list(s)|将序列 s 转换为一个列表
set(s)|将序列 s 转换为一个集合


# 五、函数

## 5.1 函数定义

**1. 常规函数**

定义一个空函数，由于没有花括号表示代码块，针对空的代码块可以用pass占位。

```
def func(x, y):
    pass
```

**2. 匿名函数**

函数定义：`lambda 参数: 表达式`

```
f = lambda x: x*x
print(f(2))
print((lambda x: x*x)(3))
```

**3. 闭包函数**

`pass`

## 5.2 函数参数 - 可变参数与关键字参数

`python`的参数传入确实是相当的方便，参数传入非常灵活。但可也可能会导致根据参数无法清楚的表达函数行为。

```
def func(x, y = 1, *args, **kwargs):
    print(type(args))
    print(type(kwargs))
    print(locals())

args = [1, 2]
kwargs = {"param": "web"}

func(1, 2, 3, '4', param="web")
# <type 'tuple'>
# <type 'dict'>
# {'y': 2, 'x': 1, 'args': (3, '4'), 'kwargs': {'param': 'web'}}

func(1, 2, *args, **kwargs)
# <type 'tuple'>
# <type 'dict'>
# {'y': 2, 'x': 1, 'args': (1, 2), 'kwargs': {'param': 'web'}}
```

- 支持设置默认参数
- 支持可变参数`*args`。函数接收到的数据类型是元祖。调用时可以以展平的方式传入，或者以列表、元祖的解引用的方式传入。
- 支持关键字参数`**kwargs`。函数接收到的数据类型是字典。

## 5.3 函数返回

**1. 多返回值**

```
def func(x, y = 1):
    return x,y

x = func(1)

print(type(x))
print(x)

# <type 'tuple'>
# (1, 1)


x, _ = func(1)
print(type(x))
print(x)

# <type 'int'>
# 1
```

`python`支持多个返回值，多个返回值实际返回的是一个元祖。多个参数时用一个参数接收时为元祖，用多个参数时可以直接将元祖解开得到具体的数据类型。

**2. 返回对象**

返回一个匿名函数。

```
def func(x, y = 1):
    # lambda 参数: 表达式
    return lambda t: x+y if t == 1 else x*y

x = func(2, 3)
print(x(1))
print(type(x))

# 5
# <type 'function'>
```

## 5.4 函数调用

`python`中**数字、字符、元组等不可变对象类型都属于值传递，而字典和列表等可变对象类型属于引用传递。**对于可变对象意味着函数内部可以修改实参的值。

```
import random

def func(x):
    x.append(random.randint(1, 100))
    print(x)

p = []
func(p)
func(p)
print(p)
# 打印[85, 86]，可见函数内部对p的修改直接影响了原有的值.
```


# 六、类 - Class

## 6.1 类定义

```
class Test(object):

    def __init__(self, x, y):
        self.x = x
        self.y = y

    def sum(self):
        return self.x + self.y


t = Test(1, 2)
print(t.sum())
```

## 6.2 访问限制

Python只有公有和私有，默认公有，当属性或方法前面有两个下划线（`__`）时，就只有内部可以访问，外部无法直接访问了。

## 6.3 成员属性

Python中并不需要在类中预定义属性，这样子会存在一个问题，并不太方便知道一个类有多少个成员属性。比如，在上面函数中增加一行：

```
    def sum(self):
        self.z = self.x + self.y
        return self.x + self.y
```

这样子在调用`sum`方法之后打印类对象：`print(vars(t))`，可以看到结果：`{'y': 2, 'x': 1, 'z': 3}`。这种不确定性在程序维护过程中会比较麻烦。所以对通用的属性可以考虑在构造函数`__init__`中做一个初始化。也方便知道类中有哪些可使用的属性。

针对属性也可以通过`@property`装饰器做一些限制。比如设置只读，对设置的值做一些检测。

```

class Test(object):

    def __init__(self):
        self._age = 18

    @property
    def age(self):
        return self._age

    @age.setter
    def age(self, value):
        if not isinstance(value, int):
            raise ValueError("age must be int.")
        self._age = value


t = Test()
print(t.age)
t.age = 20
print(t.age)
#设置非int会抛异常
#t.age = "123"
```

- 可以直接通过`示例.属性名`来获取，不需要带函数的括号。
- 不设置对应的`setter`方法时，该属性对外只读，修改会抛异常。


## 6.4 成员方法

成员方法有3中形式：实例方法、类方法、静态方法。

- 实例方法：需要显示的传入self对象参数，需要实例化类之后调用方法
- 类方法(`@classmethod`)：不需要实例化类，方法属于类
- 静态方法(`@staticmethod`): 相当于放在类里面的函数，从功能上看属于类，但不需要调用类里其他的成员属性和成员方法。

定义方式如下：

```
class Test(object):

    __n = 0

    @classmethod
    def show(cls):
        print("Hello World {0}".format(cls.__n))

    @staticmethod
    def sum(x, y):
        return x + y


Test.sum(1, 2)
Test.show()
```


# 七、错误和异常

```
try:
    t = Test(1, 'b')
    print(t.sum())
except (TypeError, NameError) as err:
    print(err)
except Exception as e:
    print(e)
else:
    raise Exception("throw exception")
```

# 八、高级用法

## 8.1 三元表达式

`[on true] if [expression] else [on false]`

等价于：

```
if [expression]:
    [on true]
else:
    [on false]
```

示例：

```
$ import random

$ x = True if random.randint(0, 10) > 5 else False
```

## 8.2. 推导式

下面是一个列表推导式，用来快速创建一个列表。可以分为列表推导式，字典推导式和集合推导式。

```
[ expression for x in X [if condition]
             for y in Y [if condition]
             ...
             for n in N [if condition] ]
```

相当于：

```
L = []
for x i X:
    [if condition]:
        for y in Y:
            [if condition]:
                L.append(x, y)

```

示例：

```
# 列表推导式
$ L = [x*y for x in range(1, 5) for y in range(6, 10) if y > 8]
# [9, 18, 27, 36]

# 字典推导式
$ L = {x: y for x in range(1, 5) for y in range(6, 10) if y > 8}
# {1: 9, 2: 9, 3: 9, 4: 9}

# 集合推导式
$ L = {x * y for x in range(1, 5) for y in range(6, 10) if y > 8}
# set([9, 18, 27, 36])
```

## 8.3 map && fliter函数

**map函数**

`map`为内置函数，用于遍历序列，然后将函数用于遍历过程中的每一个元素。函数定义：

```
map(function, sequence[, sequence, ...]) -> list
```

- function: 处理函数
- sequence：一个或多个序列
- 返回值： list

示例一(字符串也是一种序列)：

```
>>> map(ord, "Hello World")
[72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100]
```

利用上面的列表推导式也可以实现：

```
>>> [ ord(x) for x in "Hello World"]
[72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100]
```

如果是多个序列，会同时将每个序列的元素拿出一起传给函数。如果长度不一致会用None补齐。

```
>>> map(lambda x, y :  (x, y) , "Hello", "Python")
[('H', 'P'), ('e', 'y'), ('l', 't'), ('l', 'h'), ('o', 'o'), (None, 'n')]
```

**filter函数**

`filter`也为内置函数，将序列的每个元素传给`function`，然后将函数执行返回`True`的元素组成新的列表。

```
filter(function or None, sequence) -> list, tuple, or string
```

如过滤掉字母`o`

```
>>> filter(lambda x: x != "o", "Hello World")
'Hell Wrld'
```

## 8.4 sort && sorted


`sort` 与 `sorted` 区别：`sort` 是应用在 `list` 上的方法，`sorted` 可以对所有可迭代的对象进行排序操作。

`list` 的 `sort` 方法返回的是**对已经存在的列表进行操作，无返回值**，而内建函数 `sorted` 方法**返回的是一个新的 list**，而不是在原来的基础上进行的操作。

关键点：

- `sort`改变原列表，`sorted`不会改变原列表
- `sort`只用于列表，`sorted`用于所有可迭代对象

**sort**

```
>>> print(list.sort.__doc__)
L.sort(cmp=None, key=None, reverse=False) -- stable sort *IN PLACE*;
cmp(x, y) -> -1, 0,
```

示例(具体参数使用参考`sorted`)：

```
>>> a = [1, 3, 2, 4]
>>> a.sort()
>>> a
[1, 2, 3, 4]
```

**sorted**

```
>>> print(sorted.__doc__)
sorted(iterable, cmp=None, key=None, reverse=False) --> new sorted list
```

参数说明：

- iterable -- 可迭代对象。
- cmp -- 比较的函数，这个具有两个参数，参数的值都是从可迭代对象中取出，此函数必须遵守的规则为，大于则返回1，小于则返回-1，等于则返回0。
- key -- 也是指定一个函数，cmp用来指定比较方法， key用来指定该用哪个key做比较
- reverse -- 排序规则，reverse = True 降序 ， reverse = False 升序（默认）。

示例一，指定cmp函数：

```
def mycmp(x, y):
    print(x,y)
    if ord(x) > ord(y):
        return 1
    elif ord(x) == ord(y):
        return 0
    else:
        return -1

a = "b2a1c3"

b = sorted(a, cmp=mycmp)

print(a, b)

# Output
# ('2', 'b')
# ('a', '2')
# ('a', 'b')
# ('a', '2')
# ('1', 'a')
# ('1', '2')
# ('c', 'a')
# ('c', 'b')
# ('3', 'a')
# ('3', '2')
# ('b2a1c3', ['1', '2', '3', 'a', 'b', 'c'])
```

示例二，指定key函数

```
a = [{"name": "Jack", "age": 30}, {"name": "Peter", "age": 18}, {"name": "Amy", "age": 24}]

b = sorted(a, key=lambda x: x["age"])

print(a)
print(b)


# [{'age': 30, 'name': 'Jack'}, {'age': 18, 'name': 'Peter'}, {'age': 24, 'name': 'Amy'}]
# [{'age': 18, 'name': 'Peter'}, {'age': 24, 'name': 'Amy'}, {'age': 30, 'name': 'Jack'}]
```
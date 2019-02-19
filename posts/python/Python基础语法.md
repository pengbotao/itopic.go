```
{
    "url": "python-language",
    "time": "2015/12/21 19:51",
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
a = "Hello World"
b = 'iTopic.org'
```

**字符串格式化**

```
"Hello %s %d" % ("world", 2019)
"Hello {0} {1}".format("world", 2019)
```

## 1.3 类型转换

方法|说明
---|---
int(x)|将x转换为一个整数
float(x)|将x转换到一个浮点数
str(x)|将对象 x 转换为字符串
tuple(s)|将序列 s 转换为一个元组
list(s)|将序列 s 转换为一个列表
set(s)|将序列 s 转换为一个集合


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

**1. 初始化及遍历**

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

```
# 初始化元祖。只有1个元素时最后携带逗号
fruits = ("apple", "orange", "banana")

print(fruits[0])

# 遍历元祖
for val in fruits:
    print(val)
```

## 3.3 字典 - dict

```
# 初始化字典
fruits = {"apple": 1, "orange": 2, "banana": 3}

# 按KEY遍历
for idx in fruits:
    print(fruits[idx])

if "grape" in fruits:
    print(fruits["grape"])
else:
    fruits["grape"] = 4

if fruits.has_key("apple"):
    print(fruits["apple"])

# 按KEY - VALUE遍历
for idx, val in fruits.items():
    print(idx, val)
```

- update方法

## 3.4 集合 - set

- 集合里的VALUE可以是不同的类型，相同数据会去重。
- 集合无序

```
fruits = set(["apple", "orange", "banana"])
print(fruits)
```

## 3.5 总结

列表、元组和字符串都是序列。字符串是字符的序列，列表和元祖是任意类型的序列。序列的两个主要特点是索引操作符和切片操作符。索引操作符让我们可以从序列中抓取一个特定项目。切片操作符让我们能够获取序列的一个切片，即一部分序列。

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

**对比**

类型|定义|有序|可变|传引用|获取|
---|---|---|---|---|---
`list`|`[]`|是|是|是|索引
`tuple`|`()`|是|否|是|索引
`dict`|`{}`|否|是|是|键
`set`|`set()`|否|是|是|


---

关于数据结构可查看：`http://www.pythondoc.com/pythontutorial3/datastructures.html#`

# 四、函数

## 4.1 函数定义

### 4.1.1 常规函数
定义一个空函数，由于没有花括号表示代码块，针对空的代码块可以用pass占位。

```
def func(x, y):
    pass
```

### 4.1.2 匿名函数

函数定义：`lambda 参数: 表达式`

```
f = lambda x: x*x
print(f(2))
print((lambda x: x*x)(3))
```

### 4.1.3 闭包函数

`pass`

## 4.2 函数参数 - 可变参数与关键字参数

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

## 4.3 函数返回

### 4.3.1 多返回值
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

### 4.3.2 返回对象

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

## 4.4 函数调用

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

## 4.5 函数装饰器

`pass`

# 五、类 - Class

## 5.1 类定义

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

## 5.2 包与导入


# 六、错误和异常

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

# 七、语法糖

## 7.1 三元表达式

`[on true] if [expression] else [on false]`

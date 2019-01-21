```
{
    "url": "python-language",
    "time": "2015/12/10 19:51",
    "tag": "Python"
}
```

# 一、概述

- 单引号与双引号的区别

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

```
# 初始化列表
fruits = ["apple", "orange", "banana"]

# 遍历列表
for val in fruits:
    print(val)

fruits.append("grape")

# 按索引 - 值遍历，也可用于字典
for idx, val in enumerate(fruits):
    print(idx, val)

for idx in range(len(fruits)):
    print(fruits[idx])
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

## 3.4 序列

列表、元组和字符串都是序列。序列的两个主要特点是索引操作符和切片操作符。索引操作符让我们可以从序列中抓取一个特定项目。切片操作符让我们能够获取序列的一个切片，即一部分序列。

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

---

关于数据结构可查看：`http://www.pythondoc.com/pythontutorial3/datastructures.html#`

# 四、函数

# 五、Class

# 六、错误和异常

# 七、模块
url: python-start
des: 
time: 2015/12/07 06:14
tag: python
++++++++

许多人谈起Python第一感觉就是他强制缩进的语法风格，代码块用缩进来表示，而取消了经典的花括号。同一代码块中的代码需要有相同的空白字符来缩进，可以是三个空格，也可以是两个tab，还可以是空格+tab，所以下面的代码也是可以执行的：
```
def a() :
   print "#3个空格"
def b() :
        print "#两个tab"
def c() :
            print "#四个空格 + 两个TAB"
            if True :
                print '#四个TAB'
```
看来也可以用空白字符来玩点恶作剧。不过，一般情况下用4个空格来缩进，TAB替换为4个空格。对于有过其他编程经验的人来说，经典花括号的去掉会比缩进更折磨人。一个函数括号写完的瞬间，习惯性的回车加花括号，那酸爽！呵呵了。。

# Python2.x 还是 Python3.x ?

Py3K说的是比2做了一个较大的升级，有些语法不向下兼容。比如上面的程序在3下运行会报错，因为3把print改成了函数，函数调用需要用括号扩起。可能还有一些其他的高级特性等等。然后就涉及到究竟选3还是选2的问题了？以前看Python也纠结过这个问题，这次重新看貌似没有太纠结，可能是项目比较紧没有时间去考虑。不过知乎上说的好：无非就语法上的少量差异，而在编程中这些只是细枝末节的东西。思想是通的，稍微花点时间学些下就好了。

# Python 命名规范

突然觉得PHP中的美元符号挺好的，可以明确区分到这是一个变量。在Python中变量没有任何修饰符，而且可以随时替换为其他类型，看下面这个例子：
```
import time as a
 
print a.time()
a = 1
def a():
    print 'A func'
a()
 
class a() :
    def __init__(self):
        print 'A class'
a()
```
看到Python动态类型的时候觉得挺好，但看到上面这段又觉得好像挺容易出问题的，可能不小心就定义了一个time变量把导入的time给覆盖了。

# Python 类

这个可能也是比较想吐槽的，看下面几行代码：
```
class a :
 
    def b(self) :
        self.p = True
        print self.p
    def c(a) :
        print a.p
x = a()
x.b()
x.c()
print x.p
```
类的成员属性可以不用声明直接赋值，比如上面的self.p，p即为类的成员属性。可是，不声明就不知道究竟定义了多少成员属性。还有就是类中的方法第一个参数必须显式的传入类实例对象，如上面的self，a都表示当前实例。可类中为什么还要传入当前实例对象呢？不能隐式的传入么?再就是类里也没有public、private等关键词来修饰成员方法、属性，所有的方法都可以公开的访问到。

# Python基本语法
## 注释
注释分单行注释和多行注释， 单行注释在代码前面加 # 即可；多行注释则为前后加三个单引号
```
'''
#Python HelloWorld类和say方法
class HelloWorld:
    def say(self):
        return input('Please Input A Name:')
 
HelloWorld = HelloWorld();
while(True):
    print('Hello,', HelloWorld.say())
'''
```
## 变量
变量的定义和PHP类似，不用定义，直接使用即可，但别忘了Python的变量命不是以$打头的
## 函数
函数的定义如上类中的say， 使用def定义
## 字符串
字符串的操作比较好理解，和PHP类似，可以进行当数组引用某一个值，也可进行字符串的链接
## 数组
分为元组(tuple)、列表(list)和字典(dictionary)

- **元组**：固定的数组，定义后元素个数不可修改，定义方式为 (元素) 如 `("a", "b", "c")`
- **列表**：可动态的增加或删除数组中的元素，定义方式为 [元素] 如 `["a", "b", "c"]`
- **字典**：好比PHP中的关联数组，为`Key:Value`的键值对，如 `{"a":"apple", "b":"banana"}`

## 控制语句
基本的用法如下，嵌套规则跟其他语言一样
```
A = 5
# IF ... ELSE ...
if(A > 10):
    print('A>10')
elif(A > 0):
    print('A>0 && A <=10')
else:
    print('A<=0')
 
# WHILE ... (ELSE)
while(A > 0) :
    print(A)
    A -= 1
else:
    print('ELSE:', A)
 
A = ["apple", "banana"]
# for
for fruit in A:
    print(fruit)
```

最后，Python究竟有啥优点？[知乎一下](https://www.zhihu.com/question/25038841)得到的结果是：

- PHP：没有优点
- Java：库多，库多，库多
- Python：语法清楚，语法清楚，语法清楚
- C：能操纵底层，能细粒度优化性能
- C++：啥都有，啥都有，啥都有

关于PHP有句屌爆了的話：“PHP是世界上最好的语言，没有之一”

Python也有这么一句：“人生苦短，我用Python”
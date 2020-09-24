```
{
    "url": "lua-language",
    "time": "2016/01/10 12:55",
    "tag": "Lua",
    "toc": "yes"
}
```

> Lua 是一个小巧的脚本语言。是巴西里约热内卢天主教大学（Pontifical Catholic University of Rio de Janeiro）里的一个研究小组，由Roberto Ierusalimschy、Waldemar Celes 和 Luiz Henrique de Figueiredo所组成并于1993年开发。 其设计目的是为了嵌入应用程序中，从而为应用程序提供灵活的扩展和定制功能。Lua由标准C编写而成，几乎在所有操作系统和平台上都可以编译，运行。Lua并没有提供强大的库，这是由它的定位决定的。所以Lua不适合作为开发独立应用程序的语言。Lua 有一个同时进行的GIT项目，提供在特定平台上的即时编译功能。

> Lua脚本可以很容易的被C/C++ 代码调用，也可以反过来调用C/C++的函数，这使得Lua在应用程序中可以被广泛应用。不仅仅作为扩展脚本，也可以作为普通的配置文件，代替XML,ini等文件格式，并且更容易理解和维护。 Lua由标准C编写而成，代码简洁优美，几乎在所有操作系统和平台上都可以编译，运行。 一个完整的Lua解释器不过200k，在目前所有脚本引擎中，Lua的速度是最快的。这一切都决定了Lua是作为嵌入式脚本的最佳选择 

# 一、Lua注释
```
-- 单行注释：使用两个减号表示此行被注释
print("Hello Lua.")
 
--[[
多行注释。此范围内的行均被注视掉。
]]
print("Hello World.")
```

# 二、Lua数据类型
Lua是动态类型语言，变量不需要类型定义。Lua中有8个基本类型，分别是：

类型   |    说明
--- | ---
nil | Lua中的特殊类型，变量没赋值前默认值为nil；给变量赋值为nil可以删除该变量。
boolean | 布尔类型。可取值true和false
number | 数字类型
string | 字符串，Lua中字符串可以包含任何字符且字符串不可修改。
table | 表。类似其他语言中的数组、字典。
function | 函数类型。Lua可以调用Lua或者C实现的函数，Lua所有标准库都是C实现的。标准库包括string库、table库、I/O库、OS库、算术库、debug库。
userdata | 这个类型专门与Lua的宿主打交道。宿主通常是由c语言和c++语言开发的，在这种情况下，userdata可以是宿主的任何类型，常用的是结构体和指针类型
thread | 线程类型

# 三、Lua表达式
## 3.1 算术运算符
+(加) -(减) *(乘) /(除) ^(乘方) %(求模)

## 3.2 关系运算符
\>(大于) <(小于) >=(大于等于) <=(小于等于) ==(等于) ~=(不等于，这个有点不一样！) 

## 3.3 逻辑运算符
and or not. 逻辑运算符认为false和nil是假(false)，其他为真，0也是true！！！and 和 or的运算结果不是true和false，而是和它的两个操作数相关。
```
a and b -- 如果a为false，则返回a，否则返回b
a or b  -- 如果a为true，则返回a，否则返回b
```

## 3.4 连接运算符
.. 两个点表示字符串连接，如果操作数为数字，Lua将数字转换成字符串。`print("Hello " .. "Lua")`

# 四、Lua基本语法
## 4.1 变量
变量可以不用声明直接使用，给变量赋值即创建了这个变量。默认情况下Lua的所有变量都是全局变量，如果需要声明局部变量可以在前面加上local，如：
```
-- 全局变量
a = 1
 
-- 局部变量
local b = 2
 
-- 等价于x = 1;y = 2
x,y = 1,2
 
-- 交换x和y的值
x,y = y,x
```

## 4.2 if语句
```
if mon == 1 then
    print("Jan")
elseif mon == 2 then 
    print("Feb")
else 
    print("Other")
end
```

## 4.3 while语句
Lua跟其他常见语言一样，提供了while控制结构，语法上也没有什么特别的。但是没有提供do-while型的控制结构,但是提供了功能相当的repeat。
```
a=10
while a < 20 do
   print("value of a:", a)
   a = a+1
end
```
值得一提的是，Lua 并没有像许多其他语言那样提供类似 continue 这样的控制语句用来立即进入下一个循环迭代（如果有的话）。因此，我们 需要仔细地安排循环体里的分支，以避免这样的需求。

## 4.4 repeat-until语句
Lua中的repeat控制结构类似于其他语言（如：C++语言）中的do-while，但是控制方式是刚好相反的。简单点说，执行repeat循环体后，直到until的条件为真时才结束，而其他语言（如：C++语言）的do-while则是当条件为假时就结束循环。
```
a = 10
repeat
   print("value of a:", a)
   a = a + 1
until a > 15
```

## 4.5 for语句
for语句有两种形式：数字for（numeric for）和范型for（generic for）。
### 4.5.1 数值for循环
```
for init,max/min value, increment
do
   statement(s)
end
```
下面是控制在一个循环的流程：

- 1、始化步骤首先被执行，并且仅一次。这个步骤可让您声明和初始化任何循环控制变量。
- 2、接着是max/min，这是最大或最小值，直到该循环继续执行。它在内部创建了一个条件检查的初值和最大值/最小值之间进行比较。
- 3、for循环体执行后，控制流跳回至递增/递减声明。这个语句可以更新任何循环控制变量。
- 4、条件现在重新计算评估。如果这为真则循环执行，并重复这个过程(循环体，然后增加一步，然后再条件)。如果条件为假，则循环终止。

示例：
```
-- 打印 10 ... 1
for i=10,1,-1 
do 
   print(i) 
end
```
### 4.5.2 范型for循环
泛型for循环通过一个迭代器（iterator）函数来遍历所有值：
```
-- print all values of array'a'
for i,v in iparis(a) do print(v) end
 
-- print all keys of table t
for k in paris(t) do print(k) end
```

## 4.6 break和return语句
break语句用来退出当前循环（for、repeat、while）。在循环外部不可使用。return从函数返回结果。
Lua语法要求break和return只能出现在block的结尾一句（也就是，chunk的最后一句，end之前，else前或者until前）。

# 五、Lua table
table 在 Lua 里是一种重要的数据结构，它可以说是其他数据结构的基础，通常的数组、记录、线性表、队列、集合等数据结构都可以用 table 来表示，甚至连全局变量（_G）、模块、元表（metatable）等这些重要的 Lua 元素都是 table 的结构。可以说，table 是一个强大而又神奇的东西。table通过两个花括号来构造。如：
```
days = {"Sunday", "Monday", "Tuesday"}
 
for idx, day in ipairs(days) do
    print(idx, day)
end
```
可以用任意类型的值来作数组的索引，但这个值不能是 nil

在构造函数中域分隔符逗号（","）可以用分号（";"）替代，通常我们使用分号用来分割不同类型的表元素。
```
{x=10, y=45; "one", "two", "three"}
```
所有索引值都需要用 "["和"]" 括起来；如果是字符串，还可以去掉引号和中括号； 即如果没有[]括起，则认为是字符串索引
```
day = {1,2,3,4,5,6,7}
days = {[day] = "table", sun="Sunday", ["mon"] = "Monday", 123}
print(days[day])
print(days["sun"])
print(days.mon)
print(days[1])
```
如果不写索引，则索引就会被认为是数字，并按顺序自动从 1往后编；
```
days = {"Sunday", "Monday",}
for i=1, #days do 
    print(days[i])
end
```
用table时，对于字符串，可以通过.的方式访问，也可以通过[]方式访问
```
days = {sun = "Sunday", mon = "Monday",}
print(days.sun)
print(days["mon"])
```
# 六、Lua函数

函数可以存储在变量中，可以通过参数传递给其他函数，或者作为函数的返回值（类比C/C++中的函数指针），这种特性使Lua具有极大的灵活性。Lua对函数式编程提供了良好的支持，可以支持嵌套函数。

另外，Lua既可以调用Lua编写的函数，还可以调用C语言编写的函数（Lua所有的标准库都是C语言写的）。
```
function func_name(arguments-list)
    statements-list
end
```
Lua函数实参和形参与赋值语句类似，多余部分被忽略，缺少部分用nil补足。函数调用时，当只有一个参数并且这个参数是字符串或者表构造时括号是可选的。如，
```
print "Hello Lua"
f{x = 10, y = 20}
```
## 6.1 Lua函数不支持参数默认值，可以通过or来实现。
```
function sum(a , b)
    a = a or 1.1
    b = b or 2.2
    return a + b
end
 
print(sum())
```

## Lua函数支持多个返回值
```
function str_replace()
    return "Hello Lua", 1
end
 
a, b = str_replace()
print(a,b)
```

## 6.2 Lua函数可以支持变长参数
Lua函数可以接受可变数目的参数，通过三个点(...)表示函数有可变的参数。
```
function sum(...)
    for i,v in pairs({...}) do
        print(i, v)
    end
end
 
sum(2,3,4)
```
通常在遍历变长参数的时候只需要使用{…}，然而变长参数可能会包含一些nil；那么就可以用select函数来访问变长参数了：`select('#', …)`或者 `select(n, …)`

`select('#', …)`返回可变参数的长度，`select(n,…)`用于访问n到`select('#',…)`的参数

## 6.3 Lua函数支持命名参数
Lua的函数参数是和位置相关的，调用时实参会按顺序依次传给形参。有时候用名字指定参数是很有用的，比如rename函数用来给一个文件重命名，有时候我们我们记不清命名前后两个参数的顺序了：
```
rename{old="temp.lua", new="temp1.lua"}
```
Lua可以通过将所有的参数放在一个表中，把表作为函数的唯一参数来实现上面这段伪代码的功能。因为Lua语法支持函数调用时实参可以是表的构造。

## 6.4 Lua函数支持闭包
```
tab = {1,2,3,4,5,6,7,8}
 
function iter()
    local index = 0
    return function()
        index = index + 1 
        return tab[index]
    end
end
 
for i in iter() do
   print(i)
end
```

## 6.5 Lua函数可以放在变量、表中
```
local sum = function(...) 
    local s = 0
    for _, v in pairs({...}) do
        s = s + v
    end
    return s
end
print(sum(1,2,3,5))
 
op = {add = sum}
print(op.add(1,2,3))
```
同样，函数也可以作为参数传递或者作为返回值返回。
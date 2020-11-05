```
{
    "url": "shell-start",
    "time": "2018/10/01 21:24",
    "tag": "Shell,运维",
    "toc": "yes"
}
```

# 一、基本代码

```
#! /bin/bash
num=`ps -ef |grep zookeeper | grep -v grep | wc -l`

# 判断进程是否存在，如果不存在则启动，否则输出Running
if [ $num -lt 1 ];then
	su - peng -c "/Users/peng/data/apache-zookeeper-3.6.2-1/bin/zkServer.sh start"
else
	echo "It's already running"
fi
```

## 1.1 变量定义

**变量名与等号之间不能有空格**，只能使用字母、数字、下划线，不能以数字开头，不能使用bash里的关键字。

## 1.2 IF语句

条件判断条件判断方括号与变量以及变量与运算符之间需要有空格，`then`和`if`可以换行，放一行需要用分号分隔，最后以`fi`结束，如果要跟多个`elif`的示例：

```
#! /bin/bash

idx=10

if [ $idx -gt 3 ]
then 
    echo "$idx > 3"
elif [ $idx -gt 2 ]; then echo "$idx > 2"; elif [ $idx -gt 1 ]; then
    echo "${idx} > 1"
else
    echo "${idx} <= 1"
fi
```

## 1.3 反引号

执行命令并将命令的输出赋值给变量`num`，反引号的作用和`$(命令)`是一样的，借助该功能就实现可以强大的外部程序调用。

```
num=$(ps -ef |grep zookeeper | grep -v grep | wc -l)
```

## 1.4 运算符

**1. 关系运算符：**

| 运算符 | 说明                                                  | 举例                         |
| ------ | ----------------------------------------------------- | ---------------------------- |
| `-eq`  | 检测两个数是否相等，相等返回 true。                   | `[ $a -eq $b ]` 返回 true。  |
| `-ne`  | 检测两个数是否相等，不相等返回 true。                 | `[ $a -ne $b ]` 返回 true。  |
| `-gt`  | 检测左边的数是否大于右边的，如果是，则返回 true。     | `[ $a -gt $b ]` 返回 false。 |
| `-lt`  | 检测左边的数是否小于右边的，如果是，则返回 true。     | `[ $a -lt $b ]` 返回 true。  |
| `-ge`  | 检测左边的数是否大等于右边的，如果是，则返回 true。   | `[ $a -ge $b ]` 返回 false。 |
| `-le`  | 检测左边的数是否小于等于右边的，如果是，则返回 true。 | `[ $a -le $b ]` 返回 true。  |

**2 逻辑运算符：**假定变量 a 为 10，变量 b 为 20:

| 运算符 | 说明       | 举例                                        |
| :----- | :--------- | :------------------------------------------ |
| `&&`     | 逻辑的 AND | `[[ $a -lt 100 && $b -gt 100 ]]` 返回 false |
| \|\|   | 逻辑的 OR  | |

**3. 布尔运算符：**

| 运算符 | 说明                                                | 举例                                       |
| ------ | --------------------------------------------------- | ------------------------------------------ |
| `!`    | 非运算，表达式为 true 则返回 false，否则返回 true。 | `[ ! false ]` 返回 true。                  |
| `-o`   | 或运算，有一个表达式为 true 则返回 true。           | `[ $a -lt 20 -o $b -gt 100 ]`返回 true。   |
| `-a`   | 与运算，两个表达式都为 true 才返回 true。           | `[ $a -lt 20 -a $b -gt 100 ]` 返回 false。 |

**4. 字符串运算符**

| 运算符 | 说明                                      | 举例                       |
| ------ | ----------------------------------------- | -------------------------- |
| `=`    | 检测两个字符串是否相等，相等返回 true。   | `[ $a = $b ]` 返回 false。 |
| `!=`   | 检测两个字符串是否相等，不相等返回 true。 | `[ $a != $b ] `返回 true。 |
| `-z`   | 检测字符串长度是否为0，为0返回 true。     | `[ -z $a ]` 返回 false。   |
| `-n`   | 检测字符串长度是否为0，不为0返回 true。   | `[ -z $a ]` 返回 true。    |
| `str`  | 检测字符串是否为空，不为空返回 true。     | `[ $a ]` 返回 true。       |

## 1.5 注释

注释以`#`开头，只支持单行注释。

# 二、Case语句

```
#! /bin/bash

case "$1" in
  start)
    echo "start process..."
    ;;
  stop)
    echo "stop process..."
    ;;
  *)
    echo $"Usage: $0 {start|stop}"
esac
```

## 2.1 特殊参数

`$1`为启动脚本时跟着脚本后面的第一个参数。

| 变量 | 含义                                                         |
| ---- | ------------------------------------------------------------ |
| `$0` | 当前脚本的文件名                                             |
| `$n` | 传递给脚本或函数的参数。n 是一个数字，表示第几个参数。例如，第一个参数是$1，第二个参数是$2。 |
| `$#` | 传递给脚本或函数的参数个数。                                 |
| `$*` | 传递给脚本或函数的所有参数。                                 |
| `$@` | 传递给脚本或函数的所有参数。被双引号(" ")包含时，与 `$*` 稍有不同，下面将会讲到。 |
| `$?` | 上个命令的退出状态，或函数的返回值。                         |
| `$$` | 当前Shell进程ID。对于 Shell 脚本，就是这些脚本所在的进程ID。 |

```
#!/bin/bash
echo "File Name: $0"
echo "First Parameter : $1"
echo "First Parameter : $2"
echo "Quoted Values: $@"
echo "Quoted Values: $*"
echo "Total Number of Parameters : $#"
```

执行后输出：

```
$ ./1.sh -host 127.0.0.1 -port 80
File Name: ./1.sh
First Parameter : -host
First Parameter : 127.0.0.1
Quoted Values: -host 127.0.0.1 -port 80
Quoted Values: -host 127.0.0.1 -port 80
Total Number of Parameters : 4
```

**`$*` 和 `$@` 的区别**

`$*` 和 `$@` 都表示传递给函数或脚本的所有参数，不被双引号(`" "`)包含时，都以`"$1"` `"$2"` … `"$n"` 的形式输出所有参数。

但是当它们被双引号(`" "`)包含时，`"$*"` 会将所有的参数作为一个整体，以`"$1 $2 … $n"`的形式输出所有参数；`"$@"` 会将各个参数分开，以`"$1"` `"$2"` … `"$n"` 的形式输出所有参数。

## 2.2 Case语句

看示例用法，`;;`相当于`break`跳出`case`语句。

# 三、循环语句

## 3.1 for循环

```
#! /bin/bash

for loop in 1 2 3 4 5
do 
    echo "The value is $loop"
done
```

## 3.2 while循环

```
#! /bin/bash

# 重启服务过程中查看进程是否已经关闭
count=$(ps -ef | grep elasticsearch | grep -v grep | wc -l)
times=1
while [ $count -gt 0 -a $times -lt 30 ]; do
    echo "count: $count, times: $times"
    sleep 1
    let times++
    count=$(ps -ef | grep elasticsearch | grep -v grep | wc -l)
done

if [ $times -ge 30 ]; then
    echo "process closed failed"
    exit 1
fi

echo "restart program"
```

`let times++`可以实现自增，也可以写成：

```
times=`expr $times + 1`
```

# 四、函数

```
#! /bin/bash

function start() {
    echo "Start: $1 $2"
}

stop() {
    echo "Stop: $1"
    return 0
}

case "$1" in
  start)
    start $2 $3
    ;;
  stop)
    stop $2
    ;;
  *)
    echo $"Usage: $0 {start|stop}"
esac
```

- `function`字样可有可无
- 不需要写形参定义，参数传递同章节`2.1`
- 函数调用上不用打括号，函数名后面跟参数，空格分隔

调用示例：

```
pengbotao:Desktop peng$ ./1.sh start app 01
Start: app 01
pengbotao:Desktop peng$ ./1.sh stop app 01
Stop: app
```






- [1] [Linux Shell脚本教程：30分钟玩转Shell脚本编程](http://c.biancheng.net/cpp/shell/)
- [2] [Runoob Shell教程](https://www.runoob.com/linux/linux-shell.html)
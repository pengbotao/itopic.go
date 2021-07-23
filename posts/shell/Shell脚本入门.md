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

注释以`#`开头，只支持单行注释。

## 1.1 变量定义

**变量名与等号之间不能有空格**，只能使用字母、数字、下划线，不能以数字开头，不能使用bash里的关键字。

## 1.2 if语句

注意，**条件判断方括号与变量以及变量与运算符之间需要有空格**，`then`和`if`可以换行，放一行需要用分号分隔，最后以`fi`结束，如果要跟多个`elif`的示例：

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

## 1.4 关系运算符

| 运算符 | 说明                                                  | 举例                         |
| ------ | ----------------------------------------------------- | ---------------------------- |
| `-eq`  | 检测两个数是否相等，相等返回 true。                   | `[ $a -eq $b ]` 返回 true。  |
| `-ne`  | 检测两个数是否相等，不相等返回 true。                 | `[ $a -ne $b ]` 返回 true。  |
| `-gt`  | 检测左边的数是否大于右边的，如果是，则返回 true。     | `[ $a -gt $b ]` 返回 false。 |
| `-lt`  | 检测左边的数是否小于右边的，如果是，则返回 true。     | `[ $a -lt $b ]` 返回 true。  |
| `-ge`  | 检测左边的数是否大等于右边的，如果是，则返回 true。   | `[ $a -ge $b ]` 返回 false。 |
| `-le`  | 检测左边的数是否小于等于右边的，如果是，则返回 true。 | `[ $a -le $b ]` 返回 true。  |

## 1.5 布尔运算符

| 运算符 | 说明                                                | 举例                                       |
| ------ | --------------------------------------------------- | ------------------------------------------ |
| `!`    | 非运算，表达式为 true 则返回 false，否则返回 true。 | `[ ! false ]` 返回 true。                  |
| `-o`   | 或运算，有一个表达式为 true 则返回 true。           | `[ $a -lt 20 -o $b -gt 100 ]`返回 true。   |
| `-a`   | 与运算，两个表达式都为 true 才返回 true。           | `[ $a -lt 20 -a $b -gt 100 ]` 返回 false。 |

> 注：示例中多个条件为一个中括号

## 1.6 逻辑运算符

假定变量 a 为 10，变量 b 为 20:

| 运算符 | 说明       | 举例                                        |
| :----- | :--------- | :------------------------------------------ |
| `&&`     | 逻辑的 AND | `[[ $a -lt 100 && $b -gt 100 ]]` 返回 false |
| \|\| | 逻辑的 OR  | |

> 注：示例中多个条件为2个中括号

这里有两种写法：

```
if [[ $a < 20 || $b > 100 ]]; then
    echo "$a < 20 or $b > 100"
fi


if [ $a -lt 20 -o $b -gt 100 ]; then
    echo "$a < 20 or $b > 100"
fi
```

在`[[`中使用`&&`和`||`表示逻辑与和逻辑或。`[`中使用`-a` 和`-o` 表示逻辑与和逻辑或。

## 1.7 字符串运算符

| 运算符 | 说明                                      | 举例                       |
| ------ | ----------------------------------------- | -------------------------- |
| `=`    | 检测两个字符串是否相等，相等返回 true。   | `[ $a = $b ]` 返回 false。 |
| `!=`   | 检测两个字符串是否相等，不相等返回 true。 | `[ $a != $b ] `返回 true。 |
| `-z`   | 检测字符串长度是否为0，为0返回 true。     | `[ -z $a ]` 返回 false。   |
| `-n`   | 检测字符串长度是否为0，不为0返回 true。   | `[ -n $a ]` 返回 true。    |
| `str`  | 检测字符串是否为空，不为空返回 true。     | `[ $a ]` 返回 true。       |

## 1.8 文件判断

使用不同的条件标志测试不同的文件系统属性。

| 操作符             | 意义                                                         |
| :----------------- | :----------------------------------------------------------- |
| `[ -f $file_var ]` | 变量 $file_var 是一个正常的文件路径或文件名 (file)，则返回真 |
| `[ -x $var ]`      | 变量 $var 包含的文件可执行 (execute)，则返回真               |
| `[ -d $var ]`      | 变量 $var 包含的文件是目录 (directory)，则返回真             |
| `[ -e $var ]`      | 变量 $var 包含的文件存在 (exist)，则返回真                   |
| `[ -c $var ]`      | 变量 $var 包含的文件是一个字符设备文件的路径 (character)，则返回真 |
| `[ -b $var ]`      | 变量 $var 包含的文件是一个块设备文件的路径 (block)，则返回真 |
| `[ -w $var ]`      | 变量 $var 包含的文件可写(write)，则返回真                    |
| `[ -r $var ]`      | 变量 $var 包含的文件可读 (read)，则返回真                    |
| `[ -L $var ]`      | 变量 $var 包含是一个符号链接 (link)，则返回真                |

使用方法如下：

```
fpath="/etc/passwd"
if [ -e $fpath ]; then
  echo File exits;
else
  echo Does not exit;
fi
```

## 1.9 特殊参数

`$1`为启动脚本时跟着脚本后面的第一个参数。

| 变量 | 含义                                                         |
| ---- | ------------------------------------------------------------ |
| `$0` | 当前脚本的文件名                                             |
| `$n` | 传递给脚本或函数的参数。`n` 是一个数字，表示第几个参数。例如，第一个参数是`$1`，第二个参数是`$2`。 |
| `$#` | 传递给脚本或函数的参数个数。                                 |
| `$*` | 传递给脚本或函数的所有参数。                                 |
| `$@` | 传递给脚本或函数的所有参数。被双引号(`" "`)包含时，与 `$*` 稍有不同，下面将会讲到。 |
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

> 注意：函数参数传递也是通过这些特殊变量

# 二、代码块

## 3.1 Case语句

常用脚本中根据传入参数做启动与重启。

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

看示例用法，`;;`相当于`break`跳出`case`语句。

## 3.2 for循环

```
#! /bin/bash

for loop in 1 2 3 4 5
do 
    echo "The value is $loop"
done
```

`in`后的内容是一组值（数字、字符串等）组成的序列，每个值通过空格分隔。每循环一次，就将列表中的下一个值赋给变量。循环内部也可以使用`break`或`continue`跳出循环。

## 3.2 while循环

服务发布过程中判断服务进程是否已退出，如果没有退出则等待30s。

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

# 五、字符串与数组

## 5.1 字符串

```
#! /bin/bash

str1='Hello, $1'
str2="Hello $1"

echo $str1
echo $str2
```

执行后可以看到：

```
$ ./1.sh World
Hello, $1
Hello World
```

- 单引号、双引号都行，单引号不会转义，双引号会转义
- 单引号中不能有单引号，对单引号转义也不行
- 双引号中可以变量

## 5.2 字符串拼接

```
echo $str1 $str2
```

## 5.3 字符串长度 

```
string="abcd"
echo ${#string} #输出 4
```

## 5.4 提取子字符串

```
string="alibaba is a great company"
echo ${string:1:4} #输出liba
```

## 5.5 数组

```
#! /bin/bash

im=(qq weixin)

# 查看元素内容
echo ${im[0]}
# 修改或新增元素
im[1]="dingding"
echo ${im[1]}
```

- 用括号来表示数组，元素之间用空格分隔
- 数组长度可调整

## 5.6 获取数组长度

```
echo ${#im[*]}
echo ${#im[@]}
```

## 5.7 遍历数组

```
#! /bin/bash

im=(qq weixin)

for loop in ${im[*]}
# for loop in ${im[@]}
do 
echo $loop
done
```

# 六、输入输出重定向

`Unix`命令默认从标准输入设备(`stdin`)获取输入，将结果输出到标准输出设备(`stdout`)显示。一般情况下，标准输入设备就是键盘，标准输出设备就是终端，即显示器。一般情况下，每个`Unix/Linux`命令运行时都会打开三个文件：

- 标准输入文件(`stdin`)：stdin的文件描述符为0，Unix程序默认从stdin读取数据。
- 标准输出文件(`stdout`)：stdout 的文件描述符为1，Unix程序默认向stdout输出数据。
- 标准错误文件(`stderr`)：stderr的文件描述符为2，Unix程序会向stderr流中写入错误信息。

## 6.1 输出重定向

将标准输出内容重定向到文件，`>` 执行会覆盖原文件内容，`>>`会以追加的方式写到文件。

```
$ command > file
```

如果希望将stderr重定向到file文件：

```
$ command 2 > file
```

如果希望将 stdout 和 stderr 合并后重定向到 file，可以这样写：`/data/shell/backup.sh > /dev/null 2>&1`

```
$ command > file 2>&1
```

以追加的方式将内容写到文件，只是多了一个`>`

```
$ command >> file
```

## 6.2 输入重定向

```
$ command < file
```

## 6.3 Here Document

```
$ command << delimiter
    document
delimiter
```

它的作用是将两个 delimiter 之间的内容(document) 作为输入传递给 command。注意：

- 结尾的delimiter 一定要顶格写，前面不能有任何字符，后面也不能有任何字符，包括空格和 tab 缩进。
- 开始的delimiter前后的空格会被忽略掉。

示例（可以保持多行的输入）：

```
$ cat << EOF | curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/demo?pretty' -d @- 
{
    "mappings": {
        "_doc": {
            "properties": {
                "title": {
                    "type" : "keyword"
                },
                "content": {
                    "type": "text"
                }
            }
        }
    }
}
EOF
```

## 6.4 /dev/null文件

如果希望执行某个命令，但又不希望在屏幕上显示输出结果，那么可以将输出重定向到`/dev/null`：

```
$ command > /dev/null
```

`/dev/null` 是一个特殊的文件，写入到它的内容都会被丢弃；如果尝试从该文件读取内容，那么什么也读不到。但是`/dev/null` 文件非常有用，将命令的输出重定向到它，会起到”禁止输出“的效果。

如果希望屏蔽`stdout`和 `stderr`，可以这样写：

```
$ command > /dev/null 2>&1
```






- [1] [Linux Shell脚本教程：30分钟玩转Shell脚本编程](http://c.biancheng.net/cpp/shell/)
- [2] [Runoob Shell教程](https://www.runoob.com/linux/linux-shell.html)
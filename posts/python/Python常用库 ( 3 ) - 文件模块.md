```
{
    "url": "python-file",
    "time": "2016/03/15 23:35",
    "tag": "Python"
}
```

# 一、文件路径相关(os.path)

路径相关的操作主要在`os.path`库。

## 1.1 判断文件、目录是否存在


**1. os.path.exists:** 判断文件/目录是否存在

```
>>> import os
>>> print os.path.exists.__doc__
Test whether a path exists.  Returns False for broken symbolic links

# build目录
>>> os.path.exists("build")
True
# 1.txt文件
>>> os.path.exists("1.txt")
True
```

**2. os.path.isdir:** 判断是否是目录

```
>>> os.path.isdir("build")
True
>>> os.path.isdir("1.txt")
False
```

**3. os.path.isfile:** 判断是否是文件

```
>>> os.path.isfile("build")
False
>>> os.path.isfile("1.txt")
True
```

**4. os.path.islink:** 判断是否是软连

> $ ln -s 1.txt v.txt

```
>>> os.path.islink("1.txt")
False
>>> os.path.islink("v.txt")
True
```

## 1.2 获取路径、目录、文件名

**1. os.path.abspath:** 获取所在文件的绝对路径

```
>>> os.path.abspath("1.txt")
'/Users/peng/workspace/python/demo/1.txt'
```

**2. os.path.dirname:** 获取目录

```
>>> os.path.dirname(os.path.abspath("1.txt"))
'/Users/peng/workspace/python/demo'
```

**3. os.path.basename:** 获取文件名

```
>>> os.path.basename(os.path.abspath("1.txt"))
'1.txt'
```

**4. os.path.isabs:** 是否是绝对路径

```
>>> os.path.isabs("1.txt")
False
>>> os.path.isabs("/Users/peng/workspace/python/demo/1.txt")
True
```

**5. os.path.walk():** 遍历路径

## 1.3 路径处理

**1. os.path.split:** 切分路径

```
>>> os.path.split(os.path.abspath("1.txt"))
('/Users/peng/workspace/python/demo', '1.txt')
```

**2. os.path.splitext:** 

```
>>> os.path.splitext(os.path.abspath("1.txt"))
('/Users/peng/workspace/python/demo/1', '.txt')
```

**3. os.path.join:** 合并路径

```
>>> os.path.join("/Users/peng/workspace/python/demo", "1.txt")
'/Users/peng/workspace/python/demo/1.txt'
```

# 二、操作文件目录(os)

本章节主要描述文件的创建、移动、拷贝、删除等操作。

## 2.1 添加/删除目录

**1. 创建目录** 

- os.mkdir: 创建目录
- os.makedirs: 递归创建目录

```
>>> print os.mkdir.__doc__
mkdir(path [, mode=0777])

>>> os.mkdir("test")
# 递归创建时会抛异常
>>> os.mkdir("testa/testb")
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
OSError: [Errno 2] No such file or directory: 
```

**递归创建目录**

```
'testa/testb'
# 递归创建目录
>>> os.makedirs("testa/testb")
```

**2. 删除目录** 

- os.rmdir: 删除单个目录
- os.removedirs: 递归删除目录
- os.remove: 删除文件

```
$ touch 1.md
>>> os.rmdir("test")
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
OSError: [Errno 66] Directory not empty: 'test'

$ rm 1.md
>>> os.rmdir("test")

# 递归删除目录
>>> os.removedirs("testa/testb")

# 删除文件
>>> os.remove("1.md")
>>> os.remove("1.md")
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
OSError: [Errno 2] No such file or directory: '1.md'
```

## 2.2 获取目录

**1. os.listdir:** 列出目录下的文件

```
>>> os.listdir("./test")
['1.md', 'folder']
```

**2. os.getcwd:** 获取当前工作目录

- os.getcwd: 获取当前工作目录
- os.chdir：切换工作目录

```
>>> print os.getcwd.__doc__
getcwd() -> path

>>> os.getcwd()
'/Users/peng/workspace/python/demo'

>>> os.chdir("./test")
>>> print os.getcwd()
/Users/peng/workspace/python/demo/test
```

## 2.3 重命名

**1. os.rename:** 获取当前工作目录

```
>>> print os.rename.__doc__
rename(old, new)

Rename a file or directory.

>>> os.rename("1.md", "2.md")
>>> os.rename("1.md", "2.md")
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
OSError: [Errno 2] No such file or directory
```


# 三、读写文件操作

`python`内置文件相关操作，本章节主要描述对文件的读写操作。

## 3.1 操作流程

操作文件分为3个步骤：1. 打开文件 - 2. 读写文件 - 3. 关闭文件句柄

### 3.1.1 打开文件

调用open()函数打开文件并返回一个file对象

```
>>> print open.__doc__
open(name[, mode[, buffering]]) -> file object

Open a file using the file() type, returns a file object.  This is the
preferred way to open a file.  See file.__doc__ for further information

# Py3中第三个参数encoding可以指定文件编码。
f = open("1.txt", "r")

# 打印文件名
>>> print(f.name)
1.txt
# 打印文件对象是否已经关闭
>>> print(f.closed)
False
# 打印打开模式
>>> print(f.mode)
r
```

**常用模式说明**

![](../../static/uploads/pythom-mode.png)

mode|说明
---|---
b|二进制模式。
+|打开一个文件进行更新(可读可写)。
r|以只读方式打开文件。文件的指针将会放在文件的开头。这是默认模式。
r+|打开一个文件用于读写。文件指针将会放在文件的开头。
w|打开一个文件只用于写入。如果该文件已存在则打开文件，并从开头开始编辑，即原有内容会被删除。如果该文件不存在，创建新文件。
w+|打开一个文件用于读写。如果该文件已存在则打开文件，并从开头开始编辑，即原有内容会被删除。如果该文件不存在，创建新文件。
a|打开一个文件用于追加。如果该文件已存在，文件指针将会放在文件的结尾。也就是说，新的内容将会被写入到已有内容之后。如果该文件不存在，创建新文件进行写入。
a+|打开一个文件用于读写。如果该文件已存在，文件指针将会放在文件的结尾。文件打开时会是追加模式。如果该文件不存在，创建新文件用于读写。

### 3.1.2 读写文件

调用read()或write()等方法读写文件

```
>>> print file.read.__doc__
read([size]) -> read at most size bytes, returned as a string.

If the size argument is negative or omitted, read until EOF is reached.
Notice that when in non-blocking mode, less data than what was requested
may be returned, even if no size parameter was given.

f.read()
```

`read`函数接收一个读取的字节数，如果不指定默认读取到结束符为止。

### 3.1.3 关闭句柄

```
f.close()
```

操作使用完后需要手动调用`close()`方法，`Python`里引入`with`语句可以在`with`代码块作用域结束后自动关闭文件，如：

```
with open("1.txt", "r") as f:
    t = f.read()
    print(t)

print(f.closed)

# Output:
The Zen of Python, by Tim Peters

Beautiful is better than ugly.
Explicit is better than implicit.
Simple is better than complex.
...
Namespaces are one honking great idea -- let's do more of those!
True
```

## 3.2 文件读取

常用的读取函数有：

- `read()` 默认读取整个文本内容
- `readline()` 读取一行内容
- `readlines()` 读取整个文件行返回到`list`中，相当于调用`readline()`直到`EOF`，并将结果放在`list`中。

**调用readline读取**

通过`readlines`可以看到每一行即便看起来是空白行，其实也是有空白字符`\n`，所以可以用`if not t`来判断是否读取到文件尾。

```
f = open("1.txt", "r")
while True:
    t = f.readline()
    print(t)
    if not t:
        break

f.close()

# Output:
The Zen of Python, by Tim Peters
...
```

**调用readlines读取**

可方便对读取到的每行内容做数据处理，省去读取整个文本后在切割的过程。可以看到输出的行尾都有换行符。可以批量去除：`[l.strip() for l in t]`

```
import json

f = open("1.txt", "r")
t = f.readlines()
print(json.dumps(t))


f.close()

# output:
[
    "The Zen of Python, by Tim Peters\n",
    "\n",
    "Beautiful is better than ugly.\n",
    "Explicit is better than implicit.\n",
    "Simple is better than complex.\n"
]
```

## 3.3 文件写入

写入主要有`write()`和`writelines()`方法。熟悉前面文件的操作流程和打开方式后，写入就比较好理解了。常用的写入模式（`w`、`a+`）

**通过write写入**

```
f = open("x.txt", "a+")

while True:
    s = input("Pls input something:")
    f.write(s)
    f.flush()

f.close()
```

**通过writelines写入**

```
f3 = open("3.txt", "w")
f3.writelines(t)
f3.close()
```

## 3.4 其他方法

### 3.4.1 `truncate()`

```
>>> print file.truncate.__doc__
truncate([size]) -> None.  Truncate the file to at most size bytes.

Size defaults to the current file position, as returned by tell().
```

用于清空文本内容或者只保留前N个字符。

```
f = open("1.txt", "r+")
# 只保留前10个字符，不传则清空了文本内容
f.truncate(10)
f.close()
```

### 3.4.2 `tell()` 和 `seek()`

`tell() -> current file position, an integer (may be a long integer).`用于返回文件指针所在的位置。

`seek(offset[, whence]) -> None.  Move to new file position.`用于移动文件指针。

比如在行一行处理大文件时，可能前面N行已经处理成功，记录一下此时的文件位置，下次启动时可以直接移动文件指针接着继续处理。

读取2行后打印文件位置，并继续读取一行：

```
f = open("1.txt", "r+")

f.readline()
f.readline()
print(f.tell())
print(f.readline())

f.close()

# Output:
34
Beautiful is better than ugly.


```

直接移动文件指针到`34`的位置读取一行。可以看到输出内容是一致。

```
f = open("1.txt", "r+")

f.seek(34)
print(f.readline())

f.close()

# Output:
Beautiful is better than ugly.
```

### 3.4.3 `flush()`

`flush() -> None.  Flush the internal I/O buffer.`
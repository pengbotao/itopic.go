```
{
    "url": "grep-sed-awk",
    "time": "2017/11/18 21:24",
    "tag": "运维",
    "toc": "yes"
}
```

# 一、grep

grep用来查找文件里符合条件的字符串。

## 1.1 用法

| 参数                         | 说明                   | 示例                                                   |
| ---------------------------- | ---------------------- | ------------------------------------------------------ |
| `-v`                         | 排除掉匹配的结果       | 不匹配127打头的host：<br />`$ grep -v ^127 /etc/hosts` |
| `-i`                         | 忽略大小写             | `$ grep -i LOCAL /etc/hosts`                           |
| `-n`                         | 显示行号               | `$ grep -n local /etc/hosts`                           |
| `-m`                         | 最大匹配次数           | `$ grep local /etc/hosts -m 1`                         |
| `-r`                         | 递归查询目录           |                                                        |
| `--color=auto`               | 匹配字符串**高亮**显示 |                                                        |
| <BR>                         |                        |                                                        |
| `-l` `--files-with-matches`  | 只打印匹配的文件名     | `$ grep -l local /etc/hosts`                           |
| `-L` `--files-without-match` | 只打印不匹配的文件名   |                                                        |
| `-c` `--count`               | 只统计匹配的次数       |                                                        |
| `-I`                         | 忽略二进制人间         |                                                        |

## 1.2 查找字符串

`hosts`文件中包含`local`的字符串，并高亮显示。

```
$ grep local /etc/hosts --color=auto
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
127.0.0.1 api.local
```

显示行数，不区分大小写：

- `-i`：忽略大小写
- `-n`：显示行号

```
$ grep -in LOCAL /etc/hosts --color=auto
1:127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
2:::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
3:127.0.0.1 api.local
```

## 1.3 排除字符串

前一步结果中排除掉包含`localhost`字样的行。

```
$ grep local /etc/hosts | grep -v "localhost"
127.0.0.1 api.local
```

## 1.4 递归查找

查找博客下包含`"tag": "Lua"`的字样

```
$ grep -r '"tag": "Lua"' .
./posts/lua/Nginx+Lua入门知识.md:    "tag": "Lua"
./posts/lua/Lua实现类与继承.md:    "tag": "Lua"
./posts/lua/Lua基础语法.md:    "tag": "Lua"
```

通过`--include`，只查找文件包含Nginx的，也可以通过`--exclude`指定不包含。

```
$ grep -r '"tag": "Lua"' . --include=Lua* --exclude=*基础*
./posts/lua/Lua实现类与继承.md:    "tag": "Lua"

$ grep -r '"tag": "Lua"' . --include *.{php,py}
```

## 1.5 正则匹配

| 正则表达式 | 说明                       | 示例                                                         |
| ---------- | -------------------------- | ------------------------------------------------------------ |
| `^`        | 匹配开头字符串             | 匹配127开头的hosts：<br />`$ grep ^127 /etc/hosts`           |
| `$`        | 匹配结尾字符串             | 匹配local结尾的hosts：<br />`$ grep local$ /etc/hosts`       |
| `^$`       | 匹配空行                   | 打印配置文件，不包含`#`注释，不包含空行<br />`$ cat /etc/salt/master | grep -v ^# | grep -v ^$` |
| `.`        | 匹配单个字符               |                                                              |
| `.*`       | 匹配任意字符               |                                                              |
| `\`        | 转义字符                   |                                                              |
| `[abc]`    | 匹配字符集合内任意一字符串 |                                                              |
| `[^abc]`   |                            |                                                              |
| `{n,m}`    |                            |                                                              |

## 

# 二、find

用grep类似，grep查找的是文件内容，而find查找的是文件，可以通过文件名、文件类型、修改时间等各种条件进行文件查找。

## 2.1 用法

```
用法: find [-H] [-L] [-P] [-Olevel] [-D help|tree|search|stat|rates|opt|exec] [path...] [expression]
```

表达式，各种条件可以组合使用：

| 表达式      | 说明                                          | 示例                                                         |
| ----------- | --------------------------------------------- | ------------------------------------------------------------ |
| `-name`     | 指定文件名                                    | 查找当前目录和子目录：`find . -name "*.py"`                  |
| `-iname`    | 指定文件名，不区分大小写                      | 查找`root`目录和`usr`目录：`find /root /usr/ -iname readme.md` |
| `-user`     | 指定文件所有者                                |                                                              |
| `-group`    | 指定文件所属组                                |                                                              |
| `-perm`     | 指定文件权限                                  | 查找当前目录下600权限的文件：<br />`find . -maxdepth 1 -type f -perm 600` |
| `-type`     | 指定文件类型.<br /> `f` ：普通文件，`d`：目录 | 查找`usr`下的`doc`目录：`find /usr/ -type d -name doc`       |
| `-atime`    | 指定文件最后访问时间                          | 如：`find . -type f -name "*.log" -atime -6`                 |
| `-ctime`    | 指定文件修改时间                              | 如：`find . -type f -name "*.tar" -ctime +6`                 |
| `-mtime`    | 指定文件内容被修改时间                        |                                                              |
| `-size`     | 指定文件大小                                  | 查找当前目录下文件大小为0的：`find . -size 0`                |
| `-maxdepth` | 指定最大查找层级                              | 查找空文件，最深2层：`find / -size 0 -maxdepth 2`            |
| `-empty`    | 查找空白文件、文件夹                          | 查找空文件：`find / -maxdepth 1 -empty`                      |
| <br />      |                                               |                                                              |
| `-exec`     | 执行shell命令                                 | `find . -name "*.md" -exec ls -lh {} \;` <br />查找执行ls命令。`{} \;`相当于参数传给`ls`，`{}`和`\`之间有个空格。 |
| `-print`    | 输出，以换行方式                              |                                                              |
| `print0`    | 每一个输出结果后以NULL结束，而非换行          |                                                              |

## 2.2 查找过期日志

**查找日志目录下的文件并删除14天前的日志文件：**

```
find /data/logs -type f -ctime +14 -exec rm -rf {} \; > /dev/null 2>&1
```

也可以写成：

```
find /data/logs -type f -ctime +14 -print0 | xargs -0 rm -rf
```



# 三、sed

## 3.1 用法

用来处理文本文件，处理方式：

- 每次读取一行内容
- 根据规则匹配并修改数据。注意，sed默认不会修改源文件。

基本格式如下：

```
sed [选项]... {脚本(如果没有其他脚本)} [输入文件]...
```

**选项说明：**

| 选项 | 说明                                                         | 示例                                                         |
| ---- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `-e` | `-e`脚本,` --expression`=脚本，添加“脚本”到程序的运行列表。  | `sed '2d' /etc/hosts`<br />默认是`-e`，即后面跟一串脚本`'2d'` |
| `-f` | `-f` 脚本文件, `--file`=脚本文件，添加“脚本文件”到程序的运行列表。 |                                                              |
| `-n` | 取消自动打印模式空间                                         |                                                              |
| `-i` | 将修改后的内容写入源文件                                     |                                                              |

## 3.2 字符串替换

用法：

- `sed -i 's/oldstr/newstr/' file` (替换每行第一个匹配)
- `sed -i 's/oldstr/newstr/g' file` (替换所有匹配)

**示例**：将`/etc/hosts`中`api.local`替换为`test.local`，不做替换查看输出。

```
$ sed 's/api.local/test.local/' /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
127.0.0.1 test.local
```

> 注意：上面示例只会替换每行第一个匹配到的结果，如果需要替换所有匹配，可以增加`/g`

```
$ sed 's/localhost/local/g' /etc/hosts
127.0.0.1   local local.localdomain local4 local4.localdomain4
::1         local local.localdomain local6 local6.localdomain6
127.0.0.1 test.local
```

如果想将替换后的内容会写到文件，可以增加`-i`参数：

```
$ sed -i 's/api.local/test.local/' /etc/hosts
```

## 3.3 删除匹配行

将`hosts`文件中`127.0.0.1 test.local`进行删除，查看输出结果

```
$ sed '/127.0.0.1 test.local/d' /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
```

将输出结果保存到文件：

```
$ sed -i '/127.0.0.1 test.local/d' /etc/hosts
```

## 3.4 删除指定行

1、删除最后一行：

```
$ sed '$d' /etc/hosts
```

2、删除第二行

```
$ sed '2d' /etc/hosts
```

3、删除第二行到最后一行：

```
$ sed '2,$d' /etc/hosts
```

同上，如果需要写入到文件，需要增加`-i`参数。

## 3.5 查找指定行

```
$ sed -n '2p' /etc/hosts
$ sed -n '2,3p' /etc/passwd
$ sed -n '3,$p' /etc/passwd
```



# 四、awk

AWK 是一种处理文本文件的语言，是一个强大的文本分析工具。之所以叫 AWK 是因为其取了三位创始人 Alfred Aho，Peter Weinberger, 和 Brian Kernighan 的 Family Name 的首字符。

## 4.1 用法

和`sed`一样，`awk`也是逐行扫描文件，如果匹配成功则执行操作，反之则不做处理。

基本格式：

```
awk [选项] '脚本命令' 文件名
```

| 选项    | 说明                           | 示例                          |
| ------- | ------------------------------ | ----------------------------- |
| `-F fs` | 指定分隔符，默认是空格或制表符 | `awk '{print $1}' /etc/hosts` |

**内置变量：**

| 变量  | 说明                                |
| ----- | ----------------------------------- |
| `$0`  | 整行文本                            |
| `$1`  | 分隔后的第一段，`$2`则表示第二段    |
| `$NF` | 最后一段，`$(NF-1)`则表示倒数第二段 |
| `NR`  | 行数                                |

## 4.2 字符串切割

切割`/etc/passwd`打印用户名

```
$ awk -F ':' '{print $1}' /etc/passwd
```

切割后最后一段等于`/sbin/nologin`，打印整行内容

```
$ awk -F ':' '$NF=="/sbin/nologin"{print $0}' /etc/passwd
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
```

## 4.3 输出指定行

```
$ awk -F ':' 'NR>1 && NR<10{print NR " " $0}' /etc/passwd
2 bin:x:1:1:bin:/bin:/sbin/nologin
3 daemon:x:2:2:daemon:/sbin:/sbin/nologin
```

## 4.4 变量计算

```
$ free -m
             total       used       free     shared    buffers     cached
Mem:         16080      15786        294       1052        538       8276

# 内存空闲：free + buffers + cached / total
$ free -m | sed -n '2p' | awk '{print ($4+$6+$7)*100/$2}'
56.6418
```

# 五、xargs



# 六、管道与重定向

`>`

`>>`

`&`




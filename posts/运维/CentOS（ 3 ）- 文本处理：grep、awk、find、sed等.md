```
{
    "url": "grep-sed-awk",
    "time": "2017/11/18 21:24",
    "tag": "运维",
    "toc": "yes"
}
```

日常日志搜索查询，文件查找的工作少不了，grep、find等命令也经常使用，这里做一个整体的整理。主要解决以下文本、文件的查找及处理问题，比如：

- 查询满足条件的日志内容
- 根据创建时间查找文件并清理
- 对文本内容进行批量替换等操作

包含的命令：

| 命令 | 说明         |
| ---- | ------------ |
| grep | 文本内容查找 |
| find | 文件查找     |

# 一、grep

grep用来查找文件里符合条件的字符串。

```
# grep --help
用法: grep [选项]... PATTERN [FILE]...
在每个 FILE 或是标准输入中查找 PATTERN。
默认的 PATTERN 是一个基本正则表达式(缩写为 BRE)。
例如: grep -i 'hello world' menu.h main.c

正则表达式选择与解释:
  -E, --extended-regexp     PATTERN 是一个可扩展的正则表达式(缩写为 ERE)
  -F, --fixed-strings       PATTERN 是一组由断行符分隔的定长字符串。
  -G, --basic-regexp        PATTERN 是一个基本正则表达式(缩写为 BRE)
  -P, --perl-regexp         PATTERN 是一个 Perl 正则表达式
  -e, --regexp=PATTERN      用 PATTERN 来进行匹配操作
  -f, --file=FILE           从 FILE 中取得 PATTERN
  -i, --ignore-case         忽略大小写
  -w, --word-regexp         强制 PATTERN 仅完全匹配字词
  -x, --line-regexp         强制 PATTERN 仅完全匹配一行
  -z, --null-data           一个 0 字节的数据行，但不是空行

Miscellaneous:
  -s, --no-messages         suppress error messages
  -v, --invert-match        select non-matching lines
  -V, --version             display version information and exit
      --help                display this help text and exit

输出控制:
  -m, --max-count=NUM       NUM 次匹配后停止
  -b, --byte-offset         输出的同时打印字节偏移
  -n, --line-number         输出的同时打印行号
      --line-buffered       每行输出清空
  -H, --with-filename       为每一匹配项打印文件名
  -h, --no-filename         输出时不显示文件名前缀
      --label=LABEL         将LABEL 作为标准输入文件名前缀
  -o, --only-matching       show only the part of a line matching PATTERN
  -q, --quiet, --silent     suppress all normal output
      --binary-files=TYPE   assume that binary files are TYPE;
                            TYPE is 'binary', 'text', or 'without-match'
  -a, --text                equivalent to --binary-files=text
  -I                        equivalent to --binary-files=without-match
  -d, --directories=ACTION  how to handle directories;
                            ACTION is 'read', 'recurse', or 'skip'
  -D, --devices=ACTION      how to handle devices, FIFOs and sockets;
                            ACTION is 'read' or 'skip'
  -r, --recursive           like --directories=recurse
  -R, --dereference-recursive
                            likewise, but follow all symlinks
      --include=FILE_PATTERN
                            search only files that match FILE_PATTERN
      --exclude=FILE_PATTERN
                            skip files and directories matching FILE_PATTERN
      --exclude-from=FILE   skip files matching any file pattern from FILE
      --exclude-dir=PATTERN directories that match PATTERN will be skipped.
  -L, --files-without-match print only names of FILEs containing no match
  -l, --files-with-matches  print only names of FILEs containing matches
  -c, --count               print only a count of matching lines per FILE
  -T, --initial-tab         make tabs line up (if needed)
  -Z, --null                print 0 byte after FILE name

文件控制:
  -B, --before-context=NUM  打印以文本起始的NUM 行
  -A, --after-context=NUM   打印以文本结尾的NUM 行
  -C, --context=NUM         打印输出文本NUM 行
  -NUM                      same as --context=NUM
      --group-separator=SEP use SEP as a group separator
      --no-group-separator  use empty string as a group separator
      --color[=WHEN],
      --colour[=WHEN]       use markers to highlight the matching strings;
                            WHEN is 'always', 'never', or 'auto'
  -U, --binary              do not strip CR characters at EOL (MSDOS/Windows)
  -u, --unix-byte-offsets   report offsets as if CRs were not there
                            (MSDOS/Windows)

‘egrep’即‘grep -E’。‘fgrep’即‘grep -F’。
直接使用‘egrep’或是‘fgrep’均已不可行了。
若FILE 为 -，将读取标准输入。不带FILE，读取当前目录，除非命令行中指定了-r 选项。
如果少于两个FILE 参数，就要默认使用-h 参数。
如果有任意行被匹配，那退出状态为 0，否则为 1；
如果有错误产生，且未指定 -q 参数，那退出状态为 2。
```

## 1.1 常用参数

| 参数           | 说明               | 示例                                                   |
| -------------- | ------------------ | ------------------------------------------------------ |
| `-v`           | 排除掉匹配的结果   | 不匹配127打头的host：<br />`$ grep -v ^127 /etc/hosts` |
| `-i`           | 忽略大小写         | `$ grep -i LOCAL /etc/hosts`                           |
| `--color=auto` | 匹配字符串高亮显示 |                                                        |
| `-n`           | 显示行号           |                                                        |
| `-m`           | 最打匹配次数       | `grep local /etc/hosts -m 1`                           |
| `-r`           | 递归查询目录       |                                                        |

## 1.2 正则说明

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

# 二、find

用grep类似，grep查找的是文件内容，而find查找的是文件。

# 三、sed

```
# sed --help
用法: sed [选项]... {脚本(如果没有其他脚本)} [输入文件]...

  -n, --quiet, --silent
                 取消自动打印模式空间
  -e 脚本, --expression=脚本
                 添加“脚本”到程序的运行列表
  -f 脚本文件, --file=脚本文件
                 添加“脚本文件”到程序的运行列表
  --follow-symlinks
                 follow symlinks when processing in place; hard links
                 will still be broken.
  -i[SUFFIX], --in-place[=SUFFIX]
                 edit files in place (makes backup if extension supplied).
                 The default operation mode is to break symbolic and hard links.
                 This can be changed with --follow-symlinks and --copy.
  -c, --copy
                 use copy instead of rename when shuffling files in -i mode.
                 While this will avoid breaking links (symbolic or hard), the
                 resulting editing operation is not atomic.  This is rarely
                 the desired mode; --follow-symlinks is usually enough, and
                 it is both faster and more secure.
  -l N, --line-length=N
                 指定“l”命令的换行期望长度
  --posix
                 关闭所有 GNU 扩展
  -r, --regexp-extended
                 在脚本中使用扩展正则表达式
  -s, --separate
                 将输入文件视为各个独立的文件而不是一个长的连续输入
  -u, --unbuffered
                 从输入文件读取最少的数据，更频繁的刷新输出
      --help     打印帮助并退出
      --version  输出版本信息并退出

如果没有 -e, --expression, -f 或 --file 选项，那么第一个非选项参数被视为
sed脚本。其他非选项参数被视为输入文件，如果没有输入文件，那么程序将从标准
输入读取数据。
```



# 四、awk

```
# awk --help
用法: awk [POSIX 或 GNU 风格选项] -f 脚本文件 [--] 文件 ...
用法: awk [POSIX 或 GNU 风格选项] [--] '程序' 文件 ...
POSIX 选项:     		GNU 长选项:
	-f 脚本文件		--file=脚本文件
	-F fs			--field-separator=fs
	-v var=val		--assign=var=val
	-m[fr] val
	-O			--optimize
	-W compat		--compat
	-W copyleft		--copyleft
	-W copyright		--copyright
	-W dump-variables[=file]	--dump-variables[=file]
	-W exec=file		--exec=file
	-W gen-po		--gen-po
	-W help			--help
	-W lint[=fatal]		--lint[=fatal]
	-W lint-old		--lint-old
	-W non-decimal-data	--non-decimal-data
	-W profile[=file]	--profile[=file]
	-W posix		--posix
	-W re-interval		--re-interval
	-W source=program-text	--source=program-text
	-W traditional		--traditional
	-W usage		--usage
	-W use-lc-numeric	--use-lc-numeric
	-W version		--version

提交错误报告请参考“gawk.info”中的“Bugs”页，它位于打印版本中的“Reporting
Problems and Bugs”一节

翻译错误请发信至 translation-team-zh-cn@lists.sourceforge.net

gawk 是一个模式扫描及处理语言。缺省情况下它从标准输入读入并写至标准输出。

范例:
	gawk '{ sum += $1 }; END { print sum }' file
	gawk -F: '{ print $1 }' /etc/passwd
```


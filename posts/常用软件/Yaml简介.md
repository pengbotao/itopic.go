```
{
    "url": "yaml",
    "time": "2019/02/09 20:46",
    "tag": "常用软件"
}
```



```
%YAML 1.2
---
YAML: YAML Ain't Markup Language

What It Is: YAML is a human friendly data serialization standard for all programming languages.
```

# 一、概述

`YAML`是"`YAML Ain't a Markup Language`"（`YAML`不是一种标记语言）的递归缩写。在开发的这种语言时，`YAML` 的意思其实是："`Yet Another Markup Language`"（仍是一种标记语言），但为了强调这种语言以数据做为中心，而不是以标记语言为重点，而用反向缩略语重命名。

它的基本语法规则如下：

- 大小写敏感
- 使用缩进表示层级关系，缩进时不允许使用`Tab`键，只允许使用空格。缩进的空格数目不重要，只要相同层级的元素左侧对齐即可（类似`Python`）
- 键值对以冒号分隔，冒号后面需留有一个空格
- 用 `#` 表示注释，只支持单行注释
- 一个文件可以包含多个文件的内容
  - 用 `---` 表示一份新内容的开始
  - 用 `... `表示一部分内容的结束（非必需）

频率使用上没有JSON高，如果书写有困难，可以直接使用 Json to YAML 互转工具 [<sup>[4]</sup>](#refer) 。

# 二、基础类型

## 2.1 字符串

```
str: This is a String

{
  "str": "This is a String"
}
```

## 2.2 布尔值

- `true`、`True`、`TRUE`、`yes`、`Yes` 和 `YES`皆为  **真**
- `false`、`False`、`FALSE`、`no`、`No` 和 `NO`皆为  **假**

## 2.3 整数

```
int: 4
```

## 2.4 浮点数

```
number: 12.30
```

## 2.5 空

`null`、`Null` 和 `~`都是空，不指定值默认也是空

## 2.6 日期

日期采用复合 iso8601 格式的年、月、日表示。

```
date: '2019-02-19'
```

# 三、数据结构

## 3.1 对象

定义：`KEY + : + 空格 + 值`，如`env: qa`，别忘记后面的空格。示例：

```
browser:
  ie: Internet Explore
  chrome: Chrome
  other:
    safari: Safari
    ff: Firefox
```

对应JSON:

```
{
    "browser": {
        "ie": "Internet Explore",
        "chrome": "Chrome",
        "other": {
            "safari": "Safari",
            "ff": "Firefox"
        }
    }
}
```

说明：上面键值均为字符串，一般情况下不需要加引号，只有在存在反斜杠需要转义的字符时才需要引号。

支持**流式风格（ Flow style）**的语法（用花括号包裹，用逗号加空格分隔，类似 JSON）

```
browser:
  ie: Internet Explore
  chrome: Chrome
  other: {safari: Safari, ff: Firefox}
```

## 3.2 数组

一组连词线开头的行，构成一个数组。

```
- Internet Explore
- Chrome
- Other
```

转为JSON为：

```
[
  "Internet Explore",
  "Chrome",
  "Other"
]
```

数据也可以用行内表示法：

```
browser: [Internet Explore, Chrome, other]


{
    "browser": [
        "Internet Explore",
        "Chrome",
        "other"
    ]
}
```

## 3.3 复合结构

接下来就是数组和对象进行组合的情况。

```
code: 0
msg: success
data:
- 
  name: Chrome
  percent: 62.5
-
  name: Internet Explore
  percent: 1.4
```

`-`下的第一个键可以提到`-`后面，如：

```
code: 0
msg: success
data:
- name: Chrome
  percent: 62.5
- name: Internet Explore
  percent: 1.4
```

等同于

```
{
    "code": 0,
    "msg": "success",
    "data": [
        {
            "name": "Chrome",
            "percent": 62.5
        },
        {
            "name": "Internet Explore",
            "percent": 1.4
        }
    ]
}
```



---

<div id="refer"></div>

- [1] [Yaml.org](https://yaml.org/)
- [2] [YAML 语言教程](http://www.ruanyifeng.com/blog/2016/07/yaml.html)
- [3] [一文看懂 YAML](https://chenpipi.cn/post/yaml-all-in-one/)
- [4] [JSON to YAML](https://www.json2yaml.com/)


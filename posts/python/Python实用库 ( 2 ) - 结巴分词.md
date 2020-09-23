```
{
    "url": "python-jieba",
    "time": "2017/01/06 20:46",
    "tag": "Python"
}
```


# 安装方式

```
pip install jieba
```

# 基本使用

## 分词模式

- 精确模式：试图将句子最精确地切开，适合文本分析；
- 全模式：把句子中所有的可以成词的词语都扫描出来, 速度非常快，但是不能解决歧义；
- 搜索引擎模式：在精确模式的基础上，对长词再次切分，提高召回率，适合用于搜索引擎分词。

**功能说明**

- `jieba.cut` 方法接受三个输入参数: 需要分词的字符串；`cut_all`参数用来控制是否采用全模式；`HMM`参数用来控制是否使用 HMM 模型
- `jieba.cut_for_search` 方法接受两个参数：需要分词的字符串；是否使用`HMM`模型。该方法适合用于搜索引擎构建倒排索引的分词，粒度比较细
- `jieba.cut` 以及 `jieba.cut_for_search`返回的结构都是一个可迭代的`generator`，可以使用`for`循环来获得分词后得到的每一个词语(unicode)，或者用`jieba.lcut`以及`jieba.lcut_for_search`直接返回`list`
- `jieba.Tokenizer(dictionary=DEFAULT_DICT)`新建自定义分词器，可用于同时使用不同词典。`jieba.dt`为默认分词器，所有全局分词相关函数都是该分词器的映射。

**示例**

```
# encoding=utf-8
import jieba
import jieba.posseg as pseg

wd = "我想要带你去浪漫的土耳其，然后一起去北京清华大学。"

# 精确模式
print(" / ".join(jieba.cut(wd)))
# 我 / 想要 / 带你去 / 浪漫 / 的 / 土耳其 / ， / 然后 / 一起 / 去 / 北京 / 清华大学 / 。

# 全模式
print(" / ".join(jieba.cut(wd, cut_all=True)))
# 我 / 想要 / 带你去 / 浪漫 / 的 / 土耳其 /  /  / 然后 / 一起 / 去 / 北京 / 清华 / 清华大学 / 华大 / 大学 /  /

# 搜索引擎模式
print(" / ".join(jieba.cut_for_search(wd)))
# 我 / 想要 / 带你去 / 浪漫 / 的 / 土耳其 / ， / 然后 / 一起 / 去 / 北京 / 清华 / 华大 / 大学 / 清华大学 / 。
```

## 自定义词典

开发者可以指定自己自定义的词典，以便包含`jieba`词库里没有的词。虽然`jieba`有新词识别能力，但是自行添加新词可以保证更高的正确率

- 用法： `jieba.load_userdict(file_name)` # file_name 为文件类对象或自定义词典的路径
- 词典格式和`dict.txt`一样，一个词占一行；每一行分三部分：词语、词频（可省略）、词性（可省略），用空格隔开，顺序不可颠倒。`file_name`若为路径或二进制方式打开的文件，则文件必须为`UTF-8`编码。
- 词频省略时使用自动计算的能保证分出该词的词频。

```
jieba.load_userdict("./dict.txt") # 添加一行 北京清华大学
print(" / ".join(jieba.lcut(wd)))
# 我 / 想要 / 带你去 / 浪漫 / 的 / 土耳其 / ， / 然后 / 一起 / 去 / 北京清华大学 / 。

# 动态修改词典
jieba.add_word("浪漫的土耳其", freq=None, tag=None)
jieba.del_word("北京清华大学")
print(" / ".join(jieba.lcut(wd)))
# 我 / 想要 / 带你去 / 浪漫的土耳其 / ， / 然后 / 一起 / 去 / 北京 / 清华大学 / 。
```

## 词性标注

```
words = pseg.cut(wd)
for word, flag in words:
    print("%s %s" % (word, flag))

我 r
想要 v
带你去 n
浪漫的土耳其 x
， x
然后 c
一起 m
去 v
北京 ns
清华大学 nt
。 x
```

# 词性对照表

```
- a 形容词  
	- ad 副形词  
	- ag 形容词性语素  
	- an 名形词  
- b 区别词  
- c 连词  
- d 副词  
	- df   
	- dg 副语素  
- e 叹词  
- f 方位词  
- g 语素  
- h 前接成分  
- i 成语 
- j 简称略称  
- k 后接成分  
- l 习用语  
- m 数词  
	- mg 
	- mq 数量词  
- n 名词  
	- ng 名词性语素  
	- nr 人名  
	- nrfg    
	- nrt  
	- ns 地名  
	- nt 机构团体名  
	- nz 其他专名  
- o 拟声词  
- p 介词  
- q 量词  
- r 代词  
	- rg 代词性语素  
	- rr 人称代词  
	- rz 指示代词  
- s 处所词  
- t 时间词  
	- tg 时语素  
- u 助词  
	- ud 结构助词 得
	- ug 时态助词
	- uj 结构助词 的
	- ul 时态助词 了
	- uv 结构助词 地
	- uz 时态助词 着
- v 动词  
	- vd 副动词
	- vg 动词性语素  
	- vi 不及物动词  
	- vn 名动词  
	- vq 
- x 非语素词  
- y 语气词  
- z 状态词  
	- zg 
```

看起来效果还不错，更多用法参考：`https://github.com/fxsjy/jieba`
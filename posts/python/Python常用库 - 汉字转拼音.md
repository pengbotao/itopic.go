```
{
    "url": "python-pinyin",
    "time": "2017/03/02 21:53",
    "tag": "Python"
}
```

日常中偶然会碰到将中文转拼音的需求，`mozillazg/python-pinyin`库应该可以满足日常的需求。

# 安装方式

```
$ pip install pypinyin

$ pypinyin 你好
nǐ hǎo
```


# 基本使用

```
from pypinyin import pinyin, lazy_pinyin, Style

sentence = "你好，世界。"

print(pinyin(sentence))
# [['nǐ'], ['hǎo'], ['，'], ['shì'], ['jiè'], ['。']]

print(pinyin(sentence, style=Style.NORMAL))
# [['ni'], ['hao'], ['，'], ['shi'], ['jie'], ['。']]

print(pinyin(sentence, style=Style.FIRST_LETTER))
# [['n'], ['h'], ['，'], ['s'], ['j'], ['。']]

print(lazy_pinyin(sentence))
# ['ni', 'hao', '，', 'shi', 'jie', '。']
```


# 翻译示例

> 《围城》一九四七年在上海初版，一九四八年再版，一九四九年三版，以后国内没有重印过。偶然碰见它的新版，那都是香港的“盗印”本。没有看到台湾的“盗印”，据说在那里它是禁书。美国哥伦比亚大学夏志清教授的英文著作里对它作了过高的评价，导致了一些西方语言的译本。日本京都大学荒井健教授很久以前就通知我他要翻译，近年来也陆续在刊物上发表了译文。现在，人民文学出版社建议重新排印，以便原著在国内较易找着，我感到意外和忻辛。

> 《 wéi chéng 》 yī jiǔ sì qī nián zài shàng hǎi chū bǎn ， yī jiǔ sì bā nián zài bǎn ， yī jiǔ sì jiǔ nián sān bǎn ， yǐ hòu guó nèi méi yǒu chóng yìn guò 。 ǒu rán pèng jiàn tā de xīn bǎn ， nà dōu shì xiāng gǎng de “ dào yìn ” běn 。 méi yǒu kàn dào tái wān de “ dào yìn ”， jù shuō zài nà lǐ tā shì jìn shū 。 měi guó gē lún bǐ yà dà xué xià zhì qīng jiào shòu de yīng wén zhù zuò lǐ duì tā zuò le guò gāo de píng jià ， dǎo zhì le yī xiē xī fāng yǔ yán de yì běn 。 rì běn jīng dū dà xué huāng jǐng jiàn jiào shòu hěn jiǔ yǐ qián jiù tōng zhī wǒ tā yào fān yì ， jìn nián lái yě lù xù zài kān wù shàng fā biǎo le yì wén 。 xiàn zài ， rén mín wén xué chū bǎn shè jiàn yì chóng xīn pái yìn ， yǐ biàn yuán zhù zài guó nèi jiào yì zhǎo zháo ， wǒ gǎn dào yì wài hé xīn xīn 。


看起来效果还不错，更多用法参考：`https://github.com/mozillazg/python-pinyin`
```
{
    "url": "simhash",
    "time": "2018/03/04 22:36",
    "tag": "Python"
}
```

# 一、简介

- https://github.com/leonsim/simhash

# 二、安装

```
pip install simhash -i  https://pypi.douban.com/simple
```

# 三、示例

## 3.1 计算距离
```
s1 = 'aaa'
s2 = 'aba'
print(Simhash(s1).distance(Simhash(s2)))#output 26
```

## 3.2 打印64位hash值
```
print(bin(Simhash(s1).value).replace('0b',''))
# 110011111011011110101010111111010011100101010011111100000001000
print(bin(Simhash(s2).value).replace('0b',''))
# 101101111011011001001110111101001100001000110001001101000101010
```
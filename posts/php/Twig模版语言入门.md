```
{
    "url": "twig-template-language",
    "time": "2013/06/02 19:24",
    "tag": "PHP,Twig"
}
```

Twig中有两种定界符 `{% ... %}` 和 `{{ ... }}` ， 前一种用来执行语句，比如for循环、IF判断、过滤等操作，后一种模版中用来显示变量。Twig常用用法记录如下：

# 可用符号

== != < > >= <= + - ~ * / // % ** | [] . .. and or not in is b-and b-or b-xor

部分符号的含义如下

- ~：连接两字符串，相当于 PHP 中的点号
- //：整除
- **：乘方，相当于 PHP 中的 ^
- b-and、b-or、b-xor：按位与、按位或、按位异或
- -：减法以及去除空白的简写用法，如 {{- 数据 }} => 去除左边的空白、{{ 数据 -}} => 去除右边的空白、{{- 数据 -}} => 去除两边的空白

# 注释

`{# This is Twig Comment #}`

# 变量

foo 为PHP赋值变量， 则模版中可以使用 `{{ foo }}`。如果该变量为数组或者对象，则可以使用`{{ foo.bar }}`来表示，也可以直接调用类中的方法，如 `{{ foo.getName }}`或当有参数传递时使用 `{{ foo.getName(p1, p2, p3) }}`

当模版中使用{{ foo.bar }}时，PHP会做如下判断来检测对应变量的值：

- 1. 检测foo是否是数组，并且bar是一个有效的KEY
- 2. 如果foo为对象，则检测bar是否为有效的属性
- 3. 如果foo为对象，但bar不是一个有效的属性，则检测bar是否为有效方法
- 4. 如果foo为对象，但bar不是一个有效的方法，则检测getBar是否为有效方法
- 5. 如果foo为对象，但getBar不是有效方法，则检测isBar是否为有效方法
- 6. 如果没有，则返回null

当foo为数组时，模版也可以这样使用：`{{ foo['bar'] }}`

变量未定义默认值设置：
`{{ var|default('var is not defined') }}`

可使用 `{% set 变量名=变量值 %}` 声明变量，也可写成 `{% set 变量名 %}` 变量值 `{% endset %}`

PHP 中非关联数组被映射成 `[元素1, 元素2, ...]`，关联数组则被映射成 `{键1: 值1, 键2: 值2, ...}`

# 数组遍历

1. 按数组value遍历：

```
{% for value in foo %}
    {{ value }}
{% endfor %}
```

2. 按数组key遍历：

```
{% for key in foo|keys %}
    {{ key }}
{% endfor %}
```

3. 按key，value遍历

```
{% for key, value in foo %}
    {{ key}}:{{value }}
{% endfor %}
```

4. 如果 foo 非数组，还可以使用else语句，如：

```
{% for key, value in foo %}
    {{ key}}:{{value }}
{% else %}  
    foo is Not a Array
{% endfor %}
```

5. 也可以直接带条件，遍历二维数组时比较有用，可用value.field来判断：

```
{% for key, value in foo if value == 1%}
    {{ key}}:{{value }}
{% endfor %}
```

循环体内部变量：

- loop.index	循环的次数（从1开始）
- loop.index0	循环的次数（从0开始）
- loop.revindex	循环剩余次数（最小值为1）
- loop.revindex0	循环剩余次数（最小值为0）
- loop.first	当第一次循环的时候返回true
- loop.last	当最后一次循环的时候返回true
- loop.length	循环的总数
- loop.parent	被循环的数组

**条件语句**

1. 需要使用 or 和 and 代替 ||、&&

```
{% if a == '1' or b == '2' %}
a = 1 or b = 2
{% endif %}
```

2. 判断变量是否定义

```
{% if var is not defined %}
    {# do something #}
{% endif %}
```

3. 是否为NULL

```
{% if var is null %}
    {# do something #}
{% endif %}
```

4. 是否为false

```
{% if var is sameas(false) %}
    {# do something %}
{% endif %}
```

5. 解析定界符

```
{{ '{{' }}
 
{% raw %}
    <ul>
    {% for item in seq %}
        <li>{{ item }}</li>
    {% endfor %}
    </ul>
{% endraw %}
```

# 控制结构

- {% if aaa %} xxx {% elseif bbb %} yyy {% else %} zzz：判断语句
- {% for %} xxx {% endfor %}：迭代变量
- {% do %}：没什么其他含义，{% do 1+2 %} 等同于 {{ 1+2 }}
- {% flush %}：刷新输出缓冲，等同于 flush
- {% include %}：包含模板
- {% extends %}：扩展模板
- {% embed %} xxx {% endembed %}：包含模板并扩展该模板的内容，相当于 include 和 extends 的结合体
- {% use %}：包含模板，近似于多重继承
- {% from aaa import bbb as ccc %}：从指定模板导入宏并设置别名
- {% macro %} xxx {% endmacro %}：定义宏以便多次调用，与定义 PHP 函数无异
- {% sandbox %} {% include xxx %} {% endsandbox %}：对导入的模板指定沙箱模式，只对 include 语句有效，只在沙箱模式已开启的情况下生效
- {% block xxx %} 或 {% block %} xxx {% endblock %}：定义代码块或覆盖代码块
- {% set xxx %} 或 {% set %} xxx {% endset %}：在模板内定义变量
- {% filter %} xxx {% endfilter %}：多行过滤器
- {% spaceless %} xxx {% endspaceless %}：去除 HTML 片段中的空格
- {% autoescape %} xxx {% endautoescape %}：将字符串安全地处理成合法的指定数据
- {% verbatim %} xxx {% endverbatim %}：阻止模板引擎的编译，是 raw 的新名字

# 内建过滤器

过滤器用来修饰数据，各过滤器可以用竖线分隔进行链式调用，用括号传递参数。也可以将过滤器当成单独的函数来用，形式如下：

```
{% filter 过滤器名 %}
待处理的数据
{% endfilter %}
```

- batch：将数组按指定的个数分割成更小的数组，可选的第二个参数用来在元素不够的情况下进行填充。如 {{ [1, 2, 3, 4, 5]|batch(2, 'NoItem') }} => [[1, 2], [3, 4], [5, 'NoItem']]
- date_modify：修改时间，常与 date 联用。如 {{ ''|date_modify('+3 days')|date('Y-m-d') }} => 将当前时间加3天后显示
- default：当所修饰的数据不存在或为空时，提供默认值。如 {{ ''|default('Ruchee') }} => 'Ruchee'
- escape：将字符串安全地处理成合法的指定数据，可简写为 e，支持多种转换模式，默认模式为 html，其他可选模式有 html_attr、js、css、url
- first：返回数组的第一个元素或字符串的第一个字符。如 {{ {a: 1, b: 2, c: 3}|first }} => 1
- last：返回数组的最后一个元素或字符串的最后一个字符。如 {{ {a: 1, b: 2, c: 3}|last }} => 3
- replace：替换一个字符串中的指定内容。如 {{ '%s1 love %s2'|replace({'%s1': 'Ruchee', '%s2': 'Vim'}) }} => 'Ruchee love Vim'
- raw：让数据在 autoescape 过滤器里失效

# 借用自PHP自带函数的过滤器

- abs：取绝对值
- nl2br：将字符串里的 `\n` 替换成 `<br/>`
- join：将数组的各个元素按指定分隔符组成字符串
- sort：对数组排序
- trim：去除字符串首尾的指定字符，默认为空格
- date：格式化时间，可处理与 strtotime 兼容的字符串，或 DateTime/DateInterval 的实例，可选的第二个参数用于指定时区，如果所修饰的数据为空则默认为当前时间
- reverse：反转一个数组或字符串，在 array_reverse 的基础上增加了对字符串的处理
- slice：截取数组或字符串的一部分，在 array_slice 的基础上增加了对字符串的处理
- keys：将数组的全部键名提取成一个数组，等同于 array_keys
- merge：合并两数组，近似于 array_merge 。如 {{ 数组1|merge(数组2) }}
- length：返回数组元素的个数或字符串的长度，等同于 count 和 strlen 的结合体
- capitalize：将字符串的首字母大写，等同于 ucfirst
- title：将字符串中每个单词的首字母大写，等同于 ucwords
- lower：将字符串所有字母全部变成小写，等同于 strtolower
- upper：将字符串所有字母全部变成大写，等同于 strtoupper
- split：将字符串分割成数组，等同于 str_split
- striptags：去除字符串中的 HTML/PHP 标记，等同于 strip_tags
- url_encode：编码链接字符串，等同于 urlencode
- json_encode：编码 JSON 格式，等同于 json_encode
- format：格式化一个字符串，近似于 printf 。如 {{ 'My name is %s, and I love %s'|format('Ruchee', 'Vim') }} => 'My name is Ruchee, and I love Vim'
- number_format：格式化数值，等同于 number_format
- convert_encoding：编码转换，第一个参数指定转换后的编码，第二个参数指定转换前的编码，近似于 iconv

# 内建函数

- even：是否为偶数
- odd：是否为奇数
- empty：是否为空
- null：是否为 null
- defined：是否已定义
- sameas：目标变量与指定值是否指向的是内存中的同一个地址，使用形式 if 变量值 is sameas(指定值)
- divisibleby：目标数值是否能够被指定值整除，使用形式 if 目标数值 divisibleby(指定值)，其中指定值不能为 0
- iterable：目标变量是否是数组或者是否可迭代，使用形式 if 变量值 is iterable
- attribute：动态获取变量属性值，两种使用形式为 attribute(数组, '元素名') 和 attribute(对象, '方法名', 可选参数)
- block：重复引用指定代码块，如 {{ block('title') }}
- constant：从字符串或对象取得常量值
- cycle：循环显示一个数组的元素，调用形式为 cycle(数组, 一个循环变量)
- date：格式化时间
- dump：在开启调试模式的情况下显示详细的变量信息，等同于 var_dump
- include：包含其他模板文件
- parent：在覆盖代码片段时用于引用父片段的内容
- random：制造一个随机数
- range：返回一个指定区间的数组，可指定步长，Twig 使用 .. 作为其简用法，等同于 range
- template_from_string：根据字符串加载模板
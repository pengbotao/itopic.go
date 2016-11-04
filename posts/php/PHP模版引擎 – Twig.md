```
{
    "url": "php-moban-yinqing-twig",
    "time": "2013/01/26 03:25",
    "tag": "php"
}
```

在网站开发过程中模版引擎是必不可少的，PHP中用的最多的当属Smarty了。目前公司系统也是用的Smarty，如果要新增一个页面只需把网站的头、尾和左侧公共部分通过Smarty的include方式引入进来，然后主体部分写内容即可，用起来也是相当方便。这也是一种比较通用的做法。但维护一段时间后发现有些凌乱了：

- 1. 公共部分内容越加越多了，不需要用的js、css在一些页面也被强制引进来了
- 2. 新页面的css只能写在网页的body内，看起来总让人不爽。
- 3. 左侧、头部、尾部若有特殊显示，操作起来不方便，只能在公共地方去做判断了。

当然这些页面问题在设计的时候可以通过合理的拆分网页来实现，当然最重要的还在于开发人员，在好的系统也经不起开发人员的折腾，一个项目经过多次转手后，接下来的维护人员那是相当痛苦的。不扯远了， 现在要说的是另一种模版开发思路。

在PHP中CLASS用过很多次了，有一个很有用的特性那就是继承，子类继承父类后可以直接调用父类的方法，也可以对父类的方法进行重写，同样PHP的模版引擎Twig也实现了这一点，模版的书写方式可以更方便。

Twig是开源框架Symfony2的默认模版引擎，主页是http://twig.sensiolabs.org/ 当前版本为Stable: `1.12.1 `，其他模版引擎能做的它都能做，这里主要整理下使用Twig的BLOCK方式编写模版页面。

以一个常见的排版为例，有三个链接，分别是首页、关于、联系三个页面，然后头部共用，尾部共用，中间部分分成左右两部分，左边共用，右边显示具体内容，貌似很多后台都是这种布局。
先看看首页 twig_index.php ， 和Smarty差不多，初始化设置，然后设置变量并显示。

```
<?php
require './Twig-1.12.1/lib/Twig/Autoloader.php';
Twig_Autoloader::register();
 
$loader = new Twig_Loader_Filesystem('./view/twig/templates');
$twig = new Twig_Environment($loader, array(
    'cache' => './view/twig/templates_c',
    'auto_reload' => true
));
 
$twig->display('index.html', array('name' => 'Bobby'));
```
其他页面的PHP内容除了模版名称不一样外，其他内容完全一样，所以后面的PHP页面就不写了。
那接下来的主要工作就是写模版了，既然支持继承，那应该有一个父类，其他页面来继承这个类。base.html就是模版的父类，其内容如下：
```
{# /view/twig/templates/base.html #}
<!DOCTYPE html>
<html>
<head>
<title>{% block title %}Home{% endblock %} - Twig</title>
<style type="text/css">
{% block stylesheet %}
#header, #main, #footer{width:600px; margin:0px auto; border:1px solid #333;}
#main{border:0px;}
#header { height:120px;margin-bottom:10px;}
#footer{ height:80px;clear:both;margin-top:10px;}
#header h1{margin-left:20px;}
#header span{margin-left:20px;}
#leftsider{width:125px; float:left; border:1px solid #333; height:200px;
padding-top:10px;}
#leftsider span{width:100%; height:30px; line-height:30px; clear:both; 
padding-left:15px; display:block;}
#rightsider{width:460px; float:right; border:1px solid #333; height:250px;}
.clear{clear:both;}
{% endblock %}
</style>
</head>
<body>
    <div id="header">
        <h1>Twig Template Test</h1>
        <span><a href="twig_index.php">Home</a></span>
        <span><a href="twig_about.php">About</a></span>
        <span><a href="twig_contact.php">Contact</a></span>
    </div>
    <div id="main">
        <div id="leftsider">
            {% block leftsider %}
            <span>系统模块1</span>
            <span>系统模块2</span>
            {% endblock %}
        </div>
        <div id="rightsider">
            {% block rightsider %}Hello {{name}}{% endblock %}
        </div>
    </div>
    <div class="clear"></div>
    <div id="footer">
        {% block footer %}<h1>Twig Footer</h1>{% endblock %}
    </div>
</body>
</html>
```
基本的页面框架没太多说的，主要看看中间有5个block - {% block blockname%}{% endblock %}  每个BLOCK代表一个块， 这里的块可以理解成PHP父类中的一个方法。

基本的html框架搭好后，index.html该如何来写呢？首先该继承base页面，然后再考虑是否要重写base页面的内容，先只做继承看看效果。
```
{# /view/twig/templates/index.html #}
{% extends "base.html" %}
```
第一行为注释部分，可以省去，第二行表示index.html继承base.html， 未重写的情况下将直接使用base.html中的内容进行显示，也就是该页面将看到如下效果：

![](/static/uploads/twig-1.png)

效果比较简单，但是很神奇，index.html只是继承了base.html，没写其他内容呢？对，不用写了，在未重写父类方法时。子类是可以直接调用父类方法的。

那接着看看about， 假设about页面和index页面除了右边区域不同外，其他部分完全相同，也就是只需要重写rightsider这个BLOCK：
```
{# /view/twig/templates/about.html #}
{% extends "base.html" %}
{% block title %}About{% endblock %}
{% block rightsider %} {% include 'about_content.html' %} {% endblock %}
```
标题的内容改成了 About， rightsider的内容从about_content.html文件中读取，其他部分保留原有。也就是除了Hello Bobby的内容不同外，其他部分与首页都是相同的，是不是觉得很方便了？

再来看一下Contact页面怎么写？我么需要在leftsider里增加一个菜单，以及rightsider里显示其他block的内容。看看下面：
```
{# /view/twig/templates/contact.html #}
{% extends "base.html" %}
{% block title %}Contact{% endblock %}
{% block leftsider %}{{ parent() }}<span>系统模块3</span>{% endblock %}
{% block rightsider %} {{ block('footer') }} {% endblock %}
```
调用parent即可显示基类的内容，通过block('footer')则可获取footer中的Twig Footer内容。所以图片效果如下：

![](/static/uploads/twig-2.png)

很神奇吧！这种排版方式值得一试，等待机会中...

使用block后子页面不可以按照html的方式在任意地方加html， 也就是在block外写任何内容都会报错，所以需要base里去合理的设置block，block设置的越多就越灵活。具体的还得到实际项目中去尝试。

至于Twig的具体语法有时间在整理下，不过这种写模版的方式确实很让人喜欢，好像Smarty3也支持该功能了，有时间也看看。

看到Twig后联想到了 lesscss, 动态样式语言，主页http://www.lesscss.net 有兴趣的朋友可以看看。
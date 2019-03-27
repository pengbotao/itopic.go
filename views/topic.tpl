<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.topic.Title}} - 老彭的博客</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<h1 class="title">{{.topic.Title}}</h1>

<a href="{{.domain}}/"><img src="/static/img/arrow-back.png" class="title_arrow_back" /></a>

{{.topic.Content}}
<div class="eof">-- EOF --</div>
<div class="eof_arrow">
    <a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;" /></a>
</div>
<div class="eof_tag">
    发表于：
    <code style="border:0px;background:none;"><a href="/{{.topic.Time.Format "2006-01"}}.html">{{.topic.Time.Format "2006-01-02 15:04"}}</a></code>
</div>
<div class="eof_tag">
    标签：{{range .topic.Tag}}
    <code style="border:0px;background:none;">{{if .TagID}}<a href="/tag/{{.TagID}}.html">{{.TagName}}</a>{{else}}{{.TagName}}{{end}}</code>{{end}}
</div>

<div id="footer">
    <ul>
        <li>
        @2013-{{.time.Format "2006"}} 老彭的博客&nbsp;[Hosted by <a href="javascript:;" style="font-weight: bold" target="_blank">Github Pages</a>]</li>
        <li><b>Github地址</b>：<a href="{{.githubURL}}/blob/master/{{.topic.TopicPath}}">{{.githubURL}}/blob/master/{{.topic.TopicPath}}</a><li>
    </ul>
</div>

<div id="top"><a href="#"><img src="/static/img/arrow-top.png" style="width:40px;height:40px;" /></a></div>
<script>
ready(fn);
function ready(fn){  
    if(document.addEventListener){
        document.addEventListener('DOMContentLoaded',function(){
            document.removeEventListener('DOMContentLoaded',arguments.callee,false);  
            fn();
        },false);
    } else if(document.attachEvent) {
        document.attachEvent('onreadystatechange',function(){
            if(document.readyState=='complete'){
                document.detachEvent('onreadystatechange',arguments.callee);  
                fn();
            }
        });
    }  
}  
function fn(){
    if(document.getElementsByTagName("nav")[0].innerText == "") {
        document.getElementsByTagName("body")[0].style.marginLeft = "0";
    }
}
</script>
</body>
</html>
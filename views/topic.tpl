<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.topic.Title}} - 老彭的博客</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
    <script>
    var _hmt = _hmt || [];
    (function() {
        var hm = document.createElement("script");
        hm.src = "https://hm.baidu.com/hm.js?0f0111c99240380ee020030f3be990f5";
        var s = document.getElementsByTagName("script")[0]; 
        s.parentNode.insertBefore(hm, s);
    })();
    </script>
</head>
<body>
<h1 class="title">{{.topic.Title}}</h1>
<a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;float:right;margin-top:-30px;" /></a>

{{.topic.Content}}
<div style="padding: 0 10px;float:left;margin-bottom:20px;color:#aaa;">-- EOF --</div>
<div style="float:right;margin-right:50px;">
    <a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;" /></a>
</div>
<div style="padding: 0 10px;text-align:right;float:right;">
    发表于：
    <code style="border:0px;background:none;"><a href="/{{.topic.Time.Format "2006-01"}}.html">{{.topic.Time.Format "2006-01-02 15:04"}}</a></code>
</div>
<div style="padding: 0 10px;text-align:right;float:right;">
    标签：{{range .topic.Tag}}
    <code style="border:0px;background:none;">{{if .TagID}}<a href="/tag/{{.TagID}}.html">{{.TagName}}</a>{{else}}{{.TagName}}{{end}}</code>{{end}}
</div>

<div id="disqus_thread" style="margin-bottom:35px;"></div>
<script>

var disqus_config = function () {
    this.page.url = "http://itopic.org/{{.topic.TopicID}}.html";
    this.page.identifier = "{{.topic.TopicID}}";
};

setTimeout(function() { // DON'T EDIT BELOW THIS LINE
    var d = document, s = d.createElement('script');
    s.src = 'https://itopic.disqus.com/embed.js';
    s.setAttribute('data-timestamp', +new Date());
    (d.head || d.body).appendChild(s);
}, 3000);

</script>

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
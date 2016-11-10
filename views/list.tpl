<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>iTopic - 老彭的博客</title>
    <link rel="stylesheet" href="{{.domain}}/static/css/markdown.css">
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
<h1 style="font-weight:600;width:90%;margin-bottom:0px;">{{.title}}存档</h1>
<a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;float:right;margin-top:-15px;" /></a>
<hr />
<ul>{{range .topics}}
    <li>[{{.Time.Format "06-01-02"}}] <a href="{{$.domain}}/{{.TopicID}}.html">{{.Title}}</a></li>{{end}}
</ul>
</body>
</html>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.title}} - 老彭的博客</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<a href="/"><img src="/static/img/arrow-back.png" style="width:30px;height:30px;float:right;" /></a>
<ul>{{range .topics}}
    <li><a href="/{{.TopicID}}.html">{{.Title}}</a></li>{{end}}
</ul>
</body>
</html>
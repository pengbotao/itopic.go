<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>iTopic - 老彭的博客</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<h1 style="font-weight:600;width:90%;margin-bottom:0px;">{{.title}}存档</h1>
<a href="/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;float:right;margin-top:-15px;" /></a>
<hr />
<ul>{{range .topics}}
    <li>[{{.Time.Format "06-01-02"}}] <a href="/{{.TopicID}}.html">{{.Title}}</a></li>{{end}}
</ul>
</body>
</html>
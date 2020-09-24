<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.title}}</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<div class="container">
    <h1 class="title">{{.title}}</h1>
    <a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;float:right;margin-top:-30px;" /></a>

    <ul>{{range .topics}}
        <li>[{{.Time.Format "06-01-02"}}] <a href="{{$.domain}}/{{.TopicID}}.html">{{.Title}}</a></li>{{end}}
    </ul>
</div>
</body>
</html>
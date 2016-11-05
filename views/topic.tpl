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
<h1 style="font-weight:600;">{{.topic.Title}}</h1>
<hr />
{{.topic.Content}}
<div style="padding: 0 10px;margin-bottom:20px;text-align:right;">
    <span>标签：</span>{{range .topic.Tag}}
    <code>{{if .TagID}}<a href="/tag/{{.TagID}}.html">{{.TagName}}</a>{{else}}{{.TagName}}{{end}}</code>{{end}}
    <span style="margin-left:30px;">发表于：</span>
    <code><a href="/{{.topic.Time.Format "2006-01"}}.html">{{.topic.Time.Format "2006-01-02 15:04"}}</a></code>
</div>
</body>
</html>
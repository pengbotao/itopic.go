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
<h1 style="font-weight:600;width:90%;margin-bottom:0px;">{{.topic.Title}}</h1>
<a href="/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;float:right;margin-top:-15px;" /></a>
<hr />
{{.topic.Content}}
<div style="padding: 0 10px;float:left;margin-bottom:20px;color:#aaa;">-- EOF --</div>
<div style="float:right;">
    <a href="/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;" /></a>
</div>
<div style="padding: 0 10px;text-align:right;float:right;">
    发表于：
    <code style="border:0px;background:none;"><a href="/{{.topic.Time.Format "2006-01"}}.html">{{.topic.Time.Format "2006-01-02 15:04"}}</a></code>
</div>
<div style="padding: 0 10px;text-align:right;float:right;">
    标签：{{range .topic.Tag}}
    <code style="border:0px;background:none;">{{if .TagID}}<a href="/tag/{{.TagID}}.html">{{.TagName}}</a>{{else}}{{.TagName}}{{end}}</code>{{end}}
</div>
</body>
</html>
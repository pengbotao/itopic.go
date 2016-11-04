<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.topic.Title}}</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<h1>{{.topic.Title}}</h1>
<code>{{.topic.Time}}</code>
{{range .topic.Category}}
<code><a href="/category/{{.CategoryID}}.html">{{.Title}}</a></code>
{{end}}
{{.topic.Content}}
</body>
</html>
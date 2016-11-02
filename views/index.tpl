<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>Homepage</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<div id="wrapper">
{{range .topics_m}}
    <h1><a href="/{{.Month}}.html">{{.Month}}</a></h1>
    <ul>
    {{range .Topics}}
        <li><a href="/{{.TopicId}}.html">{{.Title}}</a></li>
    {{end}}
    </ul>
{{end}}
</div>
</body>
</html>
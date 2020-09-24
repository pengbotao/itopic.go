<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>iTopic - 老彭的博客</title>
    <meta name="description" content="记录和分享程序开发过程中的经验和体会。" />
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<div class="container">
    <h1 class="title">iTopic.org</h1>

    <div id="left-sider">
        {{range .topics_l}}
        <ul>{{range .Topics}}
            <li>[{{.Time.Format "06-01-02"}}] <a href="{{$.domain}}/{{.TopicID}}.html">{{.Title}}</a></li>{{end}}
        </ul>
        {{end}}
    </div>

    <div id="right-sider">
        {{range .topics_r}}
        <ul>{{range .Topics}}
            <li>[{{.Time.Format "06-01-02"}}] <a href="{{$.domain}}/{{.TopicID}}.html">{{.Title}}</a></li>{{end}}
        </ul>
        {{end}}
    </div>

    <div id="footer">
        <ul>
            <li>- 古人学问无遗力，少壮工夫老始成</li>
            <li>
            @2013-{{.time.Format "2006"}} 老彭的博客&nbsp;[Hosted by <a href="https://pages.github.com/" style="font-weight: bold" target="_blank">Github Pages</a>] <b>Github地址</b>：{{.githubURL}}</li>
        </ul>
    </div>
</div>
</body>
</html>
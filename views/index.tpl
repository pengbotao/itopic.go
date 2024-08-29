<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>首页 - IT驿站</title>
    <meta name="description" content="记录和分享IT开发过程中的知识点、技术实践和编程经验，内容涵盖编程语言、框架应用、开发工具使用以及项目管理等方面。" />
    <link rel="stylesheet" href="/static/css/markdown.css">
</head>
<body>
<div class="container">
    <h1 class="title">因上努力，果上随缘</h1>

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
            @2013-{{.time.Format "2006"}} 老彭的博客&nbsp; [鄂ICP备2024070731号-1] <b>Github地址</b>：{{.githubURL}}</li>
        </ul>
    </div>
</div>

</body>
</html>
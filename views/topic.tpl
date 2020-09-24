<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.title}}</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
    {{ if .topic.IsToc }}<style type="text/css">
        .content{padding-left:280px;width:100%;}
    </style>{{ end }}
</head>
<body>
<div class="content">
    <h1 class="title">{{.topic.Title}}</h1>

    <a href="{{.domain}}/"><img src="/static/img/arrow-back.png" class="title_arrow_back" /></a>

    {{.topic.Content}}
    <div class="eof">-- EOF --</div>
    <div class="eof_arrow">
        <a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;" /></a>
    </div>
    {{ if not .topic.LastModifyTime.IsZero }}
    <div class="eof_tag">
        最后更新于：
        <code style="border:0px;background:none;"><a href="{{$.domain}}/{{.topic.Time.Format "2006-01"}}.html">{{.topic.LastModifyTime.Format "2006-01-02 15:04"}}</a></code>
    </div>
    {{ end }}
    <div class="eof_tag">
        发表于：
        <code style="border:0px;background:none;"><a href="{{$.domain}}/{{.topic.Time.Format "2006-01"}}.html">{{.topic.Time.Format "2006-01-02 15:04"}}</a></code>
    </div>
    <div class="eof_tag">
        标签：{{range .topic.Tag}}
        <code style="border:0px;background:none;">{{if .TagID}}<a href="{{$.domain}}/tag/{{.TagID}}.html">{{.TagName}}</a>{{else}}{{.TagName}}{{end}}</code>{{end}}
    </div>

    <div id="footer">
        <ul>
            {{if .topic_left.Title}}<li>
            <b>上一篇</b>：<a href="{{$.domain}}/{{.topic_left.TopicID}}.html">{{.topic_left.Title}}</a>
            </li>
            {{end}}{{if .topic_right.Title}}
            <li>
            <b>下一篇</b>：<a href="{{$.domain}}/{{.topic_right.TopicID}}.html">{{.topic_right.Title}}</a>
            </li>{{end}}
            <li>
                <b>Github地址</b>：<a href="{{.githubURL}}/blob/master/{{.topic.TopicPath}}">{{.githubURL}}/blob/master/{{.topic.TopicPath}}</a>
            <li>
            <li>
                @2013-{{.time.Format "2006"}} 老彭的博客&nbsp;[Hosted by <a href="https://pages.github.com/" style="font-weight: bold" target="_blank">Github Pages</a>]
            </li>
        </ul>
    </div>
</div>
<div id="top"><a href="#"><img src="/static/img/arrow-top.png" style="width:40px;height:40px;" /></a></div>
</body>
</html>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.title}}</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
    {{ if .topic.IsToc }}<style type="text/css">
    @media (min-width: 1200px) {
        .content{padding-left:280px;width:100%;max-width:100%;}
    }
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

<a href="{{.githubURL}}/blob/master/{{.topic.TopicPath}}" target="_blank"  class="github-corner">
<svg width="60" height="60" viewBox="0 0 250 250" style="fill: #2D8CF0; color:#fff; position: absolute;top: 0;border: 0;right: 0;" aria-hidden="true">
    <path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path>
    <path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2" fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path>
    <path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z" fill="currentColor" class="octo-body"></path>
    </svg>
</a>
</body>
</html>
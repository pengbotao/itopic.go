<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.topic.Title}} - 老彭的博客</title>
    <link rel="stylesheet" href="{{.domain}}/static/css/markdown.css">
    <script>
    var _hmt = _hmt || [];
    (function() {
        var hm = document.createElement("script");
        hm.src = "https://hm.baidu.com/hm.js?0f0111c99240380ee020030f3be990f5";
        var s = document.getElementsByTagName("script")[0]; 
        s.parentNode.insertBefore(hm, s);
    })();
    </script>
</head>
<body>
<h1 style="font-weight:600;width:90%;margin-bottom:0px;">{{.topic.Title}}</h1>
<a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;float:right;margin-top:-15px;" /></a>
<hr />
{{.topic.Content}}
<div style="padding: 0 10px;float:left;margin-bottom:20px;color:#aaa;">-- EOF --</div>
<div style="float:right;">
    <a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;" /></a>
</div>
<div style="padding: 0 10px;text-align:right;float:right;">
    发表于：
    <code style="border:0px;background:none;"><a href="/{{.topic.Time.Format "2006-01"}}.html">{{.topic.Time.Format "2006-01-02 15:04"}}</a></code>
</div>
<div style="padding: 0 10px;text-align:right;float:right;">
    标签：{{range .topic.Tag}}
    <code style="border:0px;background:none;">{{if .TagID}}<a href="/tag/{{.TagID}}.html">{{.TagName}}</a>{{else}}{{.TagName}}{{end}}</code>{{end}}
</div>

<!-- 多说评论框 start -->
	<div class="ds-thread" data-thread-key="{{.topic.TopicID}}" data-title="{{.topic.Title}}" data-url="http://itopic.org/{{.topic.TopicID}}.html"></div>
<!-- 多说评论框 end -->
<!-- 多说公共JS代码 start (一个网页只需插入一次) -->
<script type="text/javascript">
var duoshuoQuery = {short_name:"itopic"};
	(function() {
		var ds = document.createElement('script');
		ds.type = 'text/javascript';ds.async = true;
		ds.src = (document.location.protocol == 'https:' ? 'https:' : 'http:') + '//static.duoshuo.com/embed.js';
		ds.charset = 'UTF-8';
		(document.getElementsByTagName('head')[0] 
		 || document.getElementsByTagName('body')[0]).appendChild(ds);
	})();
	</script>
<!-- 多说公共JS代码 end -->
</body>
</html>
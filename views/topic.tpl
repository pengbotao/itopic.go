<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <title>{{.topic.Title}} - 老彭的博客</title>
    <link rel="stylesheet" href="/static/css/markdown.css">
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
<h1 style="font-weight:400;width:90%;margin-bottom:0px;border:0px;">{{.topic.Title}}</h1>
<a href="{{.domain}}/"><img src="/static/img/arrow-back.png" style="width:25px;height:25px;float:right;margin-top:-30px;" /></a>

<div class="main-topic-content">
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

<BR>

<div id="cloud-tie-wrapper" class="cloud-tie-wrapper"></div>
<script>
  var cloudTieConfig = {
    url: document.location.href, 
    sourceId: "{{.topic.TopicID}}",
    productKey: "47a1277aece74470855c0e74c1208eaf",
    target: "cloud-tie-wrapper"
  };
</script>
<script src="http://img1.cache.netease.com/f2e/tie/yun/sdk/loader.js"></script>
</div>

<div id="top"><a href="#"><img src="/static/img/arrow-top.png" style="width:40px;height:40px;" /></a></div>
</body>
</html>
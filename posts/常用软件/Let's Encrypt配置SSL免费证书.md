```
{
    "url": "lets-encrypt",
    "time": "2018/01/09 18:52",
    "tag": "常用软件"
}
```

# 一、Let's Encrypt介绍

Let's Encrypt是一个于2015年三季度推出的数字证书认证机构，旨在以自动化流程消除手动创建和安装证书的复杂流程，并推广使万维网服务器的加密连接无所不在，为安全网站提供免费的SSL/TLS证书。

官网：https://letsencrypt.org/

# 二、Let's Encrypt安装

官网上推荐使用`certbot`，安装起来也是相当方便。官网地址：`https://certbot.eff.org`

## 2.1 下载certbot-auto

选择软件和系统后会有安装提示：

```
wget https://dl.eff.org/certbot-auto
chmod a+x certbot-auto
```

## 2.2 申请证书
```
./certbot-auto certonly --webroot -w /data/gopath/src/github.com/pengbotao/itopic.go -d itopic.org  --agree-tos --email pengbotao@vip.qq.com
```

说明：`certbot-auto`需要`python3`支持，可参考[Python多环境及包安装](./python-environment.html)安装。


安装好虚拟环境后执行命令提示：`To use Certbot, packages from the EPEL repository need to be installed.`

目前参考该页面解决：`https://unix.stackexchange.com/questions/165916/trying-to-enable-epel-on-centos-6-and-it-wont-show-in-repolist`

```
The CentOS Extras repository should be enabled by default, so you can just run:

sudo rpm -e epel-release
to remove the existing package and then:

sudo yum install epel-release
to enable EPEL.
```

证书安装之后提示安装路径：

```
IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/itopic.org/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live/itopic.org/privkey.pem
   Your cert will expire on 2019-06-07. To obtain a new or tweaked
   version of this certificate in the future, simply run certbot-auto

   again. To non-interactively renew *all* of your certificates, run
   "certbot-auto renew"
 - Your account credentials have been saved in your Certbot
   configuration directory at /etc/letsencrypt. You should make a
   secure backup of this folder now. This configuration directory will
   also contain certificates and private keys obtained by Certbot so
   making regular backups of this folder is ideal.
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le
```


## 2.3 Nginx使用证书

本站点下Nginx配置方式如下，配置之后访问`https`访问正常即可。

```
server
{
    server_name itopic.org;
    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/itopic.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/itopic.org/privkey.pem;
    root /data/gopath/src/github.com/pengbotao/itopic.go;

    location / {
        if ($request_uri = "/") {
            proxy_pass http://127.0.0.1:8001;
        }
        if (!-e $request_filename) {
            proxy_pass http://127.0.0.1:8001;
        }
    }
}
```

同时将原站点跳转到`https`:

```
server
{
    listen 80;
    server_name itopic.org;
    return 301 https://itopic.org$request_uri;
}
```

# 三、刷新证书

免费证书的有效期为90天，做一次自动刷新即可(每2个月的0点过5分刷新)。

```
5 0 * */2 * /root/letsencrypt/certbot-auto renew --renew-hook "/usr/local/server/nginx1.15.3/sbin/nginx -s reload"
```

参考站点：`https://segmentfault.com/a/1190000012343679`
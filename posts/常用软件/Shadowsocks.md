```
{
    "url": "shadowsocks",
    "time": "2016/11/21 22:45",
    "tag": "常用软件",
    "toc": "no"
}
```

- 客户端：https://shadowsocks.org/en/download/clients.html
- 服务端：https://shadowsocks.org/en/download/servers.html

`Github`上有各种版本：`Python`, `Go`, `Rust`等

# 服务端安装

## Golang服务端安装

```
go get github.com/shadowsocks/shadowsocks-go/cmd/shadowsocks-server
```

或者直接从发布版本中选择下载即可：`https://github.com/shadowsocks/shadowsocks-go/releases`

**config.json**

```
{
    "server":"127.0.0.1",
    "server_port":8388,
    "local_port":1080,
    "local_address":"127.0.0.1",
    "password":"barfoo!",
    "method": "aes-128-cfb",
    "timeout":600
}
```


详细安装见： `https://github.com/shadowsocks/shadowsocks-go`

## 启动

```
./shadowsocks-server -c config.json
```

# 客户端安装

`https://github.com/shadowsocks/shadowsocks-windows/releases`

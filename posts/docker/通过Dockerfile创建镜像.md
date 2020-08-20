```
{
    "url": "dockerfile",
    "time": "2020/07/09 14:08",
    "tag": "Docker"
}
```

# 一、Dockerfile

## 1.1 概述

`Docker`通过读取`Dockerfile`中的指令来自动构建镜像，`Dockerfile`是一个文本文件，它包含了构建镜像的所有命令。可以通过`docker build`读取`Dockerfile`文件来构建镜像。`docker build`提交构建镜像请求给`Docker daemon`，同时会将当前目录递归传过去。所以最好是将`Dockerfile`和需要文件放到一个空目录，再在这个目录构建。

- 常规构建: `$ docker build -t pengbotao/itopic.go .`
- 指定`Dockerfile`: `$ docker build -f /path/to/a/Dockerfile .`
- 标记为多个镜像：`$ docker build -t shykes/myapp:1.0.2 -t shykes/myapp:latest .`

## 1.2 文件格式

可以通过`Dockerfile`文件来定义构建过程，通过`.dockerignore`用来忽略文件，构建时不会讲ignore的文件传给`Docker`服务端。

1、 **`Dockerfile`**：

```
# Comment
INSTRUCTION arguments
```
`INSTRUCTION`不区分大小写，建议大写。`Dockerfile`中的指令第一个指令必需是`FROM`用来指定基础镜像。`#`用来注释。

2、 **`.dockerignore`**:

```
# comment
*/temp*
*/*/temp*
temp?
```

此文件导致以下构建行为：

规则|行为
---|---
`#` comment|忽略
`*/temp*`|不包含子目录中`temp`打头的文件和目录。如：文件 `/somedir/temporary.txt`，目录 `/somedir/temp`
`*/*/temp*`|不包含子目录下的子目录中`temp`打头的文件和目录。 如：`/somedir/subdir/temporary.txt`
`temp?`|不包含根目录中`temp`打头的文件和目录。 如：`/tempa` 和 `/tempb`

## 1.3 指令集合

编号|指令|参数|必须|示例
---|---|---|---|---
1|FROM|`FROM [--platform=<platform>] <image> [AS <name>]` <BR> `FROM [--platform=<platform>] <image>[:<tag>] [AS <name>]` <BR> `FROM [--platform=<platform>] <image>[@<digest>] [AS <name>]`|是|`FROM nginx`
2|RUN|`RUN <command>` <BR> `RUN ["executable", "param1", "param2"]`|否|`RUN ["/bin/bash", "-c", "echo hello"]`
3|CMD|`CMD ["executable","param1","param2"]` <BR> `CMD ["param1","param2"]` <BR> `CMD command param1 param2`|否|`CMD ["/usr/bin/wc","--help"]`
4|LABEL|`LABEL <key>=<value> <key>=<value> <key>=<value> ...`|否|`LABEL version="1.0"`
5|EXPOSE|`EXPOSE <port> [<port>/<protocol>...]`|否|`EXPOSE 80/tcp`
6|ENV|`ENV <key> <value>` <BR> `ENV <key>=<value> ...`|否|`ENV myCat fluffy`
7|ADD|`ADD [--chown=<user>:<group>] <src>... <dest>` <BR> `ADD [--chown=<user>:<group>] ["<src>",... "<dest>"]`|否|`ADD hom?.txt /mydir/`
8|COPY|`COPY [--chown=<user>:<group>] <src>... <dest>` <BR> `COPY [--chown=<user>:<group>] ["<src>",... "<dest>"]`|否|`COPY hom?.txt /mydir/`
9|ENTRYPOINT|`ENTRYPOINT ["executable", "param1", "param2"]` <BR> `ENTRYPOINT command param1 param2`|否|`ENTRYPOINT ["top", "-b"]`
10|VOLUME|`VOLUME ["/data"]`|否|`VOLUME /var/log` <BR> `VOLUME ["/var/log"]`
11|USER|`USER <user>[:<group>]` <BR> `USER <UID>[:<GID>]`|否|`USER patrick`
12|WORKDIR|`WORKDIR /path/to/workdir`|否|`WORKDIR /data`
13|ARG|`ARG <name>[=<default value>]`|否|`ARG user1=someuser`
14|ONBUILD|`ONBUILD <INSTRUCTION>`|否|`ONBUILD ADD . /app/src`
15|STOPSIGNAL|`STOPSIGNAL signal`|否|
16|HEALTHCHECK|`HEALTHCHECK [OPTIONS] CMD command` <BR> `HEALTHCHECK NONE`|否|
17|SHELL|`SHELL ["executable", "parameters"]`|否|`SHELL ["powershell", "-command"]`


# 二、指令详解

@todo

# 三、Dockerfiles最佳实践

## 3.1 多阶段构建

前一篇文章中构建的镜像大小有901M，Go的基础镜像就比较大。可以尝试将编译后的二进制放在新的容器中执行，解除对golang环境的依赖。

```
REPOSITORY            TAG        IMAGE ID       CREATED             SIZE
pengbotao/itopic.go   latest     36e405364fd8   25 hours ago        901MB
```

如，将编译后的`itopic`拷贝到`alpine`系统中，在看大小就只有`27M`了。

```
FROM golang:1.14 AS build

COPY . /go/src/github.com/pengbotao/itopic.go/
RUN cd /go/src/github.com/pengbotao/itopic.go/ \
&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o itopic


FROM alpine
COPY --from=build /go/src/github.com/pengbotao/itopic.go/ /www/itopic.go/
EXPOSE 8001
WORKDIR /www/itopic.go
CMD [ "-host", "0.0.0.0:8001" ]
ENTRYPOINT [ "/www/itopic.go/itopic" ]
```


- [1] [Dockerfile reference](https://docs.docker.com/engine/reference/builder/)
- [2] [Best practices for writing Dockerfiles](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [3] [Dockerfile 最佳实践及示例](http://www.dockerone.com/article/9551)
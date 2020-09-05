```
{
    "url": "dockerfile",
    "time": "2020/07/09 14:08",
    "tag": "Docker,容器化"
}
```

# 一、Dockerfile

## 1.1 概述

`Docker`通过读取`Dockerfile`中的指令来自动构建镜像，`Dockerfile`是一个文本文件（可以方便做版本管理），它包含了构建镜像的所有命令。可以通过`docker build`读取`Dockerfile`文件来构建镜像。`docker build`提交构建镜像请求给`Docker`守护进程，同时会将当前目录递归传过去。所以最好是将`Dockerfile`和需要文件放到一个空目录，再在这个目录构建。

## 1.2 构建方法

- 常规构建: 
```
$ docker build -t pengbotao/itopic.go .
```

- 指定`Dockerfile`: 
```
$ docker build -f /path/to/a/Dockerfile .
```
- 标记为多个镜像：
```
$ docker build -t shykes/myapp:1.0.2 -t shykes/myapp:latest .
```

## 1.3 文件格式

可以通过`Dockerfile`文件来定义构建过程，通过`.dockerignore`用来忽略文件，构建时不会将`ignore`的文件传给`Docker`服务端。

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

## 1.4 指令集合

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

前一篇文章中构建的镜像大小有901M，`Go`的基础镜像就比较大。可以尝试将编译后的二进制放在新的容器中执行，解除对`golang`环境的依赖。

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

## 3.2 通过.dockerignore忽略无关文件

构建时会将该指定目录传到`Docker`守护进程，所以一般建议在空目录或者只包含需要的文件，对于项目中不需要的文件可以通过`.dockerignore`文件忽略。

## 3.3 减少镜像层数

**a. 镜像与层**

Docker镜像由一系列层组成。每层代表Dockerfile中的一条指令。除最后一层外的每一层都是只读的。看以下Dockerfile：

```
FROM ubuntu:18.04
COPY . /app
RUN make /app
CMD python /app/app.py
```

包含了4条指令，每个命令会创建一层，`FROM`语句首先从`ubuntu:18.04`镜像创建一层。然后`COPY`命令从当前目录拷贝一些文件。`RUN`命令执行`make`操作，最后一层指定在容器中运行的命令。每个层仅包含了前一层的差异部分。 当我们启动一个容器的时候，`Docker`会加载镜像层并在其上添加一个可写层。容器上所做的任何更改，譬如新建文件、更改文件、删除文件，都将记录与可写层上。容器层与镜像层的结构如下图所示。

![](../../static/uploads/container-layers.jpg)

**b. 容器与层**

容器和镜像之间的主要区别是最顶层的可写层。在容器中所有写操作都存储在可写层中。删除容器后，可写层也会被删除，基础镜像保持不变。因为每个容器都有自己的可写容器层，并且所有更改都存储在该容器层中，所以多个容器可以共享对同一基础镜像。下图显示了共享相同Ubuntu 18.04镜像的多个容器。

![](../../static/uploads/sharing-layers.jpg)


## 3.4 清理无用数据

通常为命令执行过程中的一些缓存或者中间数据，可以随镜像创建完成时执行删除操作。


# 四、构建示例

## 4.1 php5.6

```
docker build -t pengbotao/php:5.6-fpm-alpine .
```

安装了常用PHP扩展：`redis` `gd` `bcmath` `bz2` `pdo_mysql` `mysqli` `opcache` `sockets` `pcntl` `xsl` `soap` `dom` `shmop` `zip` `mcrypt`

```
FROM php:5.6-fpm-alpine

WORKDIR "/data/www"

RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
&& apk update \
&& apk add autoconf libmcrypt-dev bzip2-dev libxslt-dev freetype-dev libjpeg-turbo-dev libpng-dev zlib-dev libzip-dev\
&& docker-php-ext-install bcmath bz2 pdo_mysql mysqli opcache sockets pcntl xsl soap dom shmop zip mcrypt\
&& docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ --with-png-dir=/usr/include/ \
&& docker-php-ext-install gd \
&& mkdir -p /usr/src/php/ext/redis \
&& curl -L https://github.com/phpredis/phpredis/archive/3.1.6.tar.gz | tar xvz -C /usr/src/php/ext/redis --strip 1 \
&& docker-php-ext-install redis \
&& docker-php-ext-enable redis 

CMD [ "php-fpm" ]
```

推送到Docker Hub

```
docker push pengbotao/php:5.6-fpm-alpine
```

## 4.2 php7.4.8

```
docker build -t pengbotao/php:7.4.8-fpm-alpine .
```

和上面类似，通过`PECL`安装了`Redis`扩展（http://pecl.php.net/package/redis）

```
FROM php:7.4.8-fpm-alpine

WORKDIR "/data/www"

RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
&& apk update \
&& apk add build-base autoconf libmcrypt-dev bzip2-dev libxslt-dev freetype-dev libjpeg-turbo-dev libpng-dev zlib-dev libzip-dev\
&& docker-php-ext-install bcmath bz2 pdo_mysql mysqli opcache sockets pcntl xsl soap dom shmop zip \
&& docker-php-ext-configure gd --with-freetype=/usr/include/ --with-jpeg=/usr/include/ \
&& docker-php-ext-install gd \
&& pecl install redis-5.3.1 \
&& docker-php-ext-enable redis \
&& apk del build-base 

CMD [ "php-fpm" ]
```

推送到`Docker Hub`

```
docker push pengbotao/php:7.4.8-fpm-alpine
```

# 五、构建调试

通过`RUN`可以在容器内的系统执行一些命令，但因为并非在容器内直接执行，调试上可能会打一些折扣，所以可以先在容器中跑一跑命令看看，跑通在整理到`Dockerfile`中，如果`Dockerfile`构建时失败了，我们可以用`docker run`命令来基于这次构建到目前为止已经成功的最后创建的一个容器，比如下面示例：

**Dockerfile**

```
FROM php:7.4.8-fpm-alpine
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories 
RUN docker-php-ext-install redis
```

无法通过`docker-php-ext-install redis`来安装`redis`扩展，所以可以看到到第三步时报错了。

```
$ docker build -t test1 .
Sending build context to Docker daemon  9.216kB
Step 1/3 : FROM php:7.4.8-fpm-alpine
 ---> c8aada1d51a4
Step 2/3 : RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
 ---> Using cache
 ---> e101c0717b4f
Step 3/3 : RUN docker-php-ext-install redis
 ---> Running in 3455205e4c52
error: /usr/src/php/ext/redis does not exist

...

Some of the above modules are already compiled into PHP; please check
the output of "php -i" to see which modules are already loaded.
The command '/bin/sh -c docker-php-ext-install redis' returned a non-zero code: 1
```

这是我们可以进入第二步创建成功的容器，然后运行并进入该容器(我这里是alpine所以进入后终端用的是`/bin/sh`)，然后直接在容器内做进一步调试，运行成功后就可以调整`Dockerfile`重新尝试构建

```
$ docker run -it e101c0717b4f /bin/sh
/var/www/html # more /etc/apk/repositories
http://mirrors.aliyun.com/alpine/v3.12/main
http://mirrors.aliyun.com/alpine/v3.12/community
```

# 六、私有镜像仓库

## 6.1 通过registry镜像搭建
通过Docker的镜像很容易创建自己的镜像仓库。

1、安装。这里直接用官方的`registry`镜像来构建仓库，安装和启动方式和前面类似：

```
$ docker pull registry
$ docker run -d -v ~/docker/registry:/var/lib/registry -p 5000:5000 --restart=always --name registry registry:latest
d18931b1bd1b4b3ce7fe151c6f241c4dd1eebbfab8287dfed422ab5427a7f6a6
```

会启动`5000`端口，这里本机测试就不绑`host`了，也可以绑个`host`。

2、打TAG

```
$ docker tag pengbotao/php:7.4.8-fpm-alpine 127.0.0.1:5000/php:7.4.8-fpm-alpine
```

`docker tag`类似别名的操作，推送到远程仓库需要类似这样子的格式，`仓库/镜像名称:版本号`，当然仓库默认是官方仓库，版本号是`latest`。

3、推送镜像

```
$ docker push 127.0.0.1:5000/php:7.4.8-fpm-alpine
The push refers to repository [127.0.0.1:5000/php]
9c8d79537668: Pushed
106be351bd42: Pushed
7d85ada8d1c7: Pushed
767b1a2e8377: Pushed
9c55f034ee50: Pushed
6d3a6e58ccca: Pushed
fc5c686c158f: Pushed
d274b3a9c362: Pushed
81c338ff74a3: Pushed
308ef7bef157: Pushed
d94df04fea90: Pushed
50644c29ef5a: Pushed
7.4.8-fpm-alpine: digest: sha256:364931dd4bf7c5b159f93b458ab64bc3914218452c82d06d9359f82e3476c9bc size: 2830
```

4、查看镜像

```
# 查看镜像
$ curl http://127.0.0.1:5000/v2/_catalog
{"repositories":["php"]}

# 查看镜像tags列表
$ curl http://127.0.0.1:5000/v2/php/tags/list
{"name":"php","tags":["7.4.8-fpm-alpine"]}
```

也可以在`~/docker/registry/docker/registry/v2/repositories/`中查看仓库里的镜像信息。

5、镜像拉取(尝试删除后再从镜像仓库拉取)

```
$ docker rmi 127.0.0.1:5000/php:7.4.8-fpm-alpine
Untagged: 127.0.0.1:5000/php:7.4.8-fpm-alpine
Untagged: 127.0.0.1:5000/php@sha256:364931dd4bf7c5b159f93b458ab64bc3914218452c82d06d9359f82e3476c9bc

$ docker pull 127.0.0.1:5000/php:7.4.8-fpm-alpine
7.4.8-fpm-alpine: Pulling from php
Digest: sha256:364931dd4bf7c5b159f93b458ab64bc3914218452c82d06d9359f82e3476c9bc
Status: Downloaded newer image for 127.0.0.1:5000/php:7.4.8-fpm-alpine
127.0.0.1:5000/php:7.4.8-fpm-alpine
```

## 6.2 通过harbor搭建

下载地址：`https://github.com/goharbor/harbor/releases`


# 七、小结

借助`Dockerfile`和官方的基础镜像，基本可以编译出需要的环境，`Dockerfile`里的`RUN`命令和往常没有太大区别，但对比虚拟机或者`ECS`，好处是配置一次之后便可以以文本的方式存储起来或者将镜像推送到镜像仓库，轻量很多，后续配起来比较方便。

这一篇主要介绍了`Dockerfile`创建镜像，到这里我们就基本掌握了`Docker`的基本使用以及镜像的构建、推送。但像`PHP`的运行环境需要服务同时协作，按目前的理解就需要一个服务一个服务进行启动，启停上比较麻烦。而`docker-compose`基本可以解决开发环境下的这个问题，将各个容器打包启动。

下篇来看看`docker-compose`的用法。

---

- [1] [Dockerfile reference](https://docs.docker.com/engine/reference/builder/)
- [2] [Best practices for writing Dockerfiles](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [3] [Dockerfile 最佳实践及示例](http://www.dockerone.com/article/9551)
- [4] [About storage drivers](https://docs.docker.com/storage/storagedriver/)
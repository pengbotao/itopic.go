```
{
    "url": "docker-compose",
    "time": "2020/07/21 15:19",
    "tag": "Docker"
}
```

# 一、Docker Compose

## 1.1 概述


## 1.2 参数说明

# 二、参数详解

# 三、容器编排

容器如何编排？

一个机器多个环境
多个环境组合


# 四、示例


```
version: "3.0"

services:

    nginx:
      image: nginx:1.19.2-alpine
      container_name: nginx
      restart: always
      ports:
        - 80:80
      working_dir: /data/www
      volumes:
        - ~/workspace:/data/www
        - ./nginx/conf/nginx.conf:/etc/nginx/nginx.conf
        - ./nginx/conf.d:/etc/nginx/conf.d
        - ./nginx/logs:/var/log/nginx
      networks:
        - server

    php5.6:
      #build: ./php/dockerfile/5.6
      image: pengbotao/php:5.6-fpm-alpine
      container_name: 5.6-fpm
      restart: always
      working_dir: /data/www
      volumes:
        - ~/workspace:/data/www
        - ./php/5.6:/usr/local/etc
      networks:
        - server

    php7.4.8:
      #build: ./php/dockerfile/7.4.8
      image: pengbotao/php:7.4.8-fpm-alpine
      container_name: 7.4.8-fpm
      restart: always
      working_dir: /data/www
      volumes:
        - ~/workspace:/data/www
        - ./php/7.4.8:/usr/local/etc
      networks:
        - server

    mysql:
      image: mysql:5.7
      container_name: mysql5.7
      restart: always
      ports:
        - 3406:3306
      volumes:
        - ./mysql/conf:/etc/mysql
        - ./mysql/data:/var/lib/mysql
        - ./mysql/logs:/var/log/mysql
      environment:
        MYSQL_ROOT_PASSWORD: 123456
      networks:
        - server
    
    redis:
      image: redis:6.0.6-alpine
      container_name: redis
      restart: always
      ports:
        - 6479:6379
      volumes:
        - ./redis/conf:/usr/local/etc/redis
        - ./redis/data:/data
      command:
        redis-server /usr/local/etc/redis/redis.conf
      networks:
        - server

networks:
  server:
     driver: bridge
```

启动后效果如下：

![](../../static/uploads/docker-compose-server-php.png)

# 五、小结

本章侧重点在容器编排问题。

---
- [1] [Compose file version 3 reference](https://docs.docker.com/compose/compose-file/)
- [2] [docker-compose编排参数详解](https://www.cnblogs.com/wutao666/p/11332186.html)
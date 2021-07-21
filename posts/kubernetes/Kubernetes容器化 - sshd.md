```
{
    "url": "sshd-in-k8s",
    "time": "2021/08/01 21:10",
    "tag": "sshd,Kubernetes,容器化"
}
```

通常数据库只能在内网访问，当需要从外部访问时可以通过在有公网IP的节点上创建账号，从而在本地实现对数据库的访问。这里还是需要打个镜像：

**Dockerfile**

```
$ cat Dockerfile
FROM debian:buster

MAINTAINER pengbotao "pengbotao@vip.qq.com"
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list && sed -i s@/security.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list && apt-get update

RUN apt-get install -y locales
RUN sed -i 's/^# *\(zh_CN.UTF-8\)/\1/' /etc/locale.gen && locale-gen && echo "export LANG=zh_CN.UTF-8" >> /etc/bash.bashrc

RUN apt-get install -y ssh
RUN apt-get -y clean && rm -rf /var/lib/apt/lists/*

COPY entrypoint.sh /sbin/entrypoint.sh
RUN chmod 755 /sbin/entrypoint.sh

ENTRYPOINT ["/sbin/entrypoint.sh"]
```

**entrypint.sh**

```
$ cat entrypoint.sh
#!/bin/bash
set -e

if [[ ! -d "/run/sshd" ]]; then
    mkdir -p /run/sshd
fi

if [[ -z "$SSH_USER_FILE" ]]; then
    SSH_USER_FILE=/root/ssh_user_list
fi

if [[ -f "$SSH_USER_FILE" ]]; then
    user_list=$(cat $SSH_USER_FILE | awk -F ':' '{print $1}')
    for user in $user_list; do
        useradd -M -s /usr/sbin/nologin -N $user
    done
    chpasswd < $SSH_USER_FILE
fi

/usr/sbin/sshd -D
```

可以通过环境变量`SSH_USER_FILE`来定义账号文件，该账号只能用于SSH代理而不能直接登录节点（需要登录服务器可以通过下一篇中的JumpServer来实现）。文件格式：

```
userA:passwdA
userB:passwdB
```

镜像打好之后部署示例：

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ssh
  labels:
    app: ssh
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ssh
  template:
    metadata:
      labels:
        app: ssh
    spec:
      containers:
      - name: ssh
        image: pengbotao/sshd
        imagePullPolicy: IfNotPresent
        ports:
        - name: ssh
          containerPort: 22
        volumeMounts:
        - name: ssh-config
          mountPath: /root/ssh_user_list
          subPath: ssh_user_list
      volumes:
      - name: ssh-config
        configMap:
          name: ssh-config
```

通过ConfigMap来进行账号管理，添加账号之后需要重启服务。最后公网需要访问通过SLB作为入口，也可以绑定弹性IP并将pod调度到对应节点。
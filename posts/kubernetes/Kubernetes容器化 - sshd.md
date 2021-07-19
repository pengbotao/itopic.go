```
{
    "url": "sshd-in-k8s",
    "time": "2021/08/01 21:10",
    "tag": "sshd,Kubernetes,容器化",
    "public": "no"
}
```

通常数据库只能在内网访问，当需要从外部访问时可以通过在有公网IP的节点上创建账号，从而在本地实现对数据库的访问。这里还是需要打个镜像：

**Dockerfile**

```
$ cat Dockerfile
FROM centos:7

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo

RUN yum install -y openssh-server openssh-clients 

RUN ssh-keygen -t rsa -f /root/.ssh/id_rsa -N ""
RUN ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key -N ""
RUN ssh-keygen -t ecdsa -f /etc/ssh/ssh_host_ecdsa_key -N ""
RUN ssh-keygen -t ed25519 -f /etc/ssh/ssh_host_ed25519_key -N ""
RUN cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys
RUN sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/g' /etc/ssh/sshd_config

COPY entrypoint.sh /sbin/entrypoint.sh
RUN chmod 755 /sbin/entrypoint.sh

ENTRYPOINT ["/sbin/entrypoint.sh"]

```

**entrypint.sh**

```
$ cat entrypoint.sh
#!/bin/bash
set -e

if [[ -z "$SSH_USER_FILE" ]]; then
    SSH_USER_FILE=/root/ssh_user_list
fi

if [[ -f "$SSH_USER_FILE" ]]; then
    user_list=$(cat $SSH_USER_FILE | awk -F ':' '{print $1}')
    for user in $user_list; do
        useradd -M -s /sbin/nologin -n  $user
    done
    chpasswd < $SSH_USER_FILE
fi

/usr/sbin/sshd
/usr/bin/sleep infinity
```

由于启动时直接后台了，所以加了个sleep操作。可以通过环境变量`SSH_USER_FILE`来定义账号文件，该账号只能用于SSH代理而不能直接登录节点。文件格式：

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
        image: ssh:centos7-20210716
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

通过ConfigMap来进行账号管理，添加账号之后需要重启服务，暴露服务则是通过阿里云的SLB来进行暴露。
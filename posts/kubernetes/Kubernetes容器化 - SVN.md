```
{
    "url": "svn-in-k8s",
    "time": "2021/07/18 09:22",
    "tag": "SVN,Kubernetes,容器化"
}
```

首先需要做个镜像：

```
FROM debian:buster

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list && sed -i s@/security.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list && apt-get update

RUN apt-get install -y subversion
RUN apt-get -y clean && rm -rf /var/lib/apt/lists/*

COPY entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod 755 /usr/bin/entrypoint.sh

EXPOSE 3690

CMD ["/usr/bin/entrypoint.sh"]
```

其中`entrypoint.sh`：通过`SUBVERSION_REPOS`来配置需要启动的目录，如果目录不存在则创建。

```
#!/usr/bin/env bash

if [[ -z "${SUBVERSION_REPOS}" ]]; then
    SUBVERSION_REPOS=/var/svn/repos
fi
if [[ ! -d ${SUBVERSION_REPOS} ]]; then
    mkdir -p ${SUBVERSION_REPOS}
    /usr/bin/svnadmin create ${SUBVERSION_REPOS}
fi

/usr/bin/svnserve --daemon --foreground --root=${SUBVERSION_REPOS}
```

镜像打好之后的部署就是常规流程，创建pv/service等，这里以Deployment的方式来部署SVN，示例如下：

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: svn
  labels:
    app: svn
spec:
  replicas: 1
  selector:
    matchLabels:
      app: svn
  template:
    metadata:
      labels:
        app: svn
    spec:
      containers:
      - name: svn
        image: svn:1.10-20210713
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 3690
          name: svn
        env:
        - name: SUBVERSION_REPOS
          value: /data
        volumeMounts:
        - name: svn-pv
          mountPath: /data
      volumes:
      - name: svn-pv
        persistentVolumeClaim:
          claimName: svn-pvc
```

> 注：需要单独创建pv进行svn目录的持久化处理。
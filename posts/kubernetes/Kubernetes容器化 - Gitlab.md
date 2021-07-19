```
{
    "url": "gitlab-in-k8s",
    "time": "2021/07/24 20:10",
    "tag": "Gitlab,Kubernetes,容器化"
}
```

Gitlab的镜像（`docker pull gitlab/gitlab-ce:latest`）里已经包含了各种依赖环境，如果不考虑拆分postgresql和redis搭建还是挺简单，而且Gitlab上的备份和恢复目前用起来都还挺顺利，所以这里是直接使用的镜像里的数据库和缓存。

首先，第一步还是准备好pv/pvc存储。

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: gitlab-pv
  labels:
    project: gitlab-pv
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: gitlab-pv
    volumeAttributes:
      server: "*.cn-hangzhou.nas.aliyuncs.com"
      path: "/gitlab"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: gitlab-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 20Gi
  selector:
    matchLabels:
      project: gitlab-pv
```

第二步，创建需要的Service

```
apiVersion: v1
kind: Service
metadata:
  name: gitlab
spec:
  selector:
    project: gitlab
  ports:
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  - name: ssh
    port: 22
    targetPort: 22
    protocol: TCP
  clusterIP: None
```

第三步，创建StatefulSet，这里将配置文件也进行了挂载，就没有配置ConfigMap。

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: gitlab
  labels:
    project: gitlab
spec:
  serviceName: gitlab
  replicas: 1
  selector:
    matchLabels:
      project: gitlab
  template:
    metadata:
      labels:
        project: gitlab
    spec:
      containers:
      - name: gitlab
        image: gitlab/gitlab-ce:9.1.3-ce.0
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          containerPort: 80
        - name: ssh
          containerPort: 22
        env:
        - name: TZ
          value: Asia/Shanghai
        - name: GITLAB_ROOT_PASSWORD
          value: password
        - name: GITLAB_ROOT_EMAIL
          value: mail@gitlab.com
        volumeMounts:
        - name: gitlab-pv
          mountPath: /var/opt/gitlab
          subPath: gitlab
        - name: gitlab-pv
          mountPath: /etc/gitlab
          subPath: etc
      volumes:
      - name: gitlab-pv
        persistentVolumeClaim:
          claimName: gitlab-pvc
```

通过前面几步创建之后服务就创建好了，但上面使用到的是内网Service，由于需要暴露多个端口，这里直接通过阿里云同一个SLB来暴露2个默认端口，这样在拉取仓库的时候都不需要带额外的端口。

```
apiVersion: v1
kind: Service
metadata:
  name: gitlab-slb
  annotations:
    service.beta.kubernetes.io/alibaba-cloud-loadbalancer-id: slbid
    service.beta.kubernetes.io/alicloud-loadbalancer-force-override-listeners: 'true'
spec:
  selector:
    app: gitlab
  ports:
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  - name: ssh
    port: 22
    targetPort: 22
    protocol: TCP
  type: LoadBalancer
```

上面配置之后Gitlab仓库地址还不是以域名的方式在呈现，这时需要调整一下gitlab.rb配置参数：`external_url`

```
$ cat /etc/gitlab/gitlab.rb
prometheus['enable'] = false
prometheus_monitoring['enable'] = false
external_url 'http://gitlab.x.com'
unicorn['worker_processes'] = 2
```


```
{
    "url": "k8s-example",
    "time": "2020/11/01 17:15",
    "tag": "Kubernetes,容器化",
    "toc": "yes"
}
```

# 一、概述

前面章节中都会有一些`Demo`，但不够整体，这里从运维角度看看该如何配置日常服务，后面想找各个资源对象的`Yaml`文件示例写法，看这里应该就够了。

> k8s环境： Mac下Docker Desktop启用Kubernetes

#  二、PHP + Nginx

## 2.1 环境说明

`PHP`、`Python`、`Go`这几种语言中，`PHP`的部署算是最麻烦的了，他需要依赖`Nginx`，`PHP`和`Nginx`之间还需要文件共享，静态页面由`Nginx`处理，`PHP`页面交给`php-fpm`解析，所以要配置`PHP+Nginx`需要先理一理`PHP`和`Nginx`的交互方式，大概有两种方式可供选择：

- `PHP` 和 `Nginx`在同一个`Pod`中
- `PHP` 和 `Nginx`属于不同的`Pod`，文件通过`Volumes`挂载到同一个目录实现共享

**这里选择在同一个Pod中，即同一个Pod中的多个容器**，`PHP`和`Nginx`容器都需要能够读取到源代码文件，同一个`Pod`中挂载的目录各个容器都可以读到，我们可以直接挂个空目录，应用镜像只打代码文件，然后在`Pod`的`initController`容器里将代码都拷贝到容器去。

另外，常规项目配置上的要求：

- 位置文件由运维管理。通常`database.php`、`config.php`、日志级别等。
- `Nginx`和`PHP`的配置文件由运维管理。

如果配置文件打到镜像中，则修改后需要重新构建镜像，如果通过`ConfigMap`管理配置文件，则需要将配置在运行时挂载到容器中。**这里选择通过ConfigMap来控制配置文件**。

还有就是日志文件的问题，我们先通过`hostPath`的方式实现挂载`Nginx`日志。通过`ingress`实现7层代理。数据库这个场景我们先暂时不配置，可以使用本机的`mysql`。

上面就是配置`PHP`环境的需求，接下来看看怎么配置：

## 2.2 配置镜像

我们会使用到3个镜像，分别是`PHP`镜像、`Nginx`镜像以及代码镜像。

- `PHP`镜像：选择我们[前面](dockerfile.html)创建好的`pengbotao/php:7.4.8-fpm-alpine`
- `Nginx`镜像：我们选择`nginx:1.19.2-alpine`
- 代码镜像：我们选择`busybox:1.32.0`为基础镜像，负责将源文件打入镜像。

还有一种说法是代码文件不进镜像，直接通过文件挂载方式共享文件，但这是否意味着滚动更新的作用就削弱了，只需要更新共享的代码文件即可，而此过程有可能引起服务的短暂不可用。

我们模拟一个简单的项目，包含以下3个文件：

```
$ ls
Dockerfile api.php config.php index.php

# index.php
<?php phpinfo();

# api.php
<?php
include "config.php";
echo json_encode($config);

# config.php
<?php
$config = [
    "host" => "127.0.0.1",
    "env" => "uat"
];
```

为简单起见，只设置了这么3个文件，`config.php`配置文件需要通过`ConfigMap`注入。接下来创建`Dockerfile`，只需要将源代码拷贝到容器中即可。

```
FROM busybox:1.32.0

WORKDIR /src
COPY . /src
```

然后在目录中创建镜像:

```
$ docker build -t pengbotao/project-php:v1 .
```

这样子一个简单的镜像就创建好了，代码镜像里只有纯代码，无法直接运行应用。没有设置`.dockerignore`，配置文件`config.php`也写入到镜像中了，后面我们会用线上配置文件覆盖掉，也可以打包的时候就忽略掉。

## 2.3 创建ConfigMap

通过正式的配置文件创建`config.php`

```
$ kubectl create configmap phpdemo-config --from-file=config.php
configmap/phpdemo-config created

$ kubectl describe cm phpdemo-config
Name:         phpdemo-config
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
config.php:
----
<?php
$config = [
    "host" => "0.0.0.0",
    "env" => "prod"
];

Events:  <none>
```

创建`Nginx`配置文件，我们也可以用同样的方法创建`Nginx.conf`，`PHP`的`php.ini`等`Web`服务的配置文件。

```
$ kubectl create configmap phpdemo-nginx --from-file=phpdemo.local.com.conf
configmap/phpdemo-nginx created

$ kubectl describe cm phpdemo-nginx
Name:         phpdemo-nginx
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
phpdemo.local.com.conf:
----
server {
    listen       80;
    listen  [::]:80;
    server_name  phpdemo.local.com;
    index index.html index.php;
    root /data/www;
    charset utf-8;

    location ~ \.php$ {
        fastcgi_pass   127.0.0.1:9000;
        fastcgi_index  index.php;
        fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
        include        fastcgi_params;
    }
}
Events:  <none>
```

同一个`Pod`内，所以`php`可以直接设置为本地`9000`端口。域名暂定为：`phpdemo.local.com`

## 2.4 创建 Deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: phpdemo
  labels:
    project: phpdemo
    env: prod
spec:
  replicas: 2
  selector:
    matchLabels:
      project: phpdemo
      env: prod
  template:
    metadata:
      labels:
        project: phpdemo
        env: prod
    spec:
      initContainers:
      - name: init-phpdemo-src
        image: pengbotao/project-php:v2
        imagePullPolicy: IfNotPresent
        command: ['sh', '-c', "cp -rf /src/* /src-www && cp /src-config/* /src-www/ "]
        volumeMounts:
        - name: wwwroot
          mountPath: /src-www
        - name: phpdemo-config
          mountPath: /src-config
      containers:
      - name: php
        image: pengbotao/php:7.4.8-fpm-alpine
        imagePullPolicy: IfNotPresent
        resources:
          limits:
            memory: "64Mi"
            cpu: "250m"
          requests:
            memory: "64Mi"
            cpu: "250m"
        volumeMounts:
        - name: wwwroot
          mountPath: /data/www
        ports:
        - containerPort: 9000
      - name: nginx
        image: nginx:1.19.2-alpine
        imagePullPolicy: IfNotPresent
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
          requests:
            memory: "128Mi"
            cpu: "500m"
        volumeMounts:
        - name: wwwroot
          mountPath: /data/www
        - name: phpdemo-nginx
          mountPath: /etc/nginx/conf.d
        - name: nginx-log-path
          mountPath: /var/log/nginx
        ports:
        - containerPort: 80
      volumes:
      - name: wwwroot
        emptyDir: {}
      - name: phpdemo-config
        configMap:
          name: phpdemo-config
      - name: phpdemo-nginx
        configMap:
          name: phpdemo-nginx
      - name: nginx-log-path
        hostPath: 
          path: /Users/peng/k8s/logs
```

说明：

- 配置了2个副本
- 通过`initController`将文件拷贝到`/data/www`，如果php-fpm非root执行需要考虑修改对应权限：`&& chown -R www-data:www-data /src-www`，否则php-fpm可能没有写入权限。
- `PHP`设定了最低内存为64M，`CPU`为0.25
- `Nginx`配置文件通过`ConfigMap`挂载
- `Nginx`日志文件通过`hostPath`挂载本机目录

执行之后我们可以进容器看看代码文件是否正常，如果执行正常容器里应该可以看到源代码和线上的`config.php`。

## 2.5 创建 Service

```
apiVersion: v1
kind: Service
metadata:
  name: phpdemo-svc
  labels:
    project: phpdemo
    env: prod
spec:
  selector:
    project: phpdemo
    env: prod
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
  clusterIP: None
```

创建成功后`kubectl describe svc phpdemo-svc`应该可以看到`Endpoints`已经关联上了`Pod`。

## 2.6 创建Ingress

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: phpdemo.local.com
spec:
  rules:
  - host: phpdemo.local.com
    http:
      paths:
      - path: /
        backend:
          serviceName: phpdemo-nginx-svc
          servicePort: 80
```

通过下面命令可以看到，当前`ingress`暴露的是宿主机`80`端口，但`80`已经使用了，把`ingress-nginx`绑定的端口调整为`30080`

```
$ kubectl get svc -n ingress-nginx
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.109.107.221   localhost     80:31526/TCP,443:30328/TCP   7d9h

$ kubectl edit svc ingress-nginx-controller -n ingress-nginx

  ports:
  - name: http
    nodePort: 31526
    port: 30080
    protocol: TCP
    targetPort: http
    
    
$ kubectl get svc -n ingress-nginx
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                         AGE
ingress-nginx-controller             LoadBalancer   10.109.107.221   localhost     30080:31526/TCP,443:30328/TCP   7d9h
```

然后在宿主机`hosts`绑定`127.0.0.1 phpdemo.local.com`后访问 `http://phpdemo.local.com:30080/api.php` 就可以看到输出了，输出的是我们线上配置的`config.php`。

```
{
    "host": "0.0.0.0",
    "env": "prod"
}
```

本地挂载的日志目录也可以看到`Nginx`日志，到这里配置就基本完成了，接下来就是跟后期日常维护相关的操作。

## 2.7 版本更新

调整`index.php`文件内容模拟更新版本

```
<?php print_r($_SERVER);
```

重新打镜像

```
$ docker build -t pengbotao/project-php:v2 .
```

更新镜像和回滚只需要指定镜像版本即可。

```
kubectl set image deployment phpdemo init-phpdemo-src=pengbotao/project-php:v2
```

如果前面我们通过不同的`Pod`来组合`PHP`环境，`Nginx`和`PHP`里都有代码文件，镜像更新则需要执行2个`Pod`的更新：

```
$ kubectl set image deployment phpdemo init-phpdemo-src=pengbotao/project-php:v2
$ kubectl set image deployment phpdemo-nginx init-phpdemo-src=pengbotao/project-php:v2
```

## 2.8 增加/修改配置文件

前面我们是直接通过`kubectl create configmap`命令来创建，如果要增加文件则相对麻烦，我们可以调整为通过`Yaml`文件来创建，

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: phpdemo-config
  namespace: default
data:
  config.php: |
    <?php
    $config = [
        "host" => "0.0.0.0",
        "env" => "prod"
    ];
  database.php: |
    <?php
    $database = [
        "host" => "127.0.0.1",
        "port" => 3306,
    ];
```

这样子就增加了`database.php`的`Key`，更新`Deployment`后就会看到源代码目录增加了`database.php`文件。

## 2.9 重启服务

比如像上面场景更新了配置文件想重启`Pod`，或者某些情况下尝试重启`Pod`。如果`Deployment`没变更的话，重新`kubectl apply`不会触发滚动更新。手动删除`Pod`会重建，但一个个去删除也太累了。我们可以这么操作：

```
$ kubectl rollout restart deploy phpdemo
```

`kubectl rollout`包含以下功能：

```
$ kubectl rollout -h
Manage the rollout of a resource.

 Valid resource types include:

  *  deployments
  *  daemonsets
  *  statefulsets

Examples:
  # Rollback to the previous deployment
  kubectl rollout undo deployment/abc

  # Check the rollout status of a daemonset
  kubectl rollout status daemonset/foo

Available Commands:
  history     显示 rollout 历史
  pause       标记提供的 resource 为中止状态
  restart     Restart a resource
  resume      继续一个停止的 resource
  status      显示 rollout 的状态
  undo        撤销上一次的 rollout
```

可以通过`undo`做回滚操作，比如回退到前一版本：

```
# 设置为v1版本
$ kubectl set image deployment phpdemo init-phpdemo-src=pengbotao/project-php:v1
# 升级为v2版本
$ kubectl set image deployment phpdemo init-phpdemo-src=pengbotao/project-php:v2
# 回滚到前一版本，即v1版本
$ kubectl rollout undo deploy phpdemo
```

也可以指定回滚的版本：`kubectl rollout undo deploy phpdemo --to-revision=1`，可以通过查看`rollout`查看历史记录：

```
$ kubectl rollout history deploy phpdemo
deployment.apps/phpdemo
REVISION  CHANGE-CAUSE
1         <none>
3         <none>
4         <none>
5         <none>
8         <none>
11        <none>
12        <none>
```

但这个记录前面看过，基本看不太出差别，所以感觉直接更新镜像版本或者回退到上一版本会更实用些（也有可能是没找到`CHANGE-CAUSE`列的用法）。

## 2.10 小结

这个环境里实现了：

- 代码镜像只有纯代码，不具备运行环境。
- `PHP`和`Nginx`部署在同一`Pod`中，容器之间实现代码文件共享
- 通过`Deploy`可以实现`Pod`异常自我修复以及滚动更新
- 配置`Ingress`实现7层负载均衡
- 模拟日常操作版本更新以及服务重启

`Pod`层级还有就绪检测、存活检测可以做一做，接下来在`Python`的环境中加上这两项看看。

# 三、Python + Nginx

## 3.1 环境说明

这里主要还是出于演示目的，尽量体现出每个`Demo`的差异化，`Python`环境这边想这么做：

- 通过`StatefulSet`来配置服务（实际环境中可能和上面`PHP`类似属于无状态服务）
- `Pod`中的`Nginx`与`Python`容器隔离，很多Gunicorn应用会在上面加一层Nginx，这里就是按照这种来配置，演示不同`Pod`之间的用法
- 实现就绪检测和存活检测

## 3.2 配置镜像

一个简单的`Gunicorn`+`Flask`应用，`Github`地址：`https://github.com/pengbotao/k8s-py-demo`。也可以通过`docker pull pengbotao/k8s-py-demo:v1`拉取。源代码已经打到`Python`镜像中了，默认启动是5000端口，拿到镜像就可以运行起来。接下来来部署这个镜像：

- 配置文件，用我们管理的配置文件覆盖掉系统里的`gunicorn.py`文件。

## 3.3 创建ConfigMap

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: flask-config
  namespace: default
data:
  gunicorn.py: |
    import gevent.monkey

    gevent.monkey.patch_all()
    debug = True
    loglevel = 'debug'
    bind = '0.0.0.0:80'

    workers = 1
    threads = 2
    worker_class = 'gunicorn.workers.ggevent.GeventWorker'
    daemon = False

    pidfile = './logs/gunicorn.pid'
    logfile = './logs/debug.log'
    accesslog = './logs/gunicorn_access.log'
    errorlog = './logs/gunicorn_error.log'

    x_forwarded_for_header = 'X-FORWARDED-FOR'

---

apiVersion: v1
kind: ConfigMap
metadata:
  name: flask-nginx
  namespace: default
data:
  pydemo.local.com.conf: |
    server {
        listen       80;
        listen  [::]:80;
        server_name  pydemo.local.com;

      location / {
            proxy_pass http://flask-svc;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
    }

```

容器里默认使用的`5000`端口，这里为做演示将容器端口调整为80，配置站点。`Nginx`中`proxy_pass`关联的是`flask-svc`。

## 3.4 创建PersistentVolume

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: flask-pv001
spec:
  capacity:
    storage: 2Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: flask-pv
  hostPath:
    path: /Users/peng/k8s/pv-data/pv001
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: flask-pv002
spec:
  capacity:
    storage: 2Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: flask-pv
  hostPath:
    path: /Users/peng/k8s/pv-data/pv002
```

创建了`2`个PV，用来存储项目产生的日志文件。

## 3.5 创建StatefulSet

```
apiVersion: v1
kind: Service
metadata:
  name: flask-svc
spec:
  selector:
    app: flask-pod
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
  clusterIP: None


---

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: flask-sts
spec:
  replicas: 2
  selector:
    matchLabels:
      app: flask-pod
  serviceName: flask-svc
  template:
    metadata:
      labels:
        app: flask-pod
    spec:
      containers:
      - name: flask
        image: pengbotao/k8s-py-demo:v1
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
        volumeMounts:
        - name: flask-pvc
          mountPath: /data/www/logs
        - name: flask-config
          mountPath: /data/www/gunicorn.py
          subPath: gunicorn.py
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 60
          periodSeconds: 10
      volumes:
      - name: flask-config
        configMap:
          name: flask-config
          items:
          - key: gunicorn.py
            path: gunicorn.py
  volumeClaimTemplates:
  - metadata:
      name: flask-pvc
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: flask-pv
      resources:
        requests:
          storage: 1Gi
```

配置`Pod`，增加就绪检测、存活检测。同时将配置文件覆盖已经存在的`gunicorn.py`。

## 3.6 创建Nginx

`Nginx`并不是必须的，`Ingress`可以直接关联上面的`SVC`：`flask-svc`，这配之后倒产生了一个新的问题，后面再说。

```
apiVersion: v1
kind: Service
metadata:
  name: flask-nginx-svc
spec:
  selector:
    app: flask-nginx-pod
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
  clusterIP: None

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: flask-nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: flask-nginx-pod
  template:
    metadata:
      labels:
        app: flask-nginx-pod
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.2-alpine
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
        volumeMounts:
        - name: flask-nginx
          mountPath: /etc/nginx/conf.d
      volumes:
      - name: flask-nginx
        configMap:
          name: flask-nginx
```

## 3.7 创建Ingress

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: pydemo.local.com
spec:
  rules:
  - host: pydemo.local.com
    http:
      paths:
      - path: /
        backend:
          serviceName: flask-svc
          servicePort: 80
```

创建`Ingress`，配置`Host`后就可以访问到了：

```
$ curl http://pydemo.local.com:30080/
{
    "ClientIP": "192.168.65.3, 10.1.2.196",
    "Host": "flask-sts-1",
    "ServerIP": "10.1.3.88",
    "Time": "2020-09-30 02:31:01",
    "Version": "v1"
}
```

这里的访问流程：`Ingress -> Nginx Depolyment -> Flask StatefulSet -> Pod`。正常访问是没有问题，但当`sts`做更新的时候会存在问题，更新流程：

- 触发`sts`更新。这个时候会停掉一个，比如停掉`flask-sts-1`
- 由于对应的`Pod`停了，`flask-svc`关联的`Endpoints`会去掉`flask-sts-1`的`PodIP`，也就是只有`1`个后端IP了
- 上面配置的`Nginx`关联的是`flask-svc`，但`Nginx`并不会更新解析地址
- 没有更新解析，原`IP`失效，再访问站点就失败了，需要重启`Nginx`才行

但如果不创建`Nginx Depolyment`，直接用`Ingress`关联`flask-svc`，则不存在这个问题，重启`Nginx`可以解决。问题产生的原因是`Nginx`当`proxy_pass`为域名时，会做`DNS`缓存，因为容器的`IP`变化了，通过旧的`IP`就会访问失败了。

**第一种解决方法**：

修改`ConfigMap`中的`nginx`增加`resolver`配置，指定`DNS`（通过`kubectl get svc -n kube-system`可以看到`IP`）

```
    server {
        listen       80;
        listen  [::]:80;
        server_name  pydemo.local.com;
        resolver 10.96.0.10 valid=3s;
        resolver_timeout 5s;
        set $upstream flask-svc.default.svc.cluster.local;
        location / {
            proxy_pass http://$upstream;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
    }
```

> 注：这里配置的地址是headless service的地址，用svc的名称没走通

**第二种解决方法**：

由于`StatefulSet`需要用到`Headless Service`，`Service`的`IP`变化感知很快，可以在建立一个`Service`设置相同的`Label/Selector`，`Nginx`里在设置`VIP`的地址即可。

两种方法感觉都一般，这个应用可以直接走`Ingress -> Flask StatefulSet`，从而去掉`Nginx`。

## 3.8 小结

这个环境更多的是为了演示，区分`Deploy`和`StatefulSet`，了解下`STS`的使用。

---



- [1] [Nginx doesn't work after restarting upstream pod without restarting nginx itself.](https://github.com/ubuntu/microk8s/issues/224)
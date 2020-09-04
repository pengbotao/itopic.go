```
{
    "url": "k8s-controller",
    "time": "2020/09/12 12:08",
    "tag": "Kubernetes,容器化",
    "public": "no"
}
```

# 一、概述

我们在前面概念章节提到了控制器以及常用控制器的职责。所谓可控制器就是用来控制容器的状态和行为，不过控制器并不都是同级别的，比如`Deployment`管理`ReplicaSet`，`HPA`可以对`RS`、`Deployment`设置扩容方案，但终归目的是用来控制容器。

| 编号 | 控制器                                | 说明                                                         | 应用场景   |
| ---- | ------------------------------------- | ------------------------------------------------------------ | ---------- |
| 1    | Deployment                            | 部署无状态应用                                               | Web应用    |
| 2    | StatefulSet                           | 部署有状态应用                                               | 数据库     |
| 3    | DaemonSet                             | 在每一个Node上面运行一个Pod；新加入的Node也同样会自动运行一个Pod | Agent      |
| 4    | Job/CronJob                           | 一次性任务/周期任务                                          | 脚本、备份 |
| 5    | ReplicaSet<br />ReplicationController | 控制Pod的副本数量                                            |            |
| 6    | Horizontal Pod Autoscaling            | Pod水平自动缩放                                              |            |

`k8s`能创建的最小粒度资源是`Pod`，`Pod`可以分为两类，一类是直接创建的`Pod`，没有管理者，退出后不会有重新拉起操作。另一类就是控制器控制的`Pod`，可以根据控制器的不同规则对`Pod`做不同的行为控制。

概念在前面章节已经提到，所以本章节主要从示例维度来对`Pod`和控制器进行说明。

# 二、 Pod

## 2.1 Pod示例

```
apiVersion: v1
kind: Pod
metadata:
  name: mypod
  labels:
    name: mypod
spec:
  containers:
  - name: c1
    image: busybox
    command:
      - "/bin/sh"
      - "-c"
      - "top"
    resources:
      limits:
        memory: "128Mi"
        cpu: "500m"
  - name: c2
    image: busybox
    command:
      - "/bin/sh"
      - "-c"
      - "while true ; sleep 1; do date ; done"
    resources:
      limits:
        memory: "128Mi"
        cpu: "500m"
  - name: c3
    image: busybox
    command:
      - "/bin/sh"
      - "-c"
      - "date"
    resources:
      limits:
        memory: "128Mi"
        cpu: "500m"
```

说明：

- 定义了一个`Pod`，`Pod`的名字是`mypod`，打了一个标签也是`mypod`
- `Pod`里定义了3个容器，容器都从`busybox`拉取镜像，但没有指定`busybox`的版本，默认`latest`版本
- 第一个容器执行`Top`命令，第二个容器每`1s`打印时间，第三个容器只打印了当前时间

来看看执行情况：

```
# 查看Pod
$ kubectl get pod
NAME                           READY   STATUS    RESTARTS   AGE
mypod                          2/3     Running   3          74s

# 查看c2容器日志，-f参数类似 tail -f。
$ kubectl logs mypod c2 -f
Fri Sep  4 03:46:11 UTC 2020
Fri Sep  4 03:46:12 UTC 2020
Fri Sep  4 03:46:13 UTC 2020
Fri Sep  4 03:46:14 UTC 2020
Fri Sep  4 03:46:15 UTC 2020
Fri Sep  4 03:46:16 UTC 2020
Fri Sep  4 03:46:17 UTC 2020
Fri Sep  4 03:46:18 UTC 2020
Fri Sep  4 03:46:19 UTC 2020
```

容器要求程序在前台执行，执行完了容器就退出了。因为前两个容器都会一直执行，第三个容器执行完就退出了。通过`describe`也可以看到。

```
$ kubectl describe pod mypod
Events:
  Type    Reason     Age        From                     Message
  ----    ------     ----       ----                     -------
  Normal  Scheduled  <unknown>  default-scheduler        Successfully assigned default/mypod to docker-desktop
  Normal  Pulling    7s         kubelet, docker-desktop  Pulling image "busybox"
  Normal  Pulled     6s         kubelet, docker-desktop  Successfully pulled image "busybox"
  Normal  Created    6s         kubelet, docker-desktop  Created container c1
  Normal  Started    6s         kubelet, docker-desktop  Started container c1
  Normal  Pulling    6s         kubelet, docker-desktop  Pulling image "busybox"
  Normal  Pulled     3s         kubelet, docker-desktop  Successfully pulled image "busybox"
  Normal  Created    3s         kubelet, docker-desktop  Created container c2
  Normal  Started    3s         kubelet, docker-desktop  Started container c2
  Normal  Pulling    3s         kubelet, docker-desktop  Pulling image "busybox"
  Normal  Pulled     1s         kubelet, docker-desktop  Successfully pulled image "busybox"
  Normal  Created    1s         kubelet, docker-desktop  Created container c3
  Normal  Started    1s         kubelet, docker-desktop  Started container c3
  Warning  BackOff    15s (x4 over 48s)  kubelet, docker-desktop  Back-off restarting failed container
  Normal   Pulling    3m3s (x4 over 3m57s)   kubelet, docker-desktop  Pulling image "busybox
  Warning  BackOff    83s (x23 over 6m19s)   kubelet, docker-desktop  Back-off restarting failed container

$ kubectl get pod mypod -o wide
NAME    READY   STATUS             RESTARTS   AGE   IP          NODE             NOMINATED NODE   READINESS GATES
mypod   2/3     CrashLoopBackOff   7          14m   10.1.2.46   docker-desktop   <none>           <none>
```

- 由于是`latest`版本，默认每次创建容器都会有`Pulling image "busybox"`操作，线上应避免直接拉`latest`镜像，如果要用重新打个`tag`也行。
- `c3`容器执行后就退出了，所以`Pod`会尝试重启容器，但会一直退出，并不是执行失败了，而是执行完退出了。
- 重新查看容器状态为`CrashLoopBackOff`，`READY`状态为`2/3`，`RESTARTS`重启次数为7，重启的粒度是容器层级，从日志上看`c1`、 `c2`容器的重启次数没有增多。因为重启策略默认为`Always`，所以会一直重启下去。

**重启策略：**

```
$ kubectl explain pod.spec.restartPolicy
KIND:     Pod
VERSION:  v1

FIELD:    restartPolicy <string>

DESCRIPTION:
     Restart policy for all containers within the pod. One of Always, OnFailure, Never. Default to Always.
```

看描述

- 重启策略针对`Pod`里的所有容器
- 值按字面意思理解就行
  - `Always`：容器退出时总是重启容器，默认策略。
  - `OnFailure`：容器异常退出时重启（状态码非0），如果增加`spec.restartPolicy = OnFailure`，则c3容器执行后不会在重启。
  - `Never`：容器退出时不重启容器



## 2.2 Pod生命周期

![](../../static/uploads/k8s-pod-lifecycle.png)

### 2.2.1 Pause容器

`Pause容器`：每个`Pod`里运行着一个特殊的被称之为`Pause`的容器，其他容器则为业务容器，这些业务容器共享`Pause`容器的网络栈和`Volume`挂载卷，因此他们之间通信和数据交换更为高效。在设计时可以充分利用这一特性，将一组密切相关的服务进程放入同一个`Pod`中；同一个`Pod`里的容器之间仅需通过`localhost`就能互相通信。

**Pause容器提供以下功能：**

- PID命名空间：Pod中的不同应用程序可以看到其他应用程序的进程ID。

- 网络命名空间：Pod中的多个容器能够访问同一个IP和端口范围。

- IPC命名空间：Pod中的多个容器能够使用System V IPC或POSIX消息队列进行通信。

- UTS命名空间：Pod中的多个容器共享一个主机名；Volumes（共享存储卷）。
- Pod中的各个容器可以访问在Pod级别定义的Volumes。

### 2.2.2 InitContainer

在主容器（`Main Container`）启动之前执行的容器，串行执行，只有前一个`InitContainer`正常退出，下一个才会继续。如果`InitContainer`失败，则会根据策略重启`Pod`。相当于在主容器启动之前，可以通过`Init Conntainer`容器做一些准备工作。

### 2.2.3 声明周期钩子函数

`Kubernetes `为容器提供了两种生命周期钩子：

- `Poststart`：于容器创建完成之后立即运行的钩子程序。
- `preStop`：容器终止之前立即运行的程序，是以同步方式的进行，因此其完成之前会阻塞 删除容器的调用

备注：钩子程序的执行方式有`Exec`和`HTTP`两种。

### 2.2.4 探针 - readiness

就绪性探测:用于判定容器中的主进程是否准备就绪以及能否对外提供服务。

### 2.2.5 存活 - liveness

存活性探测:用于判定主容器是否处于存活状态。

## 2.3 完整示例

@todo

# 三、ReplicaSet

## 3.1 RS / RC区别

`ReplicaSet`/`ReplicationController` 确保系统中的Pod数量永远等于设置的个数。

在新版的`Kubernetes`中建议使用`ReplicaSet (RS)`来取代`ReplicationController(RC)`。`ReplicaSet`跟`ReplicationController`没有本质的不同，只是名字不一样，但`ReplicaSet`支持集合式`selector`，k8s里通过对资源对象打标签，然后可以按不同的规则来筛选这些标签。具体的用法在前一篇`Label / Selector`中有说到，到具体的差异可以从描述文档上看到：

```
$ kubectl explain rc.spec
KIND:     ReplicationController
VERSION:  v1

FIELDS:
   selector	<map[string]string>


$ kubectl explain rs.spec
KIND:     ReplicaSet
VERSION:  apps/v1

FIELDS:
   selector	<Object> -required-

$ kubectl explain rs.spec.selector
KIND:     ReplicaSet
VERSION:  apps/v1

FIELDS:
   matchExpressions	<[]Object>
     matchExpressions is a list of label selector requirements. The requirements
     are ANDed.

   matchLabels	<map[string]string>
     matchLabels is a map of {key,value} pairs. A single {key,value} in the
     matchLabels map is equivalent to an element of matchExpressions, whose key
     field is "key", the operator is "In", and the values array contains only
     "value". The requirements are ANDed.
```

对比它俩的文档，可以看到结果略有差异，`RC`的的`selector`是一个`map`，而`RS`的`selector`为一个`object`，必填，下面有两个`key`: `matchExpressions`、`matchLabels`。而`matchExpressions`支持的筛选规则似乎更灵活些，就是上面说到的集合式`selector`。

```
apiVersion: v1
kind: ReplicationController
metadata:
  name: myapp
spec:
  selector:
      app: myapp



apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: myapp
spec:
  selector:
    matchLabels:
      app: myapp
    matchExpressions；
      - {key: tier, operator: In, values: [frontend]}
      - {key: environment, operator: NotIn, values: [dev]}
```

## 3.2 ReplicaSet示例

```
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: rsname
spec:
  replicas: 1
  selector:
    matchLabels:
      name: rspod
  template:
    metadata:
      labels: 
        name: rspod
    spec:
      containers: 
      - name: nginx
        image: nginx:1.19.2-alpine
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
        ports:
        - containerPort: 80


$ kubectl apply -f rs-nginx.yaml
replicaset.apps/rsname created

$ kubectl scale rs rsname --replicas=2
```

可以看到`ReplicaSet`的定义和`Deployment`定义是一样的，而且也支持扩容，然后来看看`Deployment`。

# 四、Deployment

## 4.1 Deployment 与 RS

定义`Deployment`会创建`Pod`和`ReplicaSet`，创建关系大致是`Deployment`创建`ReplicaSet`，`ReplicaSet`创建`Pod`。文档上对它俩的描述：

- ReplicaSet ensures that a specified number of pod replicas are running at any given time.
- Deployment enables declarative updates for Pods and ReplicaSets.

`Deployment`不直接管理`Pod`，而是通过`ReplicaSet`来进行管理，他们的功能差不多，都支持自动扩容、缩容，但`Deployment`支持滚动更新和回滚，这个是`ReplicaSet`不支持的，所以一般建议是通过`Deployment`来管理`Pod`。

## 4.2 Deployment示例

我们先来创建`Nginx`的`Deployment`，可能是本地的网络原因，需要创建`service`且指定`type`为`LoadBalancer`本机端口才能访问到，先忽略，后续在探究这个问题。多个资源清单可以写在一个文件，通过`---`进行分割即可。

```
apiVersion: v1
kind: Service
metadata:
  name: nginx-svc
  labels:
    project: nginx
spec:
  selector:
    app: nginx-pod
  type: LoadBalancer
  ports:
  - port: 38000
    targetPort: 80
    protocol: TCP

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deploy
  labels:
    project: nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx-pod
  template:
    metadata:
      labels:
        app: nginx-pod
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.2-alpine
        imagePullPolicy: IfNotPresent
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
        ports:
        - containerPort: 80

```

通过apply创建，并且在最后增加了`--record`，用来记录历史版本号，可以回滚时指定某个版本。

> --record=false: Record current kubectl command in the resource annotation. If set to false, do not record the
> command. If set to true, record the command. If not set, default to updating the existing annotation value only if one already exists.

```
$ kubectl apply -f nginx.yaml --record
service/nginx-svc created
deployment.apps/nginx-deploy created
```

这个配置前面也配过好几次了，接下来本机就可以访问`http://localhost:38000/`看到`Nginx`的欢迎页。命令行也可以看到`deployment`、`rs`、`pod`信息。

```
$ kubectl get deploy
NAME           READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deploy   2/2     2            2           25s

$ kubectl get rs
NAME                     DESIRED   CURRENT   READY   AGE
nginx-deploy-9bdb559b9   2         2         2       27s

$ kubectl get pod
NAME                           READY   STATUS    RESTARTS   AGE
nginx-deploy-9bdb559b9-hlrjq   1/1     Running   0          28s
nginx-deploy-9bdb559b9-rw2gh   1/1     Running   0          28s
```

## 4.3 滚动更新

既然`Deployment`和`RS`的主要差别在于滚动更新，我们来看看滚动更新操作。

**4.3.1 版本更新**

尝试更新nginx镜像到`1.18.0-alpine`版本。更新方式通资源清单中所讲，可以通过更新yaml文件，或者通过`set image`命令操作就会触发rollout。

```
$ kubectl set image deployment nginx-deploy nginx=nginx:1.18.0-alpine
```

再次查看资源信息，可以看到增加了一个`rs`，`DESIRED`、`CURRENT`、`READY`对应为1、1、0，也新增了一个`pod`，状态是`ImagePullBackOff`。但这里为啥是期望的`DESIRED`数量1呢？，猜测这是一个中间状态，跟滚动更新流程有关系，rs先启动一个，启动成功后再扩容一个，但由于第一个现在卡主了，所以看到了这个状态。

```
$ kubectl get deploy
NAME           READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deploy   2/2     1            2           3m8s

$ kubectl get rs
NAME                      DESIRED   CURRENT   READY   AGE
nginx-deploy-6d57fb5cfd   1         1         0       2m9s
nginx-deploy-9bdb559b9    2         2         2       3m11s


$ kubectl get pod
NAME                            READY   STATUS             RESTARTS   AGE
nginx-deploy-6d57fb5cfd-ktw6w   0/1     ImagePullBackOff   0          2m14s
nginx-deploy-9bdb559b9-hlrjq    1/1     Running            0          3m16s
nginx-deploy-9bdb559b9-rw2gh    1/1     Running            0          3m16s
```

查看Pod信息，可以到镜像获取失败。

```
$ kubectl describe pod nginx-deploy-6d57fb5cfd-ktw6w

Events:
  Type     Reason     Age                  From                     Message
  ----     ------     ----                 ----                     -------
  Normal   Scheduled  <unknown>            default-scheduler        Successfully assigned default/nginx-deploy-6d57fb5cfd-ktw6w to docker-desktop
  Warning  Failed     2m58s                kubelet, docker-desktop  Failed to pull image "nginx:1.18.0-alpine": rpc error: code = Unknown desc = Get https://registry-1.docker.io/v2/library/nginx/manifests/1.18.0-alpine: net/http: TLS handshake timeout
  Warning  Failed     2m12s                kubelet, docker-desktop  Failed to pull image "nginx:1.18.0-alpine": rpc error: code = Unknown desc = Get https://registry-1.docker.io/v2/: net/http: TLS handshake timeout
  Warning  Failed     63s                  kubelet, docker-desktop  Failed to pull image "nginx:1.18.0-alpine": rpc error: code = Unknown desc = Get https://registry-1.docker.io/v2/library/nginx/manifests/sha256:8853c7e938c2aa5d9d7439e698f0e700f058df8414a83134a09fcbb68bb0707a: net/http: TLS handshake timeout
  Warning  Failed     63s (x3 over 2m58s)  kubelet, docker-desktop  Error: ErrImagePull
  Normal   BackOff    33s (x4 over 2m58s)  kubelet, docker-desktop  Back-off pulling image "nginx:1.18.0-alpine"
  Warning  Failed     33s (x4 over 2m58s)  kubelet, docker-desktop  Error: ImagePullBackOff
  Normal   Pulling    22s (x4 over 3m37s)  kubelet, docker-desktop  Pulling image "nginx:1.18.0-alpine"
```

查看版本记录，需要前面开启`--record`，可以通过`.spec.revisionHistoryLimit`指定`deployment`最多保留多少`revision`记录。貌似这个版本记录并不太清晰，不同版本之间看不出区别。

```
$ kubectl rollout history deployment nginx-deploy
deployment.apps/nginx-deploy
REVISION  CHANGE-CAUSE
1         kubectl apply --filename=nginx.yaml --record=true
2         kubectl apply --filename=nginx.yaml --record=true
```

**4.3.2 版本回滚**

看起来这一次更新镜像异常，正好可以操作下回滚：

```
$ kubectl rollout undo deployment nginx-deploy
deployment.apps/nginx-deploy rolled back
```

也可以指定某个历史版本

```
$ kubectl rollout undo deployment nginx-deploy --to-revision=1
```

再来查看资源信息，状态已经恢复，`Deployment`增加了另一个rs用来做滚动升级。

```
$ kubectl get deploy
NAME           READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deploy   2/2     2            2           7m19s
pengbotao:k8s peng$ kubectl get rs
NAME                      DESIRED   CURRENT   READY   AGE
nginx-deploy-6d57fb5cfd   0         0         0       6m19s
nginx-deploy-9bdb559b9    2         2         2       7m21s
pengbotao:k8s peng$ kubectl get pod
NAME                           READY   STATUS    RESTARTS   AGE
nginx-deploy-9bdb559b9-hlrjq   1/1     Running   0          7m26s
nginx-deploy-9bdb559b9-rw2gh   1/1     Running   0          7m26s
```

**4.3.3 版本更新**

我们尝试手动下载镜像，然后再次升级看看。刷`rs`可以看到中间数量的变更。

```
$ kubectl get rs
NAME                      DESIRED   CURRENT   READY   AGE
nginx-deploy-6d57fb5cfd   2         2         1       8h
nginx-deploy-9bdb559b9    1         1         1       8h

$ kubectl get rs
NAME                      DESIRED   CURRENT   READY   AGE
nginx-deploy-6d57fb5cfd   2         2         2       8h
nginx-deploy-9bdb559b9    0         0         0       8h
```

通过`describe`可以看到现在的版本以及日志信息

```
  $ kubectl describe deploy nginx-deploy
  ...
  Pod Template:
  Labels:  app=nginx-pod
  Containers:
   nginx:
    Image:      nginx:1.18.0-alpine

  ...
  
  Events:
  Type    Reason             Age                  From                   Message
  ----    ------             ----                 ----                   -------
  Normal  ScalingReplicaSet  63s (x3 over 8h)     deployment-controller  Scaled up replica set nginx-deploy-6d57fb5cfd to 1
  Normal  ScalingReplicaSet  62s (x2 over 4m23s)  deployment-controller  Scaled down replica set nginx-deploy-9bdb559b9 to 1
  Normal  ScalingReplicaSet  61s (x2 over 4m23s)  deployment-controller  Scaled up replica set nginx-deploy-6d57fb5cfd to 2
  Normal  ScalingReplicaSet  59s (x2 over 4m21s)  deployment-controller  Scaled down replica set nginx-deploy-9bdb559b9 to 0
```

通过Events中`ScalingReplicaSet`看到流程是：

- 升级的1.18版本，启用一个Pod，启动成功后
- 原版本1.19停用一个Pod
- 1.18在启动一个Pod
- 1.19版本Pod都停掉，完成升级。

**4.3.4 kubectl操作**

```
# 暂停
$ kubectl rollout pause deployment nginx-deploy

# 查看状态
$ kubectl rollout status deployment nginx-deploy
deployment "nginx-deploy" successfully rolled out

# 恢复
$ kubectl rollout resume deployment nginx-deploy

# 扩容
$ kubectl scale deployment nginx-deploy --replicas 3

# 删除
$ kubectl delete service,deploy -l project=nginx
```

##  4.4 资源清单

看看资源清单中比较重要的节点：

```
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx-pod
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  minReadySeconds: 5
  template:
```

文档说明：

```
$ kubectl explain deploy.spec.strategy
KIND:     Deployment
VERSION:  apps/v1

RESOURCE: strategy <Object>

DESCRIPTION:
     The deployment strategy to use to replace existing pods with new ones.

     DeploymentStrategy describes how to replace existing pods with new ones.

FIELDS:
   rollingUpdate	<Object>
     Rolling update config params. Present only if DeploymentStrategyType =
     RollingUpdate.

   type	<string>
     Type of deployment. Can be "Recreate" or "RollingUpdate". Default is
     RollingUpdate.
```

# 五、DaemonSet

`DaemonSet`控制有下面特征：

- 每个`Node`节点上都会运行一个`Pod`的副本
- 有新节点加入时，会自动创建`Pod`，节点移除时`Pod`也会被移除
- 删除`DaemonSet`，那就创建的所有`Pod`就都会被删除

所以使用场景上一般为日志收集、监控等。查看系系统中已有的`DaemonSet`可以看到，`kube-proxy`就是通过`DS`在配置。

```
$ kubectl get ds -n kube-system
NAME         DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR                 AGE
kube-proxy   1         1         1       1            1           beta.kubernetes.io/os=linux   6d22h
```

## 5.1 DaemonSet示例



# 六、Job

负责处理一次性任务。

## 6.1 Job示例

```
apiVersion: batch/v1
kind: Job
metadata:
  name: job-demo
spec:
  template:
    spec:
      restartPolicy: OnFailure
      containers:
      - name: box
        image: busybox
        command:
        - "/bin/sh"
        - "-c"
        - "date"
```

创建后查看

```
$ kubectl apply -f job.yaml
job.batch/job-demo created

$ kubectl get job
NAME       COMPLETIONS   DURATION   AGE
job-demo   1/1           3s         11s

$ kubectl get pod
NAME                           READY   STATUS      RESTARTS   AGE
job-demo-8rl4c                 0/1     Completed   0          20s

$ kubectl logs job-demo-8rl4c
Thu Sep  3 07:19:11 UTC 2020
```

## 6.2 资源清单

- `spec.template`格式同`Pod`

# 七、 CronJob

`CronJob`就是在`Job`的基础上变成周期性的任务，可以周期性执行。周期设置的`schedule`和`crontab`一样。在后续可以看到，执行`CronJob`会产生`Job`。

## 7.1 CronJob示例

```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: cronjob-demo
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
          - name: box
            image: busybox
            command:
            - "/bin/sh"
            - "-c"
            - "date"
```

创建后查看

```
$ kubectl apply -f cronjob.yaml
cronjob.batch/cronjob-demo created

$ kubectl get cronjob
NAME           SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
cronjob-demo   */1 * * * *   False     0        <none>          16s

$ kubectl get job
NAME                      COMPLETIONS   DURATION   AGE
cronjob-demo-1599117660   1/1           3s         76s
cronjob-demo-1599117720   1/1           4s         16s
job-demo                  1/1           3s         3m7s

$ kubectl get pod
NAME                            READY   STATUS      RESTARTS   AGE
cronjob-demo-1599117660-hf6cs   0/1     Completed   0          95s
cronjob-demo-1599117720-pwrbj   0/1     Completed   0          35s
job-demo-8rl4c                  0/1     Completed   0          3m26s

$ kubectl logs cronjob-demo-1599117660-hf6cs
Thu Sep  3 07:21:02 UTC 2020


$ kubectl describe cronjob cronjob-demo
...
Events:
  Type    Reason            Age    From                Message
  ----    ------            ----   ----                -------
  Normal  SuccessfulCreate  7m34s  cronjob-controller  Created job cronjob-demo-1599117660
  Normal  SawCompletedJob   7m24s  cronjob-controller  Saw completed job: cronjob-demo-1599117660, status: Complete
  Normal  SuccessfulCreate  6m34s  cronjob-controller  Created job cronjob-demo-1599117720
  Normal  SawCompletedJob   6m24s  cronjob-controller  Saw completed job: cronjob-demo-1599117720, status: Complete
  Normal  SuccessfulCreate  5m34s  cronjob-controller  Created job cronjob-demo-1599117780
  Normal  SawCompletedJob   5m24s  cronjob-controller  Saw completed job: cronjob-demo-1599117780, status: Complete
  Normal  SuccessfulCreate  4m34s  cronjob-controller  Created job cronjob-demo-1599117840
  Normal  SawCompletedJob   4m24s  cronjob-controller  Saw completed job: cronjob-demo-1599117840, status: Complete
  Normal  SuccessfulDelete  4m24s  cronjob-controller  Deleted job cronjob-demo-1599117660
  Normal  SuccessfulCreate  3m34s  cronjob-controller  Created job cronjob-demo-1599117900
  Normal  SuccessfulDelete  3m24s  cronjob-controller  Deleted job cronjob-demo-1599117720
  Normal  SawCompletedJob   3m24s  cronjob-controller  Saw completed job: cronjob-demo-1599117900, status: Complete
  Normal  SuccessfulCreate  2m34s  cronjob-controller  Created job cronjob-demo-1599117960
  Normal  SawCompletedJob   2m24s  cronjob-controller  Saw completed job: cronjob-demo-1599117960, status: Complete
  Normal  SuccessfulDelete  2m24s  cronjob-controller  Deleted job cronjob-demo-1599117780
```

## 7.2 资源清单

- `.spec.schedule`：调度，必需字段，指定任务运行周期，格式同Crontab。
- `.spec.jobTemplate`：Job模板，必需字段，指定需要运行的任务，格式同Job
- `.spec.successfulJobHistoryLimit`和`.spec.failedJobHistoryLimit`：历史限制，可选字段，指定可以保留多少完成或失败的Job。默认情况下，他们分别设置为3和1。设置为0相当于完成后不被保留。所以可以看到上面执行了多次，但实际Pod最多只有3个。

# 八、StatefulSet

@todo

# 九、Horizontal Pod Autoscaling

@todo

# 十、小结

本章对`Pod`以及常用的控制器做一些演示及资源清单配置方式进行说明，初步了解各个控制器的作用，能进行日常配置。



---

- [1] [k8s 使用 Init Container 确保依赖的服务已经启动](https://www.cnblogs.com/weihanli/p/12018469.html)
- [2] [Kubernetes Pod 生命周期和重启策略](https://www.dazhuanlan.com/2019/11/04/5dbf2f7da84c5/)
- [3] [Kubernetes K8S之Pod 生命周期与init container初始化容器详解](https://www.cnblogs.com/zhanglianghhh/p/13493337.html)
- [4] [Kubernetes学习之路（十二）之Pod控制器--ReplicaSet、Deployment](https://www.cnblogs.com/linuxk/p/9578211.html)
- [5] [k8s job的使用](https://www.cnblogs.com/benjamin77/p/9903280.html)


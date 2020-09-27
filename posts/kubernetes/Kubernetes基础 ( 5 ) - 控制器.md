```
{
    "url": "k8s-controller",
    "time": "2020/10/01 06:24",
    "tag": "Kubernetes,容器化",
    "toc": "yes"
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
| 6    | Horizontal Pod Autoscaling            | Pod水平自动缩放                                              | 弹性收缩   |

概念在前面章节已经提到，所以本章节主要从示例维度来对控制器进行说明。

# 二、ReplicaSet

## 2.1 RS / RC区别

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

## 2.2 ReplicaSet示例

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

# 三、Deployment

## 3.1 Deployment 与 RS

定义`Deployment`会创建`Pod`和`ReplicaSet`，创建关系大致是`Deployment`创建`ReplicaSet`，`ReplicaSet`创建`Pod`。文档上对它俩的描述：

- ReplicaSet ensures that a specified number of pod replicas are running at any given time.
- Deployment enables declarative updates for Pods and ReplicaSets.

`Deployment`不直接管理`Pod`，而是通过`ReplicaSet`来进行管理，他们的功能差不多，都支持自动扩容、缩容，但`Deployment`支持滚动更新和回滚，这个是`ReplicaSet`不支持的，所以一般建议是通过`Deployment`来管理`Pod`。

##  3.2 资源清单

| 参数名                                     | 字段类型             | 说明                                                 |
| ------------------------------------------ | -------------------- | ---------------------------------------------------- |
| spec.minReadySeconds                       | Integer              |                                                      |
| spec.paused                                | boolean              |                                                      |
| spec.progressDeadlineSeconds               | integer              |                                                      |
| spec.replicas                              | integer              | 控制Pod的副本数量                                    |
| spec.revisionHistoryLimit                  | integer              |                                                      |
| **spec.selector**                          | **Object[required]** | Pod筛选                                              |
| spec.selector.matchExpressions             | []Object             |                                                      |
| spec.selector.matchExpressions[].key       | string[required]     |                                                      |
| spec.selector.matchExpressions[].operator  | string[required]     | 可选值有：`In`, `NotIn`, `Exists` and `DoesNotExist` |
| spec.selector.matchExpressions.values      | []string             |                                                      |
| spec.selector.matchLabels                  | map[string]string    |                                                      |
| **spec.strategy**                          | **Object**           |                                                      |
| spec.strategy.rollingUpdate                | Object               |                                                      |
| spec.strategy.rollingUpdate.maxSurge       | string               |                                                      |
| spec.strategy.rollingUpdate.maxUnavailable | string               |                                                      |
| spec.strategy.type                         | string               |                                                      |
| **spec.template**                          | **Object[required]** | Pod模板                                              |
| spec.template.metadata                     | Object               |                                                      |
| spec.template.spec                         | Object               |                                                      |

## 3.3 Deployment示例

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

## 3.4 滚动更新

既然`Deployment`和`RS`的主要差别在于滚动更新，我们来看看滚动更新操作。

**3.3.1 版本更新**

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

**3.3.2 版本回滚**

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

**3.3.3 版本更新**

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

**3.3.4 kubectl操作**

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

# 四、DaemonSet

## 4.1 关于DaemonSet

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

## 4.2 资源清单

| 参数名                                           | 字段类型             | 说明                                                         |
| ------------------------------------------------ | -------------------- | ------------------------------------------------------------ |
| spec.minReadySeconds                             | Integer              |                                                              |
| spec.revisionHistoryLimit                        | integer              |                                                              |
| **spec.selector**                                | **Object[required]** | Pod筛选                                                      |
| **spec.template**                                | **Object[required]** | Pod模板                                                      |
| **spec.updateStrategy**                          | **Object**           |                                                              |
| spec.updateStrategy.rollingUpdate                | Object               |                                                              |
| spec.updateStrategy.rollingUpdate.maxUnavailable | string               |                                                              |
| spec.updateStrategy.type                         | string               | 更新方式，可选值：RollingUpdate、OnDelete，默认RollingUpdate. |

## 4.3 DaemonSet示例

@todo


# 五、Job

## 5.1 关于Job

负责处理一次性任务。

## 5.2 资源清单

| 参数名                       | 字段类型   | 说明 |
| ---------------------------- | ---------- | ---- |
| spec.activeDeadlineSeconds   | integer    |      |
| spec.backoffLimit            | integer    |      |
| spec.completions             | integer    |      |
| spec.manualSelector          | boolean    |      |
| spec.parallelism             | integer    |      |
| **spec.selector**            | **Object** |      |
| **spec.template**            | **Object** |      |
| spec.ttlSecondsAfterFinished | integer    |      |

## 5.3 Job示例

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

# 六、 CronJob

## 6.1 关于CronJob

`CronJob`就是在`Job`的基础上变成周期性的任务，可以周期性执行。周期设置的`schedule`和`crontab`一样。在后续可以看到，执行`CronJob`会产生`Job`。

## 6.2 资源清单

| 参数名                          | 字段类型         | 说明                         |
| ------------------------------- | ---------------- | ---------------------------- |
| spec.concurrencyPolicy          | string           |                              |
| spec.failedJobsHistoryLimit     | integer          | 保留是失败的Job记录，默认为1 |
| **spec.jobTemplate**            | **Object**       | Job模板，格式同Job           |
| spec.schedule                   | string[required] | Job运行周期，格式同`Crontab` |
| spec.startingDeadlineSeconds    | integer          |                              |
| spec.successfulJobsHistoryLimit | integer          | 保留完成的Job记录，默认为3   |
| spec.suspend                    | boolean          |                              |

## 6.3 CronJob示例

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

# 七、StatefulSet

## 7.1 关于StatefulSet

`StatefulSet`主要解决的是有状态服务的部署问题，前面使用`Deployment+Service`创建的是无状态的`Pod`。就好比`Nginx`的，通过反向代理后面可以挂多个节点，每个节点都是平等的，更换机器只需要调整反向代理的服务即可。但有些服务不行，比如`Redis`，他涉及到数据的存储，数据之间是有状态的，如果切换机器需要将对应的数据存储同步迁移关联上才行，所以一般有状态的服务会配合`PVC`使用。

`StatefulSet`应用特点：

- 稳定且唯一的网络标识符，当节点挂掉后，`Pod`重新调度后`PodName`和`HostName`不变，基于`Headless Service`实现。
- 稳定且持久的存储，当节点挂掉后，`Pod`重新调度后访问相同的持久化存储，基于`PVC`实现。
- 有序、平滑的扩展、部署。`Pod`是有序的，部署或者扩展的时候根据定义的顺序依次进行，从0到N-1，下一个`Pod`运行之前所有的`Pod`必须是`Running`和`Ready`状态，基于`init containers`实现
- 有序、平滑的收缩、删除。根据定义的顺序倒序收缩，及从N-1到0

所以`StatefulSet`的核心功能在于解决稳定的网络表示和持久的存储、服务启停顺序也是确定的，通常包含以下几部分：

- `Headless Service`解决网络表示问题
- `volumeClaimTemplates`创建pvc，关联`pv`解决存储问题
- `StatefulSet`用于定义具体应用

## 7.2 资源清单

| 参数名                                      | 字段类型             | 说明                 |
| ------------------------------------------- | -------------------- | -------------------- |
| spec.podManagementPolicy                    | string               |                      |
| spec.replicas                               | integer              |                      |
| spec.revisionHistoryLimit                   | integer              |                      |
| **spec.selector**                           | **Object[required]** | Pod筛选              |
| spec.serviceName                            | string[required]     | Service 名称         |
| **spec.template**                           | **Object[required]** | Pod模板              |
| **spec.updateStrategy**                     | **Object**           |                      |
| spec.updateStrategy.rollingUpdate           | Object               |                      |
| spec.updateStrategy.rollingUpdate.partition | integer              |                      |
| spec.updateStrategy.type                    | string               | 默认值RollingUpdate. |
| **spec.volumeClaimTemplates**               | **[]Object**         |                      |

## 7.3 StatefulSet示例

**1. 创建Headless Service**

只是将`Service`中的`clusterIP`指定为`None`，不会分配`VIP`。

```
apiVersion: v1
kind: Service
metadata:
  name: ss-nginx-svc
spec:
  ports:
  - port: 38002
    targetPort: 80
  clusterIP: None
  selector:
    app: ss-nginx-pod
```

**2. 创建pv**

这里是存储在本机，也可以通过`nfs`挂载，用来持久化存储数据。

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nginx-pv001
spec:
  capacity:
    storage: 2Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: pv-nginx
  hostPath:
    path: /Users/peng/k8s/pv-data/pv001
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nginx-pv002
spec:
  capacity:
    storage: 2Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: pv-nginx
  hostPath:
    path: /Users/peng/k8s/pv-data/pv002
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nginx-pv003
spec:
  capacity:
    storage: 2Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: pv-nginx
  hostPath:
    path: /Users/peng/k8s/pv-data/pv003
```

**3. 创建StatefulSet应用**

这里为尽量简单，先以`Nginx`做测试。其中`serviceName`关联`Headless Service`。其他的配置同`Deployment`一样。

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: ss-nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ss-nginx-pod
  serviceName: ss-nginx-svc
  template:
    metadata:
      labels:
        app: ss-nginx-pod
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.2-alpine
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
        volumeMounts:
        - name: pvc-nginx
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: pvc-nginx
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: pv-nginx
      resources:
        requests:
          storage: 1Gi

```

查看`statefulset`、`pv`、`pvc`、`pod`、`svc`的状态

```
$ kubectl get statefulset -o wide
NAME       READY   AGE    CONTAINERS   IMAGES
ss-nginx   3/3     2m4s   nginx        nginx:1.19.2-alpine

$ kubectl get pv
NAME          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                          STORAGECLASS   REASON   AGE
nginx-pv001   2Gi        RWO            Recycle          Bound       default/pvc-nginx-ss-nginx-0   pv-nginx                2m28s
nginx-pv002   2Gi        RWO            Recycle          Bound       default/pvc-nginx-ss-nginx-1   pv-nginx                2m28s
nginx-pv003   2Gi        RWO            Recycle          Bound       default/pvc-nginx-ss-nginx-2   pv-nginx                2m28s

$ kubectl get pvc
NAME                   STATUS   VOLUME        CAPACITY   ACCESS MODES   STORAGECLASS   AGE
pvc-nginx-ss-nginx-0   Bound    nginx-pv001   2Gi        RWO            pv-nginx       6m18s
pvc-nginx-ss-nginx-1   Bound    nginx-pv002   2Gi        RWO            pv-nginx       2m24s
pvc-nginx-ss-nginx-2   Bound    nginx-pv003   2Gi        RWO            pv-nginx       2m21s

$ kubectl get pod -o wide
ss-nginx-0                      1/1     Running   0          6m59s   10.1.2.100   docker-desktop   <none>           <none>
ss-nginx-1                      1/1     Running   0          6m56s   10.1.2.101   docker-desktop   <none>           <none>
ss-nginx-2                      1/1     Running   0          6m53s   10.1.2.102   docker-desktop   <none>           <none>

$ kubectl get svc
ss-nginx-svc   ClusterIP      None            <none>        38002/TCP         8m41s

$ kubectl describe svc ss-nginx-svc
Name:              ss-nginx-svc
Namespace:         default
Labels:            <none>
Annotations:       kubectl.kubernetes.io/last-applied-configuration:
                     {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"name":"ss-nginx-svc","namespace":"default"},"spec":{"clusterIP":"None","...
Selector:          app=ss-nginx-pod
Type:              ClusterIP
IP:                None
Port:              <unset>  38002/TCP
TargetPort:        80/TCP
Endpoints:         10.1.2.100:80,10.1.2.101:80,10.1.2.102:80
Session Affinity:  None
Events:            <none>
```

3个`Pod`都创建成功了，每个`Pod`都可以通过以下域名进行访问，通信方式都是通过此域名来访问而非`IP`。

```
<PodName>.<ServiceName>.<NamespaceName>.svc.cluster.local
```

`Pod`故障后可能飘逸到其他节点上，`PodIP`可能会变，但`StatefulSet`会确保`PodName`以及这个域名不变。进入容器之后可以访问及查看。

```
# more /etc/hosts
# Kubernetes-managed hosts file.
127.0.0.1       localhost
::1     localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
fe00::0 ip6-mcastprefix
fe00::1 ip6-allnodes
fe00::2 ip6-allrouters
10.1.2.103      ss-nginx-0.ss-nginx-svc.default.svc.cluster.local       ss-nginx-0

/ # curl ss-nginx-0.ss-nginx-svc.default.svc.cluster.local
pv001
/ # curl ss-nginx-1.ss-nginx-svc.default.svc.cluster.local
pv002
/ # curl ss-nginx-2.ss-nginx-svc.default.svc.cluster.local
pv003

/ # curl ss-nginx-svc.default.svc.cluster.local
pv002
/ # curl ss-nginx-svc.default.svc.cluster.local
pv002
/ # curl ss-nginx-svc.default.svc.cluster.local
pv003
/ # curl ss-nginx-svc.default.svc.cluster.local
pv001
```

当我们删掉`ss-nginx-0`后，`pod`会重建，此时名称会和之前一样，`pvc`会关联同一个。实现原`Pod`相同的功能，达到有状态服务重启后保持相同状态的目的。

```
$ kubectl get pod
ss-nginx-0                      1/1     Running   0          3s
ss-nginx-1                      1/1     Running   0          24m
ss-nginx-2                      1/1     Running   0          24m
```

# 八、Horizontal Pod Autoscaling

## 8.1 关于HPA

应用在日常运行过程中会有高峰也会有低谷的情况，如何做到削峰填谷，提高集群资源的可利用率？HPA就是为了解决此问题。类似阿里云的弹性收缩功能。

## 8.2 资源清单

| 参数名                              | 字段类型             | 说明                               |
| ----------------------------------- | -------------------- | ---------------------------------- |
| spec.maxReplicas                    | integer[required]    | 最大副本数量，不能小于最小副本数量 |
| spec.minReplicas                    | integer              | 最小副本数量，默认值为1.           |
| spec.targetCPUUtilizationPercentage | integer              | 平均CPU使用率，百分比              |
| **spec.scaleTargetRef**             | **Object[required]** | 关联的资源对象                     |
| spec.scaleTargetRef.apiVersion      | string               | 关联资源的api版本                  |
| spec.scaleTargetRef.kind            | string[required]     | 关联资源的类型                     |
| spec.scaleTargetRef.name            | string[required]     | 关联资源的名称                     |

## 8.3 HPA示例

对前一节的StatefulSet应用增加一个HPA控制

```
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: hpa-demo
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: StatefulSet
    name: ss-nginx
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
```

设定副本数量最小3个，最大10个，定义`CPU`指标达到`80%`触发。

```
$ kubectl apply -f hpa.yaml
horizontalpodautoscaler.autoscaling/hpa-demo created
$ kubectl get hpa
NAME           REFERENCE                 TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
hpa-demo       StatefulSet/ss-nginx      <unknown>/80%   3         10        0          4s
```

当前从文档上看到只定义了`CPU`字段，但很显然后续还有会内存、请求数等等指标。创建起来不复杂，但需要配合监控收集数据。这个留待后面在做整体测试。

# 九、小结

本章对`Pod`以及常用的控制器做一些演示及资源清单配置方式进行说明，初步了解各个控制器的作用，能进行日常配置。



---

- [1] [Kubernetes学习之路（十二）之Pod控制器--ReplicaSet、Deployment](https://www.cnblogs.com/linuxk/p/9578211.html)
- [2] [k8s job的使用](https://www.cnblogs.com/benjamin77/p/9903280.html)
- [3] [Kubernetes之StatefulSet](https://www.cnblogs.com/xzkzzz/p/9871837.html)


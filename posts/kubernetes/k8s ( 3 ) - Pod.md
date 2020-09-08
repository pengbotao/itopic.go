```
{
    "url": "k8s-pod",
    "time": "2020/09/06 23:18",
    "tag": "Kubernetes,容器化"
}
```

# 一、概述

## 1.1 Pod简介

![](../../static/uploads/k8s-pod-struct.png)

`k8s`管理的最小粒度资源是`Pod`，它是`k8s`基础资源。一个`Node`里可以运行多个`Pod`，一个`Pod`里可运行多个容器，这些容器共享存储、网络。

- **网络**：每个`Pod`被分配一个独立的`IP`地址，`Pod`中的每个容器共享网络命名空间，包括IP地址和网络端口。`Pod`内的容器可以使用`localhost`相互通信。当`Pod`中的容器与`Pod`外部通信时，他们必须协调如何使用共享网络资源（如端口）。

- **存储**：`Pod`可以指定一组共享存储`volumes`。Pod中的所有容器都可以访问共享`volumes`，允许这些容器共享数据。`volumes`还用于`Pod中`的数据持久化，以防其中一个容器需要重新启动而丢失数据。

## 1.2 Pod分类

`Pod`可以分为两类：

- 自主式Pod：直接创建的`Pod`，没有管理者，退出后不会有重新拉起操作。
- 控制器Pod：可以根据控制器的不同规则对`Pod`做不同的行为控制。

控制器Pod我们会在下一章节介绍控制器时来说明，本章主要介绍自主式Pod。先来看一个示例，了解下Pod的基本用法。

```
apiVersion: v1
kind: Pod
metadata:
  name: mypod
  labels:
    name: mypod
spec:
  restartPolicy: Always
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

- 由于是`latest`版本，默认每次创建容器都会有`Pulling image "busybox"`操作，线上应避免直接拉`latest`镜像，如果要用可以重新打个`tag`。
- `c3`容器执行后就退出了，所以`Pod`会尝试重启容器，但会一直退出，并不是执行失败了，而是执行完退出了。
- 重新查看容器状态为`CrashLoopBackOff`，`READY`状态为`2/3`，`RESTARTS`重启次数为7，重启的粒度是容器层级，从日志上看`c1`、 `c2`容器的重启次数没有增多。因为重启策略默认为`Always`，所以会一直重启下去。

# 二、Pod资源清单

`apiVersion`和`kind`通过`kubectl explain`可以看到

```
$ kubectl explain pod
KIND:     Pod
VERSION:  v1
```

所以，主要介绍下`metadata`和`spec`下的字段说明。

## 2.1 metadata

| 参数名               | 字段类型          | 说明     |
| -------------------- | ----------------- | -------- |
| metadata.name        | string            | 资源名称 |
| metadata.namespace   | string            | 名称空间 |
| metadata.labels      | map[string]string |          |
| metadata.annotations | map[string]string |          |

## 2.2 spec

| 参数名   | 字段类型 | 说明       |
| -------- | -------- | ---------- |
| hostname | string   | 指定主机名 |

## 2.3 spec.containers

`Pod`中包含的容器列表，单个`Pod`可以配置多个容器。需要注意的是因为他们共享网络所以容器里不能都启动同一个端口，否则就端口冲突了，不同`Pod`则没有限制。一般一个容器负责处理一件事情，比如要配置`PHP+Nginx`，则可以`PHP`起一个容器，`Nginx`起一个容器，他们在同一个`Pod`中，通过`volume`的方式共享磁盘。所以`Pod`的概念就类似一台虚拟机，里面可以起多个容器相互协作。

在看`Docker`的时候有一个疑问是把容器理解成虚拟机包含整套服务（比如`CentOS`镜像里安装了各种环境），还是只是一个进程。在`k8s`里看，更倾向于相互独立。

### 2.3.1 image

| 参数名                            | 字段类型 | 说明                                                         |
| --------------------------------- | -------- | ------------------------------------------------------------ |
| spec.containers[]                 | []Object | 定义容器列表                                                 |
| spec.containers[].name            | string   | 定义容器名称                                                 |
| spec.containers[].image           | string   | 定义镜像来源                                                 |
| spec.containers[].imagePullPolicy | string   | 定义镜像拉取策略，有Always、Never、IfNotPresent三个值。默认Always。<br />- Always: 每次都尝试重新拉取镜像<br />- Never：仅使用本地镜像<br />- IfNotPresent：本地有就使用，没有就拉取 |
| spec.containers[].workingDir      | string   | 指定容器工作目录                                             |
| spec.containers[].command         | []string |                                                              |
| spec.containers[].args            | []string |                                                              |

### 2.3.2 ports

| 参数名                                  | 字段类型 | 说明                                          |
| --------------------------------------- | -------- | --------------------------------------------- |
| spec.containers[].ports                 | []Object | 指定容器需要用到的端口列表                    |
| spec.containers[].ports[].name          | string   |                                               |
| spec.containers[].ports[].containerPort | integer  |                                               |
| spec.containers[].ports[].hostIP        | string   |                                               |
| spec.containers[].ports[].hostPort      | integer  |                                               |
| spec.containers[].ports[].protocol      | string   | Must be UDP, TCP, or SCTP. Defaults to "TCP". |

### 2.3.3 volumeMounts

| 参数名                                     | 字段类型 | 说明                   |
| ------------------------------------------ | -------- | ---------------------- |
| spec.containers[].volumeMounts             | []Object | 指定容器内部存储券配置 |
| spec.containers[].volumeMounts[].name      |          |                        |
| spec.containers[].volumeMounts[].mountPath |          |                        |
| spec.containers[].volumeMounts[].readOnly  |          |                        |

### 2.3.4 env

| 参数名                        | 字段类型 | 说明                                 |
| ----------------------------- | -------- | ------------------------------------ |
| spec.containers[].env         | []Object | 指定容器运行前需要设置的环境变量列表 |
| spec.containers[].env[].name  |          |                                      |
| spec.containers[].env[].value |          |                                      |

### 2.3.5 resources

| 参数名                                      | 字段类型          | 说明                           |
| ------------------------------------------- | ----------------- | ------------------------------ |
| spec.containers[].resources                 | Object            | 指定资源限制和资源请求的值     |
| spec.containers[].resources.limits          | map[string]string | 指定容器运行的资源上限         |
| spec.containers[].resources.limits.cpu      |                   |                                |
| spec.containers[].resources.limits.memory   |                   |                                |
| spec.containers[].resources.requests        | map[string]string | 指定容器启动和调度时的限制设置 |
| spec.containers[].resources.requests.cpu    |                   |                                |
| spec.containers[].resources.requests.memory |                   |                                |

### 2.3.6 readinessProbe

| 参数名                           | 字段类型 | 说明     |
| -------------------------------- | -------- | -------- |
| spec.containers[].readinessProbe | Object   | 就绪检测 |

### 2.3.7 livenessProbe

| 参数名                          | 字段类型 | 说明     |
| ------------------------------- | -------- | -------- |
| spec.containers[].livenessProbe | Object   | 存活检测 |

## 2.4 spec.volumes

| 参数名         | 字段类型 | 说明   |
| -------------- | -------- | ------ |
| spec.volumes[] | []Object | 存储卷 |

## 2.5 spec.restartPolicy

前面示例容器会不断重启，除了容器执行完成后，主要还是跟重启参数设置有关，默认为总是重启。看描述：

```
$ kubectl explain pod.spec.restartPolicy
KIND:     Pod
VERSION:  v1

FIELD:    restartPolicy <string>

DESCRIPTION:
     Restart policy for all containers within the pod. One of Always, OnFailure, Never. Default to Always.
```

- 重启策略针对`Pod`里的所有容器
- 值按字面意思理解就行
  - `Always`：容器退出时总是重启容器，默认策略。
  - `OnFailure`：容器异常退出时重启（状态码非0），如果增加`spec.restartPolicy = OnFailure`，则c3容器执行后不会在重启。
  - `Never`：容器退出时不重启容器

## 2.6 spec.initContainers



# 三、 Pod生命周期

![](../../static/uploads/k8s-pod-lifecycle.png)

## 3.1 pause

`Pause容器`，全称infrastucture container（又叫infra）基础容器。每个`Pod`里运行着一个特殊的被称之为`Pause`的容器，其他容器则为业务容器，这些业务容器共享`Pause`容器的网络栈和`Volume`挂载卷，因此他们之间通信和数据交换更为高效。在设计时可以充分利用这一特性，将一组密切相关的服务进程放入同一个`Pod`中；同一个`Pod`里的容器之间仅需通过`localhost`就能互相通信。

## 3.2 init container

在主容器（`Main Container`）启动之前执行的容器，串行执行，只有前一个`InitContainer`正常退出，下一个才会继续。如果`InitContainer`失败，则会根据策略重启`Pod`。相当于在主容器启动之前，可以通过`Init Conntainer`容器做一些准备工作。

## 3.3 hook

`Kubernetes `为容器提供了两种生命周期钩子：

- `post start hook`：于容器创建完成之后立即运行的钩子程序。
- `pre stop hook`：容器终止之前立即运行的程序，是以同步方式的进行，因此其完成之前会阻塞 删除容器的调用

备注：钩子程序的执行方式有`Exec`和`HTTP`两种。

## 3.4 readiness

就绪检测：用于判定容器中的主进程是否准备就绪以及能否对外提供服务。

## 3.5 liveness

存活检测：用于判定主容器是否处于存活状态。



# 四、Pod示例

## 4.1 volumes



## 4.2 init container



# 五、小结

`Pod`是`k8s`的基础资源，着重介绍了`Pod`的资源清单配置方式以及对`Pod`的生命周期做了介绍。下一篇看看如果通过控制器来控制`Pod`。



---

- [1] [k8s 使用 Init Container 确保依赖的服务已经启动](https://www.cnblogs.com/weihanli/p/12018469.html)
- [2] [Kubernetes Pod 生命周期和重启策略](https://www.dazhuanlan.com/2019/11/04/5dbf2f7da84c5/)
- [3] [Kubernetes K8S之Pod 生命周期与init container初始化容器详解](https://www.cnblogs.com/zhanglianghhh/p/13493337.html)


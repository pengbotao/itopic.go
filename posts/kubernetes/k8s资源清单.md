```
{
    "url": "k8s-yaml",
    "time": "2020/09/01 22:06",
    "tag": "Kubernetes"
}
```

# 一、概述

## 1.1 描述文件

k8s常用资源：[<sup>[1]</sup>](#refer)

- 工作负载：`Pod`、` rs(ReplicasSet)`、`deploy(Deployment)`、`sts(StatefulSet)`、`ds(DaemonSet)`、`Job`、`Cronjob`（`ReplicationController`在v1.11版本被遗弃）
- 服务发现及负载均衡：`svc(Service)`、`ing(Ingress)`
- 配置与存储：`Volume`、`pv( persistentvolumes )`、`pvc`、`cm(ConfigMap)`、`Secret`、`DownwardAPI`
- 集群级：`ns(Namespace)`、`Node`、`Role`、`ClusterRole`、`RoleBinding`、`ClusterRoleBinding`
- 元数据：`HPA`、`PodTemplate`、`LimitRange`

## 1.2 资源清单

类似通过`Dockerfile`来表示容器的创建过程，k8s里的各种资源也可以通过文本的方式来创建，通常是通过`Yaml`方式来定义，也支持`Json`。本文档主要介绍资源清单的基本文档说明，至于各种资源的创建示例，再单独起篇章说明。

# 二、部署iTopic

同样，先来看示例。

前一篇中用`kubectl run`的方式启动成功了，但这并不是常用方式，通常还是通过资源清单的方式创建。先来看看示例，第一条命令：

## 2.1 创建deploy

`$ kubectl run itopic --image=pengbotao/itopic.go:alpine --replicas=3 --port=8001`

对应编写`itopic.yaml`（镜像为本博客镜像，可直接替换为`Nginx`镜像），然后执行：`kubectl apply -f itopic.yaml`

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: itopic-deploy
  namespace: default
  labels:
    app: itopic
spec:
  replicas: 3
  selector:
    matchLabels:
      app: itopic
  template:
    metadata:
      labels:
        app: itopic
    spec:
      containers:
      - name: itopic
        image: pengbotao/itopic.go:alpine
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
        ports:
        - containerPort: 8001
```

## 2.2 创建service

第二条命令，创建`service`：

`$ kubectl expose deployment itopic --type=LoadBalancer --port=38001 --target-port=8001 `

对应编写`itopic-svc.yaml`，然后执行：`kubectl apply -f itopic-svc.yaml`

```
apiVersion: v1
kind: Service
metadata:
  name: itopic-svc
  labels:
    app: itopic-svc
spec:
  type: LoadBalancer
  selector:
    app: itopic
  ports:
  - port: 38002
    targetPort: 8001
```

然后访问`http://localhost:38002`就可以访问到了。

# 三、Yaml基础字段

通过`kubectl explain pod`的方式可以查看对应资源的说明文档，要看子节点可以用：`kubectl explain pod.spec`，所有Yaml中的定义的字段都可以通过该方法查看文档。

```
$ kubectl explain pod
KIND:     Pod
VERSION:  v1

DESCRIPTION:
     Pod is a collection of containers that can run on a host. This resource is
     created by clients and scheduled onto hosts.

FIELDS:
   apiVersion	<string>
     APIVersion defines the versioned schema of this representation of an
     object. Servers should convert recognized schemas to the latest internal
     value, and may reject unrecognized values. More info:
     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources

   kind	<string>
     Kind is a string value representing the REST resource this object
     represents. Servers may infer this from the endpoint the client submits
     requests to. Cannot be updated. In CamelCase. More info:
     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds

   metadata	<Object>
     Standard object's metadata. More info:
     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata

   spec	<Object>
     Specification of the desired behavior of the pod. More info:
     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status

   status	<Object>
     Most recently observed status of the pod. This data may not be up to date.
     Populated by the system. Read-only. More info:
     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
```

也可以看看其他的资源，比如`deployment`，大同小异。通过上面看到，基础的`Yaml`格式如下，可以对照前面的示例进行查看，接下来主要介绍`Yaml`里的语法规则：

```
apiVersion: group/apiversion  #API版本
kind:         # 资源类型
metadata：    # 元数据对象
  name        # 如Pod名称
  namespace   # 命名空间默认default
  lables      # 标签
spec：        # 详细对象
  containers：# 容器列表
  - name      # 容器名称
    image     # 镜像
status:       # 当前状态，本字段由 Kubernetes 自身维护，用户不能定义
```

## 3.1 apiVersion

用来指定调用资源接口的版本，通过`$ kubectl api-versions`可以获取到所有支持的版本。通过`explain`也可以看到版本号。

```
$ kubectl explain deployment
KIND:     Deployment
VERSION:  apps/v1
```

## 3.2 kind

定义的资源类型。比如： `Deployment`、`Pod`，通过`explain`也可以看到。查询的时候有些简写可以直接使用，比如：`kubectl get deploy`，但`kind`这里不行，

```
$ kubectl apply -f itopic.yaml
error: unable to recognize "itopic.yaml": no matches for kind "deploy" in version "apps/v1"
```

确保值和`kubectl explain`查看出来的一致就行。

## 3.3 metadata



## 3.4 spec



| 参数名                                             | 字段类型 | 说明                                                         |
| -------------------------------------------------- | -------- | ------------------------------------------------------------ |
| apiVersion                                         | String   | k8s API的版本，可配合前面的`kubectl api-versions`查看        |
| kind                                               | String   | 定义的资源类型，比如`Deployment`、`Pod`                      |
| metadata                                           | Object   | 元数据对象                                                   |
| metadata.name                                      | String   | 元数据对象的名称                                             |
| medadata.namespace                                 | String   | 元数据对象的名称空间                                         |
| spec                                               | Object   | 详细定义对象                                                 |
| spec.containers[]                                  | List     | 定义容器列表                                                 |
| spec.containers[].name                             | String   | 定义容器名称                                                 |
| spec.containers[].image                            | String   | 定义镜像来源                                                 |
| spec.containers[].imagePullPolicy                  | String   | 定义镜像拉取策略，有Always、Never、IfNotPresent三个值。默认Always。<br />- Always: 每次都尝试重新拉取镜像<br />- Never：仅使用本地镜像<br />- IfNotPresent：本地有就使用，没有就拉取 |
| spec.containers[].command[]                        | List     | 指定容器启动命令                                             |
| spec.containers[].workingDir                       | String   | 指定容器工作目录                                             |
| spec.containers[].volumeMounts[]                   | List     | 指定容器内部存储券配置                                       |
| spec.containers[].volumeMounts[].name              |          |                                                              |
| spec.containers[].volumeMounts[].mountPath         |          |                                                              |
| spec.containers[].volumeMounts[].readOnly          |          |                                                              |
| spec.containers[].ports[]                          | List     | 指定容器需要用到的端口列表                                   |
| spec.containers[].ports[].name                     |          |                                                              |
| spec.containers[].ports[].containerPort            |          |                                                              |
| spec.containers[].ports[].hostPort                 |          |                                                              |
| spec.containers[].ports[].protocol                 |          |                                                              |
| spec.containers[].env[]                            | List     | 指定容器运行前需要设置的环境变量列表                         |
| spec.containers[].env[].name                       |          |                                                              |
| spec.containers[].env[].value                      |          |                                                              |
| spec.containers[].resources                        | Object   | 指定资源限制和资源请求的值                                   |
| spec.containers[].resources.limits                 | Object   | 指定容器运行的资源上限                                       |
| spec.containers[].resources.limits.cpu             |          |                                                              |
| spec.containers[].resources.limits.memory          |          |                                                              |
| spec.containers[].resources.limits.requests        | Object   | 指定容器启动和调度时的限制设置                               |
| spec.containers[].resources.limits.requests.cpu    |          |                                                              |
| spec.containers[].resources.limits.requests.memory |          |                                                              |
| spec.restartPolicy                                 | String   | 定义Pod重启策略，可选值：Always、Onfailure、Never，默认为Always<br />- Always：Pod一旦终止则理解重启<br />- Onfailure：非正常退出才重启（Code非0）<br />- Never：不重启 |
| spec.nodeSelector                                  | Object   |                                                              |
| spec.imagePullSecrets                              | Object   |                                                              |
| spec.hostNetwork                                   | Boolean  |                                                              |

# 四、Yaml详细解读

## 4.1 名称空间



## 4.2 Label/Selector



## 4.3 容器定义







# 五、Yaml常用操作

## 5.1 创建资源

```
kubectl apply -f x.yaml
```

或

```
kubectl create -f x.yaml
```

区别：都可以创建资源，如果存在则`create`报错，`apply`会根据新的文件进行更新。

## 5.2 修改/删除资源

**修改资源**：





**删除资源：**

可以通过`kubectl delete deploy itopic-deploy`的方式删除`deploy`，通修改类似，也可以直接指定`yaml`文件的方式来删除。

```
$ kubectl delete -f itopic.yaml
deployment.apps "itopic-deploy" deleted

$ kubectl delete -f itopic-svc.yaml
service "itopic-svc" deleted
```



## 5.3 容器调试

创建过程中可能会出现一些问题，提供一些调试方法：

### 5.3.1 describe

通过`describe`命令查看资源信息，各个资源都可以通过`describe`查看。可以看到状态以及`Events`。

### 5.3.2 查看容器日志

```
$ kubectl logs itopic-6f9dd4f4cd-6n6lp -c itopic
The topic server is running at http://0.0.0.0:8001
Quit the server with Control-C
```

如果Pod里只有一个容器可以省略`-c`参数，创建资源的时候可以通过日志看容器是否成功。

### 5.3.3 service

如果创建的service对外无法提供访问，可以通过`describe`查看`svc`信息

```
$ kubectl describe svc itopic-svc
Name:                     itopic-svc
Namespace:                default
Labels:                   app=itopic
Annotations:              kubectl.kubernetes.io/last-applied-configuration:
                            {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app":"itopic-svc"},"name":"itopic-svc","namespace":"default"},...
Selector:                 app=itopic
Type:                     LoadBalancer
IP:                       10.108.2.60
LoadBalancer Ingress:     localhost
Port:                     <unset>  38002/TCP
TargetPort:               8001/TCP
NodePort:                 <unset>  31744/TCP
Endpoints:                10.1.0.93:8001,10.1.0.94:8001,10.1.0.95:8001
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```

- `Selector: app=itopic` 确保`Selector`是正确的
- `Endpoints：10.1.0.93:8001,10.1.0.94:8001,10.1.0.95:8001` 确保是正确的，如果Selector有错，则可能找不到后端Pod，就无法访问了

## 5.4 查看完整Yaml

不同资源的完整的yaml信息可以这么看：

```
$ kubectl get svc itopic-svc -o yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: "2020-08-31T09:17:17Z"
  labels:
    app: itopic
  name: itopic-svc
  namespace: default
  resourceVersion: "418216"
  selfLink: /api/v1/namespaces/default/services/itopic-svc
  uid: 7198a068-2573-4e37-b0c6-20191940890e
spec:
  clusterIP: 10.97.234.18
  externalTrafficPolicy: Cluster
  ports:
  - nodePort: 31846
    port: 38002
    protocol: TCP
    targetPort: 8001
  selector:
    app: itopic
  sessionAffinity: None
  type: LoadBalancer
status:
  loadBalancer:
    ingress:
    - hostname: localhost
```



---

<div id="refer"></div>

- [1] [1.k8s.资源清单](https://www.cnblogs.com/elvi/p/11755617.html)
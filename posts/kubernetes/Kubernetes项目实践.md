```
{
    "url": "k8s-online",
    "time": "2021/01/22 20:46",
    "tag": "Kubernetes,容器化"
}
```

最近考虑尝试下Kubernetes，将部分项目切换到容器中运行，向容器化再迈进一步，试错成本比较低，阿里云提供了ACK托管版、ACK专有版等各种版本，节点也可以按量付费，大大降低了容器化的配置成本，只需要通过可视化的配置就可以完成一个Kubernetes集群的搭建。

![](../../static/uploads/k8s-ship.jpg)

集群版本：

- Kubernetes版本：1.18.8
- Docker版本：19.3.5
- CentOS：7.8

ACK托管版选项说明：

- Pod网络和Service的网络CIDR需要不同，且与VPC不重复。可根据集群大小填写合适的CIDR。
- API Server：K8s集群的API Server部分，配置连接信息后可以通过kubectl来访问集群。如果只是内网访问可以不用勾选`使用EIP暴露API Server`，也可以后续再配置。
- Ingress：负载均衡类型也选择了内网。开启Ingress会自动创建一个SLB，后面在这个SLB上绑EIP就可以对公网暴露。

配置好之后就可以登录`Node`节点进行查看：

- Node 和 Pod的网络是互通的，这点很重要，解决了通信问题
- 集群自动创建了2个SLB，一个用来访问API Server，一个是勾选的Ingress创建的Service，通过LoadBalancer方式暴露服务
- 按照连接信息配置config，即可使用kubectl访问集群

拿到集群后做了各种尝试发现都比较通畅，碰到一个问题是，通过Service来暴露Pod时，从Pod里访问Service的IP会出现时好时坏的情况，如果调度到自己就出问题。原因是Flannel默认设置不允许回环访问。如果有这种需求考虑用headless Service来暴露服务或者集群使用Terway网络组件。

>  注：headless svc的地址只能在容器内部访问，可以通过ingress对外。

容器部署上还没发现什么问题，省去了部署过程的繁琐，只需要侧重在应用的部署上即可。阿里云后台上提供了可视化的操作，当然也可以直接编写Yaml文件。

整体上还比较顺利，迁移了几个项目，大体就是前面知识的线上实践，对部署过程中初次可能碰到的问题做下整理：

**1. 网络**

网络是基础资源，集群已经打通了容器与VPC，提供服务时如果是内部服务可以用Service，如果是外部可以通过Ingress来暴露服务。暴露之后访问流程如下：

- DNS服务解析域名到Ingress绑定的SLB公网IP上
- 配置对应的Ingress，域名与DNS上配置一致，会连到Ingress的容器上
- 后续就是Ingress到Service到Pod的流程

**2. 镜像**

阿里云提供了镜像仓库，可以直接使用。对于镜像的操作可以参考前面示例章节。

项目可以按无状态部署，PHP项目镜像只打源文件，通过多个容器组合提供服务，拷贝源代码时需要注意下php-fpm运行账号的权限，以免目录无法写入。Python项目源文件打入镜像，基础镜像先搭建好。

镜像上可以在CentOS这样的基础镜像上安装环境，也可以直接使用nginx、php的镜像组合提供服务。目前使用后面这种方式，大部分运行在Debian系统上。

**3. 配置文件**

配置文件可以通过`ConfigMap`来配置，然后挂载到对应的项目里去，支持挂载到指定目录，或者挂载到具体的配置文件。但基本都以文件的方式注入到容器里。用法可参考前面存储章节。

**4. 会话保持**

这可能项目运行中碰到的第一个问题，Ingress会随机解析到后台一个Pod上，如果Session存储在本地则刷新页面可能出现跳出登录的情况，原因是解析到另外一个Pod上了，而该Pod上没有会话信息。可参考`https://kubernetes.github.io/ingress-nginx/examples/affinity/cookie/`在Ingress上配置：

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: test
  namespace: default
  annotations:
    nginx.ingress.kubernetes.io/affinity: "cookie"
    nginx.ingress.kubernetes.io/affinity-mode: "persistent"
    nginx.ingress.kubernetes.io/session-cookie-name: "route"
spec:
```

作用是同一个客户端解析到同一个后端容器，以达到保持会话的目的。但如果Pod重启还是会出现退出登录的情况，这就需要会话支持分布式，可以将会话信息存储到缓存中。

> Service 也可以配置会话保持 sessionAffinity: ClientIP

**5. Crontab**

项目中总会出现些需要配置Cron的地方，用Cronjob的方式可以配置，但有两点不太爽的地方：

- 计划任务比较多的时候要配置的Job比较多，
- 计划任务需要有运行脚本的环境，更新时也需要同步更新

这里用了两种方式来实现：

一种是Cronjob的镜像只需支持curl请求， 配置Cron定时给项目发请求，项目收到请求后做后续处理。

一种是单独起一个Deploy，包含运行环境，同时也安装了Crontab服务，Cron列表可以通过ConfigMap配置，容器启动时启动下Cron服务，加载任务列表并不退出容器，最后就和ecs里跑Crontab的流程一样。没直接打在线上服务的镜像中是因为线上服务往往有多个Pod，而脚本大部分只需要启动一次即可。

> 在发版过程中，POD会替换，POD的退出可能影响正在执行的程序

**6. 日志**

这部分目前还没考虑，部署的是无状态应用，日志存储在容器内部，随着pod生命周期的结束，容器内的日志就丢失了。在部署之前有考虑使用有状态还是无状态，相比而言还是更倾向于无状态，应用本身也不需要管理状态，所以日志这块是后续需要考虑的地方。

**7. 监控与部署**

集群上有`Prometheus`监控，可以观察到Node、Pod的运行情况，集群装了metrics，也通过`kubectl top node/pod`来查看资源使用情况。项目部署上需要做2个动作，一个是打新的镜像，然后就是通知集群更新镜像即可。要添加Node也很方便，通过后台直接添加节点即可，会自动装好Node需要的环境，Node下线做个节点排水设置不调度就自动在其他Node上创建了。
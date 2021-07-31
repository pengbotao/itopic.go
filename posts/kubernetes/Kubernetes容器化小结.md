```
{
    "url": "k8s",
    "time": "2021/07/31 10:20",
    "tag": "Kubernetes,容器化"
}
```

从疫情开始到目前已经一年半了，公司基本处于停滞状态，运维同学离职临时接手运维，机器资源已经做过一轮缩减但依然还有很大的空间，虽然可以当一天和尚撞一天钟，奈何性格不是如此，于是开始找寻优化的方法。本不是专业运维，在原有基础上继续混提升的空间应该很有限，而且系统会弄的越发复杂；而容器化呢，听说之前做过调研对一些服务不太可行就没有实施。行不行不知道，但通过一些初步的调研评估潜意识里觉得这是一条能对原有运维方式进行降维打击的路。

当时的调研报告：

> 现有问题：
>
> 1. 各个服务对资源的需求不同，有CPU型有内存型，按服务类型来部署利用率低。
> 2. 压缩空间有限，资源利用率难进一步提升，混合部署麻烦
> 3. CentOS-6生命周期结束，无赖系统还比较老。
>
> 容器化改造解决的问题：
>
> 1. 提升资源利用率，至少还有30%的空间
> 2. 简化部署、快速交付
> 3. 减少运维成本

总结下来主要是两点：机器利用率 与 运维人力成本。当前的运维方式还处于纯手工的方式，还没有到脚本化、系统化的程度；服务的部署利用率也不高。而目前处于业务低点，迁移改造影响相对小些，正是容器化改造的好时机。

所以决定试一试。

学习的路自然不会顺畅，Docker算是比较完善，概念也相对单一，所以学起来也比较容易写；相比而言Kubernetes的入门门槛要高很多，概念多，安装的组件也多。好在本地的环境搭建还是比较简单，算是解决了第一步，后面就是一点点的抽丝剥茧，找资料、做实验。所以有了下面一些列的Docker和Kubernetes文章，入门而言应该是足够了。

| 标题                                                        | 说明                                                    |
| ----------------------------------------------------------- | ------------------------------------------------------- |
| [Docker ( 1 ) - 小试Docker容器环境](docker-start.html)      | Docker基本介绍                                          |
| [Docker ( 2 ) - Dockerfile](dockerfile.html)                | 尝试通过Dockerfile来构建镜像                            |
| [Docker ( 3 ) - Docker Compose](docker-compose.html)        | 通过Docker Compose来组合容器                            |
| [Docker ( 4 ) - Docker Api](docker-api.html)                | Docker的API接口与基本架构介绍                           |
| [Docker ( 5 ) - Docker小结](docker-summary.html)            | Docker小结                                              |
| [Kubernetes基础 ( 1 ) - 环境搭建及概念说明](k8s-start.html) | k8s基础概念介绍，能通过minukube或者docker搭建一套环境。 |
| [Kubernetes基础 ( 2 ) - 资源清单](k8s-yaml.html)            | 通过2、3、4三个章节对基础的Pod和资源清单、存储做介绍。  |
| [Kubernetes基础 ( 3 ) - Pod](k8s-pod.html)                  | Pod的写法及生命周期介绍                                 |
| [Kubernetes基础 ( 4 ) - 存储](k8s-storage.html)             | Pod依赖的存储介绍                                       |
| [Kubernetes基础 ( 5 ) - 控制器](k8s-controller.html)        | 各种控制器的初步介绍，通过控制器来管理Pod               |
| [Kubernetes基础 ( 6 ) - Service](k8s-service.html)          | 通过Service暴露Pod                                      |
| [Kubernetes基础 ( 7 ) - Ingress](k8s-ingress.html)          | 通过7层负载暴露服务                                     |
| [Kubernetes基础 ( 8 ) - 调度器](k8s-scheduler.html)         | 介绍k8s调度器的逻辑                                     |
| [Kubernetes基础 ( 9 ) - 示例](k8s-example.html)             | 通过示例来演示基本的用法                                |
| [Kubernetes基础 ( 10 ) - 小结](k8s-summary.html)            | Kubernetes小结                                          |

大概经过了半年的准备时间，中间也在不断的做混排优化，感觉到了可操作的程度，预期在现有资源上再优化30%。

接下来就是生成环境的部署问题，首先想到的还是阿里云上有没有现成的服务可用，毕竟目标是来优化机器资源，最后决定使用的是阿里云Kubernetes托管版。

起初并没有急着去批量上线，而是探索了k8s集群内网的通信问题与兼容问题，所以有了这两篇：

| 标题                                                | 说明                           |
| --------------------------------------------------- | ------------------------------ |
| [Kubernetes项目实践 - ACK集群](k8s-online.html)     | 阿里云Kubernetes集群尝试       |
| [阿里云ACK集群中的网络问题](ack-network-issue.html) | 网络问题以及新旧的一些兼容问题 |

了解之后再开始进行服务迁移，基本步骤：

- 运行环境，打造服务应用的环境。比如Python的可通过pipres导出包，选择个Python版本打好基础镜像。
- 编写项目的Dockerfile
- 编写Yaml文件，配置、PV/PVC、Deployment、Service、Ingress
- 调测容器运行情况

集中开始迁移是今年二月底开始，大概3个月左右的时间就将线上100多个服务迁移到了容器。然后就是相关的中间件以及一些工具类的组件迁移：

| 标题                                          | 说明                     |
| ----------------------------------------------- | ------------------------ |
|[Kubernetes容器化 - Jenkins](jenkins-in-k8s.html)|CI/CD|
|[Kubernetes容器化 - JumpServer](jumpserver-in-k8s.html)|跳板机|
|[Kubernetes容器化 - sshd](sshd-in-k8s.html)|通过SSHD代理访问内网服务|
|[Kubernetes容器化 - Gitlab](gitlab-in-k8s.html)|版本控制|
|[Kubernetes容器化 - SVN](svn-in-k8s.html)|版本控制|
|[Kubernetes容器化 - Kafka](kafka-in-k8s.html)|消息队列|
|[Kubernetes容器化 - Zookeeper](zookeeper-in-k8s.html)|Zookeeper|
|[Kubernetes容器化 - MongoDB](mongo-in-k8s.html)|Mongo|
|[Kubernetes容器化 - Mysql](mysql-in-k8s.html)|Mysql|
|[Kubernetes容器化 - Redis](redis-in-k8s.html)|Redis|
|[Kubernetes容器化 - Confluence](confluence-in-k8s.html)|文档系统|
|[Kubernetes容器化 - Elasticsearch](elasticsearch-in-k8s.html)|Elasticsearch集群|

虽然大部分都是单机运行模式，比如Redis，但基本都迁移进了容器，先进来在优化。到这里开发流程就基本覆盖了原有的流程。

项目监控报警上还不够，基本只是将监控指标接入到了本地：[通过自定义Grafana监控集群](grafana-k8s.html)。

到这里容器化就差不多了，整体上看， 机器资源比预期30%的比例要多些，这部分是容器带来的优势。

遗憾的是，公司终究没逃过疫情，在行业类也算是第一梯队的公司，不是被对手打败，也不是被自己玩坏，心多少会有点不干。最后的最后就是进行资源缩减，保持运行的最低要求，原来一两百台Node最后运行在了三个Node里。

公司是结束了，了解容器化之后更坚定这会是以后各服务基础设施，加油！

![](../../static/uploads/WechatIMG766.jpeg)
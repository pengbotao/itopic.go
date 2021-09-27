```
{
    "url": "docker-api",
    "time": "2020/08/02 21:00",
    "tag": "Docker,容器化",
    "toc": "yes"
}
```

了解API之前我们需要了解一下`Docker`的架构。

# 一、Docker概述

## 1.1 Docker平台

`Docker`提供了一种在安全隔离的容器中运行近乎所有应用的方式，这种隔离性和安全性允许你在同一主机上同时运行多个容器，而容器的这种轻量级特性，意味着你可以节省更多的硬件资源，因为你不必消耗运行`hypervisor`所需要的额外负载。

基于容器虚拟化的工具或者平台可以为你提供如下帮助：

- 将应用程序（包括支撑的组件）放入`Docker`容器中；
- 将这些容器打包并分发给你的团队，以便于后续的开发和测试；
- 将这些容器部署到生产环境中，生产环境可以是本地的数据中心，也可以在云端。

## 1.2 Docker引擎

`Docker Engine`包含以下组件：

- `Docker`守护进程（`Docker Daemon`）
- `REST API`，通过API接口与`Docker`守护进程交互。
- `Docker` 客户端

![](../../static/uploads/engine-components-flow.png)

`Docker`是`Client/Server`的架构，`Docker`客户端与`Docker daemon`通过`REST API`交互，`Docker daemon`负责构建、运行和发布 `Docker`容器。客户端可以和服务端运行在同一个系统中，也可以不在。

![](../../static/uploads/architecture.svg)


# 二、Docker API

Docker生态中一共有三种API：

- `Registry API`: 提供与存储`Docker`镜像的`Docker Registry`集成的功能
- `Docker Hub API`: 提供与`Docker Hub`集成的功能
- `Docker Remote API`: 提供与`Docker`守护进程进行集成的功能

三种都是RESTful风格的接口，主要了解的是`Docker Remote API`。

@todo



# 三、Docker原理

通过namespace 来做资源隔离，cgroup 来做资源限制。

---

- [1] [Docker overview](https://docs.docker.com/get-started/overview/)
- [2] [Develop with Docker Engine API](https://docs.docker.com/engine/api/)
- [3] [Docker Engine API reference.](https://docs.docker.com/engine/api/latest/)
- [4] [理解Docker（3）：Docker 使用 Linux namespace 隔离容器的运行环境](https://www.cnblogs.com/sammyliu/p/5878973.html)
- [5] [理解Docker（4）：Docker 容器使用 cgroups 限制资源使用](https://www.cnblogs.com/sammyliu/p/5886833.html)
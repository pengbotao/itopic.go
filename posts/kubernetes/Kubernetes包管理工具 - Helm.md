```
{
    "url": "helm",
    "time": "2021/03/01 09:49",
    "tag": "Kubernetes,容器化"
}
```



# 一、概述

![](../../static/uploads/helm-chart.png)


# 二、基本用法

## 2.1 添加Helm源

添加阿里云镜像

```
$ helm repo add apphub https://apphub.aliyuncs.com/
$ helm repo update
```

**helm repo**的基本用法：

```
Usage:
  helm repo [command]

Available Commands:
  add         add a chart repository
  index       generate an index file given a directory containing packaged charts
  list        list chart repositories
  remove      remove one or more chart repositories
  update      update information of available charts locally from chart repositories

```

## 2.2 搜索Chart

```
$ helm search repo nfs
NAME                            CHART VERSION   APP VERSION     DESCRIPTION
apphub/nfs-client-provisioner   1.2.8           3.1.0           nfs-client is an automatic provisioner that use...
apphub/nfs-server-provisioner   1.0.0           2.3.0           nfs-server-provisioner is an out-of-tree dynami...
```

## 2.3 查看Chart

**helm show**的基本用法：

```
Usage:
  helm show [command]

Aliases:
  show, inspect

Available Commands:
  all         show all information of the chart
  chart       show the chart's definition
  readme      show the chart's README
  values      show the chart's values

$ helm show chart apphub/nfs-client-provisioner
```

## 2.4 安装Chart

选择仓库里的Chart进行安装。

```
$ helm install nfs-storage apphub/nfs-client-provisioner \
--set nfs.server=192.168.88.100 \
--set nfs.path=/home/pengbotao/nfs \
--set nfs.sotrageClass.name=nfs-storage \
--set sotrageClass.defaultClass=true
```


---

- [1] [在Kubernetes上安装nfs-client-provisioner来提供StorageClass](https://knner.wang/2019/12/02/install-nfs-client-provisioner-storageclass.html)
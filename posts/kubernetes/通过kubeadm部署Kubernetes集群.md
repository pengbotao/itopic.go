```
{
    "url": "k8s-kubeadm-install",
    "time": "2020/11/30 22:00",
    "tag": "Kubernetes,容器化",
    "toc": "yes"
}
```

# 一、概述

Kubeadm 是一个工具，它提供了 `kubeadm init` 以及 `kubeadm join` 这两个命令作为快速创建 kubernetes 集群的最佳实践。

kubeadm 通过执行必要的操作来启动和运行一个最小可用的集群。它被故意设计为只关心启动集群，而不是准备节点环境的工作。同样的，诸如安装各种各样的可有可无的插件，例如 Kubernetes 控制面板、监控解决方案以及特定云提供商的插件，这些都不在它负责的范围。

# 二、 就绪检测

## 2.1 机器就位

| Hostname      | IP             | 说明     |
| ------------- | -------------- | -------- |
| peng-master-1 | 172.16.196.200 | 至少2CPU |
| peng-node-1   | 172.16.196.201 |          |
| peng-node-2   | 172.16.196.202 |          |

> 注：所以节点操作以下设置。

## 2.2 关闭防火墙

```
systemctl stop firewalld
systemctl disable firewalld
```

## 2.3 关闭selinux

```
sed -i 's/enforcing/disabled/' /etc/selinux/config 
setenforce 0
```

## 2.4 关闭swap

```
swapoff -a # 临时关闭
sed -ri 's/.*swap.*/#&/' /etc/fstab  #永久关闭
```

## 2.5 将桥接的ipv4流量传到iptables的链

```
cat > /etc/sysctl.d/k8s.conf << EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

sysctl --system
```

## 2.6 时间同步

```
yum install ntpdate -y
ntpdate time.windows.com
```

# 三、安装Docker

> 注：所有节点安装

## 3.1 安装Docker源

```
yum install -y wget && wget https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo -O /etc/yum.repos.d/docker-ce.repo
```

## 3.2 安装Docker

**1. 安装指定版本Docker**

```
yum list | grep docker-ce
containerd.io.x86_64                        1.3.7-3.1.el7              docker-ce-stable
docker-ce.x86_64                            3:19.03.13-3.el7           docker-ce-stable
docker-ce-cli.x86_64                        1:19.03.13-3.el7           docker-ce-stable
docker-ce-selinux.noarch                    17.03.3.ce-1.el7           docker-ce-stable

yum install -y docker-ce-19.03.13-3.el7.x86_64
```

**2. 开机启动与启动Docker服务**

```
systemctl enable docker && systemctl start docker
```

**3. 查看Docker，确保已启动**

```
docker version
docker ps
```

# 四、安装kubeadm

> 注：所有节点安装

## 4.1 添加阿里云yum源

```
cat > /etc/yum.repos.d/kubernetes.repo << EOF
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
```

## 4.2 安装kubeadm、kubelet、kubectl

```
# yum list | grep kubelet
kubelet.x86_64                              1.19.2-0                   kubernetes
[root@peng-master-1 ~]# yum list | grep kubeadm
kubeadm.x86_64                              1.19.2-0                   kubernetes
[root@peng-master-1 ~]# yum list | grep kubectl
kubectl.x86_64                              1.19.2-0                   kubernetes


yum install -y kubelet-1.19.2-0 kubeadm-1.19.2-0 kubectl-1.19.2-0

```

**设置开机启动**

```
systemctl enable kubelet
```

# 五、配置 Master

> 注：只在master节点执行

```
kubeadm init \
--apiserver-advertise-address=172.16.196.200 \
--image-repository registry.aliyuncs.com/google_containers \
--kubernetes-version v1.19.2 \
--service-cidr=10.1.0.0/16 \
--pod-network-cidr=10.244.0.0/16
```

安装成功后提示：

```
Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 172.16.196.200:6443 --token eflt49.q765j3u6zj7yaq0r \
    --discovery-token-ca-cert-hash sha256:18170da9910aa1a9891ca02053a90b9401650534cc485e5320eb8716ab738f29
```

**使用kubectl工具**

```
 mkdir -p $HOME/.kube
 sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
 sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

# 六、安装Flannel

github地址：https://github.com/coreos/flannel

## 6.1 下载镜像

> 注：镜像每个节点安装

```
[root@peng-master-1 ~]# wget https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

[root@peng-master-1 ~]# grep image kube-flannel.yml
        image: quay.io/coreos/flannel:v0.13.0-rc2
        image: quay.io/coreos/flannel:v0.13.0-rc2
[root@peng-master-1 ~]# docker pull quay.io/coreos/flannel:v0.13.0-rc2
```

没有翻墙本地似乎可以拉取成功，如果其有些节点太慢可以打包同步。

```
[root@peng-master-1 ~]# docker save -o flannel.tar quay.io/coreos/flannel:v0.13.0-rc2
[root@peng-master-1 ~]# scp flannel.tar root@172.16.196.201:~/
[root@peng-node-1 ~]# docker load -i flannel.tar
```

## 6.2 安装Flannel

> 注：只Master节点操作

```
kubectl apply -f kube-flannel.yml
```

# 七、添加节点

> 注：Node节点操作

## 7.1 添加Node

```
kubeadm join 172.16.196.200:6443 --token eflt49.q765j3u6zj7yaq0r \
    --discovery-token-ca-cert-hash sha256:18170da9910aa1a9891ca02053a90b9401650534cc485e5320eb8716ab738f29
```

## 7.2 查看Node

```
[root@peng-master-1 ~]# kubectl get node
NAME            STATUS   ROLES    AGE   VERSION
peng-master-1   Ready    master   54m   v1.19.2
peng-node-1     Ready    <none>   13m   v1.19.2
peng-node-2     Ready    <none>   11m   v1.19.2
```

# 八、创建示例

用我们前面的用过的示例进行测试

```
apiVersion: v1
kind: Service
metadata:
  name: go-svc
spec:
  type: NodePort
  ports:
  - port: 36001
    targetPort: 38001
  selector:
    name: k8s-go-demo

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-go-demo
spec:
  replicas: 3
  selector:
    matchLabels:
      name: k8s-go-demo
  template:
    metadata:
      labels:
        name: k8s-go-demo
    spec:
      containers:
      - name: k8s-go-demo
        image: pengbotao/k8s-go-demo:v1
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 38001
```

**创建**

```
# kubectl apply -f k8s-go-demo.yaml
service/go-svc created
deployment.apps/k8s-go-demo created
```

**查看**

```
[root@peng-master-1 k8s]# kubectl get svc
NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)           AGE
go-svc       NodePort    10.1.35.223   <none>        36001:31169/TCP   9m43s

[root@peng-master-1 k8s]# kubectl get pod -o wide
NAME                          READY   STATUS    RESTARTS   AGE     IP           NODE          NOMINATED NODE   READINESS GATES
k8s-go-demo-6bd6875cc-jw9gg   1/1     Running   0          8m21s   10.244.1.2   peng-node-1   <none>           <none>
k8s-go-demo-6bd6875cc-ktscm   1/1     Running   0          8m21s   10.244.1.3   peng-node-1   <none>           <none>
k8s-go-demo-6bd6875cc-w2cgv   1/1     Running   0          8m21s   10.244.2.2   peng-node-2   <none>           <none>
```

**访问**

```
# 虚拟机访问
[root@peng-master-1 ~]# curl 10.1.35.223:36001
{"ClientIP":"10.244.0.0","Host":"k8s-go-demo-6bd6875cc-jw9gg","ServerIP":"10.244.1.2","Time":"2020-09-25 01:42:59","Version":"v1"}

# 宿主机访问
pengbotao:~ peng$ curl http://172.16.196.200:31169
{"ClientIP":"10.244.0.0","Host":"k8s-go-demo-6bd6875cc-w2cgv","ServerIP":"10.244.2.2","Time":"2020-09-25 01:43:37","Version":"v1"}

pengbotao:~ peng$ curl http://172.16.196.201:31169
{"ClientIP":"10.244.1.1","Host":"k8s-go-demo-6bd6875cc-ktscm","ServerIP":"10.244.1.3","Time":"2020-09-25 01:43:42","Version":"v1"}
```

# 九、安装Dashboard

github地址：https://github.com/kubernetes/dashboard

## 9.1 安装

`Dashboard`没有使用`nodeport`，下载`yaml`文件，并对`Service`部分增加了`type`和`nodePort`节点。

```
wget https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.4/aio/deploy/recommended.yaml

---

kind: Service
apiVersion: v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
spec:
  type: NodePort
  ports:
    - port: 443
      targetPort: 8443
      nodePort: 32000
  selector:
    k8s-app: kubernetes-dashboard

---
```

镜像也是一样的安装方式，拉好某一台后同步到其他机器：

```
[root@peng-master-1 ~]# grep image recommended.yaml
          image: kubernetesui/dashboard:v2.0.4
          imagePullPolicy: Always
          image: kubernetesui/metrics-scraper:v1.0.4
```

## 9.2 配置token

```
cat > dashboard-adminuser.yaml << EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: aks-dashboard-admin
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: aks-dashboard-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: aks-dashboard-admin
  namespace: kube-system
EOF

kubectl apply -f dashboard-adminuser.yaml
```

查看Token

```
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep aks-dashboard-admin | awk '{print $1}')
```

期间出现的问题时获取到的Token看不到名称空间，应该是权限的问题，按上面产生的Token可以。

## 9.3 宿主机访问

```
https://peng-master-1:32000/
```

用前一步的Token登录即可，Chrome可能无法登录，可以用Firefox试试。



---



> 说明：以上操作参考[南宫乘风]的教程，表示感谢。除Dashboard部分配置Token有调整外，其他基本一致。

- [1] [kubeadm部署Kubernetes（k8s）完整版详细教程](https://blog.csdn.net/heian_99/article/details/103888459)


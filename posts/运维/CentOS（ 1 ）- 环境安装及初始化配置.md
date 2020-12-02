```
{
    "url": "centos",
    "time": "2017/10/01 08:10",
    "tag": "运维,CentOS",
    "toc": "yes"
}
```

# 一、安装CentOS

工具（都可以从官网下载）：

- VMware Fusion
- CentOS：https://www.centos.org/download/

这里CentOS版本选择的7，ISO选择的x86_64，通过网易镜像选择的`CentOS-7-x86_64-Minimal-2003.iso`

Mac下也没啥需要设置的，可视化的界面，安装过程省略，安装之后可查看到版本信息：

```
# uname -a
Linux peng-master-1 3.10.0-1127.el7.x86_64 #1 SMP Tue Mar 31 23:36:51 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux

# cat /etc/redhat-release
CentOS Linux release 7.8.2003 (Core)
```

# 二、CentOS初始化

## 2.1 调整固定IP

```
$ cd  /Library/Preferences/VMware\ Fusion/vmnet8/
$ more dhcpd.conf

subnet 172.16.196.0 netmask 255.255.255.0 {
        range 172.16.196.128 172.16.196.254;
        option broadcast-address 172.16.196.255;
        option domain-name-servers 172.16.196.2;
        option domain-name localdomain;
        default-lease-time 1800;                # default is 30 minutes
        max-lease-time 7200;                    # default is 2 hours
        option netbios-name-servers 172.16.196.2;
        option routers 172.16.196.2;
}
```

虚拟机通过NAT方式连接，还有个nat.conf文件可以看到NAT网管地址。进入`CentOS`

```
# more /etc/sysconfig/network-scripts/ifcfg-ens33
TYPE=Ethernet
PROXY_METHOD=none
BROWSER_ONLY=no
BOOTPROTO=static
DEFROUTE=yes
IPV4_FAILURE_FATAL=no
IPV6INIT=yes
IPV6_AUTOCONF=yes
IPV6_DEFROUTE=yes
IPV6_FAILURE_FATAL=no
IPV6_ADDR_GEN_MODE=stable-privacy
NAME=ens33
UUID=ce2592d2-6986-4ddc-bbc1-ef087dfc2ee1
DEVICE=ens33
ONBOOT=yes
IPADDR=172.16.196.200
GATEWAY=172.16.196.2
NETMASK=255.255.255.0
DNS1=172.16.196.2
```

- `BOOTPROTO`: dchp -> static
- `ONBOOT`: no -> yes
- 新增：`IPADDR`、`GATEWAY`、`NETMASK`、`DNS1`



保存后重启网络: `service network restart`

## 2.2 调整hostname

```
# hostnamectl set-hostname peng-master-1

# hostnamectl status
   Static hostname: peng-master-1
         Icon name: computer-vm
           Chassis: vm
        Machine ID: 8f859574377f40bc8098beaa3346e00c
           Boot ID: abbf713187bf47ad926d9c92fde55296
    Virtualization: vmware
  Operating System: CentOS Linux 7 (Core)
       CPE OS Name: cpe:/o:centos:centos:7
            Kernel: Linux 3.10.0-1127.el7.x86_64
      Architecture: x86-64
```

注：需要重启机器。

# 三、宿主机调整

## 3.1 添加授权KEY

```
$ ssh-copy-id ~/.ssh/id_rsa.pub root@172.16.196.200
```

## 3.2 配置ssh与hosts

```
172.16.196.200 peng-master-1
172.16.196.201 peng-node-1
172.16.196.202 peng-node-2
```

**Mac下配置~/.ssh/config**

```
Host peng-master-1
    HostName     172.16.196.200
    Port         22
    User         root
    IdentityFile ~/.ssh/id_rsa
    ServerAliveInterval 10

Host peng-node-1
    HostName     172.16.196.201
    Port         22
    User         root
    IdentityFile ~/.ssh/id_rsa
    ServerAliveInterval 10

Host peng-node-2
    HostName     172.16.196.202
    Port         22
    User         root
    IdentityFile ~/.ssh/id_rsa
    ServerAliveInterval 10
```

宿主机登录方式：`ssh peng-master-1`，如果使用`Xshell`就可以省略掉这步。


# 四、yum 配置

yum 的配置文件分为两部分：main 和repository

- main 部分定义了全局配置选项，整个yum 配置文件应该只有一个main。常位于/etc/yum.conf 中。
- repository 部分定义了每个源/服务器的具体配置，可以有一到多个。常位于/etc/yum.repo.d 目录下的各文件中。

yum.conf 文件一般位于/etc目录下，一般其中只包含main部分的配置选项。

# 五、镜像

## 5.1 CentOS镜像

CentOS，是基于 Red Hat Linux 提供的可自由使用源代码的企业级 Linux 发行版本；是一个稳定，可预测，可管理和可复制的免费企业级计算平台。

**1. 备份**

```
mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
```

**2. 下载镜像**

下载新的 CentOS-Base.repo 到 /etc/yum.repos.d/

```
# CentOS6
$ wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-6.repo

# CentOS7
$ wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-7.repo

# CentOS8
$ wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-8.repo
```

**3. 更新缓存**

```
yum clean all
yum makecache
```

## 5.2 epel 镜像

EPEL (Extra Packages for Enterprise Linux), 是由 Fedora Special Interest Group 维护的 Enterprise Linux（RHEL、CentOS）中经常用到的包。

**1. 备份**

备份，如有配置其他epel源

```
mv /etc/yum.repos.d/epel.repo /etc/yum.repos.d/epel.repo.backup
mv /etc/yum.repos.d/epel-testing.repo /etc/yum.repos.d/epel-testing.repo.backup
```

**2. 下载镜像**

**epel(RHEL 8)**

1）安装 epel 配置包

```
$ yum install -y https://mirrors.aliyun.com/epel/epel-release-latest-8.noarch.rpm
```

2）将 repo 配置中的地址替换为阿里云镜像站地址

```
sed -i 's|^#baseurl=https://download.fedoraproject.org/pub|baseurl=https://mirrors.aliyun.com|' /etc/yum.repos.d/epel*
sed -i 's|^metalink|#metalink|' /etc/yum.repos.d/epel*
```

**eple(RHEL5-7)**

```
$ wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo

$ wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-6.repo

$ wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-5.repo
```

**3. 更新缓存**

```
yum clean all
yum makecache
```


- [1] [阿里云官方镜像站](https://developer.aliyun.com/mirror/)
- [2] [网易开源镜像站](http://mirrors.163.com/)

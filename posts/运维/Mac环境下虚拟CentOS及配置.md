```
{
    "url": "centos",
    "time": "2019/06/01 08:10",
    "tag": "运维,Kubernetes",
    "public": "no",
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

```
{
    "url": "ansible-1",
    "time": "2021/09/12 08:50",
    "tag": "运维"
}
```

# 一、关于Ansible

Ansible是一款自动化运维工具，基于Python开发，与Salt不同，Ansible属于无Agent的实现方式。

## 1.1 安装

```
# 下载epel源
$ curl -o /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo
# 安装
$ yum install -y ansible
```

安装之后就可以通过密码访问网络内的节点：

```
$ ansible 192.168.88.100 -m ping --ask-pass
$ ansible 192.168.88.100 -m command -a 'pwd' --ask-pass
```

添加公钥之后就可以去掉密码部分。

```
$ ssh-copy-id 192.168.88.100
$ ansible 192.168.88.100 -m command -a 'pwd'
```

## 1.2 管理清单

`Inventory` - 清单， 用来配置需要管理的主机列表，默认配置文件在`/etc/ansible/hosts`，可以通过`-i`来指定使用不同的主机清单。

```
$ ansible -i www.ini all -m ping
```

下面示例定义了2个组，组的名称分别是`group1`和`group2`，每个组包含2台机器，其中`group1`是直接定义的IP，`group2`给主机定义了2个别名。

```
192.168.88.88

[group1]
192.168.88.100
192.168.88.101

[group2]
alias200 ansible_host=192.168.88.200
alias201 ansible_host=192.168.88.201 ansible_port=22 ansible_user=root ansible_ssh_pass=123456
```

系统也有两个默认组名：`all` 和 `ungrouped`，分别表示所有和没有分组的主机清单。

## 1.3 基本使用

了解了上面两点基本就可以通过ansible进行日常主机管理了,基本用法：

```
$ ansible -i /etc/ansible/hosts all -m command -a 'free -m'
```

- `-i`: 指定主机清单，使用默认配置的话可以省略，后面的`all`是匹配需要执行命令的主机，不可省略。
- `-m`: 指定使用的模块，默认是`command`，可以省略。`ansible`内置了很多的模块，熟悉下就好。
- `-a`: 传递给模块的参数，参数和模块关联，根据对应的模块而定，可以通过`ansible-doc`来查看。

如果需要查看特定模块的文档可以`ansible-doc`查看，后面的章节都是基本模块的用法。

```
$ ansible-doc -s copy
```

# 二、执行命令

## 2.1 command



## 2.2 shell

执行被控制机上的命令。

## 2.3 script

**Description:**

```
- name: Runs a local script on a remote node after transferring it
  script：
    chdir:
    cmd:
    creates:
    decrypt:
    free_form:
    removes:
```

**Usage:**

```
$ ansible test -m script -a '/root/1.sh'
```

# 三、用户模块

## 3.1 user




## 3.2 group

```
ansible test -m group -a 'name=mysql'
```


# 四、文件模块

## 4.1 file

```
ansible test -m file -a 'path=/root/file.md owner=root group=root mode=644 state=touch'
```

**state:**

- file：即使文件不存在，也不会被创建
- link：创建软连接；
- hard：创建硬连接；
- touch：如果文件不存在，则会创建一个新的文件，如果文件或目录已存在，则更新其最后修改时间
- absent：删除目录、文件或者取消链接文件

## 4.2 copy

```
$ ansible test -m copy -a 'src=./www.ini dest=/tmp/ owner=root group=root mode=644'
```

## 4.3 fetch



## 4.4 synchronize



## 4.5 unarchive



# 五、服务相关

## 5.1 yum



## 5.2 service



## 5.3 pip



## 5.4 cron




# 六、其他

## 6.1 template




---

- [1] [非常好的Ansible入门教程（超简单）](https://blog.csdn.net/pushiqiang/article/details/78126063)
- [2] [ansible常用模块](https://www.cnblogs.com/ccorz/p/ansible-chang-yong-mo-kuai.html)
- [3] [ansible - dylloveyou](https://blog.csdn.net/dylloveyou/category_7621040.html)
- [4] [ansible 批量在远程主机上执行命令](https://www.cnblogs.com/amber-liu/p/10403512.html)

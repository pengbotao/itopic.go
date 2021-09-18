```
{
    "url": "ansible-1",
    "time": "2021/09/12 08:50",
    "tag": "运维",
    "toc": "yes"
}
```

# 一、关于Ansible

`Ansible`是一款自动化运维工具，基于Python开发，与Salt不同的是Ansible属于无Agent的实现方式。

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

## 1.2 清单 - Inventory

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
- `-m`: 指定使用的模块，默认是`command`，可以省略。`ansible`内置了很多的模块，模块的用法是ansible的基础。
- `-a`: 传递给模块的参数，参数和模块关联，根据对应的模块而定，可以通过`ansible-doc`来查看。

如果需要查看特定模块的文档可以`ansible-doc`查看，后面的章节都是基本模块的用法。

```
$ ansible-doc -s [module]
```

# 二、执行命令

## 2.1 command

默认模块，可省略。

```
$ ansible test -m command -a 'ls'
```

## 2.2 shell

执行被控制机上的命令。

```
$ ansible test -m shell -a 'ls'
```

## 2.3 script

执行控制机上的脚本。

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

**Description:**

```
- name: Manage user accounts
  user:
    name:  (required) Name of the user to create, remove or modify.
    group: Optionally sets the user's primary group (takes a group name).
    comment: Optionally sets the description (aka `GECOS') of user account.
    home: Optionally set the user's home directory.
    state:  Whether the account should exist or not, taking action if the state is different from what is stated.
```

**Usage:**

```
$ ansible test -m user -a 'name=peng'
```

## 3.2 group

**Description:**

```
- name: Add or remove groups
  group:
    name: (required) Name of the group to manage.
    state: Whether the group should be present or not on the remote host.
```

**Usage:**

```
$ ansible test -m group -a 'name=mysql'
```

# 四、文件模块

## 4.1 file

**Description:**

```
- name: Manage files and file properties
  file:
    state:
      - file：即使文件不存在，也不会被创建
      - link：创建软连接；
      - hard：创建硬连接；
      - touch：如果文件不存在，则会创建一个新的文件，如果文件或目录已存在，则更新其最后修改时间
      - absent：删除目录、文件或者取消链接文件
```

**Usage:**:

```
$ ansible test -m file -a 'path=/root/file.md owner=root group=root mode=644 state=touch'
```

## 4.2 copy

**Description:**

```
- name: Copy files to remote locations
  copy:
    src: Local path to a file to copy to the remote server.
    dest: (required) Remote absolute path where the file should be
```

**Usage:**:

```
$ ansible test -m copy -a 'src=./www.ini dest=/tmp/ owner=root group=root mode=644'
```

## 4.3 fetch


**Description:**

```
- name: Fetch files from remote nodes
  fetch:
    src: (required) The file on the remote system to fetch.
    dest: (required) A directory to save the file into.
```

**Usage:**

```
$ ansible test -m fetch -a 'src=/root/file.md dest=/tmp/'
```

## 4.4 synchronize

**Description:**

```
- name: A wrapper around rsync to make common tasks in your playbooks quick and easy
  synchronize:
    src: # (required) Path on the source host that will be synchronized
    dest: # (required) Path on the destination host that will be

```

**Usage:**

```
$ ansible test -m synchronize -a 'src=/root/test dest=/home/peng'
```


## 4.5 unarchive

**Description:**

```
- name: Unpacks an archive after (optionally) copying it from the local machine.
  unarchive:
    copy：在解压文件之前，是否先将文件复制到远程主机，默认为yes。若为no，则要求目标主机上压缩包必须存在。
    src: 如果copy为yes，则需要指定压缩文件的源路径 
    dest: 远程主机上的一个路径，即文件解压的路径 
```

**Usage:**

```
$ ansible test -m unarchive -a 'src=/root/test.tar.gz dest=/home/peng copy=yes'
```

# 五、服务相关

## 5.1 yum

**Description:**

```
- name: Manages packages with the `yum' package manager
  yum:
    name:  A package name or package specifier with version, like `name-1.0'.
```

**Usage:**

```
$ ansible test -m yum -a 'name=nginx'
```

## 5.2 service

**Description:**

```
- name: Manage services
  service:
    name: (required) Name of the service.
```

**Usage:**

```
$ ansible test -m service -a 'name=nginx enabled=yes state=started'
```

## 5.3 cron

管理被管理机上的Crontab。

**Description:**

```
- name: Manage cron.d and crontab entries
  cron:
```

**Usage:**

```
$ ansible test -m cron -a 'name="Test" minute="*/10" job="/bin/echo Hello"'
```

## 5.4 pip

**Description:**

```
- name: Manages Python library dependencies
  pip:
```

## 5.5 setup

获取被管理机器的资源信息，比如CPU、内存等。

**Description:**

```
- name: Gathers facts about remote hosts
  setup:
```

**Usage:**

```
$ ansible test -m setup -a 'filter=ansible_memory_mb'
```

---

- [1] [非常好的Ansible入门教程（超简单）](https://blog.csdn.net/pushiqiang/article/details/78126063)
- [2] [ansible常用模块](https://www.cnblogs.com/ccorz/p/ansible-chang-yong-mo-kuai.html)
- [3] [ansible专栏 - dylloveyou](https://blog.csdn.net/dylloveyou/category_7621040.html)
- [4] [ansible 批量在远程主机上执行命令](https://www.cnblogs.com/amber-liu/p/10403512.html)

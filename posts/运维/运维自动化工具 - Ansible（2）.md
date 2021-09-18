```
{
    "url": "ansible-2",
    "time": "2021/09/19 08:50",
    "tag": "运维",
    "toc": "yes"
}
```

# 一、概述

前面一章简单对Ansible做了个介绍，借助强大的模块可以进行一些日常操作，但操作方式相对零碎，需要更系统化的管理方式。它提供的Playbooks、Roles就可以帮我们来做系统化的管理。


# 二、剧本 - Playbooks

PlayBook可以通过Yaml的定义方式来告诉Ansible需要调用的模块、将这些调用组合起来就形成了一个PlayBook。针对Playbook、Play、Tasks的概念可以看下图。

![](../../static/uploads/Ansible_Playbook.png)

先来看一个示例：

## 1.1 剧本示例

- 示例场景：安装Nginx服务并从本地模板拷贝个配置文件过去。
- 操作步骤：首先创建`main.yaml`文件和模板文件`./templates/test.conf.j2`

```
$ cat main.yaml
- name: Playbook Test
  hosts: test
  remote_user: root
  vars:
  - listen_port: 81
  tasks:
  - name: Install Nginx
    yum: name=nginx state=present
  - name: Copy Nginx Config File
    template:
      src: test.conf.j2
      dest: /etc/nginx/conf.d/test.conf
      owner: nginx
      group: nginx
    notify: 
    - restart nginx
  - name: Start Nginx
    service: name=nginx enabled=yes state=started
  handlers:
  - name: restart nginx
    service: name=nginx state=restarted

$ cat templates/test.conf.j2
server {
    listen {{listen_port}};
    server_name _;
    root /usr/share/nginx/html;
}
```

测试运行和查看输出：

```
$ ansible-playbook main.yaml

PLAY [Playbook Test] ********************************************************************************
TASK [Gathering Facts] ******************************************************************************
ok: [192.168.88.200]

TASK [Install Nginx] ********************************************************************************
ok: [192.168.88.200]

TASK [Copy Nginx Config File] ***********************************************************************
changed: [192.168.88.200]

TASK [Start Nginx] **********************************************************************************
ok: [192.168.88.200]

RUNNING HANDLER [restart nginx] *********************************************************************
changed: [192.168.88.200]

PLAY RECAP ******************************************************************************************
192.168.88.200: ok=5    changed=2    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

## 1.2 示例解读

- 定义了一个`PlayBook`，名字是`Playbook Test`，作用于分类为`test`的主机群组。
- 通过`vars`来定义了变量，用于传递给配置文件。这里的作用不明显，应该考虑加载变量文件`vars_files`，这样模块文件里能更充分发挥模板的作用。
- tasks定义了3个任务，分别是安装Nginx、通过模板`template`将配置文件传递到远端机器、启动Nginx
- handlers和notify对应，配置方式和Task一致，用于做回调处理，上面用于更新配置后重启Nginx。

剧本的基本用法如上，指定需要执行的主机、调用一组模块来实现功能，可以将模块的组合调用作用发挥出来。

# 二、角色 - Roles

角色的理解更像是可以单独使用的组件，独立于主机之上，可以将Nginx、Mysql这些组件都以角色的方式来编排，实现公用。

Roles有一套自定义的结构，按照结构来定义可以免去维护文件之间的调用关系，Role的初始化及结构：

![](../../static/uploads/Ansible_Roles.png)

接下来，把前面的例子通过Roles的方式来实现：

## 2.1 初始化

通过`ansible-galaxy`来初始化目录结构。

```
[root@peng nginx]# mkdir roles && cd roles
[root@peng nginx]# ansible-galaxy init nginx
[root@peng nginx]# ls
defaults  files  handlers  meta  README.md  tasks  templates  tests  vars
```

## 2.2 创建任务

任务在`tasks`目录，`main.yml`为入口文件，可以通过`include_tasks`来调用其他任务。

```
$ cat roles/nginx/tasks/main.yml
---
# tasks file for nginx

- name: include install.yml
  include_tasks: install.yml

- name: include conf.yml
  include_tasks: conf.yml

- name: include start.yaml
  include_tasks: start.yml

$ cat roles/nginx/tasks/install.yml
- name: Install Nginx
  yum: name=nginx state=present

$ cat roles/nginx/tasks/conf.yml
- name: Copy Nginx Config File
  template:
    src: test.conf.j2
    dest: /etc/nginx/conf.d/test.conf
    owner: nginx
    group: nginx
  notify: 
  - restart nginx

$ cat roles/nginx/tasks/start.yml
- name: Start Nginx
  service: name=nginx enabled=yes state=started
```

## 2.3 模板

模板在`templates`目录，内容同前面。

## 2.4 变量

变量在vars目录。

```
$ cat roles/nginx/vars/main.yml
---
# vars file for nginx

listen_port: 8801
```

## 2.5 回调

回调目录在`handlers`

```
$ cat roles/nginx/handlers/main.yml
---
# handlers file for nginx

- name: restart nginx
  service: name=nginx state=restarted
```

## 2.6 调用角色

最后和`roles`目录同级别创建Playbook并执行:

```
[root@peng ~]# cat nginx.yml

- hosts: test
  roles:
  - nginx

[root@peng ~]# ansible-playbook nginx.yml

PLAY [test] ********************************************************************

TASK [Gathering Facts] *********************************************************
ok: [192.168.88.200]

TASK [nginx : include install.yml] *********************************************
included: /root/roles/nginx/tasks/install.yml for 192.168.88.200

TASK [nginx : Install Nginx] ***************************************************
ok: [192.168.88.200]

TASK [nginx : include conf.yml] ************************************************
included: /root/roles/nginx/tasks/conf.yml for 192.168.88.200

TASK [nginx : Copy Nginx Config File] ******************************************
changed: [192.168.88.200]

TASK [nginx : include start.yaml] **********************************************
included: /root/roles/nginx/tasks/start.yml for 192.168.88.200

TASK [nginx : Start Nginx] *****************************************************
ok: [192.168.88.200]

RUNNING HANDLER [restart nginx] ************************************************
changed: [192.168.88.200]

PLAY RECAP *********************************************************************
192.168.88.200             : ok=8    changed=2    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

---

[1] [Ansible-playbook 运维笔记](https://www.cnblogs.com/kevingrace/p/5569648.html)
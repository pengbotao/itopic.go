```
{
    "url": "ansible-2",
    "time": "2021/09/19 08:50",
    "tag": "运维",
    "toc": "yes"
}
```


# 一、剧本 - Playbooks

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



# 二、角色 - Roles




---

[1] [Ansible-playbook 运维笔记](https://www.cnblogs.com/kevingrace/p/5569648.html)
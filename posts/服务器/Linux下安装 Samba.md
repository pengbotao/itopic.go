```
{
    "url": "linux-samba-install",
    "time": "2014/09/17 19:52",
    "tag": "Linux,Samba"
}
```

Samba可以实现Linux与Win之间的文件共享，在内部开发文档共享上还是极好的。

# 安装Samba

查看Samba是否已安装
```
# rpm -qa | grep samba
```
通过yum直接安装Samba服务端和客户端
```
#yum -y install samba samba-client
```

# 配置Samba

Samba的主配置文件为`/etc/samba/smb.conf`，这里主要达到在win下通过帐号登录linux共享即可，复杂的配置可参考最后的地址。直接在smb.conf后加上一段共享块。
```
[public]
    comment = Public Files
    path = /data/share/public
    public = yes
    writable = yes
    browseable = yes
    create mode = 0664
    directory mode = 0775
    write list = smbuser
    guest ok = no
```
配置块通过TAB来缩进，拷贝出来的可能是空格。添加系统用户并将用户添加到samba账户中
```
# useradd -s /sbin/nologin smbuser
# smbpasswd -a smbuser
```
smbpasswd命令

- smbpasswd -a 增加用户（要增加的用户必须以是系统用户）
- smbpasswd -d 冻结用户，就是这个用户不能在登录了
- smbpasswd -e 恢复用户，解冻用户，让冻结的用户可以在使用
- smbpasswd -n 把用户的密码设置成空. 注意如果设置了"NO PASSWORD"之后，要允许使用者以空口令登入到Samba服务器，管理员必须在smb.conf配置档案的[global]段中设置以下的参数：null passwords = yes
- smbpasswd -x 删除用户 

测试配置是否无误 - **testparm**

测试Samba的设置是否正确无误，如上面的配置
```
# testparm -s smb.conf 
Load smb config files from smb.conf
Processing section "[public]"
Loaded services file OK.
Server role: ROLE_STANDALONE
[global]
    workgroup = MYGROUP
    server string = Samba Server Version %v
    passdb backend = tdbsam
    log file = /var/log/samba/%m.log
    cups options = raw
 
[public]
    comment = Public Files
    path = /data/share/public
    read only = No
```
# 启动Samba

查看Samba服务状态、启动及重启。
```
# service smb status
smbd (pid  30408) is running...
nmbd (pid  30411) is running...
 
# service smb stop
Shutting down SMB services:                                [  OK  ]
Shutting down NMB services:                                [  OK  ]
 
# service smb start
Starting SMB services:                                     [  OK  ]
Starting NMB services:                                     [  OK  ]
 
# service smb restart
Shutting down SMB services:                                [  OK  ]
Shutting down NMB services:                                [  OK  ]
Starting SMB services:                                     [  OK  ]
Starting NMB services:                                     [  OK  ]
```
设置Samba服务开机自启动
```
# chkconfig --list | grep smb
smb             0:off   1:off   2:off   3:off   4:off   5:off   6:off
# chkconfig --level 35 smb on
# chkconfig --list | grep smb
smb             0:off   1:off   2:off   3:on    4:off   5:on    6:off
```

**Windows清除共享记录**

通过Samba连接成功后会在本地记录登录的帐号密码，下次可直接连接，如果需要切换帐号可手动删除连接记录。
查看访问记录
```
C:\Users\Administrator>net use
不记录新的网络连接。
 
状态       本地        远程                      网络
 
--------------------------------------------------------------------------
OK                     \\42.121.104.209\public   Microsoft Windows Network
命令成功完成。
```
**清除访问记录**
```
C:\Users\Administrator>net use \\42.121.104.209\public /delete
\\42.121.104.209\public 已经删除。
```

**说明**：如何Window下提示没有权限访问，请与管理员管理员联系请求访问权限。则可能是selinux防火墙的问题，执行下面命令关闭selinux防火墙试试：
```
setenforce 0
```
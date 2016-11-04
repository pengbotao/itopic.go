url: linux-vsftp-install
des: 
time: 2014/09/01 22:52
tag: linux
++++++++

# 安装vsftpd

查看是否已安装vsftpd
```
# rpm -qa | grep vsftpd
vsftpd-2.0.5-28.el5
```
yum直接安装vsftpd
```
# yum -y install vsftpd
```
vsftpd服务启动及停止
```
# service vsftpd start
Starting vsftpd for vsftpd:                                [  OK  ]
# service vsftpd restart
Shutting down vsftpd:                                      [  OK  ]
Starting vsftpd for vsftpd:                                [  OK  ]
# service vsftpd stop
Shutting down vsftpd:                                      [  OK  ]
# service vsftpd status
vsftpd is stopped
```
## 卸载vsftpd
```
# rpm -e vsftpd
warning: /etc/vsftpd/vsftpd.conf saved as /etc/vsftpd/vsftpd.conf.rpmsave
```
## vsftpd相关文件

vsftpd相关文件 	| 文件说明
--- | ---
/etc/vsftpd/vsftpd.conf | 主配置文件
/usr/sbin/vsftpd | Vsftpd的主程序。
/etc/pam.d/vsftpd | PAM认证文件（此文件中file=/etc/vsftpd/ftpusers字段，指明阻止访问的用户来自/etc/vsftpd/ftpusers | 文件中的用户）
/etc/vsftpd/ftpusers | 禁止使用vsftpd的用户列表文件。记录不允许访问FTP服务器的用户名单，管理员可以把一些对系统安全有威胁的用户账号记录在此文件中，以免用户从FTP登录后获得大于上传下载操作的权利，而对系统造成损坏。
/etc/vsftpd/user_list | 禁止或允许使用vsftpd的用户列表文件。这个文件中指定的用户缺省情况（即在/etc/vsftpd/vsftpd.conf中设置userlist_deny=YES）下也不能访问FTP服务器，在设置了userlist_deny=NO时,仅允许user_list中指定的用户访问FTP服务器。
/var/ftp | vsftpd的匿名用户登录根目录。

## vsftp配置
参数配置  | 默认值 | 说明
--- | --- | ---
anonymous_enable=YES | YES | 是否用于匿名用户(ftp或anonymous)登录FTP，登录后进入/var/ftp
local_enable=YES | NO | 是否允许本地用户登录FTP服务器，登录后进入用户主目录
write_enable=YES  | NO  | 是否允许写入
local_umask=022  | 077  | 默认的umask码
anon_upload_enable=YES | NO  | 是否允许匿名用户上传文件。如果此项要生效，则配置write_enable必须激活。并且匿名用户所在相关目录有写权限。
anon_mkdir_write_enable=YES  | NO  | 是否允许匿名用户创建新目录。如果此项要生效，则配置write_enable必须激活。并且匿名用户所在相关目录有写权限。
dirmessage_enable=YES  | 	NO  | 是否激活目录欢迎信息功能。.message文件可以通过更改message_file来调整。
xferlog_enable=YES  | 	NO 	 | 是否启动记录上传和下载日志。
connect_from_port_20=YES  | 20 	 | 设定PORT模式下的连接端口(只要connect_from_port_20被激活)。
chown_uploads=YES  | NO  | 	设定是否允许改变上传文件的属主，与下面一个设定项配合使用
chown_username=whoever  | 	ROOT  | 置想要改变的上传文件的属主，如果需要，则输入一个系统用户名，例如可以把上传的文件都改成root属主。whoever：任何人
xferlog_file=/var/log/xferlog  | 	/var/log/xferlog  | 设置日志文件的文件名和存储路径
xferlog_std_format=YES  | 	NO 	 | 是否使用标准的ftpd xferlog日志文件格式
idle_session_timeout=600  | 300  | 设置空闲的用户会话中断时间,默认是10分钟
data_connection_timeout=120  | 300  | 	设置数据连接超时时间,默认是120秒
nopriv_user=ftpsecure  | nobody  | 	运行vsftpd需要的非特权系统用户
async_abor_enable=YES  | NO  | 是否允许运行特殊的ftp命令async ABOR。
ascii_upload_enable=YES ascii_download_enable=YES  | NO  | 是否使用ascii码方式上传和下载文件。
deny_email_enable=YES banned_email_file=/etc/vsftpd/banned_emails  | NO  | 禁止匿名用户通过banned_email_file定义的邮件地址做密码
chroot_list_enable=YES chroot_list_file=/etc/vsftpd/chroot_list  | NO  | 设置为NO时，用户登录FTP后具有访问自己目录以外的其他文件的权限；设置为YES时，chroot_list_file中的用户列表被锁定在自己的home目录下。此时chroot_local_user=NO，如果chroot_local_user=YES则chroot_list_file中的用户将不被锁定在home目录下。
ls_recurse_enable=YES  | 	NO 	 | 是否允许递归查询
listen=YES  | 	NO  | 	vsftpd 处于独立启动模式
listen_ipv6=YES  | 	NO  | 	是否支持IPV6
pam_service_name=vsftpd  | 	ftp  | 	设定vsftpd将要用到的PAM服务的名字。
userlist_enable=YES  | 	NO 	 | 设置为YES，vsftpd将读取userlist_file参数所指定的文件中的用户列表。当列表中的用户登录FTP服务器时，该用户在提示输入密码之前就被禁止了。即该用户名输入后，vsftpd查到该用户名在列表中，vsftpd就直接禁止掉该用户，不会再进行询问密码等后续步聚
userlist_deny=YES  | 	YES  | 	决定禁止还是只允许由userlist_file指定文件中的用户登录FTP服务器。此选项在userlist_enable 选项启动后才生效。YES，默认值，禁止文件中的用户登录，同时也不向这些用户发出输入密码的提示。NO，只允许在文件中的用户登录FTP服务器
userlist_file  | 	/etc/vsftpd/user_list  | 	当userlist_enable被激活，系统将去这里调用文件。
tcp_wrappers=YES  | 	NO  | 	是否允许tcp_wrappers管理
listen_port  | 	21  | 	如果vsftpd处于独立运行模式，这个端口设置将监听的FTP连接请求。
max_clients  | 	0  | 	FTP的最大连接数，0为无限制。
max_per_ip  | 	0  | 	单个IP的最大连接数。
anon_max_rate  | 	0  | 	匿名用户允许的最大传输速度，单位：字节/秒
local_max_rate  | 	0  | 	本地认证用户允许的最大传输速度，单位：字节/秒。

vsftp提供了3中登录性质：匿名帐号、本地帐号和虚拟帐号。默认配置就可以实现匿名帐号和真实帐号登录。

# Vsftp虚拟用户配置

## 1、添加虚拟用户列表

创建一个虚拟用户列表文件，保存需要配置的虚拟帐号，格式为：第一行用户名，第二行密码，依次类推。
```
# vi /etc/vsftpd/virtual_user_list
ftp_www
123456
ftp_log
123456
```
## 2、生成虚拟用户口令认证文件

通过db_load命令生成认证文件。查询db_load命令是否已安装，未安装则先安装。
```
# rpm -qa |grep db4-utils
```
本地CentOS未安装，直接通过yum安装即可。
```
# yum -y install db4-utils
```
将前面添加的virtual_user_list虚拟用户口令文件转换成系统识别的口令认证文件。之后若需要调整该文件，需重新执行db_load动作。
```
# db_load -T -t hash -f /etc/vsftpd/virtual_user_list /etc/vsftpd/virtual_user_list.db
```
## 3、设置虚拟用户所需的PAM配置文件

vsftpd.conf中的pam_service_name参数可指定配置文件。这里修改默认的vsftpd文件。
```
# vi /etc/pam.d/vsftpd
# 注释掉其他部分
auth required /lib/security/pam_userdb.so db=/etc/vsftpd/virtual_user_list
account required /lib/security/pam_userdb.so db=/etc/vsftpd/virtual_user_list
```
## 4、创建vsftpd宿主帐号
```
# useradd vsftpd -s /sbin/nologin
```
## 5、编辑vsftpd.conf配置文件

编辑vsftp配置文件，添加虚拟用户相关配置。
```
# vi /etc/vsftpd/vsftpd.conf
 
# 设定启用虚拟用户功能。
guest_enable=YES
 
# 指定虚拟用户的宿主用户，前面创建的vsftpd用户。
guest_username=vsftpd
 
# 设定虚拟用户的权限符合他们的宿主用户。
virtual_use_local_privs=YES
 
# 设定虚拟用户个人vsftp的配置文件存放路径。配置文件须与虚拟用户同名
# 如/etc/vsftpd/vconf/ftp_www则可以定义ftp_www的特殊配置。
user_config_dir=/etc/vsftpd/vconf
```
## 6、根据虚拟用户设置不同的权限

设定ftp_www只能访问/data/ftp/www，ftp_log只能文芳/data/ftp/log，创建所需的文件和配置信息。
```
# mkdir -p /data/ftp/{www,log}
# chown -R vsftpd:vsftpd /data/ftp/*
 
# touch /etc/vsftpd/chroot_list
 
# mkdir /etc/vsftpd/vconf
# vi /etc/vsftpd/vconf/ftp_www
local_root=/data/ftp/www
max_clients=1
max_per_ip=1
local_max_rate=10000
 
# vi /etc/vsftpd/vconf/ftp_log
local_root=/data/ftp/log
```
## 7、重启服务
```
# service vsftpd restart
```
此时整个配置文件如下：
```
# more vsftpd.conf 
anonymous_enable=NO
local_enable=YES
write_enable=YES
local_umask=022
anon_upload_enable=NO
anon_mkdir_write_enable=NO
dirmessage_enable=YES
xferlog_enable=YES
connect_from_port_20=YES
#chown_uploads=YES
#chown_username=whoever
#xferlog_file=/var/log/xferlog
xferlog_std_format=YES
#idle_session_timeout=600
#data_connection_timeout=120
#nopriv_user=ftpsecure
#async_abor_enable=YES
#ascii_upload_enable=YES
#ascii_download_enable=YES
#ftpd_banner=Welcome to blah FTP service.
#deny_email_enable=YES
#banned_email_file=/etc/vsftpd/banned_emails
chroot_list_enable=YES
chroot_list_file=/etc/vsftpd/chroot_list
chroot_local_user=YES
#ls_recurse_enable=YES
listen=YES
#listen_ipv6=YES
 
pam_service_name=vsftpd
userlist_enable=YES
tcp_wrappers=YES
 
guest_enable=YES
guest_username=vsftpd
virtual_use_local_privs=YES
user_config_dir=/etc/vsftpd/vconf
```
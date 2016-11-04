```
{
    "url": "linux-sz-rz",
    "time": "2014/01/23 10:24",
    "tag": "linux"
}
```

# 1、软件安装

**编译安装**

```
cd /tmp
wget http://www.ohse.de/uwe/releases/lrzsz-0.12.20.tar.gz
tar zxvf lrzsz-0.12.20.tar.gz && cd lrzsz-0.12.20
./configure && make && make install
```
创建软链接，并命名为rz/sz：
```
cd /usr/bin
ln -s /usr/local/bin/lrz rz
ln -s /usr/local/bin/lsz sz
```

**yum安装**

```
yum install -y lrzsz
```
# 2、使用说明

sz命令发送文件到本地：
```
# sz filename
```

rz命令本地上传文件到服务器：
```
# rz
```
执行该命令后，在弹出框中选择要上传的文件即可。

说明：打开SecureCRT软件 -> Options -> session options -> X/Y/Zmodem 下可以设置上传和下载的目录。
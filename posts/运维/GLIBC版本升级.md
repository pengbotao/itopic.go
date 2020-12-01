```
{
    "url": "glibc-upgrade",
    "time": "2017/09/12 18:08",
    "tag": "运维"
}
```

安装`Seafile-6.2.6`对`GLIBC`版本有要求：`GLIBC_2.17 not found`。于是执行了下面操作：

```
scp root@172.16.1.100:/lib64/libc.so.6 /lib64/
```

**然后，就没有然后了...**

终端立即、马上断开连接，所有服务无法运行，也无法通过`sh`、`vnc`进行远程连接，于是只好求助阿里云。修复后重新通过编译的方式进行安装成功。安装步骤如下：

```
$ wget http://ftp.gnu.org/gnu/glibc/glibc-2.18.tar.gz
$ tar zxvf glibc-2.18.tar.gz
$ cd glibc-2.18
$ mkdir build
$ cd build/
$ ../configure --prefix=/usr
$ make -j4
$ make install
```

验证是否安装成功：

```
$ strings /lib64/libc.so.6 | grep GLIBC
GLIBC_2.2.5
GLIBC_2.2.6
GLIBC_2.3
GLIBC_2.3.2
GLIBC_2.3.3
GLIBC_2.3.4
GLIBC_2.4
GLIBC_2.5
GLIBC_2.6
GLIBC_2.7
GLIBC_2.8
GLIBC_2.9
GLIBC_2.10
GLIBC_2.11
GLIBC_2.12
GLIBC_2.13
GLIBC_2.14
GLIBC_2.15
GLIBC_2.16
GLIBC_2.17
GLIBC_2.18
GLIBC_PRIVATE
```
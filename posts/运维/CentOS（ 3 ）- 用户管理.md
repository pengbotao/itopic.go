```
{
    "url": "linux-user",
    "time": "2017/11/01 20:48",
    "tag": "运维",
    "toc": "yes"
}
```

# 一、用户与用户组

## 1.1 用户列表

用户列表存储在`/etc/passwd`文件，查看用户列表：

```
$ cat /etc/passwd
root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
```

`passwd` 文件格式：

```
name:password:uid:gid:comment:home:shell
```

| 编号 | 字段     | 说明                                                         |
| ---- | -------- | ------------------------------------------------------------ |
| 1    | name     | 用户名                                                       |
| 2    | password | 用户密码，早期放这里，由于安全因素用`x`占位，实际存放在`/etc/shadow` |
| 3    | uid      | 用户ID，0是root。通过`$ id 用户名`可以查看用户ID和组ID       |
| 4    | gid      | 用户组ID，对应`/etc/group`中的第三个字段。                   |
| 5    | comment  | 注释                                                         |
| 6    | home     | Home目录或者主目录。登录后进入该目录。                       |
| 7    | shell    | 登录成功后执行的shell，为 `/sbin/nologin` 或`/bin/false`的账号表示不可登录账号 |

`UID`的范围是`0-65535`，分为超级用户、系统用户和普通用户：

- 超级用户：`root`
- 系统用户：比如`nobody`、`bin`等。`CentOS6`中系统用户`UID`取值范围为`[1, 500)`，`CentOS7`中为：`[1,1000)`
- 普通用户：非上面两种用户

查看普通用户列表

```
# 查看普通用户列表
$ awk -F ':' '$3>=500{print $0}' /etc/passwd
zabbix:x:501:501::/home/zabbix:/sbin/nologin
www:x:502:502::/home/www/:/bin/bash
```

## 1.2 密码文件

密码文件和用户是分开存储的，用户密码存储在`/etc/shadow`文件中，只有`root`用户有权限查看。

```
# more /etc/shadow
root:$6$LXFRkKpuGgx4YnCE$RCq2K/BvigVqOb5CFMZk.E/0MSUge7l7gPso32T16hZ/::0:99999:7:::
bin:*:18353:0:99999:7:::
daemon:*:18353:0:99999:7:::
ntp:!!:18529::::::
peng:!!:18546:0:99999:7:::
```

`shadow`文件格式：

```
用户名：加密密码：最后一次修改时间：最小修改时间间隔：密码有效期：密码需要变更前的警告天数：密码过期后的宽限时间：账号失效时间：保留字段
```

| 编号 | 字段                     | 说明                                                         |
| ---- | ------------------------ | ------------------------------------------------------------ |
| 1    | 用户名                   | 用户名，对应`/etc/passwd`中的用户名                          |
| 2    | 密码                     | 有几种格式：<br />- `*`<br />- `!!`：没有设置过密码<br />- 密码串：密码串会隐藏前面规则，比如`$6$`标识使用的`SHA-512`算法 |
| 3    | 最后一次修改时间         | 此字段表示最后一次修改密码的时间与1970-01-01相差的天数。     |
| 4    | 最小修改时间间隔         | 最小修改间隔时间，也就是说，该字段规定了从第 3 字段（最后一次修改密码的日期）起，多长时间之内不能修改密码。如果是 0，则密码可以随时修改；如果是 10，则代表密码修改后 10 天之内不能再次修改密码。此字段是为了针对某些人频繁更改账户密码而设计的。 |
| 5    | 密码有效期               | 这个字段可以指定距离第 3 字段（最后一次更改密码）多长时间内需要再次变更密码，否则该账户密码进行过期阶段。该字段的默认值为 99999。 |
| 6    | 密码需要变更前的警告天数 | 与第 5 字段相比较，当账户密码有效期快到时，系统会发出警告信息给此账户，提醒用户 "再过 n 天你的密码就要过期了，请尽快重新设置你的密码！"。该字段的默认值是 7。 |
| 7    | 密码过期后的宽限天数     | 也称为“口令失效日”，简单理解就是，在密码过期后，用户如果还是没有修改密码，则在此字段规定的宽限天数内，用户还是可以登录系统的；如果过了宽限天数，系统将不再让此账户登陆，也不会提示账户过期，是完全禁用。 |
| 8    | 账号失效时间             | 同第 3 个字段一样，使用自 1970 年 1 月 1 日以来的总天数作为账户的失效时间。该字段表示，账号在此字段规定的时间之外，不论你的密码是否过期，都将无法使用！ |
| 9    | 保留                     | 这个字段目前没有使用，等待新功能的加入。                     |

## 1.3 用户组

用户组存储在`/etc/group`文件，查看分组：

```
# cat /etc/group
root:x:0:
bin:x:1:
daemon:x:2:
www:x:500:peng,www,test
```

`group`文件格式：

```
组名:组密码:组ID(GID):组中的用户
```

| 编号 | 字段       | 说明                                      |
| ---- | ---------- | ----------------------------------------- |
| 1    | 组名       | 组名称                                    |
| 2    | 组密码     | `x`用来占位，真正密码存储在`/etc/gshadow` |
| 3    | 组ID       | 组ID                                      |
| 4    | 组中的用户 | 比如上面的www组下有多个用户               |

查看用户所属分组：

```
$ id peng
uid=556(peng) gid=556(peng) 组=556(peng),500(www)
```

组分为初始组和附加组：

- **初始组**：每个用户的初始组只有一个，通常是以相同的用户名创建组。`passwd`文件中的组对应的是初始组。
- **附加组**：用户可以加入多个其他的组，并拥有这些组的权限。比如上面示例拥有2个组的权限。

## 1.4 配置文件

`/etc/login.defs`文件用于在创建用户时，对用户的一些基本属性做默认设置，例如指定用户`UID`和`GID`的范围，用户的过期时间，密码的最大长度，等等。

需要注意的是，该文件的用户默认配置对`root`用户无效。并且，当此文件中的配置与 `/etc/passwd` 和 `/etc/shadow` 文件中的用户信息有冲突时，系统会以`/etc/passwd`和 `/etc/shadow` 为准。

| 设置项          | 示例            | 含义                                                         |
| --------------- | --------------- | ------------------------------------------------------------ |
| MAIL_DIR        | /var/spool/mail | 创建用户时，系统会在目录 /var/spool/mail 中创建一个用户邮箱，比如 lamp 用户的邮箱是 /var/spool/mail/lamp。 |
| PASS_MAX_DAYS   | 99999           | 密码有效期，99999 是自 1970 年 1 月 1 日起密码有效的天数，相当于 273 年，可理解为密码始终有效。 |
| PASS_MIN_DAYS   | 0               | 表示自上次修改密码以来，最少隔多少天后用户才能再次修改密码，默认值是 0。 |
| PASS_MIN_LEN    | 5               | 指定密码的最小长度，默认不小于 5 位，但是现在用户登录时验证已经被 PAM 模块取代，所以这个选项并不生效。 |
| PASS_WARN_AGE   | 7               | 指定在密码到期前多少天，系统就开始通过用户密码即将到期，默认为 7 天。 |
| UID_MIN         | 500             | 指定最小 UID 为 500，也就是说，添加用户时，默认 UID 从 500 开始。注意，如果手工指定了一个用户的 UID 是 550，那么下一个创建的用户的 UID 就会从 551 开始，哪怕 500~549 之间的 UID 没有使用。 |
| UID_MAX         | 60000           | 指定用户最大的 UID 为 60000。                                |
| GID_MIN         | 500             | 指定最小 GID 为 500，也就是在添加组时，组的 GID 从 500 开始。 |
| GID_MAX         | 60000           | 用户 GID 最大为 60000。                                      |
| CREATE_HOME     | yes             | 指定在创建用户时，是否同时创建用户主目录，yes 表示创建，no 则不创建，默认是 yes。 |
| UMASK           | 077             | 用户主目录的权限默认设置为 077。                             |
| USERGROUPS_ENAB | yes             | 指定删除用户的时候是否同时删除用户组，准备地说，这里指的是删除用户的初始组，此项的默认值为 yes。 |
| ENCRYPT_METHOD  | SHA512          | 指定用户密码采用的加密规则，默认采用 SHA512，这是新的密码加密模式，原先的 Linux 只能用 DES 或 MD5 加密。 |

# 二、用户管理

## 2.1 创建用户（useradd）

可以通过useradd来创建命名，创建过程就是向我们前面提到的文件里写入相应信息。来看示例，创建用户名：`peng`

```
$ useradd peng
```

创建`test`用户：指定UID=1024，主目录=/data/test，初始组=peng，附加组=root，bash为/bin/bash

```
$ useradd -u 1024 -d /data/test -g peng -G root -c "this is a test account" -s /bin/bash test
```

查看前面提到的文件：

```
$ cat /etc/passwd
peng:x:1000:1000::/home/peng:/bin/bash
test:x:1024:1000:this is a test account:/data/test:/bin/bash

$ cat /etc/group
root:x:0:test
peng:x:1000:

$ cat /etc/shadow
peng:!!:18547:0:99999:7:::
test:!!:18547:0:99999:7:::

$ id peng
uid=1000(peng) gid=1000(peng) 组=1000(peng)
$ id test
uid=1024(test) gid=1000(peng) 组=1000(peng),0(root)
```

也就是我们可以什么都不指定的创建一个用户，也可以完全按自己的意愿来创建用户。`useradd`命令在添加用户时参考的默认值文件主要有两个，分别是`/etc/default/useradd`和`/etc/login.defs`。

```
cat /etc/default/useradd
# useradd defaults file
GROUP=100
HOME=/home
INACTIVE=-1
EXPIRE=
SHELL=/bin/bash
SKEL=/etc/skel
CREATE_MAIL_SPOOL=yes
```

完整的流程是：

- 读取`/etc/login.defs`和`/etc/default/useradd`配置文件
- 在`/etc/passwd`中添加一条用户记录
- 在`/etc/shadow`文件中新增一条用户密码记录
- 在`/etc/group`文件中新增一条用户组记录
- 在`/etc/gshadow`文件中新增组密码记录
- 创建用户的主目录和邮箱
- 将`/etc/skel`目录中的配置文件赋值到新用户的主目录中

**注：`CentOS`下`useradd`命令和`adduser`命令没有区别。网上说的差别应该是`Ubuntu`系统下的差异。**

## 2.2 设置密码（passwd）

```
$ passwd peng
```

更多操作：

```
passwd --help
用法: passwd [选项...] <帐号名称>
  -k, --keep-tokens       保持身份验证令牌不过期
  -d, --delete            删除已命名帐号的密码(只有根用户才能进行此操作)
  -l, --lock              锁定指名帐户的密码(仅限 root 用户)
  -u, --unlock            解锁指名账户的密码(仅限 root 用户)
  -e, --expire            终止指名帐户的密码(仅限 root 用户)
  -f, --force             强制执行操作
  -x, --maximum=DAYS      密码的最长有效时限(只有根用户才能进行此操作)
  -n, --minimum=DAYS      密码的最短有效时限(只有根用户才能进行此操作)
  -w, --warning=DAYS      在密码过期前多少天开始提醒用户(只有根用户才能进行此操作)
  -i, --inactive=DAYS     当密码过期后经过多少天该帐号会被禁用(只有根用户才能进行此操作)
  -S, --status            报告已命名帐号的密码状态(只有根用户才能进行此操作)
  --stdin                 从标准输入读取令牌(只有根用户才能进行此操作)
```

比如，临时锁定用户：

```
$ passwd -l peng
锁定用户 peng 的密码 。
passwd: 操作成功

$ ssh peng@peng-node-1
peng@172.16.196.201's password:
Permission denied, please try again.
```

查看密码信息：

```
$ passwd -S peng
peng PS 2020-10-12 0 99999 7 -1 (密码已设置，使用 SHA512 算法。)
```

## 2.3 修改用户（usermod）

创建用户之后若想修改信息可以直接修改上面提供的文件，也可以通过`usermod`命令来更改。比如，

**增加用户组：**

```
# usermod -G [groupname] [username]，将用户 peng 添加进root组
$ usermod -G root peng

$ id peng
uid=1000(peng) gid=1000(peng) 组=1000(peng),0(root)
```

注意：多次操作修改会覆盖附加组的数据。

**锁定用户**

```
$ usermod -L peng
```

同上面`passwd`操作，就是在`/etc/passwd`文件的密码字段前增加了一个`!`，解锁操作则会去掉`!`

```
$ usermode -U peng
```

其他的

- `-d`: 指定新的主目录
- `-g`: 修改初始组
- `-l`:修改用户名
- `-s`: 修改登录的Shell

## 2.4 强制修改密码（chage）

```
$ chage -d 0 peng
```

通过`chage -d`命令更改密码文件中的最后一次修改时间，从而使密码失效。让用户下次登录时强制需要修改密码。

```
$ ssh peng@peng-node-1
peng@172.16.196.201's password:
You are required to change your password immediately (root enforced)
Last login: Mon Oct 12 10:29:52 2020 from 172.16.196.1
WARNING: Your password has expired.
You must change your password now and login again!
更改用户 peng 的密码 。
为 peng 更改 STRESS 密码。
（当前）UNIX 密码：
新的 密码：
重新输入新的 密码：
passwd：所有的身份验证令牌已经成功更新。
```

同样，也可以通过`chage -E`命令修改密码文件中的第8个字段，指定密码的失效日期，实现给用户临时授权的目的。`chage -l` 查看账户年龄信息：

```
$ chage -l root
最近一次密码修改时间					：从不
密码过期时间					：从不
密码失效时间					：从不
帐户过期时间						：从不
两次改变密码之间相距的最小天数		：0
两次改变密码之间相距的最大天数		：99999
在密码过期之前警告的天数	：7
```

强制修改密码的功能也可以通过:

```
$ passwd -e peng
```

## 2.5 删除用户（userdel）

```
$ userdel peng
```

只删除用户，会保留主目录，如果需要删除主目录，则可以指定`-r`参数

```
$ userdel -r peng
```

## 2.6 切换用户（su）

`su`可以实现不同登录用户之间的切换。

```
$ su peng
```

**su 和 su - 的区别**

注意，使用 su 命令时，有 - 和没有 - 是完全不同的，- 选项表示在切换用户身份的同时，连当前使用的环境变量也切换成指定用户的。我们知道，环境变量是用来定义操作系统环境的，因此如果系统环境没有随用户身份切换，很多命令无法正确执行。

# 三、用户组管理

## 3.1 添加组（groupadd）

```
$ groupadd group
```

## 3.2 修改组（groupmod）

```
$ groupmod -n newgroup group
```

修改组名称，对应的gid不变。

## 3.3 删除组（groupdel）

```
$ groupdel newgroup
```

## 3.4 加入或者移除用户组（gpasswd）

**将用户添加进组：**

```
$ gpasswd -a peng aaa
正在将用户“peng”加入到“aaa”组中

$ gpasswd -a peng bbb
正在将用户“peng”加入到“bbb”组中

$ id peng
uid=1000(peng) gid=1000(peng) 组=1000(peng),0(root),1027(aaa),1028(bbb)
```

与`usermod -G`不同的是，该命令是追加的方式，不会覆盖，而`usermod -G`会清空之前的组。

**将用户移除出组**

```
$ gpasswd -d peng aaa
正在将用户“peng”从“aaa”组中删除
```

**设置进群组管理员**

```
$ gpasswd -A peng groupname
```

群组管理员可以将用户加入组或者移除出组，具有一部分管理功能。`gpasswd`具体用法可参考手册：

```
$ gpasswd -h
用法：gpasswd [选项] 组

选项：
  -a, --add USER                向组 GROUP 中添加用户 USER
  -d, --delete USER             从组 GROUP 中添加或删除用户
  -h, --help                    显示此帮助信息并推出
  -Q, --root CHROOT_DIR         要 chroot 进的目录
  -r, --delete-password         remove the GROUP's password
  -R, --restrict                向其成员限制访问组 GROUP
  -M, --members USER,...        设置组 GROUP 的成员列表
  -A, --administrators ADMIN,...	设置组的管理员列表
除非使用 -A 或 -M 选项，不能结合使用这些选项。
```





---

- [1] [Linux用户和用户组管理详解](http://c.biancheng.net/linux_tutorial/60/)
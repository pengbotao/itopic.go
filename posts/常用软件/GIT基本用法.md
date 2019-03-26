```
{
    "url": "git",
    "time": "2016/11/01 18:25",
    "tag": "常用软件"
}
```

# 一、概述

## 1.1 Git的诞生

很多人都知道，Linus在1991年创建了开源的Linux，从此，Linux系统不断发展，已经成为最大的服务器系统软件了。

Linus虽然创建了Linux，但Linux的壮大是靠全世界热心的志愿者参与的，这么多人在世界各地为Linux编写代码，那Linux的代码是如何管理的呢？

事实是，在2002年以前，世界各地的志愿者把源代码文件通过diff的方式发给Linus，然后由Linus本人通过手工方式合并代码！

你也许会想，为什么Linus不把Linux代码放到版本控制系统里呢？不是有CVS、SVN这些免费的版本控制系统吗？因为Linus坚定地反对CVS和SVN，这些集中式的版本控制系统不但速度慢，而且必须联网才能使用。有一些商用的版本控制系统，虽然比CVS、SVN好用，但那是付费的，和Linux的开源精神不符。

不过，到了2002年，Linux系统已经发展了十年了，代码库之大让Linus很难继续通过手工方式管理了，社区的弟兄们也对这种方式表达了强烈不满，于是Linus选择了一个商业的版本控制系统BitKeeper，BitKeeper的东家BitMover公司出于人道主义精神，授权Linux社区免费使用这个版本控制系统。

安定团结的大好局面在2005年就被打破了，原因是Linux社区牛人聚集，不免沾染了一些梁山好汉的江湖习气。开发Samba的Andrew试图破解BitKeeper的协议（这么干的其实也不只他一个），被BitMover公司发现了（监控工作做得不错！），于是BitMover公司怒了，要收回Linux社区的免费使用权。

Linus可以向BitMover公司道个歉，保证以后严格管教弟兄们，嗯，这是不可能的。实际情况是这样的：

Linus花了两周时间自己用C写了一个分布式版本控制系统，这就是Git！一个月之内，Linux系统的源码已经由Git管理了！牛是怎么定义的呢？大家可以体会一下。

Git迅速成为最流行的分布式版本控制系统，尤其是2008年，GitHub网站上线了，它为开源项目免费提供Git存储，无数开源项目开始迁移至GitHub，包括jQuery，PHP，Ruby等等。

历史就是这么偶然，如果不是当年BitMover公司威胁Linux社区，可能现在我们就没有免费而超级好用的Git了。

— 来自[廖雪峰的官方网站](http://www.liaoxuefeng.com/)

## 1.2 Git和SVN对比

Git有一套本地仓库，推送到远程之前可以提交到本地，确保开发过程中的功能均能纳入版本管理。而SVN使用上通常不上线不提交到SVN，开发过程中没有版本管理。

Git的分支相比SVN更方便。使用SVN时很多时候通过拷贝代码来实现一些功能，而Git通过分支可以快速的实现。

Git的分支功能更强大，合理的使用分支可以摸索出更适合团队的工作流。

## 1.3 Git安装及图形化管理工具

- 安装Git命令行 https://git-scm.com/downloads
- 图形化工具
  - Win:TortoiseGit等
  - Mac:SourceTree等

# 二、Git基本用法

## 2.1 配置公钥

Git可以直接在本地创建仓库进行提交、回滚等操作，但如果需要跟其他人交换代码则需要连接到远程仓库，连接远程仓库一般需要账号或者公私钥验证，这里推荐以公钥的方式进行验证。首先需要在本地生成一对公私钥。一路回车后即可看到生成的秘钥地址（如果使用的非默认地址，之后使用Git的时候需要进行配置秘钥地址）。

```
$ ssh-keygen
$ ls ~/.ssh/
id_rsa		id_rsa.pub
```

`id_rsa`为私钥留在本地，`id_rsa.pub`为公钥，提供给github、gitlab等平台进行验证。这里以gitlab为例，登录之后点击账户头像 -> settings -> SSH Keys. 将公钥内容贴在里面保存即可。

![](./static/uploads/gitlab-ssh-keys.png)

## 2.2 仓库初始化

初始化一个仓库很简单，两种常用的方式：

### 2.2.1 从远程仓库克隆
将一个已存在的仓库克隆到本地，大多数情况下我们进行的也是这种操作。

```
git clone git@hostname:group/project.git
```

按照上面的地址会在本地创建一个project目录。如果想指定目录名可以在后面指定

```
# 克隆在当前目录
$ git clone git@hostname:group/project.git .
# 克隆到指定目录名
$ git clone git@hostname:group/project.git project_alias
```

### 2.2.2 初始化本地仓库
当远程仓库还未创建时，我们可以先在本地初始化仓库，在本地进行开发。这种方式唯一影响的功能就是将代码推送到远程仓库。其他的功能不受影响。

```
$ mkdir gitdemo
$ cd gitdemo
$ git init
Initialized empty Git repository in /Users/peng/workspace/gitdemo/.git/
```

这样子就创建了一个本地仓库，可以在本地仓库进行代码提交等操作，后续远程仓库创建后可以在关联远程仓库

```
git remote add origin git@hostname:group/project.git
```

## 2.3 GIT常用配置

### 2.3.1 配置用户名和邮箱

```
git config --global user.name "peng"
git config --global user.email "peng@domain.cn"
```

### 2.3.2 关闭crlf自动转换

```
git config --global core.autocrlf false
```

不同操作系统下的换行符不同，开启autocrlf会存在提交、检出自动转换换行，这就会导致对比的时候全是变化部分。所以关闭此配置同时约定文件的换行符统一采用Linux下的`\n`.

## 2.4 Git 基本流程

![](./static/uploads/git-basic-command.png)

上图展示了最常用到的几个基本Git命令。这里按照上图走一下基本流程。

### 2.4.1 克隆远程仓库

```
$ git clone git@git.hostname.net:peng/gitdemo.git
Cloning into 'gitdemo'...
warning: You appear to have cloned an empty repository.
$ cd gitdemo
```

### 2.4.2 本地开发

我们创建了一个`README.md`文件,接着通过`git status`查看下状态发现`README.md`文件没有被跟踪，同时在括号中标注`use "git add <file>..." to include in what will be committed`。此时默认是`On branch master`

```
$ touch README.md
$ git status
On branch master

Initial commit

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	README.md

nothing added to commit but untracked files present (use "git add" to track)
```

### 2.4.3 提交到暂存区

添加到暂存区，并查看状态。同时可以根据提示`use "git rm --cached <file>..." to unstage`从暂存区移除。

```
$ git add README.md
$ git status
On branch master

Initial commit

Changes to be committed:
  (use "git rm --cached <file>..." to unstage)

	new file:   README.md
```

### 2.4.4 提交到本地仓库

提交到本地仓库，如果没有配置账号、Email则会有相应提示信息，按照提示信息设置账号和Email信息。`-m`为提交注释。

```
$ git commit -m "add README"

*** Please tell me who you are.

Run

  git config --global user.email "you@example.com"
  git config --global user.name "Your Name"

to set your account's default identity.
Omit --global to set the identity only in this repository.

fatal: unable to auto-detect email address (got 'peng@pengbotao.(none)')
```

再次执行提交操作，此过程会提交数据到本地仓库。

```
$ git commit -m "add README"
[master (root-commit) bc6094e] add README
 1 file changed, 0 insertions(+), 0 deletions(-)
 create mode 100644 README.md
```

### 2.4.5 推送到远程仓库

origin为默认关联的远程地址的名称，通过`git remote -v`可以看到本地仓库已关联的远程仓库。同一个仓库可以关联多个远程地址，这里表示要推送到origin的master分支上。

```
$ git push origin master
Counting objects: 3, done.
Writing objects: 100% (3/3), 209 bytes | 0 bytes/s, done.
Total 3 (delta 0), reused 0 (delta 0)
To git.coding.net:pbt/gitdemo.git
 * [new branch]      master -> master
```

### 2.4.6 从远程仓库拉取

```
$ git pull origin master
Already up-to-date.
```

如果`push`和`pull`过程中存在和远程分支未关联的情况，按提示信息做一下关联即可。同时关联上之后可以在当前分支的推送、拉取可以直接使用`git push/pull`命令。

```
$ git push
fatal: The current branch dev has no upstream branch.
To push the current branch and set the remote as upstream, use

    git push --set-upstream origin dev
```

---

基本的操作就是重复`2.4.2`-`2.4.6`的过程。同时Git也给了足够多的提示信息，如果有相关报错信息仔细查看提示信息可以解决大部分问题。

一般本地开发会频繁执行`2.4.2`-`2.4.4`，整个过程都在本地，可以及时把新加入的功能加入到版本库。这也是`Git`和`SVN`区别比较大的一个地方。还有一个主要区别在于`Git`分支以及分支引申出来的工作流，接来下说一下分支的用法。

# 三、Git分支用法

接上面，目前处与`master`分支，修改`README.md`后并重复上面步骤进行提交。假设此时的内容为`This is master branch Demo.`，从`master`分支创建一个新的`dev`分支，修改`README.md`文件。

```
$ git checkout -b dev
Switched to a new branch 'dev'
$ vi README.md

This is dev branch Demo.
```

提交后切回`master`分支并查看文件内容

```
$ git checkout master
Switched to branch 'master'
Your branch is up-to-date with 'origin/master'.
pengbotao:gitdemo peng$ cat README.md
This is master branch Demo.

$ git checkout dev
Switched to branch 'dev'
Your branch is up-to-date with 'origin/dev'.
$ cat README.md
This is dev branch Demo.
```

来回切换发现文本内容实时就变了，从磁盘上看文件在同一个目录，同一个文件名但内容不同了。这也是`Git`和`SVN`差别比较大的一个地方，`Git`分支的不同代码切换在相同的文件目录中。以前使用`SVN`一定有拷贝代码的情况，理解`Git`的这点后只需要创建个分支并在分支上修改即可。

## 3.1 分支基本操作

分支可以先创建在切换，后者创建后直接切换。

```
# 创建一个新的分支
$ git branch bugfix
# 查看当前所有分支，当前还是dev分支
$ git branch -a
  bugfix
* dev
  master
  remotes/origin/dev
  remotes/origin/master

# 切换到bugfix分支
$ git checkout bugfix
Switched to branch 'bugfix'

# 创建feature并切换到feature分支
$ git checkout -b feature
Switched to a new branch 'feature'
pengbotao:gitdemo peng$ git branch -v
  bugfix  bb0b9db dev readme
  dev     bb0b9db dev readme
* feature bb0b9db dev readme
  master  bf06765 master readme modify

# 删除分支
$ git branch -d feature
Deleted branch feature (was bb0b9db)
```

## 3.2 分支合并

以下过程描述bugfix分支修改后合并到dev分支，这里测试直接用dev合并了bugfix分支。实际情况中可能dev分支已经被其他人修改，这个时候需要先更新dev分支在做合并。通常的做法是在bugfix分支上先合并dev，处理冲突等，处理完成之后在切回到dev合并bugfix。

```
# 修改bugfix分支README.md文件(需要先切换到bugfix分支)
$ git commit -am "bugfix"
[bugfix 164795f] bugfix
 1 file changed, 2 insertions(+)

# 切换到dev分支并合并bugfix分支
$ git checkout dev
Switched to branch 'dev'
Your branch is up-to-date with 'origin/dev'.
$ git merge bugfix
Updating bb0b9db..164795f
Fast-forward
 README.md | 2 ++
 1 file changed, 2 insertions(+)
 
# 查看文件内容
$ cat README.md
This is dev branch Demo.
write from bugfix branch
```

## 3.3 Git工作流

`Git`可以创建不同的分支，每个分支的功能可以定义为不同，比如常见的`master`分支当做发布分支，`develop`分支当开发分支，`hotfix-x`当补丁分支，`feature-x`当新功能分支。根据`Git`灵活的分支策略可以搭配出很多种工作流，这里简单减少一下几种常见的工作流。

### 3.3.1 中心化工作流

只有默认的master分支，代码编写提交都在此分支上，不需要master之外的其他分支，提交处理方式跟SVN类似。

### 3.3.2  基于功能分支的工作流

每个新功能创建一个分支，功能开发完成后合并到master分支。

### 3.3.3 Gitflow工作流

![](./static/uploads/git-model@2x.png)

### 3.3.4 Fork工作流

每个人fork一份仓库为自己的仓库，提交到自己的仓库，要发布的时候需要发送pull request，需要特定人员去处理合并。

# 四、其他用法

### 4.1 工作区状态

`git status`查看当前工作区状态，开发过程中可能多次用到。

```
$ git status
On branch master
Your branch is up-to-date with 'origin/master'.
nothing to commit, working tree clean
```

### 4.2 远程仓库相关操作

远程仓库地址的管理，更多操作可查看`git remote --help`

```
$ git remote
github
origin
oschina

# 查看远程仓库URL
$ git remote -v
github	git@github.com:peng/gitdemo.git (fetch)
github	git@github.com:peng/gitdemo.git (push)
origin	git@git.hostname.net:peng/gitdemo.git (fetch)
origin	git@git.hostname.net:peng/gitdemo.git (push)

# 关联远程仓库
$ git remote add github git@github.com:peng/gitdemo.git

# 移除远程仓库
$ git remote remove github
```

### 4.3 查看日志

```
$ git log

commit 164795fb2aa7785fcb0fe16bbf6f2ef092a0f934
Author: pengbotao <12345@qq.com>
Date:   Thu Jun 29 15:41:01 2017 +0800

    bugfix

commit bb0b9dbed62c3fd60f9a21eff7fcd660078990b2
Author: pengbotao <12345@qq.com>
Date:   Thu Jun 29 14:14:10 2017 +0800

    dev readme
```

## 4.4 版本回滚

### 4.4.1 文件添加到暂存区

文件添加到暂存区，但是没有提交到本地仓库。通过`git status`可以看到`unstage`操作命令。

```
$ git status
On branch dev
Your branch is ahead of 'origin/dev' by 1 commit.
  (use "git push" to publish your local commits)
Changes to be committed:
  (use "git reset HEAD <file>..." to unstage)

	modified:   test.php
```

通过`git reset HEAD <file>`从暂存区移除。通过`git reset HEAD`移除所有添加到暂存区的文件。

### 4.4.2 还原本地修改的文件

只在本地做了修改，想丢弃这些修改。

```
git checkout -- filename
```

如果是本地修改了很多文件，想全部恢复可使用

```
git checkout -- .
```

### 4.4.3 还原已提交到本地的文件

文件已经添加并提交到本地仓库，但是未推送到远程仓库。可通过`git log`查看`commit id`,reset后

```
git reset --hard 77e02cd8dc43940f0c817379c94935a79d476510
```

### 4.4.4 还原已经提交到远程的文件

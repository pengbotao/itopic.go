```
{
    "url": "python-environment",
    "time": "2015/12/21 19:51",
    "tag": "Python"
}
```

# 包管理工具 - pip

**pip安装：**

```
# https://pip.pypa.io/en/stable/installing/
$ curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
$ python get-pip.py
```

**方法列表**

```
$ pip

Usage:
  pip <command> [options]

Commands:
  install                     Install packages.
  download                    Download packages.
  uninstall                   Uninstall packages.
  freeze                      Output installed packages in requirements format.
  list                        List installed packages.
  show                        Show information about installed packages.
  check                       Verify installed packages have compatible dependencies.
  config                      Manage local and global configuration.
  search                      Search PyPI for packages.
  wheel                       Build wheels from your requirements.
  hash                        Compute hashes of package archives.
  completion                  A helper command used for command completion.
  help                        Show help for commands.
```

**使用示例**

```
# 安装指定包及版本
$ pip install numpy==x.xx.x
# 导出包列表到文件
$ pip freeze > requirements.txt
# 从文件安装
$ pip install -r requirements.txt
```

# 多版本管理 - pyenv

`pyenv`主要用来解决`Python`多版本共存的问题。项目地址：`https://github.com/pyenv/pyenv`

**1. 安装方式**

`github`上有安装说明，按安装说明操作即可。

```
$ git clone https://github.com/pyenv/pyenv.git ~/.pyenv
$ echo 'export PYENV_ROOT="$HOME/.pyenv"' >> ~/.bash_profile
$ echo 'export PATH="$PYENV_ROOT/bin:$PATH"' >> ~/.bash_profile
$ echo -e 'if command -v pyenv 1>/dev/null 2>&1; then\n  eval "$(pyenv init -)"\nfi' >> ~/.bash_profile
$ source ~/.bash_profile

# 显示已经安装的版本
$ pyenv versions
* system (set by /Users/peng/.pyenv/version)

# 安装python3.7.2版本
$ pyenv install 3.7.2
python-build: use openssl from homebrew
python-build: use readline from homebrew
Downloading Python-3.7.2.tar.xz...
-> https://www.python.org/ftp/python/3.7.2/Python-3.7.2.tar.xz


# 卸载版本
$ pyenv uninstall 2.7.1
```

安装错误: `zipimport.ZipImportError: can't decompress data; zlib not available`

解决方法：

```
$ xcode-select --install
$ brew install zlib
```

安装错误：`ERROR: The Python ssl extension was not compiled. Missing the OpenSSL lib?`

```
# MAC下升级openssl才得以解决
$ brew upgrade openssl
$ CFLAGS="-I$(brew --prefix openssl)/include -I$(xcrun --show-sdk-path)/usr/include" CPPFLAGS="-I$(brew --prefix openssl)/include" LDFLAGS="-L$(brew --prefix openssl)/lib" pyenv install -v 3.7.2

# CentOS
$ yum install -y zlib-devel bzip2 bzip2-devel readline-devel sqlite sqlite-devel openssl-devel xz xz-devel libffi-devel findutils
$ CFLAGS=-I/usr/include/openssl \
LDFLAGS=-L/usr/lib64 \
pyenv install -v 3.7.2
```

安装问题可查询：`https://github.com/pyenv/pyenv/wiki/Common-build-problems`


下载较慢可直接下载pip提示的文件地址：`https://www.python.org/ftp/python/3.7.2/Python-3.7.2.tar.xz`，存储到 `~/.pyenv/cache`目录后再执行上面安装命令。

```
$ pyenv versions
* system (set by /Users/peng/.pyenv/version)
  3.7.2
```


**2. pyenv使用**

```
$ pyenv
pyenv 1.2.9-16-g9baa6efe
Usage: pyenv <command> [<args>]

Some useful pyenv commands are:
```

命令|说明
---|---
commands    |List all available pyenv commands
local       |Set or show the local application-specific Python version
global      |Set or show the global Python version
shell       |Set or show the shell-specific Python version
install     |Install a Python version using python-build
uninstall   |Uninstall a specific Python version
rehash      |Rehash pyenv shims (run this after installing executables)
version     |Show the current Python version and its origin
versions    |List all Python versions available to pyenv
which       |Display the full path to an executable
whence      |List all Python versions that contain the given executable

**操作示例**

```
# 全局切换Python版本
$ pyenv global 2.7.1

# 切回系统版本
$ pyenv global system

# 当前目录及子目录切换版本
$ pyenv local 2.7.1

# 切回系统版本
$ pyenv local system
```


# 虚拟环境 - virtualenv
使用`pyenv`或可以解决多个版本的问题，但具体的特定项目下依赖包管理不太方便，使用`virtualenv`可以在已有的版本基础上创建一个新的虚拟版本，此版本下安装的包都在该环境下，方便项目层级的版本和包管理，配合`pip freeze`也可以方便导出依赖包列表。项目地址：`https://github.com/pyenv/pyenv-virtualenv`

这里是按照说明文档以`pyenv`插件的方式安装：

```
$ git clone https://github.com/pyenv/pyenv-virtualenv.git $(pyenv root)/plugins/pyenv-virtualenv
$ echo 'eval "$(pyenv virtualenv-init -)"' >> ~/.bash_profile

$ pyenv virtualenv --version
pyenv-virtualenv 1.1.5 (virtualenv unknown)
```

虚拟一个新的环境并在指定目录下使用。

```
$ pyenv virtualenv 3.7.2 my372
$ pyenv local my372
(my372) pengbotao:py3 peng$
```

删除虚拟出来的环境。

```
$ pyenv virtualenv-delete my372
pyenv-virtualenv: remove /Users/peng/.pyenv/versions/3.7.2/envs/my372? y
```

pyenv可参考：`https://blog.csdn.net/u010104435/article/details/79633067`


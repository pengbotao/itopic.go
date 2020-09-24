```
{
    "url": "cmake",
    "time": "2015/04/30 18:24",
    "tag": "C++",
    "toc": "no"
}
```

# CMake简介

CMake 是一个跨平台的自动化建构系统,它使用一个名为 CMakeLists.txt 的文件来描述构建过程,可以产生标准的构建文件,如 Unix 的 Makefile 或Windows Visual C++ 的 projects/workspaces 。文件 CMakeLists.txt 需要手工编写,也可以通过编写脚本进行半自动的生成。CMake 提供了比 autoconfig 更简洁的语法。在 linux 平台下使用 CMake 生成 Makefile 并编译的流程如下:

- 1、编写 CmakeLists.txt。
- 2、执行命令“cmake PATH”或者“ccmake PATH”生成 Makefile ( PATH 是 CMakeLists.txt 所在的目录 )。
- 3、使用 make 命令进行编译。

# 安装CMAKE
```
# wget http://www.cmake.org/files/v3.2/cmake-3.2.2.tar.gz
# tar zxvf cmake-3.2.2.tar.gz 
# cd cmake-3.2.2/
# ./bootstrap 
# make && make install
```
接下来通过cmake来构建一个简单的HelloWorld项目。

# Hello World
```
$ more main.cpp 
#include <iostream>
 
int main(int argc, char **argv) {
    std::cout << "Hello, world!" << std::endl;
    return 0;
}
```
接下来按照流程第一步编写CMakeLists.txt文件，下面为Kdevelop4.6创建project后自动生成的CMakeLists.txt。这里CMakeLists.txt和main.cpp在同一目录。
```
$ more CMakeLists.txt 
cmake_minimum_required(VERSION 2.6)
project(main)
 
add_executable(main main.cpp)
 
install(TARGETS main RUNTIME DESTINATION bin)
```
第二步，执行命令“cmake PATH”生成Makefile。
```
$ cmake .
-- The C compiler identification is GNU 4.8.2
-- The CXX compiler identification is GNU 4.8.2
-- Check for working C compiler: /usr/bin/cc
-- Check for working C compiler: /usr/bin/cc -- works
-- Detecting C compiler ABI info
-- Detecting C compiler ABI info - done
-- Detecting C compile features
-- Detecting C compile features - done
-- Check for working CXX compiler: /usr/bin/c++
-- Check for working CXX compiler: /usr/bin/c++ -- works
-- Detecting CXX compiler ABI info
-- Detecting CXX compiler ABI info - done
-- Detecting CXX compile features
-- Detecting CXX compile features - done
-- Configuring done
-- Generating done
-- Build files have been written to: /home/peng/projects/main
```
第三步，使用 make 命令进行编译，即可生成可执行程序main，直接执行main程序即可看到Hello,world!
```
$ make
Scanning dependencies of target main
[100%] Building CXX object CMakeFiles/main.dir/main.o
Linking CXX executable main
[100%] Built target main
```
最后CMakeLists也指定了install，也可执行make install操作，可将main程序安装到指定目录。
```
$ sudo make install
[100%] Built target main
Install the project...
-- Install configuration: ""
-- Installing: /usr/local/bin/main
```
至此，一个项目的编译安装就完成了。当然，也可以直接通过G++来编译，而且单文件时看起来似乎更简单，一句话就可以搞定。
```
g++ -o main main.cpp
```

**接下来看看CMake相关知识点**

# cmake基本语法规则

- 1，变量使用${}方式取值，但是在 IF 控制语句中是直接使用变量名
- 2，指令(参数 1 参数 2...)
参数使用括弧括起，参数之间使用空格或分号分开。以上面的 ADD_EXECUTABLE 指令为例，如果存在另外一个 func.cpp 源文件，就要写成：
```
ADD_EXECUTABLE(hello main.cpp func.cpp)
```
或者
```
ADD_EXECUTABLE(hello main.cpp;func.cpp)
```
- 3，指令是大小写无关的，参数和变量是大小写相关的。但推荐你全部使用大写指令。

上面的示例未内部构建（in-source build），而cmake强烈推荐的是外部构建（out-of-source build）。

# 内部构建与外部构建

在cmake中有两个预定义变量：`projectname_BINARY_DIR`以及`projectname_SOURCE_DIR`。在我们的项目中，两个变量分别为：main_BINARY_DIR和main_SOURCE_DIR。

同时cmake还预定义了PROJECT_BINARY_DIR和PROJECT_SOURCE_DIR变量。在前面项目中，PROJECT_BINARY_DIR等同于main_BINARY_DIR，PROJECT_SOURCE_DIR等同于main_SOURCE_DIR。但在有些情况下会有所不同。执行cmake有两种方式：
```
cmake .
make
```
或者
```
mkdir build
cd build
cmake ..
make
```
两种方法最大的不同在于执行cmake和make的工作路径不同。第一种方法中，cmake生成的所有中间文件和可执行文件都会存放在项目目录中；而第二种方法中，中间文件和可执行文件都将存放在build目录中。第二种方法的优点显而易见，它最大限度的保持了代码目录的整洁。同时由于第二种方法的生成、编译和安装是发生在不同于项目目录的其他目录中，所以第二种方法就叫做“外部构建”。在外部构建的情况下，PROJECT_SOURCE_DIR指向的目录同内部构建相同，仍然为~/main，而PROJECT_BINARY_DIR则有所不同，指向~/main/build目录。当然，cmake强烈推荐使用外部构建的方法。同时在实际的应用中推荐使用PROJECT_BINARY_DIR和PROJECT_SOURCE_DIR变量，这样即使项目名称发生变化也不会影响CMakeLists.txt文件。

# 安装 - INSTALL指令
这里需要引入一个新的 cmake 指令 INSTALL 和一个非常有用的变量 CMAKE_INSTALL_PREFIX。CMAKE_INSTALL_PREFIX 变量类似于 configure 脚本的 –prefix，常见的使用方法看起来是这个样子：
```
cmake -DCMAKE_INSTALL_PREFIX=/usr .
```
INSTALL 指令用于定义安装规则，安装的内容可以包括目标二进制、动态库、静态库以及文件、目录、脚本等。INSTALL 指令包含了各种安装类型，我们需要一个个分开解释：

**目标文件的安装：**
```
INSTALL(TARGETS targets...
    [[ARCHIVE|LIBRARY|RUNTIME]
        [DESTINATION <dir>]
        [PERMISSIONS permissions...]
        [CONFIGURATIONS
    [Debug|Release|...]]
        [COMPONENT <component>]
        [OPTIONAL]
       ] [...])
```
参数中的 TARGETS 后面跟的就是我们通过 ADD_EXECUTABLE 或者 ADD_LIBRARY 定义的目标文件，可能是可执行二进制、动态库、静态库。

目标类型也就相对应的有三种，ARCHIVE 特指静态库，LIBRARY 特指动态库，RUNTIME特指可执行目标二进制。

DESTINATION 定义了安装的路径，如果路径以/开头，那么指的是绝对路径，这时候CMAKE_INSTALL_PREFIX 其实就无效了。如果你希望使用 CMAKE_INSTALL_PREFIX 来定义安装路径，就要写成相对路径，即不要以/开头，那么安装后的路径就是
${CMAKE_INSTALL_PREFIX}/

举个简单的例子：
```
INSTALL(TARGETS myrun mylib mystaticlib
    RUNTIME DESTINATION bin
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION libstatic
)
```
上面的例子会将：

- 可执行二进制 myrun 安装到${CMAKE_INSTALL_PREFIX}/bin 目录
- 动态库 libmylib 安装到${CMAKE_INSTALL_PREFIX}/lib 目录
- 静态库 libmystaticlib 安装到${CMAKE_INSTALL_PREFIX}/libstatic 目录
- 特别注意的是你不需要关心 TARGETS 具体生成的路径，只需要写上 TARGETS 名称就可以了。

**普通文件的安装：**
```
INSTALL(FILES files... DESTINATION <dir>
    [PERMISSIONS permissions...]
    [CONFIGURATIONS [Debug|Release|...]]
    [COMPONENT <component>]
    [RENAME <name>] [OPTIONAL])
```
可用于安装一般文件，并可以指定访问权限，文件名是此指令所在路径下的相对路径。如果默认不定义权限 PERMISSIONS，安装后的权限为：OWNER_WRITE, OWNER_READ, GROUP_READ,和 WORLD_READ，即 644 权限。

**非目标文件的可执行程序安装(比如脚本之类)：**
```
INSTALL(PROGRAMS files... DESTINATION <dir>
    [PERMISSIONS permissions...]
    [CONFIGURATIONS [Debug|Release|...]]
    [COMPONENT <component>]
    [RENAME <name>] [OPTIONAL])
```
跟上面的 FILES 指令使用方法一样，唯一的不同是安装后权限为:OWNER_EXECUTE, GROUP_EXECUTE, 和 WORLD_EXECUTE，即 755 权限

**目录的安装：**
```
INSTALL(DIRECTORY dirs... DESTINATION <dir>
    [FILE_PERMISSIONS permissions...]
    [DIRECTORY_PERMISSIONS permissions...]
    [USE_SOURCE_PERMISSIONS]
    [CONFIGURATIONS [Debug|Release|...]]
    [COMPONENT <component>]
    [[PATTERN <pattern> | REGEX <regex>]
    [EXCLUDE] [PERMISSIONS permissions...]] [...])
```
这里主要介绍其中的 DIRECTORY、PATTERN 以及 PERMISSIONS 参数。

DIRECTORY 后面连接的是所在 Source 目录的相对路径，但务必注意：abc 和 abc/有很大的区别。如果目录名不以/结尾，那么这个目录将被安装为目标路径下的 abc，如果目录名以/结尾，代表将这个目录中的内容安装到目标路径，但不包括这个目录本身。PATTERN 用于使用正则表达式进行过滤，PERMISSIONS 用于指定 PATTERN 过滤后的文件权限。

我们来看一个例子:
```
INSTALL(DIRECTORY icons scripts/ DESTINATION share/myproj
    PATTERN "CVS" EXCLUDE
    PATTERN "scripts/*"
    PERMISSIONS OWNER_EXECUTE OWNER_WRITE OWNER_READ GROUP_EXECUTE GROUP_READ)
```
这条指令的执行结果是：

将 icons 目录安装到 /share/myproj，将 scripts/中的内容安装到 /share/myproj 不包含目录名为 CVS 的目录，对于 scripts/*文件指定权限为 OWNER_EXECUTE OWNER_WRITE OWNER_READ GROUP_EXECUTE GROUP_READ. 

# 常用命令
参数|说明
---|---
cmake_minimum_required|指定cmake最低版本
project(main)|指定项目的名称。
add_executable(main main.cpp)|生成可执行程序main，可指定多个文件，空格分隔
add_subdirectory|ADD_SUBDIRECTORY(source_dir [binary_dir] [EXCLUDE_FROM_ALL])
add_library|add_library(Hello hello.cxx) #将hello.cxx编译成静态库如libHello.a
target_link_libraries|指定链接目标文件时需要链接的外部库，相同于gcc参数-l，可以解决外部库的依赖问题。
link_directories|动态链接库或静态链接库的搜索路径，相当于gcc的-L参数；
link_libraries |  
include_directories | 指定编译过程中编译器搜索头文件的路径，当项目不在系统默认的搜索路径时需要指定头文件的搜索路径，相当于指定gcc的-I参数；
add_definitions|添加编译参数

**message** message(SEND_ERROR\|STATUS\|FATAL_ERROR “message to display”)

SEND_ERROR，产生错误，生成过程被跳过。SATUS，输出前缀为--的信息。FATAL_ERROR，立即终止所有 cmake 过程。示例如，
```
message(STATUS ${PROJECT_BINARY_DIR})
message(STATUS ${PROJECT_SOURCE_DIR})
```
**set命令** SET(VAR [VALUE] [CACHE TYPE DOCSTRING [FORCE]])
SET 指令可以用来显式的定义变量即可。比如
```
set(SRC_LIST main.c)
set(CMAKE_INSTALL_PREFIX /usr/local)
```
如果有多个源文件，也可以定义成：
```
SET(SRC_LIST main.c t1.c t2.c)
```
url: linux-c-thread
des: 
time: 2015/05/25 17:52
category: cpp
++++++++

# 进程和线程

程序是一组用计算机语言编写的命令序列的集合。程序并不能单独运行，只有将程序装载到内存中，系统为它分配资源才能运行，而这种执行的「程序」就称之为进程。

**进程**（process）是指在系统中正在运行的一个应用程序，是系统资源分配的基本单位，在内存中有其完备的数据空间和代码空间，拥有完整的虚拟空间地址。一个进程所拥有的数据和变量只属于它自己。

**线程**（thread）是进程内相对独立的可执行单元，所以也被称为轻量进程（lightweight processes）；是操作系统进行任务调度的基本单元。它与父进程的其它线程共享该进程所拥有的全部代码空间和全局变量，但拥有独立的堆栈（即局部变量对于线程来说是私有的）。

启动一个新的进程必须分配给它独立的地址空间，建立众多的数据表来维护它的代码段、堆栈段和数据段，这是一种"昂贵"的多任务工作方式。而运行于一个进程中的多个线程，它们彼此之间使用相同的地址空间，共享大部分数据，启动一个线程所花费的空间远远小于启动一个进程所花费的空间，而且，线程间彼此切换所需的时间也远远小于进程间切换所需要的时间。
同时，对不同进程来说，它们拥有独立的数据空间，要进行数据的传递只能通过通信的方式进行，这种方式不仅费时，而且很不方便。线程则不然，由于同一进程下的线程之间共享数据空间，所以一个线程的数据可以直接为其它线程所用，这不仅快捷，而且方便。当然，数据的共享也带来其他 一些问题，有的变量不能同时被两个线程所修改，有的子程序中声明为static的数据更有可能给多线程程序带来灾难性的打击，这些正是编写多线程程序时最需要注意的地方。

# 线程的使用

Linux系统下的多线程遵循POSIX线程接口，称为pthread。编写Linux下的多线程程序，需要使用头文件pthread.h ，下面看一段小程序thread.cpp：
```
#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
 
int num = 0;
 
void *sum(void * arg)
{
    for(int i=0; i<3; i++) {
        num++;
    }
}
 
int main(int argc, char **argv)
{
    pthread_t tid;
    if(pthread_create(&tid, NULL, sum, NULL)) {
        printf("create thread error\n");
        exit(1);
    }
    pthread_join(tid, NULL);
    printf("%d\n", num);
    return 0;
}
```
编译
```
g++ -o thread thread.cpp -lpthread
```
这是一段非常简单的线程程序， 启动一个线程递增num的值。下面来看看这段程序， 首先声明了pthread_t类型的变量tid。通常称为线程ID，也可以认为他是一种线程句柄。然后调用pthread_create函数创建一个线程。创建成功时返回0，创建失败返回非0值。pthread_create函数中的第一个参数是指向tid的指针，第二个参数表示线程属性，demo上未设置，使用系统默认，第三个参数表示线程启动时调用的函数名称，该函数接受void *作为参数，同时返回void * 。这表示可以用void *向新线程传递任意类型的数据，也可以返回任意类型的数据。第四个参数表示需要向线程函数中传递的参数，如不需传递则传入NULL。

创建一个新的线程后程序将存在两个线程，主线程和新创建的线程。主线程按顺序继续执行下一行程序。本程序中下一行为pthread_join，pthread_join是一个线程阻塞的函数，调用它的函数将一直等待到被等待的线程结束为止，当函数返回时，被等待线程的资源被回收。所以上面的示例将会打印3，如果去掉phtread_join这一行将可能打印0，因为可能在线程函数执行之前，主程序就执行完了。

接下来，我们在主程序也对num进行下修改，为了方便看到效果，将递增的次数增多，如下例：
```
#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
 
int num = 0;
 
void *sum(void * arg)
{
    for(int i=0; i<10000; i++) {
        num++;
    }
}
 
int main(int argc, char **argv)
{
    pthread_t tid;
    if(pthread_create(&tid, NULL, sum, NULL)) {
        printf("create thread error\n");
        exit(1);
    }
    for(int i=0; i<10000; i++) {
        num++;
    }
    pthread_join(tid, NULL);
    printf("%d\n", num);
    return 0;
}
```
编译后执行，发现每次打印的内容都不一样，线程是并发运行的。当新线程执行递增操作同时，主线程也可能在执行递增操作，而且递增之前num的值是相同的，所以最后递增出来结果要比20000少。开发的时候需要注意这类问题，如果需要限制，也可以通过加锁的方式来实现，同一时刻确保只有一个线程能访问它，如：
```
#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
 
int num = 0;
pthread_mutex_t my_mutex=PTHREAD_MUTEX_INITIALIZER;
 
void *sum(void * arg)
{
    for(int i=0; i<10000; i++) {
        pthread_mutex_lock(&my_mutex);
        num++;
        pthread_mutex_unlock(&my_mutex);
    }
}
 
int main(int argc, char **argv)
{
    pthread_t tid;
    if(pthread_create(&tid, NULL, sum, NULL)) {
        printf("create thread error\n");
        exit(1);
    }
    for(int i=0; i<10000; i++) {
        pthread_mutex_lock(&my_mutex);
        num++;
        pthread_mutex_unlock(&my_mutex);
    }
    pthread_join(tid, NULL);
    printf("%d\n", num);
    return 0;
}
```
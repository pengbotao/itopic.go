```
{
    "url": "linux-c-chatroom-demo",
    "time": "2015/05/24 11:09",
    "tag": "C++,线程"
}
```

基于socket和线程实现一个简单的群聊功能。编译方式

```
g++ -o chat_server chat_server.cpp -lpthread
g++ -o chat_client chat_client.cpp -lpthread
```
**chat_server.cpp**

监听客户端请求，当新增客户端时启动一个新的线程来通信，接收到数据后广播给其他客户端。
```
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>  
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <pthread.h>
#include <time.h>
 
#define SERVER_PORT 8801
#define MAX_CONNECT 128
#define BUFFER_SIZE 2048
 
int chat_sockfd_idx = 0, chat_user_total = 1;
int chat_sockfd_list[MAX_CONNECT] = {0};
 
struct chat_item
{
    int sockfd;
    char name[32];
    char ip[15];
};
 
void *chat_login(void * args);
void chat_send(char *msg, int sockfd = 0);
void chat_exit(chat_item * ci);
 
 
int main(int argc, char **argv)
{
    int srv_sockfd, cli_sockfd, sock_opt = 1;
    if((srv_sockfd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
        perror("socket error");
        exit(1);
    }
    struct sockaddr_in server_sockaddr, client_sockaddr;
    server_sockaddr.sin_family = AF_INET;
    server_sockaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    server_sockaddr.sin_port = htons(SERVER_PORT);
 
    setsockopt(srv_sockfd, SOL_SOCKET, SO_REUSEADDR, &sock_opt, sizeof(sock_opt));
     
    if(bind(srv_sockfd, (struct sockaddr *)&server_sockaddr, sizeof(server_sockaddr)) == -1) {
        perror("bind error");
        exit(1);
    }
 
    if(listen(srv_sockfd, MAX_CONNECT) == -1) {
        perror("listen error");
        exit(1);
    }
 
    pthread_t pt;
    socklen_t len = sizeof(client_sockaddr);
    pthread_attr_t attr;
    pthread_attr_init(& attr);
    pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
    printf("chatroom init success\n\n");
    while(true) {
        cli_sockfd = accept(srv_sockfd, (struct sockaddr *)&client_sockaddr, &len);
        if(cli_sockfd == -1) {
            continue;
        }
        chat_item ci;
        ci.sockfd = cli_sockfd;
        strcpy(ci.ip, inet_ntoa(client_sockaddr.sin_addr));
        sprintf(ci.name, "guest_%d", chat_user_total++);
        if(pthread_create(& pt, &attr, chat_login, (void *) &ci)) {
            perror("pthread error");
            continue;
        }
        chat_sockfd_list[chat_sockfd_idx++] = cli_sockfd;
        printf("%s login, ip:%s\n", ci.name, ci.ip);
    }
    pthread_attr_destroy (&attr);
    close(srv_sockfd);
    return 0;
}
 
 
void *chat_login(void *args)
{
    chat_item ci = *((chat_item *) args);
    char buffer[BUFFER_SIZE] = {0};
    sprintf(buffer, "<!--system:welcome %s join the chatroom!-->\n", ci.name);
    chat_send(buffer, 0);
    memset(buffer, 0, sizeof(buffer));
    char data[BUFFER_SIZE] = {0};
    while(true) {
        memset(buffer, 0, sizeof(buffer));
        memset(data, 0, sizeof(data));
        int len = recv(ci.sockfd, buffer, sizeof(buffer), 0);
        if(len == 0 ||  strcmp(buffer, "exit\n") == 0) {
            chat_exit(&ci);
            break;
        }
         
        if(strlen(buffer) == 0) {
            continue;
        }
        time_t t = time(0); 
        char tmp[64] = {0}; 
        strftime(tmp, sizeof(tmp), "%H:%M:%S", localtime(&t)); 
     
        sprintf(data, "%s(%s) %s\n", ci.name, ci.ip, tmp);
        strcat(data, buffer);
        chat_send(data, ci.sockfd);
    }
}
 
void chat_send(char *msg, int from_sockfd)
{
    for(int i=0; i<chat_sockfd_idx; i++) {
        if(from_sockfd != 0 && from_sockfd == chat_sockfd_list[i]) {
            continue;
        }
        if(chat_sockfd_list[i] == 0) {
            continue;
        }
        send(chat_sockfd_list[i], msg, strlen(msg), 0);
    }
}
 
void chat_exit(chat_item * ci)
{
    char buffer[BUFFER_SIZE] = {0};
    sprintf(buffer, "<!--system:%s(%s) logout-->", ci->name, ci->ip);
    chat_send(buffer, ci->sockfd);
    printf("%s logout, ip:%s\n", ci->name, ci->ip);
    int j = 0;
    for(int i=0; i<chat_sockfd_idx; i++) {
        if(ci->sockfd == chat_sockfd_list[i]) {
            close(ci->sockfd);
            ci->sockfd = 0;
            chat_sockfd_list[i] = 0;
            j = i;
            break;
        }
    }
    memmove(chat_sockfd_list, chat_sockfd_list + j + 1, chat_sockfd_idx - j - 1);
    chat_sockfd_idx--;
}
```
**chat_client.cpp**

启动两个线程一个用来接收数据，一个用来输入信息并发送。当输入exit并回车或者服务端断开时结束线程。
```
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>  
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <pthread.h>
#include <signal.h>
 
#define SERVER_HOST "127.0.0.1"
#define SERVER_PORT 8801
#define BUFFER_SIZE 2048
 
int cli_sockfd;
pthread_t trd_show, trd_push;
 
void *show_msg(void * args);
void *push_msg(void * args);
void bye(int sig);
 
int main(int argc, char **argv)
{
    signal(SIGINT, bye);
    if((cli_sockfd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
        perror("socket error");
        exit(1);
    }
    struct sockaddr_in server_sockaddr;
    server_sockaddr.sin_family = AF_INET;
    server_sockaddr.sin_addr.s_addr = inet_addr(SERVER_HOST);
    server_sockaddr.sin_port = htons(SERVER_PORT);
 
    if(connect(cli_sockfd, (struct sockaddr *)&server_sockaddr, sizeof(server_sockaddr)) == -1) {
        perror("connect error");
        exit(1);
    }
     
    pthread_attr_t attr;
    pthread_attr_init(& attr);
    pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
    pthread_create(&trd_show, &attr, show_msg, &cli_sockfd);
    pthread_create(&trd_push, &attr, push_msg, &cli_sockfd);
    pthread_attr_destroy (&attr);
    pthread_exit(NULL);
    return 0;
}
 
void *show_msg(void * args)
{
    int sockfd = *((int*)args);
    char buffer[BUFFER_SIZE] = {0};
    while(true) {
        int len = recv(sockfd, buffer, sizeof(buffer), 0);
        if(len <= 0) {
            break;
        }
        printf("%s\n", buffer);
        memset(buffer, 0, sizeof(buffer));
    }
    bye(0);
}
 
void *push_msg(void * args)
{
    int sockfd = *((int*)args);
    char sdffer[BUFFER_SIZE] = {0};
    while(fgets(sdffer, sizeof(sdffer), stdin)) {
        if(strcmp(sdffer, "\n") == 0) {
            continue;
        }
        if(strcmp(sdffer, "exit\n") == 0) {
            break;
        }
        printf("\n");
        send(sockfd, sdffer, sizeof(sdffer), 0);
        memset(sdffer, 0, sizeof(sdffer));
    }
    bye(0);
}
 
void bye(int sig)
{
    printf("\nByeBye.\n\n", sig);
    pthread_cancel(trd_push);
    pthread_cancel(trd_show);
    close(cli_sockfd);
}
```
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/epoll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <fcntl.h>
#include <errno.h>

#define MAX_EVENTS 10 // 最大事件数
#define BUF_SIZE 1024 // 缓冲区大小

// 设置非阻塞模式的函数
void setnonblocking(int sockfd) {
    int opts = fcntl(sockfd, F_GETFL);
    if (opts < 0) {
        perror("fcntl(F_GETFL)");  
        exit(EXIT_FAILURE);
    }
    opts = (opts | O_NONBLOCK);  
    if (fcntl(sockfd, F_SETFL, opts) < 0) {
        perror("fcntl(F_SETFL)");  
        exit(EXIT_FAILURE);
    }
}

int main() {
    int listen_fd, epoll_fd, n, i, conn_sock;
    struct sockaddr_in server_addr;
    struct epoll_event ev, events[MAX_EVENTS];  

    // 创建epoll事件监听池
    epoll_fd = epoll_create1(0);
    if (epoll_fd == -1) {
        perror("epoll_create1 error"); 
        exit(EXIT_FAILURE);
    }

    // 创建网络套接字
    listen_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_fd == -1) {
        perror("socket error");  
        exit(EXIT_FAILURE);
    }
    // 初始化网络地址信息
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;  
    server_addr.sin_addr.s_addr = INADDR_ANY; 
    server_addr.sin_port = htons(8080);  
    // 绑定端口8080
    if (bind(listen_fd, (struct sockaddr *)&server_addr, sizeof(server_addr)) == -1) {
        perror("bind error"); 
        exit(EXIT_FAILURE);
    }
    // 将套接字设置为监听模式
    if (listen(listen_fd, 10) == -1) {
        perror("listen error"); 
        exit(EXIT_FAILURE);
    }
    // 设置socket为非阻塞模式
    setnonblocking(listen_fd);
    // 设置监听事件类型为EPOLLIN（可读事件）
    ev.events = EPOLLIN; 
    ev.data.fd = listen_fd;
    // 将监听的文件描述符添加到epoll事件监听池
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, listen_fd, &ev) == -1) {
        perror("epoll_ctl: listener"); 
        exit(EXIT_FAILURE);
    }

    for (;;) { 
        // 等待事件的产生.
        // 事件可能是listen_fd上有新的连接，或者是其他已连接的socket有数据到来
        n = epoll_wait(epoll_fd, events, MAX_EVENTS, -1);  
        for (i = 0; i < n; i++) { 
            if (events[i].data.fd == listen_fd) { 
                // 场景一：新的连接请求
                while ((conn_sock = accept(listen_fd, NULL, NULL)) > 0) {  
                    setnonblocking(conn_sock);  
                    ev.events = EPOLLIN | EPOLLET; 
                    ev.data.fd = conn_sock;
                    // 将新的fd添加到epoll的监听队列中
                    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, conn_sock, &ev) == -1) { 
                        perror("epoll_ctl: conn_sock");  
                        exit(EXIT_FAILURE);
                    }
                }
                if (conn_sock == -1) {
                    if (errno != EAGAIN && errno != ECONNABORTED && errno != EPROTO && errno != EINTR)
                        perror("accept");  
                }
                continue;
            } else { 
                // 场景二：已连接的socket有数据到来
                ssize_t count;
                char buf[BUF_SIZE];  
                count = read(events[i].data.fd, buf, sizeof buf);  
                if (count == -1) {
                    if (errno != EAGAIN) {
                        perror("read error"); 
                        close(events[i].data.fd);  
                    }
                } else if (count == 0) {
                    close(events[i].data.fd); 
                } else {
                    // 回写数据给客户端
                    write(events[i].data.fd, buf, count);
                }
            }
        }
    }

    close(listen_fd); // 关闭监听文件描述符
    close(epoll_fd); // 关闭epoll文件描述符

    return 0;
}

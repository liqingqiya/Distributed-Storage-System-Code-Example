#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <signal.h>
#include <string.h>

#define INPUT_BUFFER_SIZE 100

// 信号处理函数
void signal_handler(int signum) {
    char buffer[INPUT_BUFFER_SIZE];
    ssize_t n;

    // 重置信号处理器，以便下次SIGIO信号再次触发
    signal(SIGIO, signal_handler);

    // 读取数据
    n = read(STDIN_FILENO, buffer, sizeof(buffer) - 1);

    if (n > 0) {
        buffer[n] = '\0'; // 确保字符串被正确终结
        printf("Received Input: %s", buffer);
    }
}

int main(void) {
    // 设置 stdin 为非阻塞模式
    fcntl(STDIN_FILENO, F_SETFL, O_NONBLOCK);

    // 设置信号处理函数（SIGIO为I/O的信号）
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_flags = 0;
    sa.sa_handler = signal_handler;
    sigaction(SIGIO, &sa, NULL);

    // 设置文件描述符的所有者为当前进程
    fcntl(STDIN_FILENO, F_SETOWN, getpid());

    // 使能异步通知
    fcntl(STDIN_FILENO, F_SETFL, fcntl(STDIN_FILENO, F_GETFL) | FASYNC);

    // 进入一个无限循环，等待信号到来
    while (1) {
        pause(); // 等待信号
    }

    return 0;
}
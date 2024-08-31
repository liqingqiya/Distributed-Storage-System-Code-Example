#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <libaio.h>

#define FILE_NAME "example.txt"
#define BUFFER_SIZE 1024

int main(int argc, char *argv[]) {
    io_context_t ctx;
    struct iocb cb;
    struct iocb *cbs[1];
    struct io_event events[1];
    int fd, ret;
    void *buf;
    
    // 分配地址按照 4KiB 对齐内存
    if (posix_memalign(&buf, 4096, BUFFER_SIZE)!= 0) {
        perror("posix_memalign error");
        return -1;
    }

    // 初始化异步I/O上下文
    memset(&ctx, 0, sizeof(ctx));
    if (io_setup(1, &ctx) != 0) {
        return -1;
    }

    // 打开文件
    fd = open(FILE_NAME, O_RDONLY | O_DIRECT);
    if (fd < 0) {
        perror("open error");
        io_destroy(ctx);
        return -1;
    }

    // 准备异步读操作的控制块
    io_prep_pread(&cb, fd, buf, BUFFER_SIZE, 0);
    cbs[0] = &cb;

    // 提交异步读请求
    if (io_submit(ctx, 1, cbs) != 1) {
        io_destroy(ctx);
        close(fd);
        return -1;
    }

    // 等待异步读操作完成
    ret = io_getevents(ctx, 1, 1, events, NULL);
    if (ret != 1) {
        io_destroy(ctx);
        close(fd);
        return -1;
    }

    // 打印读到的结果
    if (events[0].res2 == 0) {
        printf("Read %lld bytes from file: %s\n", events[0].res, buf);
    } else {
        printf("Read failed\n");
    }

out:
    io_destroy(ctx);
    close(fd);
    return 0;
}
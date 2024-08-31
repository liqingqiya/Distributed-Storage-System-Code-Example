package main

import (
	"fmt"
	"log"
	"syscall"
)

func main() {
	// 以只读的方式打开文件
	fd, err := syscall.Open("hello.txt", syscall.O_RDONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}

	// 获取文件状态信息
	stat := &syscall.Stat_t{}
	if err = syscall.Fstat(fd, stat); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1)
	// 使用 Read 函数顺序读取整个文件
	for i := 0; i < int(stat.Size); i++ {
		_, err = syscall.Read(fd, buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", buf)
	}
	// 使用 Pread 函数随机读取文件的每个字节
	for i := 0; i < int(stat.Size); i++ {
		_, err = syscall.Pread(fd, buf, int64(i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", buf)
	}
}

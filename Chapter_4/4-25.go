package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	// 每次都以截断清零的方式打开文件
	fd, err := syscall.Open("hello.txt",
		syscall.O_CREAT|syscall.O_RDWR|syscall.O_TRUNC, 0700)
	if err != nil {
		log.Fatal(err)
	}

	// 调到 4096 的偏移位置
	off, err := syscall.Seek(fd, 4096, os.SEEK_SET)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("off:%v\n", off)

	// 写一小段数据
	n, err := syscall.Write(fd, []byte("hello world"))
	if err != nil {
		log.Fatal(err)
	}

	_ = n
}

package main

import (
	"log"
	"syscall"
)

func main() {
	// 以创建、读写、截断模式打开文件
	fd, err := syscall.Open("hello.txt",
		syscall.O_CREAT|syscall.O_RDWR|syscall.O_TRUNC, 0700)
	if err != nil {
		log.Fatal(err)
	}

	// 写入一串数据
	_, err = syscall.Write(fd, []byte("hello world"))
	if err != nil {
		log.Fatal(err)
	}

	// 在文件开始位置覆盖写入一个字节
	_, err = syscall.Pwrite(fd, []byte("H"), 0)
	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"encoding/binary"
	"log"
	"os"
)

func main() {
	// 打开文件
	file, err := os.OpenFile("hello.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0700)
	if err != nil {
		log.Fatal(err)
	}
	// main 函数退出时，关闭文件
	defer file.Close()

	// 文件大小 1G
	fileSize := 1 * 1024 * 1024 * 1024
	// 每次写入 4K 的数据，顺序写入
	buf := make([]byte, 4*1024)
	for writen := 0; writen < fileSize; {
		// 装填特殊的数据
		binary.BigEndian.PutUint32(buf, uint32(writen))
		// 顺序写入
		n, err := file.Write(buf)
		if err != nil {
			log.Fatal(err)
		}
		writen += n
	}
}

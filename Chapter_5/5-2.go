package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// 打开文件
	file, err := os.OpenFile("hello.txt", os.O_RDONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	// main 函数退出时，关闭文件
	defer file.Close()

	start := time.Now()

	// 每次写入 4K 的数据，顺序写入
	buf := make([]byte, 4*1024)
	reader := io.LimitReader(file, 1*1024*1024*1024)
	for {
		// 顺序读取
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// 处理数据
		offset := binary.BigEndian.Uint32(buf)
		_ = offset
		// fmt.Printf("offset: %v\n", offset)
	}

	latency := time.Since(start)
	fmt.Printf("latency:%v\n", latency.Seconds())

}

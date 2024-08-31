package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
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

	block := 4096                                 // I/O大小以 4K 为粒度
	totalBlocks := 1 * 1024 * 1024 * 1024 / block // 1G的文件有多少个4K
	workers := 32                                 // 32 个并发 Goroutine 读取
	subTotalBlocks := totalBlocks / workers       // 每个 Goroutine 处理的4K个数

	start := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			buf := make([]byte, 4*1024)
			for n := 0; n < subTotalBlocks; n++ {
				off := (index + n*workers) * block
				// ReadAt(对应 Pread 系统调用)，支持并发
				_, err := file.ReadAt(buf, int64(off))
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				// 处理数据
				offset := binary.BigEndian.Uint32(buf)
				fmt.Printf("offset: %v\n", offset)
			}
		}(i)
	}
	wg.Wait()
	latency := time.Since(start)
	fmt.Printf("latency:%v\n", latency.Seconds())
}

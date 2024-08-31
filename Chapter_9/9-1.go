package main

import (
	"fmt"
	"hash/crc32"
)

func main() {
	// 原始数据
	data := []byte("Hello World")
	// 计算CRC32校验和
	checksum := crc32.ChecksumIEEE(data)
	// 输出结果
	fmt.Printf("CRC32 (IEEE): %x\n", checksum)
}

package main

import (
	"crypto/md5"
	"fmt"
	"io"
)

func main() {
	// 原始数据
	data := "Hello World"
	// 创建一个新的 md5 实例
	hasher := md5.New()
	// 写入数据到 md5 实例中
	io.WriteString(hasher, data)
	// 计算 MD5 哈希值
	sum := hasher.Sum(nil)
	// 输出结果
	fmt.Printf("MD5: %x\n", sum)
}

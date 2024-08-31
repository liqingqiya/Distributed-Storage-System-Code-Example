package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	// 原始数据
	data := "Hello World"
	// 创建一个 sha256 实例
	hasher := sha256.New()
	// 写入数据到 sha256 实例中
	hasher.Write([]byte(data))
	// 计算 SHA256 哈希值
	sum := hasher.Sum(nil)
	// 输出结果
	fmt.Printf("SHA256: %x\n", sum)
}

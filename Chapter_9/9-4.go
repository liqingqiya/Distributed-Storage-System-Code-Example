package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func main() {
	// 原始数据
	data := "Hello World"
	// HMAC密钥
	key := "secret"
	// 创建一个新的基于 SHA256 的 HMAC 哈希实例
	hasher := hmac.New(sha256.New, []byte(key))
	// 写入数据到哈希实例中
	hasher.Write([]byte(data))
	// 计算 HMAC 值
	sum := hasher.Sum(nil)
	// 输出结果
	fmt.Printf("HMAC-SHA256: %x\n", sum)
}

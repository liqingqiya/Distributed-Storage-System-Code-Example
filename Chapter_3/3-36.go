package main

import (
	"io"
	"os"
	"strings"
)

func main() {
	// 字符串
	str := "hello world"
	// 构建Reader
	readerFromBytes := strings.NewReader(str)
	// 把Reader的数据拷贝到标准输出
	n, err := io.Copy(os.Stdout, readerFromBytes)

	_, _ = n, err
}

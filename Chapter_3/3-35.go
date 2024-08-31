package main

import (
	"bytes"
	"io"
	"os"
)

func main() {
	// 'h', 'e', 'l', 'l', 'o', '\n'
	buffer := []byte{104, 101, 108, 108, 111, 10}
	// 构建Reader
	readerFromBytes := bytes.NewReader(buffer)
	// 把Reader的数据拷贝到标准输出
	n, err := io.Copy(os.Stdout, readerFromBytes)

	_ = n
	_ = err
}

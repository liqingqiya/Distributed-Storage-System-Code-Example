package main

import (
	"embed"
	"fmt"
)

//go:embed hello.txt
var f embed.FS

func main() {
	// 打开文件
	file, _ := f.Open("hello.txt")

	// 读文件
	buf := make([]byte, 4096)
	_, _ = file.Read(buf)

	fmt.Printf("%s\n", buf)
}

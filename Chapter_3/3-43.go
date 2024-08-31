package main

import (
	"io"
	"os"
)

func main() {
	// 一行代码实现回显
	io.Copy(os.Stdout, os.Stdin)
}

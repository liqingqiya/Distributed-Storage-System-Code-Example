package main

import (
	"io"
	"os"
)

func main() {
	f, _ := os.OpenFile("hello.txt", os.O_RDONLY, 0)
	io.Copy(io.Discard, f)
}

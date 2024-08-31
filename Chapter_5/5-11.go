package main

import (
	"bufio"
	"os"
)

func main() {
	file, err := os.OpenFile("hello.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 4*1024*1024)
	for i := 0; i < 1000000; i++ {
		file.Write([]byte("Hello, world!\n"))
	}
	writer.Flush()
}

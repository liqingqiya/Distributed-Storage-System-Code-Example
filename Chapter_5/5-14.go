package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("hello.txt", os.O_RDONLY, 0700)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 4*1024*1024)
	buf := make([]byte, 1024)
	for i := 0; i < 1000000; i++ {
		_, err := reader.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stdout, "%s", buf)
	}
	fmt.Println()
}

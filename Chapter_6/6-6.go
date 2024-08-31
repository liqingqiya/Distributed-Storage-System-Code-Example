package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	consoleReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入消息: ")
		input, _ := consoleReader.ReadString('\n')

		// 发送消息到服务器
		n, err := conn.Write([]byte(input))
		if err != nil {
			// 错误处理
		}
		_ = n

		// 读取来自服务器的响应
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("服务器响应错误:", err)
			break
		}
		fmt.Print("来自服务器的响应: " + response)
	}
}

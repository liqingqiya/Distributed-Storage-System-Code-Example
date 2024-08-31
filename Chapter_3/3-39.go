package main

import (
	"fmt"
	"log"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	// 读取网络数据（客户端发送过来）
	conn.Read(buf)
	// 处理：打印到控制台
	fmt.Printf("get: <%s>\n", buf)

	// 处理之后回复响应
	conn.Write([]byte("pong: "))
	conn.Write(buf)
}

func main() {
	// 监听端口
	server, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("err:%v", err)
	}
	for {
		// 监听获取客户端连接
		c, err := server.Accept()
		if err != nil {
			log.Fatalf("err:%v", err)
		}
		// 处理客户端网络请求
		go handleConn(c)
	}
}

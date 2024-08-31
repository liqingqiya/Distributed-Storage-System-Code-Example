package main

import (
	"io"
	"net"
	"os"
)

func main() {
	// 和服务端建立连接
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	// 发送数据给服务端
	conn.Write([]byte("hello world"))
	// 打印服务端的响应
	io.Copy(os.Stdout, conn)
}

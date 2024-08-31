package main

import (
	"fmt"
	"io"
	"os"
)

// 实现一个字符串的Reader。实际项目中更常见的是文件、网络的句柄
type MyReader struct {
	s string
	i int64
}

func (r *MyReader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// 实现一个简单的Writer。打印到标准输出。实际项目中更常见的是文件、网络的句柄
type MyWriter struct{}

func (w *MyWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("%s", string(p))
	return len(p), nil
}

func main() {
	var r io.Reader
	var w io.Writer

	r = &MyReader{s: "hello world\n"}
	w = &MyWriter{}
	w = io.MultiWriter(w, os.Stderr)

	io.CopyBuffer(w, r, make([]byte, 1))
}

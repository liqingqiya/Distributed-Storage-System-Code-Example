package main

import (
	"fmt"
	"hash/crc32"
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

	// 分流 1：数据副本（用 pipe 写转读，然后并发操作）
	pr, pw := io.Pipe()
	r = io.TeeReader(r, pw)
	repErrCh := make(chan error, 1)

	go func() {
		var repErr error
		defer func() {
			pr.CloseWithError(repErr)
			repErrCh <- repErr
		}()
		// 处理 3：处理数据流副本，写到标准错误输出（控制台）。
		// 常见的是把副本数据写到文件
		_, repErr = io.CopyBuffer(os.Stderr, pr, make([]byte, 1))
	}()

	// 分流 2：做 crc 计算
	crcW := crc32.NewIEEE()
	r = io.TeeReader(r, crcW)

	// 处理 1：读写数据
	io.CopyBuffer(w, r, make([]byte, 1))

	// 处理 2：计算 crc
	fmt.Printf("crc:%x\n", crcW.Sum32())

	pw.Close()
	<-repErrCh
}

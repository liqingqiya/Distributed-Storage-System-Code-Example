package client

import (
	"io"
	"log"
	"sync"
)

func FragWrite(r io.Reader, size int) {
	fragSize := 4096
	remain := size

	wg := sync.WaitGroup{}
	var buf []byte
	for remain > 0 {
		if remain >= fragSize {
			buf = make([]byte, fragSize)
		} else {
			buf = make([]byte, remain)
		}

		_, err := io.ReadFull(r, buf)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		// 请求下发，并发处理
		go func(buf []byte) {
			wg.Done()
			log.Printf("len()=%d", len(buf))
		}(buf)
	}

	// 汇总处理
	wg.Wait()
}

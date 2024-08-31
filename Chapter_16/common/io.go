package common

import "io"

type Writer struct {
	io.WriterAt
	Offset int64
}

func (p *Writer) Write(val []byte) (n int, err error) {
	n, err = p.WriteAt(val, p.Offset)
	p.Offset += int64(n)
	return
}

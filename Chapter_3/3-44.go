package main

import (
	"bufio"
	"io"
)

func main() {
	var r io.Reader
	var w io.Writer

	w = bufio.NewWriter(w)
	w = bufio.NewWriterSize(w, 512)
	r = bufio.NewReader(r)
	r = bufio.NewReaderSize(r, 512)
}

package main

import "fmt"

func main() {
	var twoDimensional [4][8]byte = [4][8]byte{
		{0x1, 0x2, 0x3},
		{0x4, 0x5, 0x6},
	}

	n := cap(twoDimensional)
	n1 := len(twoDimensional)

	fmt.Printf("n:%v, n1:%v\n", n, n1)

	_ = twoDimensional
}

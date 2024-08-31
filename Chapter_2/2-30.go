package main

import "fmt"

func main() {
	c := make(chan int, 1)
	v := 1

	select {
	case c <- v:
		fmt.Println("c <- v")
	default:
		fmt.Println("default")
	}
}

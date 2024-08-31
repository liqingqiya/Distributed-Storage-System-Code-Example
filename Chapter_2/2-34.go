package main

import "fmt"

func main() {
	c := make(chan int, 1)

	c <- 1

	select {
	case v, ok := <-c:
		fmt.Printf("v:%v, ok:%v\n", v, ok)
	default:
		fmt.Println("default")
	}
}

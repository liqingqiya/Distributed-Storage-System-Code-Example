package main

import "fmt"

type TestAPI interface {
	Method1() error
	Method2() error
}

type Tester struct{}

func (t *Tester) Method1() error { return nil }
func (t *Tester) Method2() error { return nil }

func main() {
	var test *Tester = nil
	var api TestAPI
	// 虽然 test 值是 nil，但它携带了类型信息
	api = test
	if api != nil {
		fmt.Printf("not nil\n")
	}
}

package main

type TestAPI interface {
	Method1() error
	Method2() error
}

type Tester struct{}

func (t *Tester) Method1() error { return nil }
func (t *Tester) Method2() error { return nil }

func main() {
	var t Tester
	var api TestAPI

	api = &t

	// 函数调用 ...
	api.Method1()
	api.Method2()
	_ = api
}

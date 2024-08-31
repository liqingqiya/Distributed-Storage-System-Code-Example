package main

// TestAPI 接口定义
type TestAPI interface {
	Method1() error
	Method2() error
}

// Tester 结构体定义
type Tester struct {
	TestAPI
	Name string
}

// Tester 的 Method1 方法的实现
func (t *Tester) Method1() error { return nil }

func main() {
	var t Tester

	t.Method1()         // 正常调用，等价于：((*Tester)t).Method1()
	t.Method2()         // 运行报告 panic 错误，等价于：((TestAPI)nil).Method2()
	t.TestAPI.Method1() // 运行报告 panic 错误，等价于：((TestAPI)nil).Method1()
	t.TestAPI.Method2() // 运行报告 panic 错误，等价于：((TestAPI)nil).Method2()
}

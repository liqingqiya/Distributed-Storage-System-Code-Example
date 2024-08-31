package main

type TestAPI interface {
	Method1() error
	Method2() error
}

type Tester1 struct{}

func (t *Tester1) Method1() error { return nil }
func (t *Tester1) Method2() error { return nil }

type Tester2 struct{}

func (t *Tester2) Method1() error { return nil }
func (t *Tester2) Method2() error { return nil }
func (t *Tester2) Method3() error { return nil }

func main() {
	var t1 Tester1
	var t2 Tester2

	var api1, api2, api3 TestAPI

	api1 = &t1
	api2 = &t2
	api3 = &t2

	api1.Method1()
	api2.Method1()
	api3.Method2()
}

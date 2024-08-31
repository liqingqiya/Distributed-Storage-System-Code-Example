package main

type Object struct{ name string }

func (s *Object) Name() string { return s.name }

type API1 interface{ Name() string }
type API2 interface{ Name() string }

func main() {
	var api1 API1
	var api2 API2

	obj := &Object{name: "concrete obj"}

	api1 = obj  // 具体类型到接口的赋值
	api2 = api1 // 接口到接口的赋值

	_, _ = api1, api2
}

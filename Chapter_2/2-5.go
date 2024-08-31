package main

import (
	"encoding/json"
)

type Person struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

func main() {
	p := &Person{Name: "Test", Age: 10}

	// 使用 JSON 规则进行序列化
	data, _ := json.Marshal(p)

	var p1 Person
	// 使用 JSON 规则进行反序列化
	json.Unmarshal(data, &p1)
}

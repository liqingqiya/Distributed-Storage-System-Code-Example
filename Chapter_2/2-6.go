package main

import (
	"encoding/binary"
)

type Person struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

// 序列化函数
// 序列化规则：|-NameLength(4B)-|--Name( 变长 )--|-Age(4B)-|
func (p *Person) Marshal() (data []byte, err error) {
	// 计算名字的长度
	nameLen := len(p.Name)
	// 计算需要分配的内存空间.
	n := 4 + nameLen + 4
	// 分配内存空间
	data = make([]byte, n)

	// 填充 length
	binary.BigEndian.PutUint32(data[:4], uint32(nameLen))
	// 填充 Name
	copy(data[4:], []byte(p.Name))
	// 填充 Age
	binary.BigEndian.PutUint32(data[4+nameLen:], uint32(p.Age))

	return data, nil
}

func (p *Person) Unmarshal(data []byte) (err error) {
	// 恢复 length
	nameLen := binary.BigEndian.Uint32(data[:4])
	// 恢复 Name
	p.Name = string(data[4 : 4+nameLen])
	// 恢复 Age
	p.Age = int32(binary.BigEndian.Uint32(data[4+nameLen:]))

	return nil
}

func main() {
	// 序列化：拆解
	p := &Person{Name: "Test", Age: 10}
	data, _ := p.Marshal()

	// 此处省略跨网络、进程等过程

	// 反序列化：恢复
	p1 := &Person{}
	p1.Unmarshal(data)
}

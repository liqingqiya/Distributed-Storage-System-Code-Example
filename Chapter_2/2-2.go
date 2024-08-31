package main

import "encoding/binary"

func main() {
	var n uint32 = 0x12345678
	bigBuf := make([]byte, 4) // 大端序
	litBuf := make([]byte, 4) // 小端序

	// 序列化为大端序和小端序
	binary.BigEndian.PutUint32(bigBuf, n)    // 大端序
	binary.LittleEndian.PutUint32(litBuf, n) // 小端序

	// 大端序和小端序的反序列化
	n1 := binary.BigEndian.Uint32(bigBuf)    // 大端序
	n2 := binary.LittleEndian.Uint32(litBuf) // 小端序

	//
	_, _ = n1, n2
}

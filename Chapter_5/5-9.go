package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
	AlignSize = 512
)

// 在 block 这个字节数组首地址，往后找，找到符合 AlignSize 对齐的地址，并返回
// 这里用到位操作，速度很快；
func alignment(block []byte, AlignSize int) int {
	return int(uintptr(unsafe.Pointer(&block[0])) & uintptr(AlignSize-1))
}

// 分配 BlockSize 大小的内存块
// 地址按照 512 对齐
func AlignedBlock(BlockSize int) []byte {
	// 分配一个，分配大小比实际需要的稍大
	block := make([]byte, BlockSize+AlignSize)

	// 计算这个 block 内存块往后多少偏移，地址才能对齐到 512
	a := alignment(block, AlignSize)
	offset := 0
	if a != 0 {
		offset = AlignSize - a
	}

	// 偏移指定位置，生成一个新的 block，这个 block 将满足地址对齐 512；
	block = block[offset : offset+BlockSize]
	if BlockSize != 0 {
		// 最后做一次校验
		a = alignment(block, AlignSize)
		if a != 0 {
			log.Fatal("Failed to align block")
		}
	}

	return block
}

func main() {
	file, err := os.OpenFile("hello.txt", syscall.O_DIRECT|os.O_RDWR, 0700)
	if err != nil {
		log.Fatal(err)
	}

	buf := AlignedBlock(4096)
	n, err := file.Write(buf)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(n)
}

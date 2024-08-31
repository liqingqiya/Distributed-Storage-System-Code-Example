package server

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"log"
	"os"
	"sync"

	"example.com/hellofs/common"
)

var (
	magic = [4]byte{0xab, 0xcd, 0xef, 0xcc}
	// header: |--magic(4)--|--id(8)--|--size(8)--|
	headerSize = 4 + 8 + 4
	// footer: |--crc(4)--|
	footerSize = 4
)

type FileStorage struct {
	lock sync.Mutex // 互斥锁
	off  int64      // 分配空间的偏移
	fd   *os.File   // 底层
	id   int64      // 文件 ID
}

func (s *FileStorage) FID() int64 {
	return s.id
}

func (s *FileStorage) Sync() {
	s.fd.Sync()
}

func (s *FileStorage) Write(id int64, reader io.Reader, size int64) (loc *common.Location, err error) {
	s.lock.Lock()
	startPos, pos := s.off, s.off
	s.off += (size + int64(headerSize) + int64(footerSize))
	s.lock.Unlock()

	crc := crc32.NewIEEE()
	reader = io.LimitReader(reader, int64(size))
	reader = io.TeeReader(reader, crc)

	header := make([]byte, headerSize)
	footer := make([]byte, footerSize)

	copy(header[:4], magic[:])
	binary.BigEndian.PutUint64(header[4:], uint64(id))
	binary.BigEndian.PutUint32(header[4+8:], uint32(size))

	// 写入头部
	s.fd.WriteAt(header, pos)
	pos += int64(headerSize)

	// 写入数据
	writer := &common.Writer{WriterAt: s.fd, Offset: pos}
	n, err := io.Copy(writer, reader)
	if err != nil {
		log.Fatal(err)
	}
	pos += n

	crc32Sum := crc.Sum32()
	binary.BigEndian.PutUint32(footer, crc32Sum)

	// 写入尾部
	s.fd.WriteAt(footer, pos)

	loc = &common.Location{
		FileID: uint64(s.id),
		Offset: startPos,
		Length: size,
		Crc:    crc32Sum,
	}

	return loc, nil
}

func (s *FileStorage) Read(loc *common.Location) (data []byte, err error) {
	header := make([]byte, headerSize)
	footer := make([]byte, footerSize)

	secReader := io.NewSectionReader(s.fd, loc.Offset, loc.Length+int64(headerSize)+int64(footerSize))

	// 读取头部
	_, err = secReader.Read(header)
	if err != nil {
		log.Fatal(err)
	}
	// 校验头部 ...
	// _magic := header[:4]
	// id := binary.BigEndian.Uint64(header[4:])
	// size := binary.BigEndian.Uint32(header[4+8:])

	// 读取数据
	crc := crc32.NewIEEE()
	reader := io.LimitReader(secReader, int64(loc.Length))
	reader = io.TeeReader(reader, crc)

	data, err = io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	// 读取尾部
	secReader.Seek(int64(headerSize)+loc.Length, io.SeekStart)
	secReader.Read(footer)
	__crcSum := binary.BigEndian.Uint32(footer)

	// 校验crc32 ...
	crcSum := crc.Sum32()
	if __crcSum != crcSum {
		log.Fatal("crc not match")
	}

	return data, nil
}

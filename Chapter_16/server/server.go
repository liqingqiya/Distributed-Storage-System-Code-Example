package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"example.com/hellofs/common"
	"github.com/gorilla/mux"
)

type StorageService struct {
	lock     sync.Mutex             // 互斥锁
	current  *FileStorage           // 当前可写的文件句柄
	files    map[int64]*FileStorage // 所有的文件句柄
	rootPath string                 // 存放数据的目录
	seqIdx   int64                  // 文件编号
}

// 服务端：写数据
func (s *StorageService) ObjectWrite(w http.ResponseWriter, r *http.Request) {
	// 参数解析
	vals := mux.Vars(r)
	idStr := vals["id"]
	sizeStr := vals["size"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// 获取文件的写句柄
	s.lock.Lock()
	current := s.current
	s.lock.Unlock()

	// 写数据
	loc, err := current.Write(int64(id), r.Body, int64(size))
	if err != nil {
		log.Printf("write error. err=%v", err)
		w.WriteHeader(500)
		return
	}

	// 返回 Location
	data, err := json.Marshal(loc)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
	log.Printf("write success. loc=%v", loc)
}

// 服务端：读数据
func (s *StorageService) ObjectRead(w http.ResponseWriter, r *http.Request) {
	// 参数解析
	vals := mux.Vars(r)
	fidStr := vals["fid"]
	offStr := vals["off"]
	lenStr := vals["size"]
	crcStr := vals["crc"]

	fileID, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	offset, err := strconv.ParseInt(offStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	length, err := strconv.ParseInt(lenStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	crcsum, err := strconv.ParseUint(crcStr, 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	loc := &common.Location{
		FileID: uint64(fileID),
		Offset: offset,
		Length: length,
		Crc:    uint32(crcsum),
	}

	// 获取文件句柄
	s.lock.Lock()
	stor, exist := s.files[fileID]
	s.lock.Unlock()
	if !exist {
		log.Printf("fileID:%v not exist", fileID)
		w.WriteHeader(400)
		return
	}

	// 读取数据，并返回
	data, err := stor.Read(loc)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
	log.Printf("len(data)=%d", len(data))
}

func (s *StorageService) Rotate() error {
	if s.current != nil {
		s.current.Sync()
	}

	newIdx := atomic.AddInt64(&s.seqIdx, 1)
	name := fmt.Sprintf("%s/idx.%d", s.rootPath, newIdx)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	s.lock.Lock()
	s.current = &FileStorage{fd: f, id: newIdx, off: 0}
	s.files[newIdx] = s.current
	s.lock.Unlock()

	return nil
}

func (s *StorageService) Init() error {
	_, err := os.Stat(s.rootPath)
	if err != nil {
		log.Fatal(err)
	}

	dirs, err := os.ReadDir(s.rootPath)
	if err != nil {
		log.Fatal(err)
	}

	var maxIdx int64
	for _, dir := range dirs {
		name := dir.Name()
		idxStr := strings.Split(name, "idx.")[1]
		idx, err := strconv.ParseInt(idxStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Open(path.Join(s.rootPath, name))
		if err != nil {
			log.Fatal(err)
		}
		if idx > maxIdx {
			maxIdx = idx
		}
		s.files[idx] = &FileStorage{fd: f, id: idx}
	}
	s.seqIdx = maxIdx

	s.Rotate()

	return nil
}

func NewStorageService(datapath string) *StorageService {
	s := &StorageService{rootPath: datapath, files: make(map[int64]*FileStorage)}
	return s
}

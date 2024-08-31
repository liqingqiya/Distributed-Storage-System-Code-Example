package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// 对应一个文件的 Sync 服务
type SyncJob struct {
	*sync.Cond                         // 引入条件变量
	holding    int32                   // 计数器，用于标识挂起的Sync请求的个数
	lastErr    error                   // 记录最后一次执行Sync操作的结果
	syncPoint  *sync.Once              // 确保Sync只执行一次
	syncFunc   func(interface{}) error // 执行 Sync 的函数
}

// Do方法执行实际的Sync操作
func (s *SyncJob) Do(job interface{}) error {
	s.L.Lock()
	if s.holding > 0 {
		// 如果有正在执行的请求则等待
		s.Wait()
	}
	s.holding += 1
	syncPoint := s.syncPoint
	s.L.Unlock()

	// 执行聚合Sync操作
	syncPoint.Do(func() {
		// 串行化执行Sync系统调用
		s.lastErr = s.syncFunc(job)
		s.L.Lock()
		// Sync调用之后，累计计数清零
		s.holding = 0
		// 重置sync.Once，准备下一次的聚合
		s.syncPoint = &sync.Once{}
		// 通知所有等待的Goroutine
		s.Broadcast()
		s.L.Unlock()
	})
	return s.lastErr
}

func NewSyncJob(fn func(interface{}) error) *SyncJob {
	return &SyncJob{
		Cond:      sync.NewCond(&sync.Mutex{}),
		syncFunc:  fn,
		syncPoint: &sync.Once{},
	}
}

func main() {
	file, err := os.OpenFile("hello.txt", os.O_RDWR, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 初始化Sync聚合服务
	syncJob := NewSyncJob(func(interface{}) error {
		fmt.Printf("performing sync...\n")
		time.Sleep(time.Second)
		return file.Sync()
	})
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// 执行写操作
			fmt.Printf("writing to file...\n")
			file.WriteAt([]byte(fmt.Sprintf("Idx:%v\n", idx)), int64(idx*1024))
			// 执行Sync操作，通过SyncJob聚合请求
			syncJob.Do(file)
		}(i)
	}
	wg.Wait()
}

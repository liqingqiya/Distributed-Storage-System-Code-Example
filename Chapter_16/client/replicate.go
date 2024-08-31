package client

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync"

	"example.com/hellofs/common"
)

// 表：server_group（复制组的配置）
type ServerGroup struct {
	GroupID int      `bson:"group_id"` // 副本组ID
	Servers []string `bson:"servers"`  // 副本组地址
}

type Err struct {
	idx int
	loc *common.Location
	err error
}

func (rep *ServerGroup) pickServerForRead() int {
	var idx int
	// 根据健康状态、均衡压力等多种维度考虑
	return idx
}

// 副本读策略：ReadOne
func (rep *ServerGroup) Read(ctx context.Context, loc *common.FileDataLocation) (content []byte, retErr error) {
	cli := common.NewClient()
	// 根据策略选择一个合适的副本
	pickIdx := rep.pickServerForRead()
	server := rep.Servers[pickIdx]
	args := &common.ReadArgs{
		Loc: *loc.Locations[pickIdx],
	}

	reader, err := cli.Read(ctx, server, args)
	if err != nil {
		return nil, err
	}

	content, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// 副本写策略：WriteAll
func (rep *ServerGroup) Write(ctx context.Context, data []byte, id int64) (retLocs []*common.Location, retErr error) {
	cli := common.NewClient()
	wg := sync.WaitGroup{}
	retResps := make([]Err, len(rep.Servers))
	for idx, sAddr := range rep.Servers {
		wg.Add(1)
		args := &common.WritedArgs{
			ID:   uint64(id),
			Size: int64(len(data)),
			Body: bytes.NewReader(data),
		}
		go func(index int, addr string) {
			defer wg.Done()
			loc, err := cli.Write(ctx, addr, args)
			retResps[index] = Err{idx: index, loc: loc, err: err}
		}(idx, sAddr)
	}
	wg.Wait()

	for _, ret := range retResps {
		if ret.err != nil {
			return nil, ret.err
		}
		retLocs = append(retLocs, ret.loc)
	}

	log.Printf("write success: %v", rep.Servers)

	return retLocs, nil
}

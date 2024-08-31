package client

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"bazil.org/fuse/fs"
	"example.com/hellofs/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	rootIno = 1
)

// hellofs 文件系统的主体
type HellofsService struct {
	metaCli      *mongo.Database // mongodb 的操作句柄
	replication  []*ServerGroup  // 复制组的集群拓扑
	inoAllocator *InodeAllocator // inode 的分配器
}

// 起点：返回根目录
func (fs *HellofsService) Root() (fs.Node, error) {
	return &Dir{
		inode: rootIno,
		files: make(map[string]*File),
		fs:    fs,
	}, nil
}

// 按照策略挑选一组可写的副本组
func (fs *HellofsService) PickWriteServerGroup() *ServerGroup {
	// 随机算法
	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(fs.replication))
	rep := fs.replication[idx]
	return rep
}

// 获取到指定的副本组
func (fs *HellofsService) PickServerGroupByID(groupID int) (*ServerGroup, error) {
	for _, s := range fs.replication {
		if s.GroupID == groupID {
			return s, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

// 分配 Inode 编号
func (fs *HellofsService) AllocInode() (uint64, error) {
	newIno := fs.inoAllocator.Alloc()
	return newIno, nil
}

// 加载超级块
func (fs *HellofsService) LoadSuperblock() error {
	var sb common.Superblock
	err := fs.metaCli.Collection(TableSuperblock).FindOne(context.TODO(), bson.D{}).Decode(&sb)
	if err != nil {
		return err
	}
	// ...

	// 加载分布式节点信息
	cursor, err := fs.metaCli.Collection(TableServerGroup).Find(context.TODO(), bson.D{})
	if err != nil {
		return err
	}

	var results []ServerGroup
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, ret := range results {
		cursor.Decode(&ret)
		var g ServerGroup
		g.GroupID = ret.GroupID
		g.Servers = append(g.Servers, ret.Servers...)
		fs.replication = append(fs.replication, &g)
	}

	return nil
}

func NewHelloFS(metadb *mongo.Database) *HellofsService {
	ialloc := &InodeAllocator{
		metaCli: metadb,
	}

	fs := &HellofsService{
		metaCli:      metadb,
		replication:  []*ServerGroup{},
		inoAllocator: ialloc,
	}

	// 从元数据中心加载
	if err := fs.LoadSuperblock(); err != nil {
		panic(err)
	}

	return fs
}

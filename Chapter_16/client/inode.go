package client

import (
	"context"
	"sync"

	"example.com/hellofs/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Inode 编号分配器
type InodeAllocator struct {
	lock     sync.RWMutex
	startIno uint64
	endIno   uint64
	step     uint64
	metaCli  *mongo.Database
}

// 申请 Inode 编号
func (a *InodeAllocator) Alloc() uint64 {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.startIno < a.endIno {
		newIno := a.startIno
		a.startIno += 1
		return newIno
	} else {
		a.fill()
	}
	newIno := a.startIno
	a.startIno += 1
	return newIno
}

// 找元数据中心申请一段可用的 inode 编号
func (a *InodeAllocator) fill() {
	ctx := context.TODO()
	coll := a.metaCli.Collection(TableSuperblock)

	filter := bson.M{"name": "hellofs"}
	var sb common.Superblock
	err := coll.FindOne(ctx, filter).Decode(&sb)
	if err != nil {
		panic(err)
	}

	startIno := sb.InodeNumBase
	step := sb.InodeStep
	nextStart := startIno + uint64(step)

	filter = bson.M{"name": "hellofs", "inode_num_base": startIno}
	update := bson.M{"$set": bson.M{"inode_num_base": nextStart}}

	var nsb common.Superblock
	err = coll.FindOneAndUpdate(ctx, filter, update).Decode(&nsb)
	if err != nil {
		panic(err)
	}

	a.startIno = startIno
	a.endIno = nextStart
}

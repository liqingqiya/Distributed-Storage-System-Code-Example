package client

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"example.com/hellofs/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TableServerGroup = "server_group"
	TableFileMeta    = "file_meta"
	TableSuperblock  = "superblock"
)

// 目录实现
type Dir struct {
	fs    *HellofsService  // 核心主体
	inode uint64           // inode 编号
	files map[string]*File // 目录的文件列表
}

// 目录的属性
func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	// Hellofs 全局只有一个目录
	a.Inode = d.inode
	a.Mode = os.ModeDir | 0555
	return nil
}

// 在目录中查找某名字的文件
func (d *Dir) Lookup(ctx context.Context, name string) (resp fs.Node, err error) {
	coll := d.fs.metaCli.Collection(TableFileMeta)
	var m common.FileMetaInfo
	filter := bson.M{"name": name}
	err = coll.FindOne(context.TODO(), filter).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, syscall.ENOENT
		}
	}
	f := NewFile(d.fs, &m.Inode, m.Name)
	return f, nil
}

// 创建文件
func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	// ...

	// 分配新的 Inode
	newIno, err := d.fs.AllocInode()
	if err != nil {
		return nil, nil, fmt.Errorf("alloc new inode failed")
	}

	meta := &common.FileMetaInfo{
		Name: req.Name,
		Inode: common.InodeInfo{
			Inode:      newIno,
			Mode:       uint32(req.Mode),
			Uid:        req.Uid,
			Gid:        req.Gid,
			Size:       0,
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
			AccessTime: time.Now(),
		},
	}

	// 写入元数据中心
	coll := d.fs.metaCli.Collection(TableFileMeta)
	_, err = coll.InsertOne(context.TODO(), meta)
	if err != nil {
		return nil, nil, fmt.Errorf("create file failed")
	}

	// 构造文件结构
	f := NewFile(d.fs, &meta.Inode, req.Name)

	// 放入缓存
	d.files[req.Name] = f.(*File)

	return f, f, nil
}

// 返回所有文件的列表
func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var dirDirs = []fuse.Dirent{}

	// 从元数据中心获取
	coll := d.fs.metaCli.Collection(TableFileMeta)
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var results []common.FileMetaInfo
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, ret := range results {
		cursor.Decode(&ret)
		dirent := fuse.Dirent{Inode: ret.Inode.Inode, Name: ret.Name, Type: fuse.DT_File}
		dirDirs = append(dirDirs, dirent)
	}

	return dirDirs, nil
}

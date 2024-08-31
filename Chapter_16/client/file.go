// 实现分布式加密的文件系统
package client

import (
	"context"
	"os"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"example.com/hellofs/common"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// 控制内核 fuse 模块的缓冲
	AttrValidDuration = 5 * time.Second
)

type File struct {
	fs           *HellofsService          // 文件系统管理结构
	fileName     string                   // 文件名称
	metaInfo     *common.InodeInfo        // 文件 inode 信息
	dataLocation *common.FileDataLocation // 位置元数据信息
}

// 获取文件属性
func (f *File) Attr(ctx context.Context, attr *fuse.Attr) error {
	// ...
	f.refreshMeta()
	attr.Inode = f.metaInfo.Inode
	attr.Mode = os.FileMode(f.metaInfo.Mode)
	attr.Size = f.metaInfo.Size
	attr.Atime = f.metaInfo.AccessTime
	attr.Mtime = f.metaInfo.ModifyTime
	attr.Ctime = f.metaInfo.CreateTime
	attr.Valid = AttrValidDuration
	return nil
}

// 读取文件的所有数据
func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	if f.dataLocation == nil {
		if err := f.refreshMeta(); err != nil {
			return nil, err
		}
	}

	repSet, err := f.fs.PickServerGroupByID(f.dataLocation.GroupID)
	if err != nil {
		return nil, err
	}

	content, err := repSet.Read(ctx, f.dataLocation)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// 文件写数据
func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	// 选节点
	repSet := f.fs.PickWriteServerGroup()

	// 副本冗余
	locs, err := repSet.Write(ctx, req.Data, int64(f.metaInfo.Inode))
	if err != nil {
		return err
	}

	dataLoc := &common.FileDataLocation{
		GroupID:   repSet.GroupID,
		Locations: locs,
	}

	// 元数据持久化到元数据中心
	err = f.writeMeta(ctx, dataLoc, len(req.Data))
	if err != nil {
		return err
	}

	resp.Size = len(req.Data)
	return nil
}

// 更新文件对应的元数据
func (f *File) writeMeta(ctx context.Context, loc *common.FileDataLocation, size int) error {
	coll := f.fs.metaCli.Collection(TableFileMeta)

	filter := bson.M{"inode.inode": f.metaInfo.Inode}
	update := bson.M{"$set": bson.M{"data_location": loc, "inode.size": size}}

	var m common.FileMetaInfo
	if err := coll.FindOneAndUpdate(ctx, filter, update).Decode(&m); err != nil {
		return err
	}
	return nil
}

// 刷新文件元数据
func (f *File) refreshMeta() error {
	coll := f.fs.metaCli.Collection(TableFileMeta)
	filter := bson.M{"inode.inode": f.metaInfo.Inode}

	var m common.FileMetaInfo
	err := coll.FindOne(context.TODO(), filter).Decode(&m)
	if err != nil {
		return err
	}

	f.metaInfo = &m.Inode
	f.dataLocation = m.DataLocs
	return nil
}

func NewFile(fs *HellofsService, meta *common.InodeInfo, fileName string) fs.Node {
	f := &File{
		fs:       fs,
		metaInfo: meta,
		fileName: fileName,
	}

	return f
}

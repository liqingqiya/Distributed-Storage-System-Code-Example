package common

import "time"

// 表：super_block（超级块信息）
type Superblock struct {
	Name         string `bson:"name"`
	InodeNumBase uint64 `bson:"inode_num_base"`
	InodeStep    uint64 `bson:"inode_step"`
}

// 表：file_meta（文件元数据）
type FileMetaInfo struct {
	Name     string            `bson:"name"`
	Inode    InodeInfo         `bson:"inode"`
	DataLocs *FileDataLocation `bson:"data_location"`
}

// 文件 Inode 信息
type InodeInfo struct {
	Inode      uint64    `bson:"inode"`
	Mode       uint32    `bson:"mode"`
	Size       uint64    `bson:"size"`
	Uid        uint32    `bson:"uid"`
	Gid        uint32    `bson:"gid"`
	Generation uint64    `bson:"gen"`
	ModifyTime time.Time `bson:"mtime"`
	CreateTime time.Time `bson:"ctime"`
	AccessTime time.Time `bson:"atime"`
}

// 文件位置信息
type FileDataLocation struct {
	GroupID   int         `bson:"group_id"`
	Locations []*Location `bson:"locations"`
}

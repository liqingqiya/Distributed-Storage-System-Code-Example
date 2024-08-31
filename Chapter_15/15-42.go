// 实现分布式加密的文件系统
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)

var (
	readonlyFiles = map[string]*File{
		"hello": {
			Inode: 20230606, Content: []byte("value: hello\n"), Mode: 0444,
		},
		"world": {
			Inode: 20230607, Content: []byte("value: world\n"), Mode: 0444,
		},
	}
)

// hellofs 文件系统的主体
type HellofsService struct{}

// 根节点
func (HellofsService) Root() (fs.Node, error) {
	return &Dir{inode: 20230101, files: readonlyFiles}, nil
}

// 定义目录
type Dir struct {
	inode uint64           // inode 编号
	files map[string]*File // 目录的文件列表
}

// GetAttr 请求
func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = d.inode
	a.Mode = os.ModeDir | 0444
	return nil
}

// Lookup 请求：根据名字查找对应的结构
func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	v, exist := readonlyFiles[name]
	if !exist {
		return nil, syscall.ENOENT
	}
	return v, nil
}

// Read 请求：对应读取目录的内容，读取全部项目
func (d *Dir) ReadDirAll(ctx context.Context) (dirents []fuse.Dirent, err error) {
	for name, file := range readonlyFiles {
		// 构造一个目录项
		entry := fuse.Dirent{Inode: uint64(file.Inode), Name: name, Type: fuse.DT_File}
		dirents = append(dirents, entry)
	}
	return dirents, nil
}

// 定义文件
type File struct {
	Inode   int64       // inode 编号
	Content []byte      // 文件内容
	Mode    os.FileMode // 文件模式
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = uint64(f.Inode)
	a.Size = uint64(len(f.Content))
	a.Mode = f.Mode
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return f.Content, nil
}

func main() {
	var mountpoint string
	// 挂载参数
	flag.StringVar(&mountpoint, "mountpoint", "", "mount point(dir)?")
	flag.Parse()
	if mountpoint == "" {
		log.Fatal("please input invalid mount point\n")
	}

	// 挂载到操作系统的目录上
	c, err := fuse.Mount(mountpoint, fuse.FSName("HelloWorld"), fuse.Subtype("HelloFS"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// 开启Daemon守护程序的处理（阻塞处理）
	err = fs.Serve(c, HellofsService{})
	if err != nil {
		log.Fatal(err)
	}
}

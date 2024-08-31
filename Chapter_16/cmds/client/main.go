package main

import (
	"context"
	"flag"
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"example.com/hellofs/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var mountpoint string
	flag.StringVar(&mountpoint, "mountpoint", "", "mount point(dir)?")
	flag.Parse()

	if mountpoint == "" {
		log.Fatal("please input invalid mount point\n")
	}
	c, err := fuse.Mount(mountpoint, fuse.FSName("helloworld"), fuse.Subtype("hellofs"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	uri := "mongodb://127.0.0.1:27017"
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := cli.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	helloFS := client.NewHelloFS(cli.Database("hellofs_db"))
	err = fs.Serve(c, helloFS)
	if err != nil {
		log.Fatal(err)
	}
}

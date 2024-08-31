package main

import (
	"context"
	"log"

	cli "example.com/hellofs/client"
	"example.com/hellofs/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://127.0.0.1:27017"
	dbCli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := dbCli.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	hellofsDB := dbCli.Database("hellofs_db")

	// 插入镜像关系
	coll := hellofsDB.Collection(cli.TableServerGroup)

	rep1 := cli.ServerGroup{
		GroupID: 1,
		Servers: []string{
			"http://127.0.0.1:37001",
			"http://127.0.0.1:37002",
			"http://127.0.0.1:37003",
		},
	}

	rep2 := cli.ServerGroup{
		GroupID: 2,
		Servers: []string{
			"http://127.0.0.1:37004",
			"http://127.0.0.1:37005",
			"http://127.0.0.1:37006",
		},
	}

	_, err = coll.InsertOne(context.TODO(), rep1)
	if err != nil {
		log.Fatal(err)
	}

	_, err = coll.InsertOne(context.TODO(), rep2)
	if err != nil {
		log.Fatal(err)
	}

	// 插入 superblock
	sb := common.Superblock{
		Name:         "hellofs",
		InodeNumBase: 1000,
		InodeStep:    1000,
	}
	coll = hellofsDB.Collection(cli.TableSuperblock)
	_, err = coll.InsertOne(context.TODO(), sb)
	if err != nil {
		log.Fatal(err)
	}
}

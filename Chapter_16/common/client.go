package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Location struct {
	FileID uint64
	Offset int64
	Length int64
	Crc    uint32
}

type ReadArgs struct {
	Loc Location
}

type WritedArgs struct {
	ID   uint64
	Size int64
	Body io.Reader
}

type WriteResp struct {
	Loc Location
}

type StorageAPI interface {
	Read(ctx context.Context, host string, args *ReadArgs) (reader io.Reader, err error)
	Write(ctx context.Context, host string, args *WritedArgs) (loc *Location, err error)
}

type client struct{}

func (c *client) Read(ctx context.Context, host string, args *ReadArgs) (reader io.Reader, err error) {
	urlStr := fmt.Sprintf("%v/object/read/fid/%d/off/%d/size/%d/crc/%d",
		host, args.Loc.FileID, args.Loc.Offset, args.Loc.Length, args.Loc.Crc)
	req, err := http.NewRequest(http.MethodPost, urlStr, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res.Body, nil
}

func (c *client) Write(ctx context.Context, host string, args *WritedArgs) (loc *Location, err error) {
	urlStr := fmt.Sprintf("%v/object/write/id/%d/size/%d",
		host, args.ID, args.Size)
	req, err := http.NewRequest(http.MethodPost, urlStr, args.Body)
	if err != nil {
		return
	}
	req.ContentLength = args.Size

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	loc = &Location{}
	if err = json.Unmarshal(resData, loc); err != nil {
		log.Fatal(err)
	}

	return loc, nil
}

func NewClient() StorageAPI {
	return &client{}
}

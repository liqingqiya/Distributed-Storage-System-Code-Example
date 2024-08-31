package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"example.com/hellofs/server"
	"github.com/gorilla/mux"
)

func main() {
	var datapath string
	var port int

	flag.StringVar(&datapath, "datapath", "./hellofs_data/", "data path")
	flag.IntVar(&port, "port", 8899, "listen port")
	flag.Parse()

	s := server.NewStorageService(datapath)
	s.Init()

	router := mux.NewRouter()

	router.HandleFunc("/object/write/id/{id}/size/{size}",
		s.ObjectWrite)
	router.HandleFunc("/object/read/fid/{fid}/off/{off}/size/{size}/crc/{crc}",
		s.ObjectRead)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	srv := http.Server{
		Handler: router,
		Addr:    address,
	}
	log.Fatal(srv.ListenAndServe())
}

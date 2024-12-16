package main

import (
	"c_cache/server"
	"flag"
	"fmt"
)

var (
	config = server.Config{}
)

func init() {
	flag.IntVar(&config.Port, "port", 8080, "Geecache server port")
}

func main() {
	flag.Parse()

	err := server.NewServer(fmt.Sprintf(":%d", config.Port)).Run()
	if err != nil {
		panic(err)
	}
}

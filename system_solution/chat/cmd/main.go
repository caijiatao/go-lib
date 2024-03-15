package main

import (
	"flag"
	"golib/system_solution/chat"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	chat.GetServer().RegisterRoutes()
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

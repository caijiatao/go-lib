package main

import (
	"flag"
	im2 "golib/system_solution/chat/im"
	"log"
	"net/http"
	"time"
)

func watchChange(w http.ResponseWriter, r *http.Request) {
	key := "test"
	err := im2.Register(w, r, key)
	if err != nil {
		return
	}
	go func() {
		for {
			after := time.After(time.Second)
			select {
			case <-after:
				im2.PushMessage(im2.NewMessage(key, []byte("test msg")))
			}
		}
	}()
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "C:\\project\\github\\go-lib\\im\\examples\\watch_change\\home.html")
}

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", watchChange)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		panic(err)
	}
}

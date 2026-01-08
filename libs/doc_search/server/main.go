package main

import (
	"fmt"
	"golib/libs/doc_search"
	"log"
	"net/http"
)

// ---- main ----

func main() {
	cfg := doc_search.LoadConfig()
	srv, err := doc_search.NewServer(cfg)
	if err != nil {
		log.Fatalf("init server failed: %v", err)
	}
	srv.StartWorkers()

	mux := http.NewServeMux()

	mux.HandleFunc("/upsert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			doc_search.WriteJSON(w, 405, map[string]any{"error": "method not allowed"})
			return
		}
		srv.UpsertHandler(w, r)
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			doc_search.WriteJSON(w, 405, map[string]any{"error": "method not allowed"})
			return
		}
		srv.SearchHandler(w, r)
	})

	mux.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			doc_search.WriteJSON(w, 405, map[string]any{"error": "method not allowed"})
			return
		}
		srv.DeleteHandler(w, r)
	})

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

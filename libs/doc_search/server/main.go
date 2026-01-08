package main

import (
	"encoding/json"
	"fmt"
	"golib/libs/doc_search"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

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
			writeJSON(w, 405, map[string]any{"error": "method not allowed"})
			return
		}
		srv.UpsertHandler(w, r)
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, 405, map[string]any{"error": "method not allowed"})
			return
		}
		srv.SearchHandler(w, r)
	})

	mux.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			writeJSON(w, 405, map[string]any{"error": "method not allowed"})
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

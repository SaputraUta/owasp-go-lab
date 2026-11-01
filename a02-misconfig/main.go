package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if os.Getenv("DB_PASSWORD") == "" {
		log.Fatal("DB_PASSWORD is required")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok \n"))
	})

	srv := &http.Server{
		Addr:              ":8081",
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("listening on :8081")
	log.Fatal(srv.ListenAndServe())
}

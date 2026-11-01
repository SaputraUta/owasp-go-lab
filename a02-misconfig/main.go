package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("DB_PASSWORD") == "" {
		log.Fatal("DB_PASSWORD is required")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok \n"))
	})

	log.Println("listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}

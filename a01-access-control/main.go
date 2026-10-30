package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("GET /me", requireAuth(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "you are %s\n", currentUser(r))
	}))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

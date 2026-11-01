package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	dbPassword := flag.String("db-password", "", "database password")
	flag.Parse()

	if *dbPassword == "" {
		log.Fatal("db-password is required")
	}

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok \n"))
	})

	log.Println("listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

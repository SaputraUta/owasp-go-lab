package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`
		CREATE TABLE products (name TEXT, category TEXT, released INTEGER);
		INSERT INTO products VALUES
			('Kopi Gayo', 'coffee', 1),
			('Teh Melati', 'tea', 1),
			('Kopi Rahasia 2027', 'coffee', 0);
	`)
	return err
}

func main() {
	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", listProducts)
	srv := &http.Server{
		Addr:              ":8085",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("listening on :8085")
	log.Fatal(srv.ListenAndServe())
}

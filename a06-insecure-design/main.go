package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Product struct {
	Name  string
	Price int
}

var products = map[string]Product{
	"kopi-gayo":  {Name: "Kopi Gayo", Price: 85000},
	"teh-melati": {Name: "Teh Melati", Price: 40000},
}

type CheckoutRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

const maxQuantity = 100

type Order struct {
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
	Total    int    `json:"total"`
}

func checkout(w http.ResponseWriter, r *http.Request) {
	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Quantity < 1 || req.Quantity > maxQuantity {
		http.Error(w, "invalid quantity", http.StatusBadRequest)
		return
	}

	p, ok := products[req.ProductID]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	order := Order{Product: p.Name, Quantity: req.Quantity, Total: p.Price * req.Quantity}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /checkout", checkout)

	srv := &http.Server{
		Addr:              ":8086",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("listening on :8086")
	log.Fatal(srv.ListenAndServe())
}

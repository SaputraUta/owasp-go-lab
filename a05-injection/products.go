package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Product struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

var db *sql.DB

func listProducts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	rows, err := db.Query("SELECT name, category FROM products WHERE released = 1 AND category = ?", category)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Name, &p.Category); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

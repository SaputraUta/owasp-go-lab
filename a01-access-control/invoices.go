package main

import (
	"encoding/json"
	"net/http"
)

type Invoice struct {
	ID     string `json:"id"`
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

var invoices = map[string]Invoice{
	"1": {ID: "1", Owner: "alice", Amount: 1500000},
	"2": {ID: "2", Owner: "bob", Amount: 9800000},
	"3": {ID: "3", Owner: "alice", Amount: 250000},
}

func getInvoice(w http.ResponseWriter, r *http.Request) {
	inv, ok := invoices[r.PathValue("id")]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

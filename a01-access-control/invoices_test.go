package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInvoiceAccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /invoices/{id}", requireAuth(getInvoice))

	cases := []struct {
		name  string
		token string
		id    string
		want  int
	}{
		{"owner reads own invoice", "alice-token", "1", http.StatusOK},
		{"other user's invoice is hidden", "alice-token", "2", http.StatusNotFound},
		{"bob reads own invoice", "bob-token", "2", http.StatusOK},
		{"no token", "", "1", http.StatusUnauthorized},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/invoices/"+tc.id, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Errorf("got %d, want %d", rec.Code, tc.want)
			}
		})
	}
}

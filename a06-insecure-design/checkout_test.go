package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckoutAbuseCases(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		wantCode  int
		wantTotal int
	}{
		{"normal order", `{"product_id":"kopi-gayo","quantity":2}`, http.StatusOK, 170000},
		{"tampered price is ignored", `{"product_id":"kopi-gayo","quantity":2,"price":1}`, http.StatusOK, 170000},
		{"negative quantity", `{"product_id":"kopi-gayo","quantity":-2}`, http.StatusBadRequest, 0},
		{"zero quantity", `{"product_id":"kopi-gayo","quantity":0}`, http.StatusBadRequest, 0},
		{"quantity above limit", `{"product_id":"kopi-gayo","quantity":101}`, http.StatusBadRequest, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/checkout", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			checkout(rec, req)

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body %q)", rec.Code, tc.wantCode, rec.Body.String())
			}
			if tc.wantCode != http.StatusOK {
				return
			}

			var order Order
			if err := json.NewDecoder(rec.Body).Decode(&order); err != nil {
				t.Fatal(err)
			}
			if order.Total != tc.wantTotal {
				t.Errorf("total %d, want %d", order.Total, tc.wantTotal)
			}
		})
	}
}

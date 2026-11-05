package main

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestListProductsRejectsInjection(t *testing.T) {
	if err := initDB(); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name     string
		category string
		want     string
		notWant  string
	}{
		{"normal search", "coffee", "Kopi Gayo", "Kopi Rahasia"},
		{"injection payload", "coffee' OR '1'='1", "[]", "Kopi Rahasia"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/products?category="+url.QueryEscape(tc.category), nil)
			rec := httptest.NewRecorder()
			listProducts(rec, req)

			body := rec.Body.String()
			if !strings.Contains(body, tc.want) {
				t.Errorf("body %q does not contain %q", body, tc.want)
			}
			if strings.Contains(body, tc.notWant) {
				t.Errorf("unreleased product leaked: %q", body)
			}
		})
	}
}

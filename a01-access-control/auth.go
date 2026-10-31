package main

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey struct{}

var tokens = map[string]string{
	"alice-token": "alice",
	"bob-token":   "bob",
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		user, ok := tokens[token]
		if !found || !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			// return
		}
		ctx := context.WithValue(r.Context(), ctxKey{}, user)
		next(w, r.WithContext(ctx))
	}
}

func currentUser(r *http.Request) string {
	return r.Context().Value(ctxKey{}).(string)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var users = map[string][]byte{}
var dummyHash []byte

func seedUsers() error {
	hash, err := bcrypt.GenerateFromPassword([]byte("sunshine"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	users["alice"] = hash

	dummyHash, err = bcrypt.GenerateFromPassword([]byte("dummy"), bcrypt.DefaultCost)
	return err
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if isLocked(req.Username) {
		http.Error(w, "too many attempts, try again later", http.StatusTooManyRequests)
		return
	}

	hash, ok := users[req.Username]
	if !ok {
		hash = dummyHash
	}

	err := bcrypt.CompareHashAndPassword(hash, []byte(req.Password))
	if !ok || err != nil {
		recordFailure(req.Username)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	resetFailures(req.Username)

	fmt.Fprintf(w, "welcome, %s\n", req.Username)
}

func main() {
	if err := seedUsers(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)

	srv := &http.Server{
		Addr:              ":8087",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("listening on :8087")
	log.Fatal(srv.ListenAndServe())
}

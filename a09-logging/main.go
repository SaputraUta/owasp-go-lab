package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var users = map[string][]byte{}
var dummyHash []byte
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

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
		logger.Warn("login blocked", "user", req.Username, "ip", r.RemoteAddr)
		http.Error(w, "too many attempts, try again later", http.StatusTooManyRequests)
		return
	}

	hash, ok := users[req.Username]
	if !ok {
		hash = dummyHash
	}
	err := bcrypt.CompareHashAndPassword(hash, []byte(req.Password))
	if !ok || err != nil {
		logger.Warn("login failed", "user", req.Username, "ip", r.RemoteAddr)
		if recordFailure(req.Username) {
			logger.Error("account locked", "user", req.Username, "ip", r.RemoteAddr)
		}
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	resetFailures(req.Username)
	logger.Info("login succeeded", "user", req.Username, "ip", r.RemoteAddr)
	fmt.Fprintf(w, "welcome, %s\n", req.Username)
}

func main() {
	if err := seedUsers(); err != nil {
		logger.Error("seed users failed", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)

	srv := &http.Server{
		Addr:              ":8089",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	logger.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

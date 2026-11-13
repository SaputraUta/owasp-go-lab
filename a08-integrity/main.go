package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Session struct {
	User string `json:"user"`
	Role string `json:"role"`
}

var sessionKey []byte

func sign(payload string) string {
	mac := hmac.New(sha256.New, sessionKey)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "user required", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(Session{User: user, Role: "user"})
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	payload := base64.RawURLEncoding.EncodeToString(data)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    payload + "." + sign(payload),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	fmt.Fprintln(w, "logged in")
}

func readSession(r *http.Request) (Session, error) {
	var s Session
	c, err := r.Cookie("session")
	if err != nil {
		return s, err
	}

	payload, sig, ok := strings.Cut(c.Value, ".")
	if !ok {
		return s, errors.New("malformed session")
	}
	if !hmac.Equal([]byte(sig), []byte(sign(payload))) {
		return s, errors.New("invalid signature")
	}

	data, err := base64.RawURLEncoding.DecodeString(payload)

	if err != nil {
		return s, err
	}

	err = json.Unmarshal(data, &s)
	return s, err
}

func admin(w http.ResponseWriter, r *http.Request) {
	s, err := readSession(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if s.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	fmt.Fprintln(w, "welcome to the admin panel")
}

func main() {
	sessionKey = []byte(os.Getenv("SESSION_KEY"))
	if len(sessionKey) < 32 {
		log.Fatal("SESSION_KEY must be at least 32 bytes")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("GET /admin", admin)

	srv := &http.Server{
		Addr:              ":8088",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("listening on :8088")
	log.Fatal(srv.ListenAndServe())
}

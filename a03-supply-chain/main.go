package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/language"
)

var matcher = language.NewMatcher([]language.Tag{language.English, language.Indonesian})

func greet(w http.ResponseWriter, r *http.Request) {
	tags, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	tag, _, _ := matcher.Match(tags...)
	if base, _ := tag.Base(); base.String() == "id" {
		fmt.Fprintln(w, "halo")
		return
	}
	fmt.Fprintln(w, "hello")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /greet", greet)

	srv := &http.Server{
		Addr:              ":8083",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("listening on :8083")
	log.Fatal(srv.ListenAndServe())
}

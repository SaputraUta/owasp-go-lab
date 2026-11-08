package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(t *testing.T) {
	t.Helper()
	if err := seedUsers(); err != nil {
		t.Fatal(err)
	}
	attempts = map[string]*attempt{}
}

func doLogin(user, pw string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, user, pw)
	rec := httptest.NewRecorder()
	login(rec, httptest.NewRequest("POST", "/login", strings.NewReader(body)))
	return rec
}

func TestLoginSuccess(t *testing.T) {
	setup(t)
	if rec := doLogin("alice", "sunshine"); rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestLoginNoUserEnumeration(t *testing.T) {
	setup(t)
	known := doLogin("alice", "wrong")
	unknown := doLogin("bob", "wrong")

	if known.Code != unknown.Code || known.Body.String() != unknown.Body.String() {
		t.Errorf("responses differ: alice=%d %q, bob=%d %q", known.Code, known.Body.String(), unknown.Code, unknown.Body.String())
	}
}

func TestLoginLockout(t *testing.T) {
	for _, user := range []string{"alice", "bob"} {
		t.Run(user, func(t *testing.T) {
			setup(t)
			for i := 0; i < maxFailures; i++ {
				if rec := doLogin(user, "wrong"); rec.Code != http.StatusUnauthorized {
					t.Fatalf("attempt %d: status = %d, want 401", i+1, rec.Code)
				}
			}
			if rec := doLogin(user, "sunshine"); rec.Code != http.StatusTooManyRequests {
				t.Errorf("after %d failures: status = %d, want 429", maxFailures, rec.Code)
			}
		})
	}
}

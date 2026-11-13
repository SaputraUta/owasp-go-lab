package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func loginCookie(t *testing.T, user string) string {
	t.Helper()
	rec := httptest.NewRecorder()
	login(rec, httptest.NewRequest("POST", "/login?user="+user, nil))
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("login set %d cookies, want 1", len(cookies))
	}
	return cookies[0].Value
}

func adminStatus(cookie string) int {
	req := httptest.NewRequest("GET", "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: cookie})
	rec := httptest.NewRecorder()
	admin(rec, req)
	return rec.Code
}

func TestAdminRejectsTamperedSession(t *testing.T) {
	sessionKey = bytes.Repeat([]byte("k"), 32)

	userCookie := loginCookie(t, "uta")
	_, userSig, _ := strings.Cut(userCookie, ".")

	adminPayload := base64.RawURLEncoding.EncodeToString([]byte(`{"user":"uta","role":"admin"}`))
	validAdmin := adminPayload + "." + sign(adminPayload)

	cases := []struct {
		name   string
		cookie string
		want   int
	}{
		{"real user cookie", userCookie, http.StatusForbidden},
		{"admin cookie signed by the server", validAdmin, http.StatusOK},
		{"admin payload without signature", adminPayload, http.StatusUnauthorized},
		{"admin payload with the user's signature", adminPayload + "." + userSig, http.StatusUnauthorized},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := adminStatus(tc.cookie); got != tc.want {
				t.Errorf("status = %d, want %d", got, tc.want)
			}
		})
	}
}

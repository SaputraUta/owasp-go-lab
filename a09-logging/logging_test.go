package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(t *testing.T) *bytes.Buffer {
	t.Helper()
	if err := seedUsers(); err != nil {
		t.Fatal(err)
	}
	attempts = map[string]*attempt{}

	var buf bytes.Buffer
	logger = slog.New(slog.NewJSONHandler(&buf, nil))
	return &buf
}

func doLogin(user, pw string) {
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, user, pw)
	login(httptest.NewRecorder(), httptest.NewRequest("POST", "/login", strings.NewReader(body)))
}

func logEntries(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var entries []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var e map[string]any
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("log line is not JSON: %q", line)
		}
		entries = append(entries, e)
	}
	return entries
}

func TestLoginEventsWithoutPassword(t *testing.T) {
	buf := setup(t)
	doLogin("alice", "Sunshine1")
	doLogin("alice", "sunshine")

	entries := logEntries(t, buf)
	if len(entries) != 2 || entries[0]["msg"] != "login failed" || entries[1]["msg"] != "login succeeded" {
		t.Fatalf("unexpected log entries: %v", entries)
	}
	if strings.Contains(buf.String(), "Sunshine1") {
		t.Errorf("password leaked into logs: %s", buf.String())
	}
}

func TestLogInjectionStaysOnOneLine(t *testing.T) {
	buf := setup(t)
	doLogin("x\n2026/11/14 17:30:00 login ok user=admin", "a")

	if entries := logEntries(t, buf); len(entries) != 1 {
		t.Errorf("got %d log lines, want 1: %s", len(entries), buf.String())
	}
}

func TestLockoutIsLogged(t *testing.T) {
	buf := setup(t)
	for range maxFailures {
		doLogin("alice", "wrong")
	}
	doLogin("alice", "sunshine")

	var locked, blocked int
	for _, e := range logEntries(t, buf) {
		switch e["msg"] {
		case "account locked":
			locked++
		case "login blocked":
			blocked++
		}
	}
	if locked != 1 || blocked != 1 {
		t.Errorf("account locked = %d, login blocked = %d, want 1 and 1", locked, blocked)
	}
}

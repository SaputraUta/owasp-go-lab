package main

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordStorage(t *testing.T) {
	users = map[string][]byte{}

	if err := register("alice", "summer2026"); err != nil {
		t.Fatal(err)
	}
	if err := register("bob", "summer2026"); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(users["alice"], users["bob"]) {
		t.Error("same password produced the same hash: no per-user salt")
	}

	cost, err := bcrypt.Cost(users["alice"])
	if err != nil || cost < bcrypt.DefaultCost {
		t.Errorf("stored hash is not bcrypt with cost >= %d (cost=%d, err=%v)", bcrypt.DefaultCost, cost, err)
	}

	if !login("alice", "summer2026") {
		t.Error("right password rejected")
	}
	if login("alice", "winter2026") {
		t.Error("wrong password accepted")
	}
}

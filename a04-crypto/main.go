package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var users = map[string][]byte{}

func register(name, pw string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	users[name] = h
	return nil
}

func login(name, pw string) bool {
	h, ok := users[name]
	if !ok {
		return false
	}
	return bcrypt.CompareHashAndPassword(h, []byte(pw)) == nil
}

func main() {
	seed := map[string]string{
		"alice": "summer2026",
		"bob":   "summer2026",
		"carol": "correct-horse-battery",
	}
	for name, pw := range seed {
		if err := register(name, pw); err != nil {
			log.Fatal(err)
		}
	}
	for name, h := range users {
		fmt.Println(name, string(h))
	}
	fmt.Println("alice, right password:", login("alice", "summer2026"))
	fmt.Println("alice, wrong password:", login("alice", "winter2026"))
}

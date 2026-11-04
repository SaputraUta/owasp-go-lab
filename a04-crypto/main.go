package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

var users = map[string]string{}

func hashPassword(pw string) string {
	sum := md5.Sum([]byte(pw))
	return hex.EncodeToString(sum[:])
}

func register(name, pw string) {
	users[name] = hashPassword(pw)
}

func main() {
	register("alice", "summer2026")
	register("bob", "summer2026")
	register("carol", "correct-horse-battery")
	for name, h := range users {
		fmt.Println(name, h)
	}
}

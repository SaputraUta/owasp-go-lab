package main

import (
	"sync"
	"time"
)

const (
	maxFailures     = 5
	lockoutDuration = 15 * time.Minute
)

type attempt struct {
	failures    int
	lockedUntil time.Time
}

var (
	mu       sync.Mutex
	attempts = map[string]*attempt{}
)

func isLocked(username string) bool {
	mu.Lock()
	defer mu.Unlock()

	a, ok := attempts[username]
	return ok && time.Now().Before(a.lockedUntil)
}

func recordFailure(username string) bool {
	mu.Lock()
	defer mu.Unlock()

	a, ok := attempts[username]
	if !ok {
		a = &attempt{}
		attempts[username] = a
	}
	a.failures++
	if a.failures >= maxFailures {
		a.lockedUntil = time.Now().Add(lockoutDuration)
		a.failures = 0
		return true
	}
	return false
}

func resetFailures(username string) {
	mu.Lock()
	defer mu.Unlock()

	delete(attempts, username)
}

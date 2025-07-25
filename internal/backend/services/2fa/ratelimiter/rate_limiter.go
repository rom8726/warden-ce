package ratelimiter

import (
	"sync"
	"time"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
)

var _ contract.TwoFARateLimiter = (*InMemoryTwoFARateLimiter)(nil)

type InMemoryTwoFARateLimiter struct {
	mu     sync.Mutex
	limits map[domain.UserID]*rateLimitEntry
	max    int
	window time.Duration
}

type rateLimitEntry struct {
	count     int
	firstTime time.Time
}

func New() *InMemoryTwoFARateLimiter {
	return &InMemoryTwoFARateLimiter{
		limits: make(map[domain.UserID]*rateLimitEntry),
		max:    5,
		window: time.Minute * 10,
	}
}

func (r *InMemoryTwoFARateLimiter) Inc(userID domain.UserID) (int, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.limits[userID]
	now := time.Now()
	if !ok || now.Sub(entry.firstTime) > r.window {
		r.limits[userID] = &rateLimitEntry{count: 1, firstTime: now}

		return 1, false
	}

	entry.count++
	if entry.count > r.max {
		return entry.count, true
	}

	return entry.count, false
}

func (r *InMemoryTwoFARateLimiter) Reset(userID domain.UserID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.limits, userID)
}

func (r *InMemoryTwoFARateLimiter) IsBlocked(userID domain.UserID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.limits[userID]
	if !ok {
		return false
	}

	if time.Since(entry.firstTime) > r.window {
		delete(r.limits, userID)

		return false
	}

	return entry.count > r.max
}

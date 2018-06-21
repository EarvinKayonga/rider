package httpx

import (
	"net/http"

	"golang.org/x/time/rate"
)

// example https://www.alexedwards.net/blog/how-to-rate-limit-http-requests

// A Limiter controls how frequently events are allowed to happen.
// It implements a "token bucket" of size b, initially full and
// refilled at rate r tokens per second.
//
// func NewLimiter(r Limit, b int) *Limiter
// the limiter permits you to consume an average of r tokens per second,
// with a maximum of b tokens in any single 'burst'.
// So in the code above our limiter allows 10 tokens to be consumed per second,
// with a maximum burst size of 20.
// var limiter = rate.NewLimiter(10, 20)

const (
	defaultLimit = 10
	defaultBurst = 20
)

// Limiter is a simple parameterized rate limiting implementation.
func Limiter(limit float64, burst int) func(next http.Handler) http.Handler {
	if limit < defaultLimit || burst < defaultBurst {
		limit = defaultLimit
		burst = defaultBurst
	}

	var limiter = rate.NewLimiter(rate.Limit(limit), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

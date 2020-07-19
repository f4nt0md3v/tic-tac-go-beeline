package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies middleware to a http.HandlerFunc
// Chain provides a cleaner interface for chaining middleware for single routes (implements Decorator pattern).
// Middleware functions are simple HTTP handlers (w http.ResponseWriter, r *http.Request).
// Here middleware functions are using custom httpx.ExtendedHandler.
// Note: order in which middleware functions supplied matters.
func Chain(f http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		f = m(f)
	}
	return f
}

package main

import (
	"net/http"
	"strings"
)

// corsMiddleware sets CORS headers for browser requests from local dev servers (any port).
// It does not short-circuit OPTIONS; chi routes handle OPTIONS explicitly so preflight gets 204, not 405.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowDevOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if req := r.Header.Get("Access-Control-Request-Headers"); req != "" {
			w.Header().Set("Access-Control-Allow-Headers", req)
		} else {
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		}
		next.ServeHTTP(w, r)
	})
}

func allowDevOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	lo := strings.ToLower(origin)
	return strings.HasPrefix(lo, "http://localhost:") ||
		strings.HasPrefix(lo, "http://127.0.0.1:")
}

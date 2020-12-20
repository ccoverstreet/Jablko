// middleware.go: General Purpose Middleware
// Cale Overstreet
// 2020/12/20
// Contains middleware not specific to MainApp functionality.

package middleware

import (
	"net/http"
	"time"
	"log"
)

func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call rest of handler stack
		next.ServeHTTP(w, r)

		end := time.Now()

		log.Printf("Request \"%s\" took %.3f ms\n", r.URL.Path, float32(end.Sub(start)) / 1000000)
	})
}

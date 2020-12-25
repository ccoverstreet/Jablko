// middleware.go: General Purpose Middleware
// Cale Overstreet
// 2020/12/20
// Contains middleware not specific to MainApp functionality.

package middleware

import (
	"net/http"
	"time"
	"github.com/ccoverstreet/Jablko/src/jlog"
)

func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call rest of handler stack
		next.ServeHTTP(w, r)

		end := time.Now()

		jlog.Printf("Request \"%s\" took %.3f ms\n", r.URL.Path, float32(end.Sub(start)) / 1000000)
	})
}

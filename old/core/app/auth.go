package app

import (
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler, core *JablkoCore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

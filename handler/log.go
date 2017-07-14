package handler

import (
	"log"
	"net/http"
)

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(IDKey).(string)
		log.Printf("[%s] %s", id, r.URL.String())
		h.ServeHTTP(w, r)
	})
}

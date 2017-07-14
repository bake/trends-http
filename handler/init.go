package handler

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

const IDKey key = "id"

func Init(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), IDKey, xid.New().String())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

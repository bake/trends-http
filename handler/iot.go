package handler

import (
	"context"
	"net/http"

	"github.com/BakeRolls/trends"
)

type key string

const IotKey key = "iot"
const NamesKey key = "qs"

func Iot(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qs, ok := r.URL.Query()["q"]
		if !ok || len(qs) == 0 {
			http.Error(w, "?q=foo&q=bar", http.StatusBadRequest)
			return
		}

		iot, err := trends.InterestOverTime(qs...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, IotKey, iot)
		ctx = context.WithValue(ctx, NamesKey, qs)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

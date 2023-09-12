package service

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type str string

func (s *service) WithMetricName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, `metricName`)

		var key str
		key = `metricName`
		ctx := context.WithValue(r.Context(), key, metricName)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *service) WithMetricValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricValue := chi.URLParam(r, `metricValue`)

		var key str
		key = `metricValue`
		ctx := context.WithValue(r.Context(), key, metricValue)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

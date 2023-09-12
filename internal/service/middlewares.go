package service

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *service) WithMetricName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, `metricName`)

		ctx := context.WithValue(r.Context(), `metricName`, metricName)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *service) WithMetricValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricValue := chi.URLParam(r, `metricValue`)

		ctx := context.WithValue(r.Context(), `metricValue`, metricValue)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

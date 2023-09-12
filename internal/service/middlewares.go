package service

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *service) WithMetricData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mdata := MetricData{
			mtype: chi.URLParam(r, `metricType`),
			name:  chi.URLParam(r, `metricName`),
			value: chi.URLParam(r, `metricValue`),
		}

		// https://gist.github.com/ww9/4ad7b2ddfb94816a30dfdf2218e02f48
		ctx := context.WithValue(r.Context(), s.metricDataContextKey, mdata)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

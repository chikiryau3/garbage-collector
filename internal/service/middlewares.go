package service

import (
	"net/http"
	"time"
)

func (s *service) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		s.log.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	})
}

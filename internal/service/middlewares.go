package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
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
			//"req hash", r.Header.Get("HashSHA256"),
			//"res hash", w.Header().Get("HashSHA256"),
		)
	})
}

type signWriter struct {
	w   http.ResponseWriter
	key string
}

func newSignWriter(w http.ResponseWriter, key string) *signWriter {
	return &signWriter{
		w,
		key,
	}
}

func (sw *signWriter) Write(p []byte) (int, error) {
	hash := hmac.New(sha256.New, []byte(sw.key))
	hash.Write(p)

	sw.w.Header().Add("HashSHA256", base64.URLEncoding.EncodeToString(hash.Sum(nil)))

	return sw.w.Write(p)
}

func (sw *signWriter) Header() http.Header {
	return sw.w.Header()
}

func (sw *signWriter) WriteHeader(statusCode int) {
	sw.w.WriteHeader(statusCode)
}

func (s *service) WithSignCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.Key == "" {
			next.ServeHTTP(w, r)
			return
		}

		hash := hmac.New(sha256.New, []byte(s.config.Key))
		var body bytes.Buffer
		_, err := body.ReadFrom(r.Body)
		if err != nil {
			s.log.Error("body parsing error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash.Write(body.Bytes())

		rSign, _ := base64.URLEncoding.DecodeString(r.Header.Get("HashSHA256"))

		s.log.Infoln("BODY", body.String(), "HEADER", r.Header.Get("HashSHA256"))
		s.log.Infoln("BODY HASH", hash.Sum(nil), "HEADER HASH", rSign)

		if !hmac.Equal(hash.Sum(nil), rSign) {
			s.log.Error("invalid request signature")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body.Bytes()))

		next.ServeHTTP(newSignWriter(w, s.config.Key), r)
	})
}

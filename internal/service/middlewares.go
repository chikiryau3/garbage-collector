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

//type gzipWriter struct {
//	http.ResponseWriter
//	Writer io.Writer
//}
//
//func (w gzipWriter) Write(b []byte) (int, error) {
//	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
//	return w.Writer.Write(b)
//}
//
//func (s *service) WithGzip(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
//			next.ServeHTTP(w, r)
//			return
//		}
//
//		// создаём gzip.Writer поверх текущего w
//		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
//		if err != nil {
//			io.WriteString(w, err.Error())
//			return
//		}
//		defer gz.Close()
//
//		w.Header().Set("Content-Encoding", "gzip")
//		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
//		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
//	})
//}

package rest

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		latency := time.Since(start).Nanoseconds()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		message := fmt.Sprintf("%v %v %v %v %v %v ns %v",
			ip,
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			latency,
			r.UserAgent())

		switch {
		case lrw.statusCode < 300:
			s.Info(message)
		case lrw.statusCode >= 300 && lrw.statusCode < 400:
			s.Warn(message)
		default:
			s.Error(message)
		}
	})
}

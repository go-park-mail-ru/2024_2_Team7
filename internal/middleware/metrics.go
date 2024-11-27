package middleware

import (
	"fmt"
	"kudago/internal/metrics"
	"net/http"
	"time"
)

func MetricsMiddleware(next http.Handler, serviceName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		fmt.Println(duration)
		metrics.RequestDuration.WithLabelValues(r.URL.Path, r.Method, serviceName, http.StatusText(ww.statusCode)).Observe(duration)
		metrics.RequestCount.WithLabelValues(r.URL.Path, r.Method, serviceName, http.StatusText(ww.statusCode)).Inc()

		if ww.statusCode >= 400 {
			metrics.ErrorCount.WithLabelValues(r.URL.Path, r.Method, serviceName, http.StatusText(ww.statusCode)).Inc()
		}
	})
}

package middleware

import (
	"net/http"
	"regexp"
	"time"

	"kudago/internal/metrics"
)

func MetricsMiddleware(next http.Handler, serviceName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		path := r.URL.Path
		re := regexp.MustCompile(`/(.*?)/\d+$`)
		if re.MatchString(r.URL.Path) {
			path = re.ReplaceAllString(r.URL.Path, "/$1/{id}")
		}

		if matched, _ := regexp.MatchString(`^/static/images/`, path); matched {
			path = "/static/images"
		}

		metrics.RequestDuration.WithLabelValues(path, r.Method, serviceName, http.StatusText(ww.statusCode)).Observe(duration)
		metrics.RequestCount.WithLabelValues(path, r.Method, serviceName, http.StatusText(ww.statusCode)).Inc()

		if ww.statusCode >= 400 {
			metrics.ErrorCount.WithLabelValues(path, r.Method, serviceName, http.StatusText(ww.statusCode)).Inc()
		}
	})
}

package middleware

import (
	"net/http"
	"time"

	"kudago/internal/http/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter, requestID string) *responseWriter {
	rw := responseWriter{w, http.StatusOK}
	rw.Header().Add("X-Request-ID", requestID)
	return &rw
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()

		r = r.WithContext(utils.SetRequestIDInContext(r.Context(), requestID))
		wrappedWriter := NewResponseWriter(w, requestID)
		next.ServeHTTP(wrappedWriter, r)
		utils.LogRequestData(r.Context(), logger, "http request", wrappedWriter.statusCode, r.Method, r.URL.Path, r.RemoteAddr, time.Since(start), nil)
	})
}

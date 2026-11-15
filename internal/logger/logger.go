package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func New() (*zap.SugaredLogger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return l.Sugar(), nil
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LogResponse(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(lrw, r)

			logger.Infow("Выполнен запрос",
				"method", r.Method,
				"uri", r.URL.String(),
				"status", lrw.statusCode,
				"size", lrw.size,
				"duration", time.Since(start),
			)
		})
	}
}

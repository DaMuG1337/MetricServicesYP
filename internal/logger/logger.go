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

var Log *zap.SugaredLogger

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Initialize() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	Log = logger.Sugar()
	return nil
}

func LogResponse(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		h.ServeHTTP(lrw, r)

		Log.Infow("Выполнен запрос",
			"Методом", r.Method,
			"URL", r.URL.String(),
			"За время", time.Since(start),
		)
	})
}

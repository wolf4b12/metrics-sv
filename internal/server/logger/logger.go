package logger

import (
    "bytes"
    "net/http"
    "time"
    "go.uber.org/zap"
)

// LoggingMiddleware для логирования запросов и ответов
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            startTime := time.Now()

            wrappedWriter := &loggingResponseWriter{
                ResponseWriter: w,
                body:           &bytes.Buffer{},
            }

            next.ServeHTTP(wrappedWriter, r)

            elapsedTime := time.Since(startTime)

            logger.Info(
                "Запрос",
                zap.String("URI", r.URL.Path),
                zap.String("Метод", r.Method),
                zap.Duration("Время выполнения", elapsedTime),
                zap.Int("Код статуса", wrappedWriter.statusCode),
                zap.Int("Размер ответа", wrappedWriter.body.Len()),
            )
        })
    }
}

// loggingResponseWriter оборачивает ResponseWriter для захвата кода статуса и тела ответа
type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
    body       *bytes.Buffer
}

// WriteHeader переопределён для захвата кода статуса
func (l *loggingResponseWriter) WriteHeader(code int) {
    l.statusCode = code
    l.ResponseWriter.WriteHeader(code)
}

// Write переопределён для записи тела ответа в буфер
func (l *loggingResponseWriter) Write(data []byte) (int, error) {
    n, err := l.body.Write(data)
    n2, err2 := l.ResponseWriter.Write(data[:n])
    if err2 != nil {
        return n2, err2
    }
    return n, err
}
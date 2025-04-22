package compress

import (
    "compress/gzip"
    "io"
    "net/http"

)

// Структура-обертка для ResponseWriter с поддержкой gzip
type GzipResponseWriter struct {
    http.ResponseWriter
    gz *gzip.Writer
}

// Write реализует интерфейс ResponseWriter для gzip
func (gw GzipResponseWriter) Write(b []byte) (int, error) {
    return gw.gz.Write(b)
}





// Middleware для обработки сжатия входящих запросов
func GzipRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Content-Encoding") == "gzip" {
            reader, err := gzip.NewReader(r.Body)
            if err != nil {
                http.Error(w, "Failed to decode gzip content", http.StatusBadRequest)
                return
            }
            defer reader.Close()
            r.Body = io.NopCloser(reader)
        }
        next.ServeHTTP(w, r)
    })
}

// Middleware для отправки сжатых ответов клиентам
func GzipResponse(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Accept-Encoding") == "gzip" {
            w.Header().Set("Content-Encoding", "gzip")
            gz := gzip.NewWriter(w)
            defer gz.Close()
            gzw := GzipResponseWriter{w, gz}
            next.ServeHTTP(gzw, r)
        } else {
            next.ServeHTTP(w, r)
        }
    })
}


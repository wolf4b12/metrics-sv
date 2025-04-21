package handlers

import (
    "compress/gzip"
    "net/http"
)

// CompressionMiddleware — middleware для сжатия ответов
func CompressionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        acceptEncodings := r.Header.Get("Accept-Encoding")
        if !contains(acceptEncodings, "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        w.Header().Set("Content-Encoding", "gzip")
        gzw := gzip.NewWriter(w)
        defer gzw.Close()

        gzwResponseWriter := gzipResponseWriter{w, gzw}
        next.ServeHTTP(gzwResponseWriter, r)
    })
}

// gzipResponseWriter — адаптер для gzip.Writer, чтобы использовать его как ResponseWriter
type gzipResponseWriter struct {
    http.ResponseWriter
    gzipWriter *gzip.Writer
}

func (gw gzipResponseWriter) Write(b []byte) (int, error) {
    return gw.gzipWriter.Write(b)
}

// contains проверяет наличие определенного значения в строке
func contains(haystack string, needle string) bool {
    return haystack == needle || (len(haystack) >= len(needle)+2 &&
        haystack[len(haystack)-len(needle)-1:] == ","+needle ||
        haystack[:len(needle)] == needle+",")
}
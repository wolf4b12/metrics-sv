package handlers

import (
//    "bytes"
    "compress/gzip"
    "io"
    "net/http"
)

// DecompressMiddleware — Middleware для автоматической декомпрессии тела запроса
func DecompressMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        contentEncoding := r.Header.Get("Content-Encoding")
        
        // Проверяем, пришло ли тело запроса в сжатом виде
        if contentEncoding == "gzip" || contentEncoding == "deflate" {
            r.Body, _ = decompressBody(r.Body, contentEncoding)
        }

        next.ServeHTTP(w, r)
    })
}

// decompressBody — декомпрессия тела запроса в зависимости от типа сжатия
func decompressBody(body io.Reader, encoding string) (io.ReadCloser, error) {
    switch encoding {
    case "gzip":
        reader, err := gzip.NewReader(body)
        if err != nil {
            return nil, err
        }
        return reader, nil
    default:
        return body.(io.ReadCloser), nil
    }
}
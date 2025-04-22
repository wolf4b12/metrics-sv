package srv

import (
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/wolf4b12/metrics-sv.git/internal/server/handlers"
    "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
    lgr "github.com/wolf4b12/metrics-sv.git/internal/server/logger" // Импортируем пакет логирования
    "go.uber.org/zap"
)

// Middleware для обработки сжатия входящих запросов
func gzipRequest(next http.Handler) http.Handler {
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
func gzipResponse(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Accept-Encoding") == "gzip" {
            w.Header().Set("Content-Encoding", "gzip")
            gz := gzip.NewWriter(w)
            defer gz.Close()
            gzw := gzipResponseWriter{w, gz}
            next.ServeHTTP(gzw, r)
        } else {
            next.ServeHTTP(w, r)
        }
    })
}

// Структура-обертка для ResponseWriter с поддержкой gzip
type gzipResponseWriter struct {
    http.ResponseWriter
    gz *gzip.Writer
}

// Write реализует интерфейс ResponseWriter для gzip
func (gw gzipResponseWriter) Write(b []byte) (int, error) {
    return gw.gz.Write(b)
}

type Server struct {
    router *chi.Mux
    server *http.Server
}

func NewServer(addr string) *Server {
    storage := storage.NewMemStorage()

    router := chi.NewRouter()

    // Инициализация логгера Zap
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Не удалось инициализировать логгер: %v", err)
    }

    // Применяем middleware для логирования
    router.Use(lgr.LoggingMiddleware(logger)) // Используем middleware из пакета logger
    router.Use(middleware.Logger)

    // Поддерживаем прием сжатых запросов
    router.Use(gzipRequest)

    // Включаем поддержку выдачи сжатых ответов
    router.Use(gzipResponse)

    // Маршруты остаются такими же
    router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateHandler(storage))
    router.Post("/update/", handlers.UpdateJSONHandler(storage))
    router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(storage))
    router.Post("/value/", handlers.PostJSONValueHandler(storage))
    router.Get("/", handlers.ListMetricsHandler(storage))

    return &Server{
        router: router,
        server: &http.Server{
            Addr:    addr,
            Handler: router,
        },
    }
}

func (s *Server) Run() error {
    log.Printf("Запуск сервера на http://%s\n", s.server.Addr)
    if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        return fmt.Errorf("не удалось запустить сервер: %w", err)
    }
    return nil
}
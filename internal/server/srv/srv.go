// internal/server/server.go
package srv

import (
    "fmt"
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
    "github.com/wolf4b12/metrics-sv.git/internal/server/handlers"
)

type Server struct {
    router *chi.Mux
    server *http.Server
}

func NewServer(addr string) *Server {
    storage := storage.NewMemStorage()

    router := chi.NewRouter()
    router.Use(middleware.Logger)
    router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateHandler(storage))
    router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(storage))
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
        return fmt.Errorf("failed to start server: %w", err)
    }
    return nil
}
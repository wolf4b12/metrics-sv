package srv

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/wolf4b12/metrics-sv.git/internal/server/handlers"
    "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
    lgr "github.com/wolf4b12/metrics-sv.git/internal/server/logger" // Импортируем пакет логирования
    "go.uber.org/zap"
    cm  "github.com/wolf4b12/metrics-sv.git/internal/server/compress"
)

type Server struct {
    router   *chi.Mux
    server   *http.Server
}

// Константа пути к файлу с метриками
const filePath = "/tmp/metrics.json"

// Конструктор сервера
func NewServer(addr string, loadOnStart bool, saveInterval time.Duration) *Server {
    storage := storage.NewMemStorage()

    // Загрузка метрик при старте
    if loadOnStart {
        err := storage.LoadFromFile(filePath)
        if err != nil {
            log.Printf("Не удалось загрузить предыдущие метрики: %v\n", err)
        } else {
            log.Println("Предыдущие метрики успешно загружены.")
        }
    }

    router := chi.NewRouter()

    // Логгер
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Не удалось инициализировать логгер: %v", err)
    }

    // Применение middleware
    router.Use(lgr.LoggingMiddleware(logger))
    router.Use(middleware.Logger)
    router.Use(cm.GzipRequest)
    router.Use(cm.GzipResponse)

    // Маршруты
    router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateHandler(storage))
    router.Post("/update/", handlers.UpdateJSONHandler(storage))
    router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(storage))
    router.Post("/value/", handlers.PostJSONValueHandler(storage))
    router.Get("/", handlers.ListMetricsHandler(storage))

    // Автосохранение метрик
    ticker := time.NewTicker(saveInterval)
    go func() {
        for range ticker.C {
            err := storage.SaveToFile(filePath)
            if err != nil {
                log.Printf("Ошибка при сохранении метрик: %v\n", err)
            } else {
                log.Println("Метрики успешно сохранены.")
            }
        }
    }()

    return &Server{
        router: router,
        server: &http.Server{
            Addr:    addr,
            Handler: router,
        },
    }
}

// Метод для плавного завершения работы сервера
func (s *Server) GracefulShutdown(ctx context.Context) error {
    log.Println("Начинаем плавное завершение работы сервера...")
    err := s.server.Shutdown(ctx)
    if err != nil {
        return fmt.Errorf("ошибка при завершении работы сервера: %w", err)
    }
    log.Println("Сервер завершил работу.")
    return nil
}

// Основной метод запуска сервера
func (s *Server) Run() error {
    log.Printf("Запуск сервера на адресе: %s\n", s.server.Addr)
    if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        return fmt.Errorf("не удалось запустить сервер: %w", err)
    }
    return nil
}
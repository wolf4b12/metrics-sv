package srv

import (
    "fmt"
    "log"
    "net/http"
//    "os"
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




// Создадим функцию, которая выбирает тип хранилища в зависимости от настроек
func chooseStorage(loadOnStart bool, filePath string) storage.Storage {
    if loadOnStart {
        fs, err := storage.NewFileStorage(filePath)
        if err != nil {
            log.Fatalf("Ошибка при создании FileStorage: %v", err)
        }
        return fs
    }
    return storage.NewMemStorage()
}

func NewServer(addr string, loadOnStart bool, saveInterval time.Duration, filePath string) *Server {
    // Выберите подходящее хранилище исходя из настроек
    storage := chooseStorage(loadOnStart, filePath)

    // Если выбрано FileStorage, восстановим старые метрики
    if loadOnStart {
        err := storage.LoadFromFile(filePath)
        if err != nil {
            log.Printf("Не удалось загрузить предыдущие метрики: %v\n", err)
        } else {
            log.Println("Предыдущие метрики успешно загружены.")
        }
    }










    router := chi.NewRouter()

    // Инициализируем логгер Zap
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Не удалось инициализировать логгер: %v", err)
    }

    // Применяем middleware для логирования
    router.Use(lgr.LoggingMiddleware(logger)) // Используем middleware из пакета logger
    router.Use(middleware.Logger)

    // Поддерживаем прием сжатых запросов
    router.Use(cm.GzipRequest)

    // Включаем поддержку выдачи сжатых ответов
    router.Use(cm.GzipResponse)

    // Маршруты остаются теми же
    router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateHandler(storage))
    router.Post("/update/", handlers.UpdateJSONHandler(storage))
    router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(storage))
    router.Post("/value/", handlers.PostJSONValueHandler(storage))
    router.Get("/", handlers.ListMetricsHandler(storage))

    // Начнем периодический процесс сохранения метрик
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

func (s *Server) Run() error {
    log.Printf("Запуск сервера на http://%s\n", s.server.Addr)
    if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        return fmt.Errorf("не удалось запустить сервер: %w", err)
    }
    return nil
}
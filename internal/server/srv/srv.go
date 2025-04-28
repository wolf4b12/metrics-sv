package srv

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "sync"
    "syscall"
    "os/signal"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/wolf4b12/metrics-sv.git/internal/server/handlers"
    "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
    lgr "github.com/wolf4b12/metrics-sv.git/internal/server/middlewares/logger" // Импортируем пакет логирования
    "go.uber.org/zap"
    cm  "github.com/wolf4b12/metrics-sv.git/internal/server/compress"

)

type Server struct {
    router   *chi.Mux
    server   *http.Server
    storage  *storage.MetricStorage
    ticker   *time.Ticker
    wg       sync.WaitGroup
}

// Запуск сервера
func NewServer(addr string, restore bool, storeInterval time.Duration, filePath string, ) *Server {
    // Создание KV-хранилища

    // Создание адаптера для работы с метриками
    metricStorage, err :=  storage.NewMetricStorage(restore, filePath) 

    if err != nil {
        log.Fatalf("Не удалось создать хранилище метрик: %v", err)
    }

    // Загрузка данных из файла при старте, если указано



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
    router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateHandler(metricStorage))
    router.Post("/update/", handlers.UpdateJSONHandler(metricStorage))
    router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(metricStorage))
    router.Post("/value/", handlers.PostJSONValueHandler(metricStorage))
    router.Get("/", handlers.ListMetricsHandler(metricStorage))



    // Создание сервера
    srv := &Server{
        router:   router,
        server: &http.Server{
            Addr:    addr,
            Handler: router,
        },
        storage: metricStorage,
    }

    // Запуск периодического сохранения, если интервал больше 0
    if storeInterval > 0 {
        srv.ticker = time.NewTicker(time.Duration(storeInterval) * time.Second)
        srv.wg.Add(1)
        go srv.periodicSave(filePath)
    }

    // Обработка сигналов для сохранения при завершении
    srv.wg.Add(1)
    go srv.handleSignals(filePath)

    return srv
}

// Периодическое сохранение метрик
func (s *Server) periodicSave(filePath string) {
    defer s.wg.Done()
    for range s.ticker.C {
        err := s.storage.SaveToFile(filePath)
        if err != nil {
            log.Printf("Ошибка при сохранении метрик: %v\n", err)
        } else {
            log.Println("Метрики успешно сохранены.")
        }
    }
}

// Обработка сигналов для сохранения при завершении
func (s *Server) handleSignals(filePath string) {
    defer s.wg.Done()
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    log.Println("Получен сигнал завершения. Сохранение метрик...")
    err := s.storage.SaveToFile(filePath)
    if err != nil {
        log.Printf("Ошибка при сохранении метрик: %v\n", err)
    } else {
        log.Println("Метрики успешно сохранены.")
    }
    s.server.Shutdown(nil)
}

func (s *Server) Run() error {
    log.Printf("Запуск сервера на http://%s\n", s.server.Addr)
    if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        return fmt.Errorf("не удалось запустить сервер: %w", err)
    }
    return nil
}
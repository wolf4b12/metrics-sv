package srv

import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
//    "fmt"
    "log"
    "net/http"
    "time"
    "sync"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/wolf4b12/metrics-sv/internal/server/handlers"
    "github.com/wolf4b12/metrics-sv/internal/server/storage"
    lgr "github.com/wolf4b12/metrics-sv/internal/server/middlewares/logger" // Импортируем пакет логирования 
    "go.uber.org/zap"
    cm  "github.com/wolf4b12/metrics-sv/internal/server/compress"

)

type Server struct {
    router   *chi.Mux
    server   *http.Server
    storage  *storage.MetricStorage
    ticker   *time.Ticker
    db       *sql.DB
    wg       sync.WaitGroup
}

// Запуск сервера
func NewServer(addr string, restore bool, storeInterval time.Duration, filePath string, dbDSN string) (*Server) {
    // Создание KV-хранилища
    kv := storage.NewKVStorage()
    // Создание адаптера для работы с метриками
    metricStorage, err :=  storage.NewMetricStorage(kv, restore, storeInterval, filePath) 

    if err != nil {
        log.Fatalf("Не удалось создать хранилище метрик: %v", err)
    }
    // Создание подключения к базе данных
    db, _ := sql.Open("postgres", dbDSN)
//    if err != nil {
//        fmt.Printf("не удалось подключиться к базе данных")
//    }
    // Проверка соединения с базой данных
//    if err := db.Ping(); err != nil {
//        fmt.Printf("не удалось проверить соединение с базой данных")
//    }
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
    router.Get("/ping", handlers.PingHandler(db))

    // Создание сервера
    srv := &Server{
        router:   router,
        server: &http.Server{
            Addr:    addr,
            Handler: router,
        },
        storage: metricStorage,
        db:      db,
    }
    return srv
}


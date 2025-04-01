package main

import (
    "flag"
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/wolf4b12/metrics-sv.git/internal/storage"
    handler "github.com/wolf4b12/metrics-sv.git/internal/handlers"
)

func main() {
    // Определение флага для адреса сервера
    addr := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")

    // Парсинг флагов
    flag.Parse()

    storage := storage.NewMemStorage()

    // Создание нового роутера с использованием chi
    router := chi.NewRouter()

    // Настройка middleware для журналирования запросов
    router.Use(middleware.Logger)

    // Маршрут для обновления метрик
    router.Post("/update/{metricType}/{metricName}/{metricValue}", handler.UpdateHandler(storage))

    // Маршрут для получения конкретной метрики
    router.Get("/value/{metricType}/{metricName}", handler.ValueHandler(storage))

    // Маршрут для получения списка всех метрик
    router.Get("/", handler.ListMetricsHandler(storage))

    server := &http.Server{
        Addr:    *addr,
        Handler: router,
    }

    log.Printf("Starting server on http://%s\n", *addr)
    log.Fatal(server.ListenAndServe())
}

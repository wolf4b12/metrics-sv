package main

import (
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/wolf4b12/metrics-sv.git/internal/storage"
    handler "github.com/wolf4b12/metrics-sv.git/internal/handlers"
)

func main() {
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
        Addr:    "localhost:8080",
        Handler: router,
    }

    log.Printf("Starting server on http://localhost:8080\n")
    log.Fatal(server.ListenAndServe())
}
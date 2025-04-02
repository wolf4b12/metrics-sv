package main

import (
    "flag"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    storage  "github.com/wolf4b12/metrics-sv.git/internal/storage"
    handler "github.com/wolf4b12/metrics-sv.git/internal/handlers"

)

func main() {
    // Определение флага для адреса сервера
    addr := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")

    // Парсинг флагов
    flag.Parse()

    // Проверка наличия неизвестных флагов
    if flag.NArg() > 0 {
        log.Fatalf("Неизвестный флаг: %s\n", os.Args[flag.NArg()-1])
    }

    storage := storage.NewMemStorage()

    // Создание нового роутера с использованием chi
    router := chi.NewRouter()

    // Настройка middleware для журналирования запросов
    router.Use(middleware.Logger)

    // Маршруты для обработки запросов
    router.Post("/update/{metricType}/{metricName}/{metricValue}", handler.UpdateHandler(storage))
    router.Get("/value/{metricType}/{metricName}", handler.ValueHandler(storage))
    router.Get("/", handler.ListMetricsHandler(storage))

    server := &http.Server{
        Addr:    *addr,
        Handler: router,
    }

    log.Printf("Запуск сервера на http://%s\n", *addr)
    log.Fatal(server.ListenAndServe())
}
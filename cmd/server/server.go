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
    var addr string

    // Проверяем наличие переменной окружения ADDRESS
    if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
        addr = envAddr
        log.Println("Используется переменная окружения ADDRESS:", addr)
    } else {
        // Если переменная окружения не найдена, проверяем флаг командной строки
        flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
        defaultAddr := "localhost:8080"
        flagSet.StringVar(&addr, "a", defaultAddr, "адрес эндпоинта HTTP-сервера")

        // Парсим флаги
        err := flagSet.Parse(os.Args[1:])
        if err != nil {
            log.Fatalf("Ошибка парсинга флагов: %v", err)
        }

        // Проверяем наличие неизвестных флагов
        if flagSet.NArg() > 0 {
            log.Fatalf("Неизвестный флаг: %s\n", flagSet.Arg(flagSet.NArg()-1))
        }

        if addr == "" {
            // Если флаг не был указан, используем значение по умолчанию
            addr = defaultAddr
            log.Println("Переменная окружения ADDRESS не найдена, используется значение по умолчанию:", addr)
        } else {
            log.Println("Используется флаг командной строки -a:", addr)
        }
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
        Addr:    addr,
        Handler: router,
    }

    log.Printf("Запуск сервера на http://%s\n", addr)
    log.Fatal(server.ListenAndServe())
}
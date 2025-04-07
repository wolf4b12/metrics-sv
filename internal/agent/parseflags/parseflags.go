package parseflags

import (
    "flag"
    "log"
    "os"
    "strconv"
    "time"
)


func ParseFlags() (time.Duration, time.Duration, string) {
    // Определяем дефолтные значения для переменных
    var (
        addr           string
        reportInterval int
        pollInterval   int
    )

    // Читаем переменные окружения
    if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
        addr = envAddr
        log.Println("Использована переменная окружения ADDRESS:", addr)
    } else {
        // Переменная окружения не найдена, читаем аргумент командной строки
        flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")
    }

    if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
        var err error
        reportInterval, err = strconv.Atoi(envReportInterval)
        if err != nil {
            log.Fatalf("Неверное значение REPORT_INTERVAL: %v", err)
        }
        log.Println("Использована переменная окружения REPORT_INTERVAL:", reportInterval)
    } else {
        // Переменная окружения не найдена, читаем аргумент командной строки
        flag.IntVar(&reportInterval, "r", 10, "Интервал отправки метрик в секундах")
    }

    if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
        var err error
        pollInterval, err = strconv.Atoi(envPollInterval)
        if err != nil {
            log.Fatalf("Неверное значение POLL_INTERVAL: %v", err)
        }
        log.Println("Использована переменная окружения POLL_INTERVAL:", pollInterval)
    } else {
        // Переменная окружения не найдена, читаем аргумент командной строки
        flag.IntVar(&pollInterval, "p", 2, "Интервал сбора метрик в секундах")
    }

    // Парсим флаги
    flag.Parse()

    // Проверка наличия неизвестных флагов
    if flag.NArg() > 0 {
        log.Fatalf("Неизвестный флаг или аргумент: %v", flag.Args())
    }

    return time.Duration(pollInterval) * time.Second,
        time.Duration(reportInterval) * time.Second,
        addr
}

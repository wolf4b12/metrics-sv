package main

import (
    "log"
    "time"

    server   "github.com/wolf4b12/metrics-sv.git/internal/server/srv"
    config   "github.com/wolf4b12/metrics-sv.git/internal/server/config"
)

func main() {
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("ошибка при создании конфигурации: %v", err)
    }

    // Настройки загрузки и интервала сохранения
    loadPrevious := true           // Нужно ли подгружать старые метрики при старте
    interval := 5 * time.Minute    //  Интервал сохранения метрик (каждый 5 минут)

    srv := server.NewServer(cfg.GetAddr(), loadPrevious, interval)

    err = srv.Run()
    if err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
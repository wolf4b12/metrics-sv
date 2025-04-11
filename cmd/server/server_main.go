package main

import (
    "log"
    server   "github.com/wolf4b12/metrics-sv.git/internal/server/srv"
    config   "github.com/wolf4b12/metrics-sv.git/internal/server/config"
)



func main() {

    cfg, err := config.NewConfig()

    if err != nil {
        log.Fatalf("ошибка при создании конфигурации: %v", err)
    }

    
    
    srv := server.NewServer(cfg.GetAddr()) // Используем метод GetAddr()

    err = srv.Run()
    if err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
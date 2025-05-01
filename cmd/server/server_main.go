package main

import (
    "log"

    srv "github.com/wolf4b12/metrics-sv/internal/server/srv"
    config "github.com/wolf4b12/metrics-sv/internal/server/config"
)

func main() {
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("ошибка при создании конфигурации: %v", err)
    }

    // Обращаемся непосредственно к полю конфигурации
    server  := srv.NewServer(
        cfg.GetAddr(),                   // Адрес прослушивания
        cfg.IsRestoreEnabled(),          // Нужна ли загрузка предыдущих метрик
        cfg.GetStoreInterval(),          // Интервал сохранения метрик
        cfg.GetFileStoragePath(),        // путь до файла с метриками
    )
    if err != nil {
        log.Fatalf("ошибка при создании сервера: %v", err)
    }

    err = server.Run()
    if err != nil {
        log.Fatalf("ошибка при запуске сервера: %v", err)
    }
}
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    srv "github.com/wolf4b12/metrics-sv.git/internal/server/srv"
    config "github.com/wolf4b12/metrics-sv.git/internal/server/config"
    "time"
)

func main() {
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("ошибка при создании конфигурации: %v", err)
        os.Exit(1)
    }

    server := srv.NewServer(
        cfg.GetAddr(),
        cfg.IsRestoreEnabled(),
        cfg.GetStoreInterval(),
    )
    if err != nil {
        log.Fatalf("ошибка при создании сервера: %v", err)
 os.Exit(1)
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-done
        log.Println("получил сигнал завершения, пытаюсь закрыть аккуратно...")
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        if err := server.GracefulShutdown(ctx); err != nil {
            log.Fatalf("ошибка graceful shutdown: %v", err)
        }
    }()

    if err := server.Run(); err != nil {
        log.Fatalf("ошибка при запуске сервера: %v", err)
        os.Exit(1)
    }
}
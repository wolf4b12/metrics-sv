package main

import (
    "flag"
    "log"
//    "net/http"
    "os"

//    "github.com/go-chi/chi/v5"
//    "github.com/go-chi/chi/v5/middleware"

//    storage  "github.com/wolf4b12/metrics-sv.git/internal/storage"
//    handler "github.com/wolf4b12/metrics-sv.git/internal/handlers"
    server   "github.com/wolf4b12/metrics-sv.git/internal/server"
)

func main() {
    var addr string

    if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
        addr = envAddr
        log.Println("Используется переменная окружения ADDRESS:", addr)
    } else {
        flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
        defaultAddr := "localhost:8080"
        flagSet.StringVar(&addr, "a", defaultAddr, "адрес эндпоинта HTTP-сервера")

        err := flagSet.Parse(os.Args[1:])
        if err != nil {
            log.Fatalf("Ошибка парсинга флагов: %v", err)
        }

        if flagSet.NArg() > 0 {
            log.Fatalf("Неизвестный флаг: %s\n", flagSet.Arg(flagSet.NArg()-1))
        }

        if addr == "" {
            addr = defaultAddr
            log.Println("Переменная окружения ADDRESS не найдена, используется значение по умолчанию:", addr)
        } else {
            log.Println("Используется флаг командной строки -a:", addr)
        }
    }


    srv := server.NewServer(addr)
    err := srv.Run()
    if err != nil {
        log.Fatalf("Failed to run server: %v", err)

    }
}
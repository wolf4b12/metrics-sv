package main

import (
    "log"
    "net/http"
    storage "github.com/wolf4b12/metrics-sv.git/internal/storage"
    handler "github.com/wolf4b12/metrics-sv.git/internal/handlers"
)

func main() {
    storage := storage.NewMemStorage()
    mux := http.NewServeMux() // создаем новый мультиплексор
    
    // Инициализация маршрутов
    mux.Handle("/update/", handler.UpdateHandler(storage))

    server := &http.Server{
        Addr:    "localhost:8080",
        Handler: mux,
    }

    log.Printf("Starting server on http://localhost:8080\n")
    log.Fatal(server.ListenAndServe())
}

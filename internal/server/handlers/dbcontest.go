package handlers

import (
    "database/sql"
//    _ "github.com/jackc/pgx/v4/stdlib" // Подключение драйвера pgx
    "log"
    "net/http"
 //   "fmt"
)

// PingHandler обработчик для проверки соединения с базой данных
func PingHandler(s string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        
        // Формирование строки подключения

    
        // Открытие подключения к базе данных
        db, err := sql.Open("pgx", s)
        defer db.Close() // Закрываем подключение после завершения обработки запроса

        if err != nil {
            log.Printf("Ошибка открытия соединения: %v\n", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // Проверка доступности базы данных путем ping-запроса
        err = db.Ping()
        if err != nil {
            log.Printf("Ошибка ping'а базы данных: %v\n", err)
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }

        // Если всё прошло хорошо — возвращаем статус 200 OK
        w.WriteHeader(http.StatusOK)
    }
}
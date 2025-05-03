package handlers

import (
    "database/sql"
    "log"
    "net/http"
)

// PingHandler обработчик для проверки соединения с базой данных
func PingDataBase(s string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        
    
        // Открытие подключения к базе данных
        db, err := sql.Open("pgx", s)

        if err != nil {
            log.Printf("Ошибка открытия соединения: %v\n", err)
            // Если соединения не установлено возвращаем статус 500 OK
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
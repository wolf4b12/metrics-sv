package handlers

import (
    "database/sql"
    "net/http"
)

// PingHandler обработчик для проверки соединения с базой данных
func PingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            http.Error(w, "Не удалось проверить соединение с базой данных", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}

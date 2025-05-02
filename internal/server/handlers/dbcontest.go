package handlers

import (
    "database/sql"
    "net/http"
	"fmt"
)

// PingHandler обработчик для проверки соединения с базой данных
func PingHandler(s string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

           // Создание подключения к базе данных
   

    ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
    s, `test`, `XXXXXXXX`, `test`)
 
	db, err := sql.Open("pgx", ps)
    if err != nil {
    panic(err)
   }

    if err := db.Ping(); err != nil {
            http.Error(w, "Не удалось проверить соединение с базой данных", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}

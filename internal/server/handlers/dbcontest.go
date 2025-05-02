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
 	_, err := sql.Open("pgx", ps)
	
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
    }  else  {w.WriteHeader(http.StatusOK)
	}
	} 
}

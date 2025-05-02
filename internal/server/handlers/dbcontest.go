package handlers

import (
    "database/sql"
    "net/http"
//	"fmt"
)

// PingHandler обработчик для проверки соединения с базой данных
func PingHandler(s string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

    // Создание подключения к базе данных
   
   
 	_, err := sql.Open("pgx", s)
	
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
    }  else  {w.WriteHeader(http.StatusOK)
	}
	} 
}

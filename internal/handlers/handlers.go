package handlers

import (
//   "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "github.com/go-chi/chi/v5"
    storage "github.com/wolf4b12/metrics-sv.git/internal/storage" // Импортируем пользовательский пакет storage
)





// UpdateHandler — обработчик для обновления метрик
func UpdateHandler(storage storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        pathParts := strings.Split(r.URL.Path, "/")[2:]
        if len(pathParts) != 3 {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        metricType, metricName, metricValue := pathParts[0], pathParts[1], pathParts[2]

        switch metricType {
        case "gauge":
            value, err := strconv.ParseFloat(metricValue, 64)
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            storage.UpdateGauge(metricName, value)

        case "counter":
            value, err := strconv.ParseInt(metricValue, 10, 64)
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            storage.UpdateCounter(metricName, value)

        default:
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}

// ValueHandler обработчик для получения значения метрики
func ValueHandler(storage storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        metricType := chi.URLParam(r, "metricType")
        metricName := chi.URLParam(r, "metricName")

        var value any
        var err error

        switch metricType {

        case "gauge":
            value, err = storage.GetGauge(metricName)
            if err == nil {
                w.WriteHeader(http.StatusOK)
                return
            }
            

        case "counter":
            
            value, err = storage.GetCounter(metricName)

            if err == nil {
                w.WriteHeader(http.StatusOK)
                return
            }
            
          
        }



           if err != nil {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "Error getting metric: %s/%s", metricType, metricName)
            return
        }

   //     w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, value)
    }
}

// ListMetricsHandler обработчик для получения списка всех метрик
func ListMetricsHandler(storage storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")

        allMetrics := storage.AllMetrics()

        html := "<html><body>"
        for metricType, metrics := range allMetrics {
            html += fmt.Sprintf("<h2>%s Metrics:</h2>", metricType)
            for metricName, value := range metrics {
                html += fmt.Sprintf("<p>%s: %v</p>", metricName, value)
            }
        }
        html += "</body></html>"

        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, html)
    }
}
package handlers

import (
    "fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/wolf4b12/metrics-sv.git/internal/server/storage" // Импортируем пользовательский пакет storage
)



// ValueHandler обработчик для получения значения метрики
func ValueHandler(storage storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        metricType := chi.URLParam(r, "metricType")
        metricName := chi.URLParam(r, "metricName")

        var value interface{}
        var err error

        switch metricType {
        case "gauge":
            value, err = storage.GetGauge(metricName)
        case "counter":
            value, err = storage.GetCounter(metricName)
        default:
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "Unknown metric type: %s", metricType)
            return
        }

        if err == storage.ErrMetricNotFound() {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "Metric not found: %s/%s", metricType, metricName)
            return
        } else if err != nil {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "Error getting metric: %s/%s", metricType, metricName)
            return
        }

        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, value)
    }
}


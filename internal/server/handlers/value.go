package handlers

import (
    "fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/wolf4b12/metrics-sv.git/internal/constant" // Импортируем константы
)

// GetStorage интерфейс для получения метрик
type GetStorage interface {
    GetGauge(name string) (float64, error)
    GetCounter(name string) (int64, error)
}

// ValueHandler обработчик для получения значения метрики
func ValueHandler(storage GetStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        metricType := chi.URLParam(r, "metricType")
        metricName := chi.URLParam(r, "metricName")

        var value interface{}
        var err error

        switch metricType {
        case constant.MetricTypeGauge:
            value, err = storage.GetGauge(metricName)
        case constant.MetricTypeCounter:
            value, err = storage.GetCounter(metricName)
        default:
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "Unknown metric type: %s", metricType)
            return
        }

        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "Error getting metric: %s/%s", metricType, metricName)
            return
        }

        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, value)
    }
}
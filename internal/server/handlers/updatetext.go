package handlers

import (
    "net/http"
    "strings"
    "strconv"
    "github.com/wolf4b12/metrics-sv/internal/constant" // Импортируем константы
)

// UpdateStorage интерфейс для обновления метрик
type UpdateStorage interface {
    UpdateGauge(name string, value float64)
    UpdateCounter(name string, value int64)
}

// UpdateHandler — обработчик для обновления метрик
func UpdateHandler(storage UpdateStorage) http.HandlerFunc {
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
        case constant.MetricTypeGauge:
            value, err := strconv.ParseFloat(metricValue, 64)
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            storage.UpdateGauge(metricName, value)

        case constant.MetricTypeCounter:
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






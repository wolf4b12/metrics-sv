package handlers

import (
    "fmt"
    "net/http"

)

// AllMetricsStorage интерфейс для получения всех метрик
type allMetricsStorage interface {
    AllMetrics() map[string]map[string]interface{}
}

// ListMetricsHandler обработчик для получения списка всех метрик
func ListMetricsHandler(storage allMetricsStorage) http.HandlerFunc {
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
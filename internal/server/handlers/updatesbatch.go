package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/wolf4b12/metrics-sv/internal/server/metricssrv" // Импортируем структуру с метриками
)



// BatchMetrics структура для представления набора метрик
type BatchMetrics struct {
    Metrics []metricssrv.Metrics `json:"metrics"`
}

// UpdatesHandler обработчик для обновления множества метрик
func UpdatesHandler(storage UpdateStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Чтение тела запроса
        var batch BatchMetrics
        if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
            http.Error(w, "Не удалось декодировать тело запроса", http.StatusBadRequest)
            return
        }

        // Обходим каждую метрику в батче
        for _, metric := range batch.Metrics {
            switch metric.MType {
            case "gauge":
                if metric.Value == nil {
                    http.Error(w, "Отсутствует значение для gauge-метрики", http.StatusBadRequest)
                    return
                }
                storage.UpdateGauge(metric.ID, *metric.Value)
            case "counter":
                if metric.Delta == nil {
                    http.Error(w, "Отсутствует дельта для counter-метрики", http.StatusBadRequest)
                    return
                }
                storage.UpdateCounter(metric.ID, *metric.Delta);
            default:
                http.Error(w, "Недопустимый тип метрики", http.StatusBadRequest)
                return
            }
        }

        // Всё прошло успешно
        w.WriteHeader(http.StatusOK)
    }
}
package handlers

import (
    "net/http"
    "io"
    "encoding/json"
    "github.com/wolf4b12/metrics-sv/internal/constant" // Импортируем константы
    "github.com/wolf4b12/metrics-sv/internal/server/metricssrv" //Импортируем серуктуру с метриками
)

// UpdateStorage интерфейс для обновления метрик

func UpdateJSONHandler(storage UpdateStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if r.Method != http.MethodPost {
            http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
            return
        }

        // Читаем тело запроса
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
            return
        }

        // Декодируем JSON как одиночная метрика (а не массив)
        var receivedMetric metricssrv.Metrics
        err = json.Unmarshal(body, &receivedMetric)
        if err != nil {
            http.Error(w, "Неверная структура JSON", http.StatusBadRequest)
            return
        }

        // Обрабатываем одну метрику
        switch receivedMetric.MType {
        case constant.MetricTypeGauge:
            if receivedMetric.Value == nil {
                http.Error(w, "'value' отсутствует для gauge-метрики", http.StatusBadRequest)
                return
            }
            storage.UpdateGauge(receivedMetric.ID, *receivedMetric.Value)

        case constant.MetricTypeCounter:
            if receivedMetric.Delta == nil {
                http.Error(w, "'delta' отсутствует для counter-метрики", http.StatusBadRequest)
                return
            }
            storage.UpdateCounter(receivedMetric.ID, *receivedMetric.Delta)

        default:
            http.Error(w, "Тип метрики неизвестен", http.StatusBadRequest)
            return
        }

        // Формируем ответ
        respData, err := json.Marshal(receivedMetric)
        if err != nil {
            http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
            return
        }

        w.Write(respData)
        w.WriteHeader(http.StatusOK)
    }
}
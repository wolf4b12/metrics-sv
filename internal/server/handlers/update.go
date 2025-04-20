package handlers

import (
    "net/http"
    "strings"
    "strconv"
    "io"
    "encoding/json"
    "github.com/wolf4b12/metrics-sv.git/internal/constant" // Импортируем константы
    "github.com/wolf4b12/metrics-sv.git/internal/server/metricssrv" //Импортируем серуктуру с метриками
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





func UpdateJSONHandler(storage UpdateStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if r.Method != http.MethodPost {
            http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
            return
        }

        // Читаем тело запроса
        body, err := io.ReadAll(r.Body)
        defer r.Body.Close()
        if err != nil {
            http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
            return
        }

        // Декодируем JSON как массив метрик
        var receivedMetrics []metricssrv.Metrics
        err = json.Unmarshal(body, &receivedMetrics)
        if err != nil {
            http.Error(w, "Неверная структура JSON", http.StatusBadRequest)
            return
        }

        // Обрабатываем каждую метрику в массиве
        for _, metric := range receivedMetrics {
            switch metric.MType {
            case constant.MetricTypeGauge:
                if metric.Value == nil {
                    http.Error(w, "'value' отсутствует для gauge-метрики", http.StatusBadRequest)
                    return
                }
                storage.UpdateGauge(metric.ID, *metric.Value)

            case constant.MetricTypeCounter:
                if metric.Delta == nil {
                    http.Error(w, "'delta' отсутствует для counter-метрики", http.StatusBadRequest)
                    return
                }
                storage.UpdateCounter(metric.ID, *metric.Delta)

            default:
                http.Error(w, "Тип метрики неизвестен", http.StatusBadRequest)
                return
            }
        }

        // Формирование ответа
        respData, err := json.Marshal(receivedMetrics)
        if err != nil {
            http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
            return
        }

        w.Write(respData)
    }
}
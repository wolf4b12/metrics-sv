package handlers

import (
    "encoding/json"
    "net/http"
//    "github.com/go-chi/chi/v5"
    "github.com/wolf4b12/metrics-sv.git/internal/constant"     // Импортируем константы
    "github.com/wolf4b12/metrics-sv.git/internal/server/metrics_srv" // Импортируем структуру метрик
)

// Интерфейс для получения значений метрик
type GetStorage interface {
    GetGauge(name string) (float64, error)
    GetCounter(name string) (int64, error)
}

// Обработчик для получения метрики по имени и типу
func ValueHandler(storage GetStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json") // Устанавливаем Content-Type

        if r.Method != http.MethodPost {
            http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
            return
        }

        // Читаем тело запроса
        var requestMetric metrics_srv.Metrics
        err := json.NewDecoder(r.Body).Decode(&requestMetric)
        if err != nil || requestMetric.ID == "" || requestMetric.MType == "" {
            http.Error(w, "Некорректный запрос", http.StatusBadRequest)
            return
        }

        // Определяем переменную для хранения результата
        var responseMetric metrics_srv.Metrics
        responseMetric.ID = requestMetric.ID
        responseMetric.MType = requestMetric.MType

        switch requestMetric.MType {
        case constant.MetricTypeGauge:
            value, errResult := storage.GetGauge(requestMetric.ID)
            if errResult != nil {
                http.Error(w, "Ошибка получения метрики", http.StatusNotFound)
                return
            }
            responseMetric.Value = &value // Передаем реальную переменную

        case constant.MetricTypeCounter:
            value, errResult := storage.GetCounter(requestMetric.ID)
            if errResult != nil {
                http.Error(w, "Ошибка получения метрики", http.StatusNotFound)
                return
            }
            responseMetric.Delta = &value // Передаем реальную переменную

        default:
            http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
            return
        }

        // Кодируем ответ в JSON и отправляем клиенту
        jsonResponse, _ := json.Marshal(responseMetric)
        w.Write(jsonResponse)
    }
}
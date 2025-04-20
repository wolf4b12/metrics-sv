package handlers

import (
    "encoding/json"
    "net/http"
//    "strings"
    "io"
//    "github.com/go-chi/chi/v5"
    "github.com/wolf4b12/metrics-sv.git/internal/constant"     // Импортируем константы
    "github.com/wolf4b12/metrics-sv.git/internal/server/metricssrv" // Импортируем структуру метрик
)

// GetStorage интерфейс для получения метрик
type GetStorage interface {
    GetGauge(name string) (float64, error)
    GetCounter(name string) (int64, error)
}

// ValueHandler обработчик для получения значения метрики
func ValueHandler(storage GetStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Установка заголовков ответа
        w.Header().Set("Content-Type", "application/json")

        // Разрешаем только POST-запросы
        if r.Method != http.MethodPost {
            http.Error(w, "Only POST method is allowed.", http.StatusMethodNotAllowed)
            return
        }

        // Чтение тела запроса
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Error reading request body.", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        // Парсим тело запроса в структуру Metric
        var inputMetric metricssrv.Metrics
        err = json.Unmarshal(body, &inputMetric)
        if err != nil {
            http.Error(w, "Invalid JSON format.", http.StatusBadRequest)
            return
        }

        // Получаем метрики из хранилища в зависимости от типа
        var outputMetric metricssrv.Metrics
        outputMetric.ID = inputMetric.ID
        outputMetric.MType = inputMetric.MType

        switch inputMetric.MType {
        case constant.MetricTypeGauge:
            val, err := storage.GetGauge(inputMetric.ID)
            if err != nil {
                http.Error(w, "Error fetching gauge metric.", http.StatusNotFound)
                return
            }
            outputMetric.Value = &val

        case constant.MetricTypeCounter:
            val, err := storage.GetCounter(inputMetric.ID)
            if err != nil {
                http.Error(w, "Error fetching counter metric.", http.StatusNotFound)
                return
            }
            outputMetric.Delta = &val

        default:
            http.Error(w, "Unsupported metric type.", http.StatusBadRequest)
            return
        }

        // Кодируем ответ в JSON и отправляем клиенту
        respBody, err := json.Marshal(outputMetric)
        if err != nil {
            http.Error(w, "Error encoding response.", http.StatusInternalServerError)
            return
        }

        w.Write(respBody)
    }
}
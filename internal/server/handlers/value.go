package handlers

import (
    "encoding/json"
    "net/http"
    "io"
    "fmt"
    "github.com/go-chi/chi/v5"
    "github.com/wolf4b12/metrics-sv.git/internal/constant"     // Импортируем константы
    "github.com/wolf4b12/metrics-sv.git/internal/server/metricssrv" // Импортируем структуру метрик
)

// GetStorage интерфейс для получения метрик
type GetStorage interface {
    GetGauge(name string) (float64, error)
    GetCounter(name string) (int64, error)
}

// Old ValueHandler обработчик для получения значения метрики (оставляется без изменений)
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

// New JSON-based ValueHandler
func PostJSONValueHandler(storage GetStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Установим заголовок ответа Content-Type
        w.Header().Set("Content-Type", "application/json")

        // Прочитаем тело запроса
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        // Расшифровываем JSON-данные в структуру Metrics
        var inputMetric metricssrv.Metrics
        err = json.Unmarshal(body, &inputMetric)
        if err != nil {
            http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
            return
        }

        // Проверим, были ли переданы оба необходимых параметра
    //    if inputMetric.ID == "" || inputMetric.MType == "" {
     //        http.Error(w, "Параметры id и type обязательны", http.StatusBadRequest)
     //       return
     //   }

        // Запрашиваем нужную метрику исходя из типа
        var outputMetric metricssrv.Metrics
        outputMetric.ID = inputMetric.ID
        outputMetric.MType = inputMetric.MType

        switch inputMetric.MType {
        case constant.MetricTypeGauge:
            value, err := storage.GetGauge(inputMetric.ID)
            if err != nil {
                http.Error(w, "Ошибка получения gauge-метрики", http.StatusNotFound)
                return
            }
            outputMetric.Value = &value

        case constant.MetricTypeCounter:
            value, err := storage.GetCounter(inputMetric.ID)
            if err != nil {
                http.Error(w, "Ошибка получения counter-метрики", http.StatusNotFound)
                return
            }
            outputMetric.Delta = &value

        default:
            http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
            return
        }

        // Конвертируем результат в JSON и отправляем клиенту
        respBytes, err := json.Marshal(outputMetric)
        if err != nil {
            http.Error(w, "Ошибка конвертации в JSON", http.StatusInternalServerError)
            return
        }

        w.Write(respBytes)
    }
}
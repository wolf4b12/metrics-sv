package handlers

import (
//   "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "github.com/go-chi/chi/v5"
    storage "github.com/wolf4b12/metrics-sv.git/internal/storage" // Импортируем пользовательский пакет storage
  // "errors"
)





// UpdateHandler — обработчик для обновления метрик
func UpdateHandler(storage storage.Storage) http.HandlerFunc { // UpdateHandler создает новый HTTP HandlerFunc с использованием Storage
    return func(w http.ResponseWriter, r *http.Request) { // Возвращаемый handlerFunc получает writer и request
        w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Устанавливаем заголовок Content-Type для правильного рендеринга текста

        if r.Method != http.MethodPost { // Проверяем, является ли метод запроса POST
            w.WriteHeader(http.StatusMethodNotAllowed) // Если метод не POST, возвращаем ошибку 405 Method Not Allowed
            return
        }

        pathParts := strings.Split(r.URL.Path, "/")[2:] // Разбиваем URL по слэшам и берем последние три сегмента
        if len(pathParts) != 3 { // Проверяем, что количество сегментов равно трем
            w.WriteHeader(http.StatusNotFound) // Если сегментов меньше трех, возвращаем ошибку 404 Not Found
            return
        }

        metricType, metricName, metricValue := pathParts[0], pathParts[1], pathParts[2] // Извлекаем тип метрики, название и значение

        switch metricType { // Анализируем тип метрики
        case "gauge": // Если тип metriсType равен "gauge"
            value, err := strconv.ParseFloat(metricValue, 64) // Преобразуем строковое значение в число с плавающей точкой
            if err != nil { // Если произошла ошибка при преобразовании
                w.WriteHeader(http.StatusBadRequest) // Возвращаем ошибку 400 Bad Request
                return
            }
            storage.UpdateGauge(metricName, value) // Обновляем gauge метрику

        case "counter": // Если тип metriсType равен "counter"
            value, err := strconv.ParseInt(metricValue, 10, 64) // Преобразуем строковое значение в целое число
            if err != nil { // Если произошла ошибка при преобразовании
                w.WriteHeader(http.StatusBadRequest) // Возвращаем ошибку 400 Bad Request
                return
            }
            storage.UpdateCounter(metricName, value) // Обновляем counter метрику

        default: // Если тип метрики не совпадает ни с одним известным типом
            w.WriteHeader(http.StatusBadRequest) // Возвращаем ошибку 400 Bad Request
            return
        }

        w.WriteHeader(http.StatusOK) // Все прошло успешно, возвращаем 200 OK
    }
}

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
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "Unknown metric type: %s", metricType)
            return
        }

        // Обрабатываем ошибки
        if err != nil {
  {
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, "MetricNotFound")
            }
            return
        }

        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, value)
    }
}

// ListMetricsHandler обработчик для получения списка всех метрик
func ListMetricsHandler(storage storage.Storage) http.HandlerFunc {
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
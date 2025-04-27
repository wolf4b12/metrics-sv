package agentmethods

import (
    "bytes"
//    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    metrics "github.com/wolf4b12/metrics-sv.git/internal/agent/metricsagent"
)


// sendMetric отправляет одну метрику в виде POST-запроса
func (a *Agent) sendMetric(metric metrics.Metrics) {
    data, err := metric.MarshalJSON()
    if err != nil {
        a.handleErrorAndContinue("маршализации метрики в JSON", err)
        return
    }

    // Формируем URL для отправки метрики
    url := fmt.Sprintf("http://%s/update/", a.addr)

    // Сжимаем данные
    compressedData, err := a.compressPayload(data)
    if err != nil {
        a.handleErrorAndContinue("сжатия метрики", err)
        return
    }

    // Формируем запрос с Gzip-данными
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(compressedData))
    if err != nil {
        a.handleErrorAndContinue("формирования запроса", err)
        return
    }

    // Устанавливаем заголовки
    a.SetHeaders(req, "application/json")

    // Выполняем запрос
    resp, err := a.client.Do(req)
    if err != nil {
        a.handleErrorAndContinue("отправки метрики", err)
        return
    }

    // Обрабатываем ответ
    if err := a.handleResponse(resp); err != nil {
        a.handleErrorAndContinue("обработки ответа", err)
    }
}

// SendJSONCollectedMetrics отправляет собранные метрики в формате JSON
func (a *Agent) SendJSONCollectedMetrics() {
    for {
        a.mu.Lock()

        // Отправляем Gauges
        for _, gauge := range a.Gauges {
            if gauge.Value == nil {
                log.Printf("Отсутствует обязательный параметр 'Value' для сенсора '%s'\n", gauge.ID)
                continue
            }
            a.sendMetric(gauge)
        }

        // Отправляем Counters
        for _, counter := range a.Counters {
            if counter.Delta == nil {
                log.Printf("Отсутствует обязательный параметр 'Delta' для счетчика '%s'\n", counter.ID)
                continue
            }
            a.sendMetric(counter)
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}
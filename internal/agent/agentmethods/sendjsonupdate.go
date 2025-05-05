package agentmethods

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    "bytes"
)

// sendSingleMetric отправляет одну метрику в формате JSON
func (a *Agent) sendSingleMetric(metric interface{}, metricID string, checkRequired func() bool) {
    if !checkRequired() {
        log.Printf("Отсутствует обязательный параметр для метрики '%s'\n", metricID)
        return
    }

    // Маршализируем метрику в JSON
    data, err := json.Marshal(metric)
    if err != nil {
        a.HandleErrorAndContinue("маршализации метрики в JSON", err)
        return
    }

    // Формируем URL для отправки метрики
    url := fmt.Sprintf("http://%s/update/", a.addr)

    // Сжимаем данные
    compressedData, err := a.CompressPayload(data)
    if err != nil {
        a.HandleErrorAndContinue("сжатия метрики", err)
        return
    }

    // Формируем запрос с Gzip-данными
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(compressedData))
    if err != nil {
        a.HandleErrorAndContinue("формирования запроса", err)
        return
    }

    // Устанавливаем заголовки
    a.SetHeaders(req, "application/json")

    // Выполняем запрос
    resp, err := a.client.Do(req)
    if err != nil {
        a.HandleErrorAndContinue("Ошибка отправки метрики", err)
        return
    }

    // Обрабатываем ответ
    if err := a.HandleResponse(resp); err != nil {
        a.HandleErrorAndContinue("Ошибка обработки ответа", err)
    }
}

// SendJSONCollectedMetrics отправляет собранные метрики в формате JSON
func (a *Agent) SendJSONCollectedMetrics() {
    for {
        a.mu.Lock()

        for _, gauge := range a.Gauges {
            a.sendSingleMetric(
                gauge,
                gauge.ID,
                func() bool { return gauge.Value != nil },
            )
        }

        for _, counter := range a.Counters {
            a.sendSingleMetric(
                counter,
                counter.ID,
                func() bool { return counter.Delta != nil },
            )
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}

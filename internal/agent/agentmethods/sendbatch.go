package agentmethods

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "bytes"
    metrics "github.com/wolf4b12/metrics-sv/internal/agent/metricsagent"

)

// sendSingleMetric отправляет одну метрику в формате JSON
func (a *Agent) sendBatch(batch []metrics.Metrics) {
    if len(batch) == 0 {
        return // Нет смысла отправлять пустой батч
    }

    // Маршализация данных в JSON
    payload, err := json.Marshal(batch)
    if err != nil {
        log.Printf("Ошибка маршализации метрик в JSON: %v\n", err)
        return
    }

    // Компрессия данных
    compressedData, err := a.CompressPayload(payload)
    if err != nil {
        log.Printf("Ошибка сжатия данных: %v\n", err)
        return
    }

    // Формирование HTTP-запроса
    url := "http://" + a.addr + "/updates/"
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(compressedData))
    if err != nil {
        log.Printf("Ошибка создания запроса: %v\n", err)
        return
    }

    // Заголовки запроса
    a.SetHeaders(req, "application/json")

    // Отправка запроса
    resp, err := a.client.Do(req)
    if err != nil {
        a.HandleErrorAndContinue("Ошибка отправки batch", err)
        return
    }
    defer resp.Body.Close()

    // Обрабатываем ответ
    if err := a.HandleResponse(resp); err != nil {
        a.HandleErrorAndContinue("Ошибка обработки ответа", err)
    }
}

// CollectAndSendBatches собираем и отправляем метрики пакетами
func (a *Agent) CollectAndSendBatches() {
    for {
        a.mu.Lock()
        batch := make([]metrics.Metrics, 0, len(a.Gauges)+len(a.Counters))

        // Собираем Gauges
        for i := range a.Gauges {
            batch = append(batch, a.Gauges[i])
        }

        // Собираем Counters
        for i := range a.Counters {
            batch = append(batch, a.Counters[i])
        }

        a.mu.Unlock()

        // Отправляем пакет
        a.sendBatch(batch)

        // Пауза между отправками
        time.Sleep(a.reportInterval)
    }
}
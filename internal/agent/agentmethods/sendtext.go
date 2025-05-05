package agentmethods

import (
    "context"
    "bytes"
    "fmt"
    "log"
    "net/http"
    "time"
)

// sendSingleTextMetric отправляет одну метрику в текстовом формате
func (a *Agent) SendSingleTextMetric(ctx context.Context, urlPath string, payload string, metricID string, checkRequired func() bool) {
    if !checkRequired() {
        log.Printf("Отсутствует обязательное поле для метрики '%s'\n", metricID)
        return
    }

    // Сжимаем payload
    compressedData, err := a.CompressPayload([]byte(payload))
    if err != nil {
        a.HandleErrorAndContinue("сжатия URL", err)
        return
    }

    // Формируем POST-запрос с Gzip-данными
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewBuffer(compressedData))
    if err != nil {
        a.HandleErrorAndContinue("формирования запроса", err)
        return
    }

    // Устанавливаем заголовки
    a.SetHeaders(req, "text/plain")

    // Выполняем запрос
    resp, err := a.client.Do(req)
    if err != nil {
        a.HandleErrorAndContinue("отправки метрики", err)
        return
    }

    // Обрабатываем ответ
    if err := a.HandleResponse(resp); err != nil {
        a.HandleErrorAndContinue("обработки ответа", err)
    }
}

// SendTextCollectedMetrics отправляет собранные метрики в текстовом формате
func (a *Agent) SendTextCollectedMetrics(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            a.mu.Lock()

            baseURL := fmt.Sprintf("http://%s/update", a.addr)

            for _, gauge := range a.Gauges {
                a.SendSingleTextMetric(
                    ctx,
                    baseURL+"/gauge",
                    fmt.Sprintf("%s/gauge/%s/%f", baseURL, gauge.ID, *(gauge.Value)),
                    gauge.ID,
                    func() bool { return gauge.Value != nil },
                )
            }

            for _, counter := range a.Counters {
                a.SendSingleTextMetric(
                    ctx,
                    baseURL+"/counter",
                    fmt.Sprintf("%s/counter/%s/%d", baseURL, counter.ID, *(counter.Delta)),
                    counter.ID,
                    func() bool { return counter.Delta != nil },
                )
            }

            a.mu.Unlock()
            time.Sleep(a.reportInterval)
        }
    }
}

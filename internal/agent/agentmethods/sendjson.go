package agentmethods

import (
 
    "encoding/json"
    "fmt"
    "log"
    "time"
    "bytes"
    "net/http"
    "compress/gzip"


)


// SendJSONCollectedMetrics отправляет собранные метрики в формате JSON
// SendJSONCollectedMetrics отправляет собранные метрики в формате JSON
func (a *Agent) SendJSONCollectedMetrics() {
    for {
        a.mu.Lock()

        // Отправляем каждую метрику отдельно
        for _, gauge := range a.Gauges {
            if gauge.Value == nil {
                log.Printf("Отсутствует обязательный параметр 'Value' для сенсора '%s'\n", gauge.ID)
                continue
            }

            // Маршализируем единичную метрику в JSON
            data, err := json.Marshal(gauge)
            if err != nil {
                log.Printf("Ошибка маршализации метрики в JSON: %v\n", err)
                continue
            }

            // Формируем URL для отправки метрики
            url := fmt.Sprintf("http://%s/update", a.addr)

            // Сжимаем данные с помощью Gzip
            var buf bytes.Buffer
            zw := gzip.NewWriter(&buf)
            if _, err := zw.Write(data); err != nil {
                log.Printf("Ошибка сжатия метрики: %v\n", err)
                continue
            }
            if err := zw.Close(); err != nil {
                log.Printf("Ошибка закрытия компрессора: %v\n", err)
                continue
            }

            // Отправляем метрику
            if err := a.doRequest(http.MethodPost, url, &buf, map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
            }
        }

        // Повторяем аналогичную процедуру для счетчиков
        for _, counter := range a.Counters {
            if counter.Delta == nil {
                log.Printf("Отсутствует обязательный параметр 'Delta' для счетчика '%s'\n", counter.ID)
                continue
            }

            // Маршализируем единичную метрику в JSON
            data, err := json.Marshal(counter)
            if err != nil {
                log.Printf("Ошибка маршализации метрики в JSON: %v\n", err)
                continue
            }

            // Формируем URL для отправки метрики
            url := fmt.Sprintf("http://%s/update", a.addr)

            // Сжимаем данные с помощью Gzip
            var buf bytes.Buffer
            zw := gzip.NewWriter(&buf)
            if _, err := zw.Write(data); err != nil {
                log.Printf("Ошибка сжатия метрики: %v\n", err)
                continue
            }
            if err := zw.Close(); err != nil {
                log.Printf("Ошибка закрытия компрессора: %v\n", err)
                continue
            }

            // Отправляем метрику
            if err := a.doRequest(http.MethodPost, url, &buf, map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}
package agentmethods

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
//	"runtime/metrics"
	"time"
)




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
                a.handleErrorAndContinue("маршализации метрики в JSON", err)
                continue
            }

            // Формируем URL для отправки метрики
            url := fmt.Sprintf("http://%s/update/", a.addr)

            // Сжимаем данные с помощью Gzip
            var buf bytes.Buffer
            zw := gzip.NewWriter(&buf)
            if _, err := zw.Write(data); err != nil {
                a.handleErrorAndContinue("сжатия метрики", err)
                continue
            }
            if err := zw.Close(); err != nil {
                a.handleErrorAndContinue("закрытия компрессора", err)
                continue
            }

            // Формируем запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, url, &buf)
            if err != nil {
                a.handleErrorAndContinue("формирования запроса", err)
                continue
            }

            // Устанавливаем заголовки для сжатия
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Content-Encoding", "gzip")
            req.Header.Set("Accept-Encoding", "gzip")

            // Выполняем запрос
            resp, err := a.client.Do(req)
            if err != nil {
                a.handleErrorAndContinue("отправки метрики", err)
                continue
            }

            // Обрабатываем ответ
            if err := a.handleResponse(resp); err != nil {
                a.handleErrorAndContinue("обработки ответа", err)
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
                a.handleErrorAndContinue("маршализации метрики в JSON", err)
                continue
            }

            // Формируем URL для отправки метрики
            url := fmt.Sprintf("http://%s/update/", a.addr)

            // Сжимаем данные с помощью Gzip
            var buf bytes.Buffer
            zw := gzip.NewWriter(&buf)
            if _, err := zw.Write(data); err != nil {
                a.handleErrorAndContinue("сжатия метрики", err)
                continue
            }
            if err := zw.Close(); err != nil {
                a.handleErrorAndContinue("закрытия компрессора", err)
                continue
            }

            // Формируем запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, url, &buf)
            if err != nil {
                a.handleErrorAndContinue("формирования запроса", err)
                continue
            }

            // Устанавливаем заголовки для сжатия
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Content-Encoding", "gzip")
            req.Header.Set("Accept-Encoding", "gzip")

            // Выполняем запрос
            resp, err := a.client.Do(req)
            if err != nil {
                a.handleErrorAndContinue("отправки метрики", err)
                continue
            }

            // Обрабатываем ответ
            if err := a.handleResponse(resp); err != nil {
                a.handleErrorAndContinue("обработки ответа", err)
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}

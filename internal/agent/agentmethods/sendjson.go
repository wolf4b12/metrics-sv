package agentmethods

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    "compress/gzip"
    "io"

)


// Метод для отправки собранных метрик
// Метод для отправки собранных метрик в формате JSON с поддержкой Gzip
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
            url := fmt.Sprintf("http://%s/update/", a.addr)

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

            // Формируем запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, url, &buf)
            if err != nil {
                log.Printf("Ошибка формирования запроса: %v\n", err)
                continue
            }

            // Устанавливаем заголовки для сжатия
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Content-Encoding", "gzip")
            req.Header.Set("Accept-Encoding", "gzip")

            // Выполняем запрос
            resp, err := a.client.Do(req)
            if err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
                continue
            }
            defer resp.Body.Close()

            // Проверяем статус ответа
            if resp.StatusCode != http.StatusOK {
                log.Printf("Получен неправильный статус-код (%d)\n", resp.StatusCode)
            }

            // Если ответ приходит в сжатом виде, разархивируем его
            if resp.Header.Get("Content-Encoding") == "gzip" {
                reader, err := gzip.NewReader(resp.Body)
                if err != nil {
                    log.Printf("Ошибка разбора Gzip-ответа: %v\n", err)
                    continue
                }
                defer reader.Close()

                // Читаем ответ
                bodyBytes, err := io.ReadAll(reader)
                if err != nil {
                    log.Printf("Ошибка чтения тела ответа: %v\n", err)
                    continue
                }

                fmt.Println(string(bodyBytes))
            } else {
                // Ответ несжатый, читаем обычный
                bodyBytes, err := io.ReadAll(resp.Body)
                if err != nil {
                    log.Printf("Ошибка чтения тела ответа: %v\n", err)
                    continue
                }

                fmt.Println(string(bodyBytes))
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
            url := fmt.Sprintf("http://%s/update/", a.addr)

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

            // Формируем запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, url, &buf)
            if err != nil {
                log.Printf("Ошибка формирования запроса: %v\n", err)
                continue
            }

            // Устанавливаем заголовки для сжатия
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Content-Encoding", "gzip")
            req.Header.Set("Accept-Encoding", "gzip")

            // Выполняем запрос
            resp, err := a.client.Do(req)
            if err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
                continue
            }
            defer resp.Body.Close()

            // Проверяем статус ответа
            if resp.StatusCode != http.StatusOK {
                log.Printf("Получен неправильный статус-код (%d)\n", resp.StatusCode)
            }

            // Если ответ приходит в сжатом виде, разархивируем его
            if resp.Header.Get("Content-Encoding") == "gzip" {
                reader, err := gzip.NewReader(resp.Body)
                if err != nil {
                    log.Printf("Ошибка разбора Gzip-ответа: %v\n", err)
                    continue
                }
                defer reader.Close()

                // Читаем ответ
                bodyBytes, err := io.ReadAll(reader)
                if err != nil {
                    log.Printf("Ошибка чтения тела ответа: %v\n", err)
                    continue
                }

                fmt.Println(string(bodyBytes))
            } else {
                // Ответ несжатый, читаем обычный
                bodyBytes, err := io.ReadAll(resp.Body)
                if err != nil {
                    log.Printf("Ошибка чтения тела ответа: %v\n", err)
                    continue
                }

                fmt.Println(string(bodyBytes))
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}

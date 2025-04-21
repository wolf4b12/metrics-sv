package agentmethods

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "io"
    "bytes"
    "compress/gzip"

)


func (a *Agent) SendTextCollectedMetrics() {
    for {
        a.mu.Lock()

        // Подготовим URL для отправки метрик
        baseURL := fmt.Sprintf("http://%s/update", a.addr)

        // Отправляем измерители (Gauges)
        for _, gauge := range a.Gauges {
            if gauge.Value == nil {
                log.Printf("Отсутствует обязательное поле 'Value' для датчика '%s'\n", gauge.ID)
                continue
            }

            // Формируем URL для конкретной метрики
            textURL := fmt.Sprintf("%s/gauge/%s/%f", baseURL, gauge.ID, *(gauge.Value))

            // Сжимаем URL в Gzip
            var buffer bytes.Buffer
            writer := gzip.NewWriter(&buffer)
            if _, err := writer.Write([]byte(textURL)); err != nil {
                log.Printf("Ошибка сжатия URL: %v\n", err)
                continue
            }
            if err := writer.Close(); err != nil {
                log.Printf("Ошибка закрытия Gzip-компрессора: %v\n", err)
                continue
            }

            // Формируем POST-запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, baseURL+"/gauge", &buffer)
            if err != nil {
                log.Printf("Ошибка формирования запроса: %v\n", err)
                continue
            }

            // Добавляем заголовки для сжатия
            req.Header.Set("Content-Type", "text/plain")
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

        // Отправляем счетчики (Counters)
        for _, counter := range a.Counters {
            if counter.Delta == nil {
                log.Printf("Отсутствует обязательное поле 'Delta' для счётчика '%s'\n", counter.ID)
                continue
            }

            // Формируем URL для конкретной метрики
            textURL := fmt.Sprintf("%s/counter/%s/%d", baseURL, counter.ID, *(counter.Delta))

            // Сжимаем URL в Gzip
            var buffer bytes.Buffer
            writer := gzip.NewWriter(&buffer)
            if _, err := writer.Write([]byte(textURL)); err != nil {
                log.Printf("Ошибка сжатия URL: %v\n", err)
                continue
            }
            if err := writer.Close(); err != nil {
                log.Printf("Ошибка закрытия Gzip-компрессора: %v\n", err)
                continue
            }

            // Формируем POST-запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, baseURL+"/counter", &buffer)
            if err != nil {
                log.Printf("Ошибка формирования запроса: %v\n", err)
                continue
            }

            // Добавляем заголовки для сжатия
            req.Header.Set("Content-Type", "text/plain")
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
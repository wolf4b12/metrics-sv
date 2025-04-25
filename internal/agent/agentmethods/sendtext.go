package agentmethods

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"time"

)


// Агент для сбора и отправки метрик




// SendTextCollectedMetrics отправляет собранные метрики в текстовом формате
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
                a.handleErrorAndContinue("сжатия URL", err)
                continue
            }
            if err := writer.Close(); err != nil {
                a.handleErrorAndContinue("закрытия Gzip-компрессора", err)
                continue
            }

            // Формируем POST-запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, baseURL+"/gauge", &buffer)
            if err != nil {
                a.handleErrorAndContinue("формирования запроса", err)
                continue
            }

            // Добавляем заголовки для сжатия
            req.Header.Set("Content-Type", "text/plain")
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
                a.handleErrorAndContinue("сжатия URL", err)
                continue
            }
            if err := writer.Close(); err != nil {
                a.handleErrorAndContinue("закрытия Gzip-компрессора", err)
                continue
            }

            // Формируем POST-запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, baseURL+"/counter", &buffer)
            if err != nil {
                a.handleErrorAndContinue("формирования запроса", err)
                continue
            }

            // Добавляем заголовки для сжатия
            req.Header.Set("Content-Type", "text/plain")
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
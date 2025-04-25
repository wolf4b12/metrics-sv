package agentmethods

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

)

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

            // Сжимаем URL
            compressedData, err := a.compressPayload([]byte(textURL))
            if err != nil {
                a.handleErrorAndContinue("сжатия URL", err)
                continue
            }

            // Формируем POST-запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, baseURL+"/gauge", bytes.NewBuffer(compressedData))
            if err != nil {
                a.handleErrorAndContinue("формирования запроса", err)
                continue
            }

            // Устанавливаем заголовки
            a.SetHeaders(req, "text/plain")

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

            // Сжимаем URL
            compressedData, err := a.compressPayload([]byte(textURL))
            if err != nil {
                a.handleErrorAndContinue("сжатия URL", err)
                continue
            }

            // Формируем POST-запрос с Gzip-данными
            req, err := http.NewRequest(http.MethodPost, baseURL+"/counter", bytes.NewBuffer(compressedData))
            if err != nil {
                a.handleErrorAndContinue("формирования запроса", err)
                continue
            }

            // Устанавливаем заголовки
            a.SetHeaders(req, "text/plain")

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
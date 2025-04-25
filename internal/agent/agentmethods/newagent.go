package agentmethods

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
//	"runtime/metrics"
	"sync"
	"time"
    metrics "github.com/wolf4b12/metrics-sv.git/internal/agent/metricsagent"
)


// Агент для сбора и отправки метрик
type Agent struct {
    Gauges         []metrics.Metrics
    Counters       []metrics.Metrics
    pollCount      int64
    mu             *sync.Mutex
    pollInterval   time.Duration
    reportInterval time.Duration
    addr           string
    client         *http.Client
}

// Структура метрики

// Создание нового агента
func NewAgent(poll, report time.Duration, addr string) *Agent {
    return &Agent{
        Gauges:         make([]metrics.Metrics, 0),
        Counters:       make([]metrics.Metrics, 0),
        pollInterval:   poll,
        reportInterval: report,
        addr:           addr,
        mu:             &sync.Mutex{},
        client:         &http.Client{Timeout: 5 * time.Second},
    }
}

// Helper function для отправки сжатых данных
func (a *Agent) sendCompressedRequest(url string, payload []byte, contentType string) error {
    // Сжимаем данные с помощью Gzip
    var buf bytes.Buffer
    zw := gzip.NewWriter(&buf)
    if _, err := zw.Write(payload); err != nil {
        return fmt.Errorf("ошибка сжатия метрики: %v", err)
    }
    if err := zw.Close(); err != nil {
        return fmt.Errorf("ошибка закрытия компрессора: %v", err)
    }

    // Формируем запрос с Gzip-данными
    req, err := http.NewRequest(http.MethodPost, url, &buf)
    if err != nil {
        return fmt.Errorf("ошибка формирования запроса: %v", err)
    }

    // Устанавливаем заголовки для сжатия
    req.Header.Set("Content-Type", contentType)
    req.Header.Set("Content-Encoding", "gzip")
    req.Header.Set("Accept-Encoding", "gzip")

    // Выполняем запрос
    resp, err := a.client.Do(req)
    if err != nil {
        return fmt.Errorf("ошибка отправки метрики: %v", err)
    }
    defer resp.Body.Close()

    // Проверяем статус ответа
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("получен неправильный статус-код (%d)", resp.StatusCode)
    }

    // Если ответ приходит в сжатом виде, разархивируем его
    if resp.Header.Get("Content-Encoding") == "gzip" {
        reader, err := gzip.NewReader(resp.Body)
        if err != nil {
            return fmt.Errorf("ошибка разбора Gzip-ответа: %v", err)
        }
        defer reader.Close()

        // Читаем ответ
        bodyBytes, err := io.ReadAll(reader)
        if err != nil {
            return fmt.Errorf("ошибка чтения тела ответа: %v", err)
        }

        fmt.Println(string(bodyBytes))
    } else {
        // Ответ несжатый, читаем обычный
        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            return fmt.Errorf("ошибка чтения тела ответа: %v", err)
        }

        fmt.Println(string(bodyBytes))
    }

    return nil
}

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

            // Отправляем метрику
            if err := a.sendCompressedRequest(url, data, "application/json"); err != nil {
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

            // Отправляем метрику
            if err := a.sendCompressedRequest(url, data, "application/json"); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}

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

            // Отправляем метрику
            if err := a.sendCompressedRequest(baseURL+"/gauge", []byte(textURL), "text/plain"); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
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

            // Отправляем метрику
            if err := a.sendCompressedRequest(baseURL+"/counter", []byte(textURL), "text/plain"); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}
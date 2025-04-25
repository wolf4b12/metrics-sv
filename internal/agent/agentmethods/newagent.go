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

// helperDoRequest отправляет запрос и обрабатывает ответ
func (a *Agent) helperDoRequest(method, url string, body io.Reader, headers map[string]string) error {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return fmt.Errorf("ошибка формирования запроса: %v", err)
    }

    // Устанавливаем заголовки
    for k, v := range headers {
        req.Header.Set(k, v)
    }

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
            if err := a.helperDoRequest(http.MethodPost, url, &buf, map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}); err != nil {
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
            if err := a.helperDoRequest(http.MethodPost, url, &buf, map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}); err != nil {
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

            // Отправляем метрику
            if err := a.helperDoRequest(http.MethodPost, baseURL+"/gauge", &buffer, map[string]string{"Content-Type": "text/plain", "Content-Encoding": "gzip"}); err != nil {
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

            // Отправляем метрику
            if err := a.helperDoRequest(http.MethodPost, baseURL+"/counter", &buffer, map[string]string{"Content-Type": "text/plain", "Content-Encoding": "gzip"}); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}
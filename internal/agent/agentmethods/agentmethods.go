package agentmethods

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "runtime"
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

// Метод для сбора метрик
func (a *Agent) CollectMetrics() {
    for {
        a.mu.Lock()

        // Чистка старых коллекций перед сборкой новых данных
        a.Gauges = a.Gauges[:0]
        a.Counters = a.Counters[:0]

        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)

        // Собираем runtime-метрики и добавляем их в Gauges
        runtimeMetrics := metrics.GetRuntimeMetricsGauge(memStats)
        for key, value := range runtimeMetrics {
            a.Gauges = append(a.Gauges, metrics.Metrics{ID: key, MType: "gauge", Value: &value})
        }

        // Счётчик опроса PollCount
        a.pollCount++
        a.Counters = append(a.Counters, metrics.Metrics{ID: "PollCount", MType: "counter", Delta: &a.pollCount})

        a.mu.Unlock()
        time.Sleep(a.pollInterval)
    }
}

// Метод для отправки собранных метрик
func (a *Agent) SendCollectedMetrics() {
    for {
        a.mu.Lock()

        // Проходим по каждой метрике отдельно
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
            req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            if err != nil {
                log.Printf("Ошибка формирования запроса: %v\n", err)
                continue
            }
            req.Header.Set("Content-Type", "application/json")

            // Выполняем запрос
            resp, err := a.client.Do(req)
            if err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
                continue
            }
            
            resp.Body.Close()
            // Проверяем статус ответа
            if resp.StatusCode != http.StatusOK {
                log.Printf("Получен неправильный статус-код (%d)\n", resp.StatusCode)
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
            req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            if err != nil {
                log.Printf("Ошибка формирования запроса: %v\n", err)
                continue
            }
            req.Header.Set("Content-Type", "application/json")

            // Выполняем запрос
            resp, err := a.client.Do(req)
            if err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
                continue
            }
            
            resp.Body.Close()

            // Проверяем статус ответа
            if resp.StatusCode != http.StatusOK {
                log.Printf("Получен неправильный статус-код (%d)\n", resp.StatusCode)
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}
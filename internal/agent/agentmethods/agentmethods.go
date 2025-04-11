package agentmethods

import (
    "fmt"
    "log"
    "net/http"
    "runtime"
    "sync"
    "time"
    metrics "github.com/wolf4b12/metrics-sv.git/internal/agent/metrics"
)

type Agent struct {
    gauges         map[string]float64
    counters       map[string]int64
    pollCount      int64
    mu             *sync.Mutex
    pollInterval   time.Duration
    reportInterval time.Duration
    addr           string
}

func NewAgent(poll, report time.Duration, addr string) *Agent {
    return &Agent{
        gauges:         make(map[string]float64),
        counters:       make(map[string]int64),
        pollInterval:   poll,
        reportInterval: report,
        addr:           addr,
        mu:             &sync.Mutex{},
    }
}

func (a *Agent) CollectMetrics() { // собираем метрики
    var memStats runtime.MemStats

    for {
        runtime.ReadMemStats(&memStats)

        a.mu.Lock()

        // Используем карту для инициализации всех runtime-метрик
        runtimeMetrics := metrics.GetRuntimeMetricsGauge(memStats)
        for key, value := range runtimeMetrics {
            a.gauges[key] = value
        }

        // Используем карту для инициализации кастомных метрик
        customMetrics := metrics.GetCustomMetrics()
        for key, value := range customMetrics {
            a.counters[key] = value
        }

        // Обновляем счетчики
        a.pollCount++
        a.counters["PollCount"] = a.pollCount

        a.mu.Unlock()
        time.Sleep(a.pollInterval)
    }
}

func (a *Agent) SendCollectedMetrics() { // отправляем собранные метрики
    client := &http.Client{Timeout: 5 * time.Second}
    baseURL := fmt.Sprintf("http://%s/update", a.addr)

    for {
        a.mu.Lock()

        // Send gauge metrics
        for name, value := range a.gauges {
            url := fmt.Sprintf("%s/gauge/%s/%f", baseURL, name, value)
            go SendMetricToServer(client, url)
        }

        // Send counter metrics
        for name, value := range a.counters {
            url := fmt.Sprintf("%s/counter/%s/%d", baseURL, name, value)
            go SendMetricToServer(client, url)
        }

        a.mu.Unlock()
        time.Sleep(a.reportInterval)
    }
}

func SendMetricToServer(client *http.Client, url string) { // вспомогательная функция для отправки метрик
    resp, err := client.Post(url, "text/plain", nil)
    if err != nil {
        log.Printf("Error sending metric: %v\n", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Unexpected status code: %d\n", resp.StatusCode)
    }
}
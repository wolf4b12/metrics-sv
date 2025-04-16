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

// Определим интерфейс коллектора метрик
type MetricsCollector interface {
    CollectMetrics()
}

// Определим интерфейс отправителя метрик
type MetricsSender interface {
    SendCollectedMetrics()
}

// Объединяем оба интерфейса в общий интерфейс агента
type AgentInterface interface {
    MetricsCollector
    MetricsSender
}

// Структура нашего агента
type Agent struct {
    gauges         map[string]float64
    counters       map[string]int64
    pollCount      int64
    mu             *sync.Mutex
    pollInterval   time.Duration
    reportInterval time.Duration
    addr           string
}

// Конструктор агента
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

// Реализация метода CollectionMetrics
func (a *Agent) CollectMetrics() {
    var memStats runtime.MemStats

    for {
        runtime.ReadMemStats(&memStats)

        a.mu.Lock()

        // Получаем runtime-метрики
        runtimeMetrics := metrics.GetRuntimeMetricsGauge(memStats)
        for key, value := range runtimeMetrics {
            a.gauges[key] = value
        }

        // Получаем кастомные метрики
        customMetrics := metrics.GetCustomMetrics()
        for key, value := range customMetrics {
            a.counters[key] = value
        }

        // Обновляем счётчик поллов
        a.pollCount++
        a.counters["PollCount"] = a.pollCount

        a.mu.Unlock()
        time.Sleep(a.pollInterval)
    }
}

// Реализация метода SendCollectedMetrics
func (a *Agent) SendCollectedMetrics() {
    client := &http.Client{Timeout: 5 * time.Second}
    baseURL := fmt.Sprintf("http://%s/update", a.addr)

    for {
        a.mu.Lock()

        // Отправляем данные по гейтикам
        for name, value := range a.gauges {
            url := fmt.Sprintf("%s/gauge/%s/%f", baseURL, name, value)
            go SendMetricToServer(client, url)
        }

        // Отправляем данные по контрметрикам
        for name, value := range a.counters {
            url := fmt.Sprintf("%s/counter/%s/%d", baseURL, name, value)
            go SendMetricToServer(client, url)
        }

        a.mu.Unlock()
        time.Sleep(a.reportInterval)
    }
}

// Вспомогательная функция отправки метрик на сервер
func SendMetricToServer(client *http.Client, url string) {
    resp, err := client.Post(url, "text/plain", nil)
    if err != nil {
        log.Printf("Ошибка отправки метрики: %v\n", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Непредвиденный статус HTTP-кода: %d\n", resp.StatusCode)
    }
}

// Проверка соответствия типа структуры типу интерфейса
var _ AgentInterface = (*Agent)(nil)
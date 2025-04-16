// github.com/wolf4b12/metrics-sv.git/internal/agent/agentmethods/agent.go

package agentmethods

import (
    "fmt"
    "log"
    "net/http"
//    "runtime"
    "sync"
    "time"
    metrics "github.com/wolf4b12/metrics-sv.git/internal/agent/metrics"
)

// Основные интерфейсы и структура сохраняются прежними

type MetricsSender interface {
    SendCollectedMetrics()
}

type AgentInterface interface {
    MetricsSender
}

type Agent struct {
    gauges         map[string]float64
    counters       map[string]int64
    pollCount      int64
    mu             *sync.Mutex
    pollInterval   time.Duration
    reportInterval time.Duration
    addr           string
}

// Новый конструктор агента остаётся прежним

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

// Теперь метод делегирует выполнение сборщику метрик из пакета metrics
func (a *Agent) CollectMetrics() {
    metrics.CollectMetrics(a.mu, a.gauges, a.counters, a.pollInterval)
}

// Методы отправки метрик остаются прежними
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

// Проверка правильности интерфейсов
var _ AgentInterface = (*Agent)(nil)
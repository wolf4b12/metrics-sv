package agentmethods

import (
	"net/http"
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




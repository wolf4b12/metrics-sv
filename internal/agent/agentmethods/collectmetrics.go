package agentmethods

import (
    "runtime"
    "time"

    metrics "github.com/wolf4b12/metrics-sv/internal/agent/metricsagent"
)

// Агент собирает метрики с определенной частотой

// Метод для начала непрерывного сбора метрик
func (a *Agent) StartCollectingMetrics() {
    ticker := time.NewTicker(a.pollInterval)
    defer ticker.Stop()

    for range ticker.C {
        a.collectMetricsOnce()
    }
}

// Внутренний метод для однократного сбора метрик
func (a *Agent) collectMetricsOnce() {
    a.mu.Lock()
    defer a.mu.Unlock()

    // Чистим старые коллекции перед новым опросом 
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
}
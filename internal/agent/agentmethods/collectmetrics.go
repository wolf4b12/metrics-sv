package agentmethods

import (
    "runtime"
    "time"

    metrics "github.com/wolf4b12/metrics-sv.git/internal/agent/metricsagent"
)

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
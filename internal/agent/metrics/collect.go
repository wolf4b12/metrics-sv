package metrics

import (
    "runtime"
    "sync"
    "time"
)

// Интерфейс для методов сбора метрик
type MetricsCollector interface {
    CollectMetrics(*sync.Mutex, map[string]float64, map[string]int64)
}

// Метод сбора метрик (используется агентом)
func CollectMetrics(mu *sync.Mutex, gauges map[string]float64, counters map[string]int64, pollInterval time.Duration) {
    var memStats runtime.MemStats

    for {
        runtime.ReadMemStats(&memStats)

        mu.Lock()

        // Читаем runtime-метрики
        runtimeMetrics := GetRuntimeMetricsGauge(memStats)
        for key, value := range runtimeMetrics {
            gauges[key] = value
        }

        // Читаем кастомные метрики
        customMetrics := GetCustomMetrics()
        for key, value := range customMetrics {
            counters[key] = value
        }

        // Обновляем счётчик колла
        pollCount := counters["PollCount"]
        pollCount++
        counters["PollCount"] = pollCount

        mu.Unlock()
        time.Sleep(pollInterval)
    }
}


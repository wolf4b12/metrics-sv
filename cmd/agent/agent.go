package main

import (
    "flag"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "os"
    "runtime"
    "sync/atomic"
    "time"
)

// Metrics содержит метрики, которые собирает агент
type Metrics struct {
    PollCount   int64         // Counter, увеличивается при каждом обновлении метрик
    RandomValue float64       // Gauge, случайное значение
    MemStats    runtime.MemStats // Метрики из runtime
}

// Config содержит конфигурацию агента
type Config struct {
    PollInterval   time.Duration // Интервал обновления метрик
    ReportInterval time.Duration // Интервал отправки метрик на сервер
    ServerAddress  string        // Адрес сервера для отправки метрик
}

// CollectMetrics собирает метрики из runtime и обновляет кастомные метрики
func CollectMetrics(metrics *Metrics) {
    atomic.AddInt64(&metrics.PollCount, 1)
    metrics.RandomValue = rand.Float64()
    runtime.ReadMemStats(&metrics.MemStats)
}

// SendMetric отправляет одну метрику на сервер по HTTP
func SendMetric(client *http.Client, serverAddress, metricType, metricName string, value interface{}) error {
    url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, metricType, metricName, value)
    req, err := http.NewRequest(http.MethodPost, url, nil)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "text/plain")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
    return nil
}

// SendAllMetrics отправляет все собранные метрики на сервер
func SendAllMetrics(client *http.Client, serverAddress string, metrics *Metrics) {
    // Отправка кастомных метрик
    SendMetric(client, serverAddress, "counter", "PollCount", atomic.LoadInt64(&metrics.PollCount))
    SendMetric(client, serverAddress, "gauge", "RandomValue", metrics.RandomValue)

    // Отправка метрик из runtime.MemStats
    memStats := metrics.MemStats
    SendMetric(client, serverAddress, "gauge", "Alloc", memStats.Alloc)
    SendMetric(client, serverAddress, "gauge", "BuckHashSys", memStats.BuckHashSys)
    SendMetric(client, serverAddress, "gauge", "Frees", memStats.Frees)
    SendMetric(client, serverAddress, "gauge", "GCCPUFraction", memStats.GCCPUFraction)
    SendMetric(client, serverAddress, "gauge", "GCSys", memStats.GCSys)
    SendMetric(client, serverAddress, "gauge", "HeapAlloc", memStats.HeapAlloc)
    SendMetric(client, serverAddress, "gauge", "HeapIdle", memStats.HeapIdle)
    SendMetric(client, serverAddress, "gauge", "HeapInuse", memStats.HeapInuse)
    SendMetric(client, serverAddress, "gauge", "HeapObjects", memStats.HeapObjects)
    SendMetric(client, serverAddress, "gauge", "HeapReleased", memStats.HeapReleased)
    SendMetric(client, serverAddress, "gauge", "HeapSys", memStats.HeapSys)
    SendMetric(client, serverAddress, "gauge", "LastGC", memStats.LastGC)
    SendMetric(client, serverAddress, "gauge", "Lookups", memStats.Lookups)
    SendMetric(client, serverAddress, "gauge", "MCacheInuse", memStats.MCacheInuse)
    SendMetric(client, serverAddress, "gauge", "MCacheSys", memStats.MCacheSys)
    SendMetric(client, serverAddress, "gauge", "MSpanInuse", memStats.MSpanInuse)
    SendMetric(client, serverAddress, "gauge", "MSpanSys", memStats.MSpanSys)
    SendMetric(client, serverAddress, "gauge", "Mallocs", memStats.Mallocs)
    SendMetric(client, serverAddress, "gauge", "NextGC", memStats.NextGC)
    SendMetric(client, serverAddress, "gauge", "NumForcedGC", memStats.NumForcedGC)
    SendMetric(client, serverAddress, "gauge", "NumGC", memStats.NumGC)
    SendMetric(client, serverAddress, "gauge", "OtherSys", memStats.OtherSys)
    SendMetric(client, serverAddress, "gauge", "PauseTotalNs", memStats.PauseTotalNs)
    SendMetric(client, serverAddress, "gauge", "StackInuse", memStats.StackInuse)
    SendMetric(client, serverAddress, "gauge", "StackSys", memStats.StackSys)
    SendMetric(client, serverAddress, "gauge", "Sys", memStats.Sys)
    SendMetric(client, serverAddress, "gauge", "TotalAlloc", memStats.TotalAlloc)

    log.Println("Metrics sent successfully")
}

// RunAgent запускает агент для сбора и отправки метрик с заданными интервалами
func RunAgent(cfg Config) {
    client := &http.Client{Timeout: 5 * time.Second}
    var metrics Metrics

    tickerPoll := time.NewTicker(cfg.PollInterval)   // Таймер для обновления метрик
    tickerReport := time.NewTicker(cfg.ReportInterval) // Таймер для отправки метрик

    for {
        select {
        case <-tickerPoll.C:
            log.Println("Collecting metrics...")
            CollectMetrics(&metrics)

        case <-tickerReport.C:
            log.Println("Sending metrics...")
            go SendAllMetrics(client, cfg.ServerAddress, &metrics) // Отправка в отдельной горутине

            // Сброс счетчика PollCount после отправки
            atomic.StoreInt64(&metrics.PollCount, 0)
        }
    }
}

func main() {
    // Определение флагов
    serverAddr := flag.String("a", "localhost:8080", "Адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).")
    reportInterval := flag.Int("r", 5, "Частота отправки метрик на сервер в секундах (по умолчанию 10 секунд).")
    pollInterval := flag.Int("p", 2, "Частота опроса метрик из пакета runtime в секундах (по умолчанию 2 секунды).")

    // Парсинг флагов
    flag.Parse()

    // Проверка наличия неизвестных флагов
    if flag.NArg() > 0 {
        log.Fatalf("Неизвестный флаг: %s\n", os.Args[flag.NArg()-1])
    }

    // Преобразование значений интервалов в Duration
    reportDuration := time.Duration(*reportInterval) * time.Second
    pollDuration := time.Duration(*pollInterval) * time.Second

    // Конфигурация агента
    cfg := Config{
        PollInterval:   pollDuration,
        ReportInterval: reportDuration,
        ServerAddress:  *serverAddr,
    }

    log.Println("Starting agent...")
    RunAgent(cfg) // Запуск агента
}
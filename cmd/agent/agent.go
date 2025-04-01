package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
//	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

// Metrics содержит метрики, которые собирает агент
type Metrics struct {
	PollCount   int64            // Counter, увеличивается при каждом обновлении метрик
	RandomValue float64          // Gauge, случайное значение
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
	SendMetric(client, serverAddress, "counter", "PollCount", atomic.LoadInt64(&metrics.PollCount))
	SendMetric(client, serverAddress, "gauge", "RandomValue", metrics.RandomValue)

	memStats := metrics.MemStats
	SendMetric(client, serverAddress, "gauge", "Alloc", memStats.Alloc)
	SendMetric(client, serverAddress, "gauge", "HeapAlloc", memStats.HeapAlloc)
	SendMetric(client, serverAddress, "gauge", "TotalAlloc", memStats.TotalAlloc)

	log.Println("Metrics sent successfully")
}

// RunAgent запускает агент для сбора и отправки метрик с заданными интервалами
func RunAgent(cfg Config) {
	client := &http.Client{Timeout: 5 * time.Second}
	var metrics Metrics

	tickerPoll := time.NewTicker(cfg.PollInterval)
	tickerReport := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-tickerPoll.C:
			log.Println("Collecting metrics...")
			CollectMetrics(&metrics)

		case <-tickerReport.C:
			log.Println("Sending metrics...")
			go SendAllMetrics(client, cfg.ServerAddress, &metrics)

			atomic.StoreInt64(&metrics.PollCount, 0)
		}
	}
}

func parseFlags() Config {
	var (
		serverAddress  = flag.String("a", "http://localhost:8080", "HTTP server address")
		reportInterval = flag.String("r", "10", "Report interval in seconds")
		pollInterval   = flag.String("p", "2", "Poll interval in seconds")
	)

	flag.Parse()

	reportIntervalSec, err := strconv.Atoi(*reportInterval)
	if err != nil || reportIntervalSec <= 0 {
		log.Fatalf("Invalid report interval: %s. Must be a positive integer.", *reportInterval)
	}

	pollIntervalSec, err := strconv.Atoi(*pollInterval)
	if err != nil || pollIntervalSec <= 0 {
		log.Fatalf("Invalid poll interval: %s. Must be a positive integer.", *pollInterval)
	}

	if len(flag.Args()) > 0 {
		log.Fatalf("Unknown flags provided: %v", flag.Args())
	}

	return Config{
		ServerAddress:  *serverAddress,
		PollInterval:   time.Duration(pollIntervalSec) * time.Second,
		ReportInterval: time.Duration(reportIntervalSec) * time.Second,
	}
}

func main() {
	cfg := parseFlags()

	log.Printf("Starting agent with config: %+v\n", cfg)
	RunAgent(cfg)
}

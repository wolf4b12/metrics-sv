package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

type Metrics struct {
	PollCount   int64
	RandomValue float64
	MemStats    runtime.MemStats
}

type Config struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerAddress  string
}

func CollectMetrics(metrics *Metrics) {
	atomic.AddInt64(&metrics.PollCount, 1)
	metrics.RandomValue = rand.Float64()
	runtime.ReadMemStats(&metrics.MemStats)
}

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

func SendAllMetrics(client *http.Client, serverAddress string, metrics *Metrics) {
	// Кастомные метрики
	SendMetric(client, serverAddress, "counter", "PollCount", atomic.LoadInt64(&metrics.PollCount))
	SendMetric(client, serverAddress, "gauge", "RandomValue", metrics.RandomValue)

	// Метрики runtime
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

func parseDuration(durationStr string) (time.Duration, error) {
	if !strings.HasSuffix(durationStr, "s") && 
	   !strings.HasSuffix(durationStr, "m") && 
	   !strings.HasSuffix(durationStr, "h") {
		durationStr += "s"
	}
	return time.ParseDuration(durationStr)
}

func main() {
	addr := flag.String("a", "http://localhost:8080", "Server address")
	pollStr := flag.String("p", "2", "Poll interval (seconds)")
	reportStr := flag.String("r", "10", "Report interval (seconds)")

	flag.Parse()

	pollInterval, err := parseDuration(*pollStr)
	if err != nil {
		log.Fatalf("Invalid poll interval: %v", err)
	}

	reportInterval, err := parseDuration(*reportStr)
	if err != nil {
		log.Fatalf("Invalid report interval: %v", err)
	}

	cfg := Config{
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		ServerAddress:  *addr,
	}

	log.Println("Starting agent...")
	RunAgent(cfg)
}

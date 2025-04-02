package main

import (
    "flag"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "runtime"
    "sync"
    "time"
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

func parseFlags() (time.Duration, time.Duration, string) {
    var (
        addr           string
        reportInterval int
        pollInterval   int
    )

    flag.StringVar(&addr, "a", "localhost:8080", "HTTP server address")
    flag.IntVar(&pollInterval, "p", 2, "Poll interval in seconds")
    flag.IntVar(&reportInterval, "r", 10, "Report interval in seconds")
    
    flag.Parse()
    
    if flag.NArg() > 0 {
        log.Fatalf("Unknown flags or arguments: %v", flag.Args())
    }

    return time.Duration(pollInterval) * time.Second,
        time.Duration(reportInterval) * time.Second,
        addr
}

func (a *Agent) CollectMetrics() {
    var memStats runtime.MemStats
    
    for {
        runtime.ReadMemStats(&memStats)
        
        a.mu.Lock()
        
        // Runtime gauge metrics
        a.gauges["Alloc"] = float64(memStats.Alloc)
        a.gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
        a.gauges["Frees"] = float64(memStats.Frees)
        a.gauges["GCCPUFraction"] = memStats.GCCPUFraction
        a.gauges["GCSys"] = float64(memStats.GCSys)
        a.gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
        a.gauges["HeapIdle"] = float64(memStats.HeapIdle)
        a.gauges["HeapInuse"] = float64(memStats.HeapInuse)
        a.gauges["HeapObjects"] = float64(memStats.HeapObjects)
        a.gauges["HeapReleased"] = float64(memStats.HeapReleased)
        a.gauges["HeapSys"] = float64(memStats.HeapSys)
        a.gauges["LastGC"] = float64(memStats.LastGC)
        a.gauges["Lookups"] = float64(memStats.Lookups)
        a.gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
        a.gauges["MCacheSys"] = float64(memStats.MCacheSys)
        a.gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
        a.gauges["MSpanSys"] = float64(memStats.MSpanSys)
        a.gauges["Mallocs"] = float64(memStats.Mallocs)
        a.gauges["NextGC"] = float64(memStats.NextGC)
        a.gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
        a.gauges["NumGC"] = float64(memStats.NumGC)
        a.gauges["OtherSys"] = float64(memStats.OtherSys)
        a.gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
        a.gauges["StackInuse"] = float64(memStats.StackInuse)
        a.gauges["StackSys"] = float64(memStats.StackSys)
        a.gauges["Sys"] = float64(memStats.Sys)
        a.gauges["TotalAlloc"] = float64(memStats.TotalAlloc)
        
        // Custom metrics
        a.gauges["RandomValue"] = rand.Float64()
        a.pollCount++
        a.counters["PollCount"] = a.pollCount
        
        a.mu.Unlock()
        time.Sleep(a.pollInterval)
    }
}

func (a *Agent) SendMetrics() {
    client := &http.Client{Timeout: 5 * time.Second}
    baseURL := fmt.Sprintf("http://%s/update", a.addr)
    
    for {
        a.mu.Lock()
        
        // Send gauge metrics
        for name, value := range a.gauges {
            url := fmt.Sprintf("%s/gauge/%s/%f", baseURL, name, value)
            go sendMetric(client, url)
        }
        
        // Send counter metrics
        for name, value := range a.counters {
            url := fmt.Sprintf("%s/counter/%s/%d", baseURL, name, value)
            go sendMetric(client, url)
        }
        
        a.mu.Unlock()
        time.Sleep(a.reportInterval)
    }
}

func sendMetric(client *http.Client, url string) {
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

func main() {
    rand.New(rand.NewSource(time.Now().UnixNano())) // Create new source for random numbers

    poll, report, addr := parseFlags()
    agent := NewAgent(poll, report, addr)

    go agent.CollectMetrics()
    go agent.SendMetrics()

    select {} // Keep main goroutine alive
}
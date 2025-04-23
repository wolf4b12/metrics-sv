package metricsag

import (

"runtime"

 "math/rand"

)



// GetRuntimeMetricsMap возвращает карту с метриками runtime
func GetRuntimeMetricsGauge(memStats runtime.MemStats) map[string]float64 {
    return map[string]float64{
        "Alloc":                  float64(memStats.Alloc),
        "BuckHashSys":            float64(memStats.BuckHashSys),
        "Frees":                  float64(memStats.Frees),
        "GCCPUFraction":          float64(memStats.GCCPUFraction),
        "GCSys":                  float64(memStats.GCSys),
        "HeapAlloc":              float64(memStats.HeapAlloc),
        "HeapIdle":               float64(memStats.HeapIdle),
        "HeapInuse":              float64(memStats.HeapInuse),
        "HeapObjects":            float64(memStats.HeapObjects),
        "HeapReleased":           float64(memStats.HeapReleased),
        "HeapSys":                float64(memStats.HeapSys),
        "LastGC":                 float64(memStats.LastGC),
        "Lookups":                float64(memStats.Lookups),
        "MCacheInuse":            float64(memStats.MCacheInuse),
        "MCacheSys":              float64(memStats.MCacheSys),
        "MSpanInuse":             float64(memStats.MSpanInuse),
        "MSpanSys":               float64(memStats.MSpanSys),
        "Mallocs":                float64(memStats.Mallocs),
        "NextGC":                 float64(memStats.NextGC),
        "NumForcedGC":            float64(memStats.NumForcedGC),
        "NumGC":                  float64(memStats.NumGC),
        "OtherSys":               float64(memStats.OtherSys),
        "PauseTotalNs":           float64(memStats.PauseTotalNs),
        "StackInuse":             float64(memStats.StackInuse),
        "StackSys":               float64(memStats.StackSys),
        "Sys":                    float64(memStats.Sys),
        "TotalAlloc":             float64(memStats.TotalAlloc),
        "RandomValue":            float64(rand.Intn(1000000)),
    }
}


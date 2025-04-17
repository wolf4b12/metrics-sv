package storage

import (
    "fmt"
    "sync"
    "github.com/wolf4b12/metrics-sv.git/internal/constant" // Импортируем константы
)

// MemStorage реализация хранилища в памяти
type MemStorage struct {
    mu       sync.RWMutex
    gauges   map[string]float64
    counters map[string]int64
}

// NewMemStorage конструктор хранилища
func NewMemStorage() *MemStorage {
    return &MemStorage{
        gauges:   make(map[string]float64),
        counters: make(map[string]int64),
    }
}

// UpdateGauge обновление gauge-метрики

func (s *MemStorage) UpdateGauge(name string, value float64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.gauges[name] = value
}

// UpdateCounter обновление counter-метрики
func (s *MemStorage) UpdateCounter(name string, value int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.counters[name] += value
}

// GetGauge получение gauge-метрики
func (s *MemStorage) GetGauge(name string) (float64, error) {
    s.mu.RLock()
 defer s.mu.RUnlock()
    value, ok := s.gauges[name]
    if !ok {
        return 0, fmt.Errorf("metric not found")
    }
    return value, nil
}

// GetCounter получение counter-метрики
func (s *MemStorage) GetCounter(name string) (int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, ok := s.counters[name]
    if !ok {
        return 0, fmt.Errorf("metric not found")
    }
    return value, nil
}

// AllMetrics получение всех метрик
func (s *MemStorage) AllMetrics() map[string]map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := map[string]map[string]interface{}{
        constant.MetricTypeGauge:   make(map[string]interface{}),
        constant.MetricTypeCounter: make(map[string]interface{}),
    }
    for k, v := range s.gauges {
        result[constant.MetricTypeGauge][k] = v
    }
    for k, v := range s.counters {
        result[constant.MetricTypeCounter][k] = v
    }
    return result
}
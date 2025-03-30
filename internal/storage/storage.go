package storage

import (
    "fmt"
    "sync"
)

// Storage интерфейс для работы с хранилищем
type Storage interface {
    UpdateGauge(name string, value float64)
    UpdateCounter(name string, value int64)
    GetGauge(name string) (float64, error)
    GetCounter(name string) (int64, error)
    AllMetrics() map[string]map[string]interface{}
    ErrMetricNotFound() error // Объявление метода ErrMetricNotFound в интерфейсе
}

// MemStorage реализация хранилища в памяти
type MemStorage struct {
    mu       sync.RWMutex
    gauges   map[string]float64
    counters map[string]int64
}

// ErrMetricNotFound реализация метода ErrMetricNotFound для MemStorage
func (s *MemStorage) ErrMetricNotFound() error {
    return fmt.Errorf("metric not found")
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
        return 0, s.ErrMetricNotFound() // Использование метода ErrMetricNotFound
    }
    return value, nil
}

// GetCounter получение counter-метрики
func (s *MemStorage) GetCounter(name string) (int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, ok := s.counters[name]
    if !ok {
        return 0, s.ErrMetricNotFound() // Использование метода ErrMetricNotFound
    }
    return value, nil
}

// AllMetrics получение всех метрик
func (s *MemStorage) AllMetrics() map[string]map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := map[string]map[string]interface{}{
        "gauges":   make(map[string]interface{}),
        "counters": make(map[string]interface{}),
    }
    for k, v := range s.gauges {
        result["gauges"][k] = v
    }
    for k, v := range s.counters {
        result["counters"][k] = v
    }
    return result
}
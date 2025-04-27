// mem_storage.go
package storage

import (
    "sync"
	"errors"
)

// MemStorage хранит метрики в памяти
type MemStorage struct {
    mu       sync.RWMutex
    gauges   map[string]float64
    counters map[string]int64
}

// NewMemStorage создаёт новое хранилище метрик в памяти
func NewMemStorage() *MemStorage {
    return &MemStorage{
        gauges:   make(map[string]float64),
        counters: make(map[string]int64),
    }
}

// UpdateGauge обновляет значение gauge-метрики
func (s *MemStorage) UpdateGauge(name string, value float64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.gauges[name] = value
}

// UpdateCounter увеличивает значение counter-метрики
func (s *MemStorage) UpdateCounter(name string, value int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.counters[name] += value
}

// GetGauge получает текущее значение gauge-метрики
func (s *MemStorage) GetGauge(name string) (float64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, ok := s.gauges[name]
    if !ok {
        return 0, errors.New("metric not found")
    }
    return value, nil
}

// GetCounter получает текущее значение counter-метрики
func (s *MemStorage) GetCounter(name string) (int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, ok := s.counters[name]
    if !ok {
        return 0, errors.New("metric not found")
    }
    return value, nil
}

// AllMetrics возвращает список всех существующих метрик
func (s *MemStorage) AllMetrics() map[string]map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := map[string]map[string]interface{}{}
    result["gauges"] = make(map[string]interface{})
    result["counters"] = make(map[string]interface{})

    for k, v := range s.gauges {
        result["gauges"][k] = v
    }
    for k, v := range s.counters {
        result["counters"][k] = v
    }
    return result
}

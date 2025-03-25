package storage

import (
    "sync"
)

// Storage интерфейс для работы с хранилищем
type Storage interface {
    UpdateGauge(name string, value float64)
    UpdateCounter(name string, value int64)
}

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

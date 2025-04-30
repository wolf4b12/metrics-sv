package storage

import (
    "errors"
    "sync"
    "log"
    "time"
)

// KVStorageInterface интерфейс для работы с KV-хранилищем
type KVStorageInterface interface {
    Set(key string, value any)
    Get(key string) (any, bool)
    Delete(key string)
    All() map[string]any
}

// MetricStorage — адаптер для работы с метриками
type MetricStorage struct {
    kv       KVStorageInterface
    gauges   map[string]float64
    counters map[string]int64
    mu       sync.RWMutex
    saveTicker      *time.Ticker
    wg              sync.WaitGroup
    stopCh          chan struct{}
}

// NewMetricStorage создаёт новый адаптер для работы с метриками
func NewMetricStorage(kv KVStorageInterface, restore bool, storeInterval time.Duration, filePath string) (*MetricStorage, error) {
    s := &MetricStorage{
        kv:       kv,
        gauges:   make(map[string]float64),
        counters: make(map[string]int64),
        saveTicker: nil,
        stopCh:    make(chan struct{}),
    }

    // Установка интервала сохранения
    if storeInterval > 0 {
        s.saveTicker = time.NewTicker(storeInterval)
        s.wg.Add(1)
        go s.StartPeriodicSaving(filePath)
    }
    // Обработка сигналов завершения

    if restore {
        err := s.LoadFromFile(filePath)
        if err != nil {
            log.Printf("Не удалось загрузить предыдущие метрики: %v\n", err)
        } else {
            log.Println("Предыдущие метрики успешно загружены.")
        }
    }

    return s, nil
}

// UpdateGauge обновляет значение gauge-метрики
func (s *MetricStorage) UpdateGauge(name string, value float64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.gauges[name] = value
    s.kv.Set(name, value)
}

// UpdateCounter увеличивает значение counter-метрики
func (s *MetricStorage) UpdateCounter(name string, value int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    currentValue, exists := s.counters[name]
    if !exists {
        currentValue = 0
    }
    s.counters[name] = currentValue + value
    s.kv.Set(name, currentValue+value)
}

// GetGauge получает текущее значение gauge-метрики
func (s *MetricStorage) GetGauge(name string) (float64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, exists := s.gauges[name]
    if !exists {
        return 0, errors.New("metric not found")
    }
    return value, nil
}

// GetCounter получает текущее значение counter-метрики
func (s *MetricStorage) GetCounter(name string) (int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, exists := s.counters[name]
    if !exists {
        return 0, errors.New("metric not found")
    }
    return value, nil
}

// AllMetrics возвращает список всех существующих метрик
func (s *MetricStorage) AllMetrics() map[string]map[string]any {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := make(map[string]map[string]any)
    result["gauges"] = make(map[string]any)
    for k, v := range s.gauges {
        result["gauges"][k] = v
    }
    result["counters"] = make(map[string]any)
    for k, v := range s.counters {
        result["counters"][k] = v
    }
    return result
}

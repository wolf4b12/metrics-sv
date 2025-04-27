package storage

import (
    "encoding/json"
    "errors"
    "os"
    "sync"
    "fmt"
)

// KVStorage — базовое хранилище ключ-значение
type KVStorage struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

// NewKVStorage создаёт новое KV-хранилище
func NewKVStorage() *KVStorage {
    return &KVStorage{
        data: make(map[string]interface{}),
    }
}







// Set устанавливает значение по ключу
func (s *KVStorage) Set(key string, value interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
}

// Get получает значение по ключу
func (s *KVStorage) Get(key string) (interface{}, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, exists := s.data[key]
    return value, exists
}

// Delete удаляет значение по ключу
func (s *KVStorage) Delete(key string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.data, key)
}

// All возвращает все данные
func (s *KVStorage) All() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

// MetricStorage — адаптер для работы с метриками
type MetricStorage struct {
    kv       *KVStorage
    gauges   map[string]float64
    counters map[string]int64
    mu       sync.RWMutex
    filePath string
    data     map[string]map[string]interface{}
}





// NewMetricStorage создаёт новый адаптер для работы с метриками
func NewMetricStorage(kv *KVStorage) *MetricStorage {
    return &MetricStorage{
        kv:       kv,
        gauges:   make(map[string]float64),
        counters: make(map[string]int64),
    }
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
func (s *MetricStorage) AllMetrics() map[string]map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := make(map[string]map[string]interface{})
    result["gauges"] = make(map[string]interface{})
    for k, v := range s.gauges {
        result["gauges"][k] = v
    }
    result["counters"] = make(map[string]interface{})
    for k, v := range s.counters {
        result["counters"][k] = v
    }
    return result
}

// LoadFromFile загружает метрики из файла
func (s *MetricStorage) LoadFromFile(filePath string) error {
    rawData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    err = json.Unmarshal(rawData, &s.data)
    if err != nil {
        return fmt.Errorf("ошибка разбора JSON: %v", err)
    }
    return nil
}

// SaveToFile сохраняет текущее состояние метрик в файл
func (s *MetricStorage) SaveToFile(filePath string) error {
    rawData, err := json.MarshalIndent(s.data, "", "\t")
    if err != nil {
        return err
    }
    return os.WriteFile(filePath, rawData, 0644)
}
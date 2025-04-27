package storage

import (
    "errors"
    "sync"
    "encoding/json"
    "os"
    "fmt"
)

// MetricStorage — адаптер для работы с метриками
type MetricStorage struct {
    kv *KVStorage
    mu sync.RWMutex
}

// NewMetricStorage создаёт новый адаптер для работы с метриками
func NewMetricStorage(kv *KVStorage) *MetricStorage {
    return &MetricStorage{
        kv: kv,
    }
}

// UpdateGauge обновляет значение gauge-метрики
func (s *MetricStorage) UpdateGauge(name string, value float64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.kv.Set(name, value)
}

// UpdateCounter увеличивает значение counter-метрики
func (s *MetricStorage) UpdateCounter(name string, value int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    currentValue, exists := s.kv.Get(name)
    if !exists {
        currentValue = int64(0)
    }
    s.kv.Set(name, currentValue.(int64)+value)
}

// GetGauge получает текущее значение gauge-метрики
func (s *MetricStorage) GetGauge(name string) (float64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, exists := s.kv.Get(name)
    if !exists {
        return 0, errors.New("metric not found")
    }
    return value.(float64), nil
}

// GetCounter получает текущее значение counter-метрики
func (s *MetricStorage) GetCounter(name string) (int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, exists := s.kv.Get(name)
    if !exists {
        return 0, errors.New("metric not found")
    }
    return value.(int64), nil
}

// AllMetrics возвращает список всех существующих метрик
// AllMetrics возвращает список всех существующих метрик
func (s *MetricStorage) AllMetrics() map[string]map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := make(map[string]map[string]interface{})
    for k, v := range s.kv.All() {
        result[k] = map[string]interface{}{"value": v}
    }
    return result
}


// LoadFromFile загружает метрики из файла
func (s *MetricStorage) LoadFromFile(filePath string) error {
    rawData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    var loadedData map[string]interface{}
    if err := json.Unmarshal(rawData, &loadedData); err != nil {
        return fmt.Errorf("ошибка разбора JSON: %v", err)
    }

    s.mu.Lock()
    defer s.mu.Unlock()
    s.kv.data = loadedData
    return nil
}

// SaveToFile сохраняет текущее состояние метрик в файл
func (s *MetricStorage) SaveToFile(filePath string) error {
    rawData, err := json.MarshalIndent(s.kv.data, "", "\t")
    if err != nil {
        return err
    }
    return os.WriteFile(filePath, rawData, 0644)
}
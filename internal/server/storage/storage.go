package storage

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"

    "github.com/wolf4b12/metrics-sv.git/internal/constant" // Импортируем константы
)

// MemStorage реализует хранение метрик в памяти
type MemStorage struct {
    mu       sync.RWMutex
    gauges   map[string]float64
    counters map[string]int64
}

// NewMemStorage создает новое хранилище метрик
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
// UpdateCounter обновление counter-метрики
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
        return 0, fmt.Errorf("metric not found")
    }
    return value, nil
}

// GetCounter получает текущее значение counter-метрики
func (s *MemStorage) GetCounter(name string) (int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, ok := s.counters[name]
    if !ok {
        return 0, fmt.Errorf("metric not found")
    }
    return value, nil
}

// AllMetrics возвращает список всех существующих метрик
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

// LoadFromFile восстанавливает метрики из файла
func (s *MemStorage) LoadFromFile(filePath string) error {
    rawData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    var loadedData map[string]map[string]interface{}
    if err := json.Unmarshal(rawData, &loadedData); err != nil {
        return fmt.Errorf("ошибка разбора JSON: %v", err)
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    for typ, metrics := range loadedData {
        switch typ {
        case constant.MetricTypeGauge:
            for key, val := range metrics {
                fVal, ok := val.(float64)
                if !ok {
                    return fmt.Errorf("некорректный тип данных для Gauges: %v", val)
                }
                s.gauges[key] = fVal
            }
        case constant.MetricTypeCounter:
            for key, val := range metrics {
                var iVal int64
                switch v := val.(type) {
                case float64:
                    // Преобразуем float64 в int64 с округлением вниз
                    iVal = int64(v)
                case int64:
                    iVal = v
                default:
                    return fmt.Errorf("некорректный тип данных для Counter: %T", val)
                }
                s.counters[key] = iVal
            }
        default:
            continue
        }
    }

    return nil
}


// SaveToFile сохраняет текущее состояние метрик в файл
func (s *MemStorage) SaveToFile(filePath string) error {
    allMetrics := s.AllMetrics()

    rawData, err := json.MarshalIndent(allMetrics, "", "\t") // marshal with indentation for readability
    if err != nil {
        return err
    }

    return os.WriteFile(filePath, rawData, 0644)
}
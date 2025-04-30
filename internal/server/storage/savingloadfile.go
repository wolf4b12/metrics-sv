package storage

import (
    "encoding/json"
    "os"
    "fmt"
    "log"
)

// LoadFromFile загружает метрики из файла
func (s *MetricStorage) LoadFromFile(filePath string) error {
    rawData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    var loadedData map[string]map[string]any
    if err := json.Unmarshal(rawData, &loadedData); err != nil {
        return fmt.Errorf("ошибка разбора JSON: %v", err)
    }

    s.mu.Lock()
    defer s.mu.Unlock()
    for metricType, metrics := range loadedData {
        for name, value := range metrics {
            switch metricType {
            case "gauges":
                if v, ok := value.(float64); ok {
                    s.metrics[name] = MetricValue{Type: Gauge, Value: v}
                    s.kv.Set(name, MetricValue{Type: Gauge, Value: v})
                }
            case "counters":
                if v, ok := value.(float64); ok {
                    intValue := int64(v)
                    s.metrics[name] = MetricValue{Type: Counter, Value: intValue}
                    s.kv.Set(name, MetricValue{Type: Counter, Value: intValue})
                }
            }
        }
    }
    return nil
}

// SaveToFile сохраняет текущее состояние метрик в файл
func (s *MetricStorage) SaveToFile(filePath string) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    data := make(map[string]map[string]any)
    data["gauges"] = make(map[string]any)
    data["counters"] = make(map[string]any)
    for k, v := range s.metrics {
        if v.Type == Gauge {
            data["gauges"][k] = v.Value
        } else if v.Type == Counter {
            data["counters"][k] = v.Value
        }
    }
    rawData, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }
    return os.WriteFile(filePath, rawData, 0644)
}

func (s *MetricStorage) StartPeriodicSaving(filePath string) {
    defer s.wg.Done()
    for {
        select {
        case <-s.stopCh:
            return
        case <-s.saveTicker.C:
            err := s.SaveToFile(filePath)
            if err != nil {
                log.Printf("Ошибка при сохранении метрик: %v\n", err)
            } else {
                log.Println("Метрики успешно сохранены.")
            }
        }
    }
}

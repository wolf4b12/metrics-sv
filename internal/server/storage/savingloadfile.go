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
                    s.gauges[name] = v
                    s.kv.Set(name, v)
                }
            case "counters":
                if v, ok := value.(float64); ok {
                    intValue := int64(v)
                    s.counters[name] = intValue
                    s.kv.Set(name, intValue)
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
    for k, v := range s.gauges {
        data["gauges"][k] = v
    }
    data["counters"] = make(map[string]any)
    for k, v := range s.counters {
        data["counters"][k] = v
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


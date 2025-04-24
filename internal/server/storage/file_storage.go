// file_storage.go
package storage

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "sync"
)

// FileStorage хранит метрики в файле
type FileStorage struct {
    filePath string
    data     map[string]map[string]interface{}
    mu       sync.Mutex
}

// NewFileStorage создаёт новое хранилище метрик в файле
func NewFileStorage(filePath string) (*FileStorage, error) {
    fs := &FileStorage{
        filePath: filePath,
        data:     make(map[string]map[string]interface{}),
    }
    err := fs.LoadFromFile(fs.filePath)
    if err != nil && !os.IsNotExist(err) {
        return nil, fmt.Errorf("ошибка загрузки данных из файла: %w", err)
    }
    return fs, nil
}

// UpdateGauge обновляет значение gauge-метрики
func (fs *FileStorage) UpdateGauge(name string, value float64) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    if _, exists := fs.data["gauges"]; !exists {
        fs.data["gauges"] = make(map[string]interface{})
    }
    fs.data["gauges"][name] = value
}

// UpdateCounter увеличивает значение counter-метрики
func (fs *FileStorage) UpdateCounter(name string, value int64) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    if _, exists := fs.data["counters"]; !exists {
        fs.data["counters"] = make(map[string]interface{})
    }
    currentValue, _ := fs.data["counters"][name].(int64)
    fs.data["counters"][name] = currentValue + value
}

// GetGauge получает текущее значение gauge-метрики
func (fs *FileStorage) GetGauge(name string) (float64, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    if gauges, exists := fs.data["gauges"]; exists {
        if value, ok := gauges[name]; ok {
            return value.(float64), nil
        }
    }
    return 0, errors.New("metric not found")
}

// GetCounter получает текущее значение counter-метрики
func (fs *FileStorage) GetCounter(name string) (int64, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    if counters, exists := fs.data["counters"]; exists {
        if value, ok := counters[name]; ok {
            return value.(int64), nil
        }
    }
    return 0, errors.New("metric not found")
}

// AllMetrics возвращает список всех существующих метрик
func (fs *FileStorage) AllMetrics() map[string]map[string]interface{} {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    return fs.data
}

// LoadFromFile загружает метрики из файла
func (fs *FileStorage) LoadFromFile(filePath string) error {
    rawData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    err = json.Unmarshal(rawData, &fs.data)
    if err != nil {
        return fmt.Errorf("ошибка разбора JSON: %v", err)
    }
    return nil
}

// SaveToFile сохраняет текущее состояние метрик в файл
func (fs *FileStorage) SaveToFile(filePath string) error {
    rawData, err := json.MarshalIndent(fs.data, "", "\t")
    if err != nil {
        return err
    }
    return os.WriteFile(filePath, rawData, 0644)
}
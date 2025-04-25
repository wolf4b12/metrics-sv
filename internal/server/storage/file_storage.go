package storage

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "sync"

//    "github.com/wolf4b12/metrics-sv.git/internal/constant"
)

// Тип хранения метрик в файле
type FileStorage struct {
    filePath string
    data     struct {
        Gauges   map[string]float64 `json:"gauges,omitempty"`
        Counters map[string]int64   `json:"counters,omitempty"`
    }
    mu sync.Mutex
}

// Создание нового экземпляра хранилища метрик
func NewFileStorage(filePath string) (*FileStorage, error) {
    fs := &FileStorage{
        filePath: filePath,
        data: struct {
            Gauges   map[string]float64 `json:"gauges,omitempty"`
            Counters map[string]int64   `json:"counters,omitempty"`
        }{
            make(map[string]float64),
            make(map[string]int64),
        },
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
    fs.data.Gauges[name] = value
}

// UpdateCounter увеличивает значение counter-метрики
func (fs *FileStorage) UpdateCounter(name string, value int64) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    currentValue, exists := fs.data.Counters[name]
    if !exists {
        currentValue = 0
    }
    fs.data.Counters[name] = currentValue + value
}

// GetGauge получает текущее значение gauge-метрики
func (fs *FileStorage) GetGauge(name string) (float64, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    value, exists := fs.data.Gauges[name]
    if !exists {
        return 0, errors.New("metric not found")
    }
    return value, nil
}

// GetCounter получает текущее значение counter-метрики
func (fs *FileStorage) GetCounter(name string) (int64, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    value, exists := fs.data.Counters[name]
    if !exists {
        return 0, errors.New("metric not found")
    }
    return value, nil
}

// AllMetrics возвращает список всех существующих метрик
func (fs *FileStorage) AllMetrics() map[string]map[string]interface{} {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    result := make(map[string]map[string]interface{})
    result["gauges"] = make(map[string]interface{})
    for k, v := range fs.data.Gauges {
        result["gauges"][k] = v
    }
    result["counters"] = make(map[string]interface{})
    for k, v := range fs.data.Counters {
        result["counters"][k] = v
    }
    return result
}

// LoadFromFile загружает метрики из файла
func (fs *FileStorage) LoadFromFile(filePath string) error {
    rawData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    var loadedData struct {
        Gauges   map[string]float64 `json:"gauges,omitempty"`
        Counters map[string]int64   `json:"counters,omitempty"`
    }

    if err := json.Unmarshal(rawData, &loadedData); err != nil {
        return fmt.Errorf("ошибка разбора JSON: %v", err)
    }

    fs.mu.Lock()
    defer fs.mu.Unlock()

    fs.data.Gauges = loadedData.Gauges
    fs.data.Counters = loadedData.Counters

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
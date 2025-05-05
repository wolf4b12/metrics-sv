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

    // MetricType тип метрики
    type MetricType string

    const (
        Gauge   MetricType = "gauge"
        Counter MetricType = "counter"
    )

    // MetricValue значение метрики
    type MetricValue struct {
        Type  MetricType
        Value any
    }

    // MetricStorage — адаптер для работы с метриками
    type MetricStorage struct {
        kv       KVStorageInterface
        metrics  map[string]MetricValue
        mu       sync.RWMutex
        saveTicker      *time.Ticker
        wg              sync.WaitGroup
        stopCh          chan struct{}
    }

    // NewMetricStorage создаёт новый адаптер для работы с метриками
    func NewMetricStorage(kv KVStorageInterface, restore bool, storeInterval time.Duration, filePath string) (*MetricStorage, error) {
        s := &MetricStorage{
            kv:       kv,
            metrics:  make(map[string]MetricValue),
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
        s.metrics[name] = MetricValue{Type: Gauge, Value: value}
        s.kv.Set(name, MetricValue{Type: Gauge, Value: value})
    }

    // UpdateCounter увеличивает значение counter-метрики
    func (s *MetricStorage) UpdateCounter(name string, value int64) {
        s.mu.Lock()
        defer s.mu.Unlock()
        currentValue, exists := s.metrics[name]
        if !exists {
            currentValue = MetricValue{Type: Counter, Value: int64(0)}
        }
        if currentValue.Type != Counter {
            log.Printf("Неверный тип метрики для %s: ожидается counter, получено %s\n", name, currentValue.Type)
            return
        }
        s.metrics[name] = MetricValue{Type: Counter, Value: currentValue.Value.(int64) + value}
        s.kv.Set(name, MetricValue{Type: Counter, Value: currentValue.Value.(int64) + value})
    }

    // GetGauge получает текущее значение gauge-метрики
    func (s *MetricStorage) GetGauge(name string) (float64, error) {
        s.mu.RLock()
        defer s.mu.RUnlock()
        value, exists := s.metrics[name]
        if !exists || value.Type != Gauge {
            return 0, errors.New("metric not found")
        }
        return value.Value.(float64), nil
    }

    // GetCounter получает текущее значение counter-метрики
    func (s *MetricStorage) GetCounter(name string) (int64, error) {
        s.mu.RLock()
        defer s.mu.RUnlock()
        value, exists := s.metrics[name]
        if !exists || value.Type != Counter {
            return 0, errors.New("metric not found")
        }
        return value.Value.(int64), nil
    }

    // AllMetrics возвращает список всех существующих метрик
    func (s *MetricStorage) AllMetrics() map[string]map[string]any {
        s.mu.RLock()
        defer s.mu.RUnlock()
        result := make(map[string]map[string]any)
        result["gauges"] = make(map[string]any)
        result["counters"] = make(map[string]any)
        for k, v := range s.metrics {
            if v.Type == Gauge {
                result["gauges"][k] = v.Value
            } else if v.Type == Counter {
                result["counters"][k] = v.Value
            }
        }
        return result
    }

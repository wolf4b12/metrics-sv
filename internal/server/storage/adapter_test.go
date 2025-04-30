package storage_test

import (
     "testing"
 
 
    "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
)

// kvStorageMock is a mock implementation of the KVStorageInterface
type kvStorageMock struct {
    data map[string]any
}

func (m *kvStorageMock) Set(key string, value any) {
    m.data[key] = value
}

func (m *kvStorageMock) Get(key string) (any, bool) {
    value, exists := m.data[key]
    return value, exists
}

func (m *kvStorageMock) Delete(key string) {
    delete(m.data, key)
}

func (m *kvStorageMock) All() map[string]any {
    return m.data
}

func TestNewMetricStorage(t *testing.T) {
    kv := &kvStorageMock{data: make(map[string]any)}
    s, err := storage.NewMetricStorage(kv, false, 0, "")
    if err != nil {
        t.Errorf("NewMetricStorage() error = %v, wantErr %v", err, false)
    }
    if s == nil {
        t.Errorf("NewMetricStorage() got nil, want non-nil")
    }
}

func TestMetricStorage_UpdateGauge(t *testing.T) {
    kv := &kvStorageMock{data: make(map[string]any)}
    s, _ := storage.NewMetricStorage(kv, false, 0, "")
    s.UpdateGauge("metric1", 1.23)
    if value, exists := kv.Get("metric1"); !exists || value != 1.23 {
        t.Errorf("UpdateGauge() got %v, want %v", value, 1.23)
    }
}

func TestMetricStorage_UpdateCounter(t *testing.T) {
    kv := &kvStorageMock{data: make(map[string]any)}
    s, _ := storage.NewMetricStorage(kv, false, 0, "")
    s.UpdateCounter("metric1", 1)
    if value, exists := kv.Get("metric1"); !exists || value != int64(1) {
        t.Errorf("UpdateCounter() got %v, want %v", value, int64(1))
    }
}

func TestMetricStorage_GetGauge(t *testing.T) {
    kv := &kvStorageMock{data: make(map[string]any)}
    s, _ := storage.NewMetricStorage(kv, false, 0, "")
    s.UpdateGauge("metric1", 1.23)
    value, err := s.GetGauge("metric1")
    if err != nil {
        t.Errorf("GetGauge() error = %v, wantErr %v", err, false)
    }
    if value != 1.23 {
        t.Errorf("GetGauge() got %v, want %v", value, 1.23)
    }
}

func TestMetricStorage_GetCounter(t *testing.T) {
    kv := &kvStorageMock{data: make(map[string]any)}
    s, _ := storage.NewMetricStorage(kv, false, 0, "")
    s.UpdateCounter("metric1", 1)
    value, err := s.GetCounter("metric1")
    if err != nil {
        t.Errorf("GetCounter() error = %v, wantErr %v", err, false)
    }
    if value != int64(1) {
        t.Errorf("GetCounter() got %v, want %v", value, int64(1))
    }
}

func TestMetricStorage_AllMetrics(t *testing.T) {
    kv := &kvStorageMock{data: make(map[string]any)}
    s, _ := storage.NewMetricStorage(kv, false, 0, "")
    s.UpdateGauge("metric1", 1.23)
    s.UpdateCounter("metric2", 1)
    allMetrics := s.AllMetrics()
    if len(allMetrics) != 2 {
        t.Errorf("AllMetrics() got %v, want %v", len(allMetrics), 2)
    }
    if gaugeValue, exists := allMetrics["gauges"]["metric1"]; !exists || gaugeValue != 1.23 {
        t.Errorf("AllMetrics() got gauge value %v, want %v", gaugeValue, 1.23)
    }
    if counterValue, exists := allMetrics["counters"]["metric2"]; !exists || counterValue != int64(1) {
        t.Errorf("AllMetrics() got counter value %v, want %v", counterValue, int64(1))
    }
}

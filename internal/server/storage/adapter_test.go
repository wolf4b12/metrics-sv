package storage

import (
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/wolf4b12/metrics-sv/mocks"
    "github.com/stretchr/testify/assert" 
)



func TestNewMetricStorage(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    kv := mocks.NewMockKVStorageInterface(ctrl)
    kv.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
    kv.EXPECT().Get(gomock.Any()).AnyTimes()
    kv.EXPECT().Delete(gomock.Any()).AnyTimes()
    kv.EXPECT().All().AnyTimes()

    s, err := NewMetricStorage(kv, false, 0, "")
    if err != nil {
        t.Errorf("NewMetricStorage() error = %v, wantErr %v", err, false)
    }
    if s == nil {
        t.Errorf("NewMetricStorage() got nil, want non-nil")
    }
}

func TestMetricStorage_UpdateGauge(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    kv := mocks.NewMockKVStorageInterface(ctrl)
    kv.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()

    s, _ := NewMetricStorage(kv, false, 0, "")
    s.UpdateGauge("metric1", 1.23)
}

func TestMetricStorage_UpdateCounter(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    kv := mocks.NewMockKVStorageInterface(ctrl)
    kv.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()

    s, _ := NewMetricStorage(kv, false, 0, "")
    s.UpdateCounter("metric1", 1)
}


func TestMetricStorage_GetGauge(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    kv := mocks.NewMockKVStorageInterface(ctrl)
    
    // Логируем все вызовы Get
    kv.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (interface{}, bool) {
        t.Logf("Get called with key: %s", key)
        return nil, false
    }).AnyTimes()
    
    kv.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()

    ms, _ := NewMetricStorage(kv, false, 0, "")

    ms.UpdateGauge("metric1", 1.23)
    value, err := ms.GetGauge("metric1")
    
    t.Logf("Result: %v, %v", value, err)
    // ... assertions
}


func TestMetricStorage_GetCounter(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    kv := mocks.NewMockKVStorageInterface(ctrl)
    
    // Логируем все вызовы Get для диагностики
    kv.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (interface{}, bool) {
        t.Logf("Get called with key: %s", key)
        if key == "counter1" {
            return MetricValue{Type: Counter, Value: int64(42)}, true
        }
        return nil, false
    }).AnyTimes()
    
    kv.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()

    ms, _ := NewMetricStorage(kv, false, 0, "")

    ms.UpdateGauge("metric1", 1.23)
    value, err := ms.GetCounter("metric1")
    
    t.Logf("Result: %v, %v", value, err)
    // ... assertions
}







func TestMetricStorage_AllMetrics(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    kv := mocks.NewMockKVStorageInterface(ctrl)

    // Ожидание вызова Set()
    kv.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()

    ms, _ := NewMetricStorage(kv, false, 0, "")

    // Добавляем метрику
    ms.UpdateGauge("metric1", 1.23)

    // Получаем все метрики
    allMetrics := ms.AllMetrics()

    // Проверяем, что результат содержит обе группы
    assert.Len(t, allMetrics, 2)

    // Проверяем группу gauges
    gauges, exists := allMetrics["gauges"]
    assert.True(t, exists)
    assert.Equal(t, 1.23, gauges["metric1"])

    // Проверяем группу counters
    counters, exists := allMetrics["counters"]
    assert.True(t, exists)
    assert.Empty(t, counters)
}
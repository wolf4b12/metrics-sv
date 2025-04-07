package storage_test

import (
 //   "reflect"
    "testing"

    "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
)

// Тест для UpdateGauge
func TestUpdateGauge(t *testing.T) {
    memStorage := storage.NewMemStorage()
    name := "my_gauge"
    value := 10.5

    memStorage.UpdateGauge(name, value)

    result, _ := memStorage.GetGauge(name)
    expected := value

    if result != expected {
        t.Errorf("TestUpdateGauge failed. Expected: %.2f, Got: %.2f", expected, result)
    }
}

// Тест для UpdateCounter
func TestUpdateCounter(t *testing.T) {
    memStorage := storage.NewMemStorage()
    name := "my_counter"
    value := int64(100)

    memStorage.UpdateCounter(name, value)

    result, _ := memStorage.GetCounter(name)
    expected := value

    if result != expected {
        t.Errorf("TestUpdateCounter failed. Expected: %d, Got: %d", expected, result)
    }
}

// Тест для GetGauge с несуществующей метрикой
func TestGetGauge_NotFound(t *testing.T) {
    memStorage := storage.NewMemStorage()
    name := "non_existent_gauge"

    _, err := memStorage.GetGauge(name)

    if err == nil || err.Error() != "metric not found" {
        t.Errorf("TestGetGauge_NotFound failed. Expected error 'metric not found', got: %v", err)
    }
}

// Тест для GetCounter с несуществующим счетчиком
func TestGetCounter_NotFound(t *testing.T) {
    memStorage := storage.NewMemStorage()
    name := "non_existent_counter"

    _, err := memStorage.GetCounter(name)

    if err == nil || err.Error() != "metric not found" {
        t.Errorf("TestGetCounter_NotFound failed. Expected error 'metric not found', got: %v", err)
    }
}


// Тест для ErrMetricNotFound
func TestErrMetricNotFound(t *testing.T) {
//	var name string
    memStorage := storage.NewMemStorage()
//    name = "non_existent_metric"

    err := memStorage.ErrMetricNotFound()

    if err == nil || err.Error() != "metric not found" {
        t.Errorf("TestErrMetricNotFound failed. Expected error 'metric not found', got: %v", err)
    }
}
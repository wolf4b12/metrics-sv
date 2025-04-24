// storage.go
package storage

import (

)

// Storage определяет методы для хранилищ метрик
type Storage interface {
    UpdateGauge(name string, value float64)
    UpdateCounter(name string, value int64)
    GetGauge(name string) (float64, error)
    GetCounter(name string) (int64, error)
    AllMetrics() map[string]map[string]interface{}
    LoadFromFile(filePath string) error
    SaveToFile(filePath string) error
}
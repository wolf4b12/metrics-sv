package metricsag

import (

 "encoding/json"

)



// Метрика для отправки на сервер
type Metrics struct {
    ID    string   `json:"id"`              // имя метрики
    MType string   `json:"type"`            // тип метрики (gauge или counter)
    Delta *int64   `json:"delta,omitempty"` // изменение значения (для счётчиков)
    Value *float64 `json:"value,omitempty"` // текущее значение (для датчиков)
}


type Metric interface {
    ID() string
    MarshalJSON() ([]byte, error)
}



// Метод интерфейса для идентификации метрики
func (m *Metrics) GetID() string { return m.ID }

// Маршаллизация метрики в JSON
func (m *Metrics) MarshalJSON() ([]byte, error) {
    return json.Marshal(m)
}
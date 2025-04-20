package metricssrv

// Метрика для отправки на сервер
type Metrics struct {
    ID    string   `json:"id"`              // имя метрики
    MType string   `json:"type"`            // тип метрики (gauge или counter)
    Delta *int64   `json:"delta,omitempty"` // изменение значения (для счётчиков)
    Value *float64 `json:"value,omitempty"` // текущее значение (для датчиков)
}
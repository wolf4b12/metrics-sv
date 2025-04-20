package agentmethods

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "runtime"
    "sync"
    "time"

    metrics "github.com/wolf4b12/metrics-sv.git/internal/agent/metrics"
)

// Метрика для отправки на сервер
type Metrics struct {
    ID    string   `json:"id"`              // Имя метрики
    MType string   `json:"type"`            // Тип метрики (gauge или counter)
    Delta *int64   `json:"delta,omitempty"` // Изменение значения (для счётчиков)
    Value *float64 `json:"value,omitempty"` // Текущее значение (для датчиков)
}

// Агент для сбора и отправки метрик
type Agent struct {
    Gauges         []Metrics
    Counters       []Metrics
    pollCount      int64
    mu             *sync.Mutex
    pollInterval   time.Duration
    reportInterval time.Duration
    addr           string
    client         *http.Client
}

// Новый агент
func NewAgent(poll, report time.Duration, addr string) *Agent {
    return &Agent{
        Gauges:         make([]Metrics, 0),
        Counters:       make([]Metrics, 0),
        pollInterval:   poll,
        reportInterval: report,
        addr:           addr,
        mu:             &sync.Mutex{},
        client:         &http.Client{Timeout: 5 * time.Second},
    }
}

// Сбор метрик
func (a *Agent) CollectMetrics() {
    for {
        a.mu.Lock()

        // Чистка старых коллекций перед сборкой новых данных
        a.Gauges = a.Gauges[:0]
        a.Counters = a.Counters[:0]

        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)

        // Собираем runtime-метрики и добавляем их в Gauges
        runtimeMetrics := metrics.GetRuntimeMetricsGauge(memStats)
        for key, value := range runtimeMetrics {
            a.Gauges = append(a.Gauges, Metrics{ID: key, MType: "gauge", Value: &value})
        }

        // Кастомные метрики добавляются в Counters
        customMetrics := metrics.GetCustomMetrics()
        for key, value := range customMetrics {
            a.Counters = append(a.Counters, Metrics{ID: key, MType: "counter", Delta: &value})
        }

        // Счётчик опроса PollCount
        a.pollCount++
        a.Counters = append(a.Counters, Metrics{ID: "PollCount", MType: "counter", Delta: &a.pollCount})

        a.mu.Unlock()
        time.Sleep(a.pollInterval)
    }
}

// Отправка собранных метрик
func (a *Agent) SendCollectedMetrics() {
    for {
        a.mu.Lock()

        // Получаем объединённый список метрик
        var metricsSlice []Metrics
        for _, gauge := range a.Gauges {
            if gauge.Value == nil {
                log.Printf("Отсутствует обязательное поле 'Value' для датчика '%s'\n", gauge.ID)
                continue
            }
            metricsSlice = append(metricsSlice, Metrics{
                ID:    gauge.ID,
                MType: "gauge",
                Value: gauge.Value,
            })
        }
        for _, counter := range a.Counters {
            if counter.Delta == nil {
                log.Printf("Отсутствует обязательное поле 'Delta' для счетчика '%s'\n", counter.ID)
                continue
            }
            metricsSlice = append(metricsSlice, Metrics{
                ID:    counter.ID,
                MType: "counter",
                Delta: counter.Delta,
            })
        }

        // Если метрики отсутствуют, пропускаем этот шаг
        if len(metricsSlice) == 0 {
            log.Println("Нет собранных метрик для отправки.")
            a.mu.Unlock()
            continue
        }

        // Маршализация метрик в JSON
        data, err := json.Marshal(metricsSlice)
        if err != nil {
            log.Printf("Ошибка маршализации метрик в JSON: %v\n", err)
            a.mu.Unlock()
            continue
        }

        // Формирование запроса
        url := fmt.Sprintf("http://%s/update", a.addr)
        req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
        if err != nil {
            log.Printf("Ошибка создания запроса: %v\n", err)
            a.mu.Unlock()
            continue
        }
        req.Header.Set("Content-Type", "application/json")

        // Выполнение запроса
        resp, err := a.client.Do(req)
        if err != nil {
            log.Printf("Ошибка отправки метрик: %v\n", err)
            a.mu.Unlock()
            continue
        }
        defer resp.Body.Close()

        // Проверка статуса ответа
        if resp.StatusCode != http.StatusOK {
            log.Printf("Получен непредвиденный статус-код (%d)\n", resp.StatusCode)
        }

        a.mu.Unlock()

        // Ждем отчетный интервал перед следующей отправкой
        time.Sleep(a.reportInterval)
    }
}
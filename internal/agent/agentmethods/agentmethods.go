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

    // Структура метрики для отправки на сервер
        type Metrics struct {
            ID    string   `json:"id"`              // имя метрики
            MType string   `json:"type"`            // тип метрики (gauge или counter)
            Delta *int64   `json:"delta,omitempty"` // изменение значения (для счётчиков)
            Value *float64 `json:"value,omitempty"` // текущее значение (для датчиков)
        }

    type Agent struct {
        Gauges         []Metrics
        Counters       []Metrics
        pollCount      int64
        mu             *sync.Mutex
        pollInterval   time.Duration
        reportInterval time.Duration
        addr           string
        client         *http.Client
//        useJsonFormat  bool // новое поле для переключения формата отправки
    }

    func NewAgent(poll, report time.Duration, addr string) *Agent {
        return &Agent{
            Gauges:         make([]Metrics, 0),
            Counters:       make([]Metrics, 0),
            pollInterval:   poll,
            reportInterval: report,
            addr:           addr,
            mu:             &sync.Mutex{},
            client:         &http.Client{Timeout: 5 * time.Second}, // используем общий клиент
//            useJsonFormat:  useJson,                                // задаём формат отправки
        }
    }

    func (a *Agent) CollectMetrics() { // собираем метрики
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


    func (a *Agent) SendCollectedMetrics() { // отправляем собранные метрики
        ticker := time.NewTicker(a.reportInterval)
        defer ticker.Stop()
    
        for range ticker.C {
            a.mu.Lock()
    
 {
                // Объединяем все метрики в единый срез
                var metricsSlice []Metrics
    
                // Добавляем gauges (датчики)
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
    
                // Добавляем counters (счетчики)
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
    
                // Преобразуем метрики в JSON
                data, err := json.Marshal(metricsSlice)
                if err != nil {
                    log.Printf("Ошибка преобразования метрик в JSON: %v\n", err)
                    a.mu.Unlock()
                    continue
                }
    
                // Формируем URL для отправки метрик
                url := fmt.Sprintf("http://%s/update", a.addr)
    
                // Создаем POST-запрос
                req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
                if err != nil {
                    log.Printf("Ошибка создания запроса: %v\n", err)
                    a.mu.Unlock()
                    continue
                }
    
                // Настраиваем заголовок Content-Type
                req.Header.Set("Content-Type", "application/json")
    
                // Выполняем запрос
                resp, err := a.client.Do(req)
                if err != nil {
                    log.Printf("Ошибка отправки метрик: %v\n", err)
                    a.mu.Unlock()
                    continue
                }
                defer resp.Body.Close()
    
                // Проверяем статус ответа
                if resp.StatusCode != http.StatusOK {
                    log.Printf("Получен неожиданный статус-код (%d)\n", resp.StatusCode)
                }
            }
    
            a.mu.Unlock()
        }
    }

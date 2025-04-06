package agentmethods

import (
    "fmt"
    "net/http"
    "reflect"
 //   "runtime"
    "sync"
    "testing"
    "time"
)

// TestNewAgent проверяет создание нового агента с правильными параметрами.
func TestNewAgent(t *testing.T) {
    pollInterval := 2 * time.Second
    reportInterval := 10 * time.Second
    addr := "example.com"

    expectedAgent := &Agent{
        gauges:         make(map[string]float64),
        counters:       make(map[string]int64),
        pollInterval:   pollInterval,
        reportInterval: reportInterval,
        addr:           addr,
        mu:             &sync.Mutex{},
    }

    actualAgent := NewAgent(pollInterval, reportInterval, addr)

    if !reflect.DeepEqual(expectedAgent, actualAgent) {
        t.Errorf("Expected agent %+v, got %+v", expectedAgent, actualAgent)
    }
}

// TestCollectMetrics проверяет сбор метрик агентом.
func TestCollectMetrics(t *testing.T) {
    // Создаем новый агент
    agent := NewAgent(2*time.Second, 10*time.Second, "example.com")

    // Запускаем сбор метрик
    go agent.CollectMetrics()

    time.Sleep(100 * time.Millisecond) // Даем немного времени на выполнение сбора метрик

    // Проверяем, что хотя бы одна метрика собрана
    if len(agent.gauges) == 0 || len(agent.counters) == 0 {
        t.Error("Метрики не были собраны")
    }

    // Проверяем наличие конкретных метрик
    expectedGaugeKeys := []string{"Alloc", "BuckHashSys", "Frees"}
    for _, key := range expectedGaugeKeys {
        if _, ok := agent.gauges[key]; !ok {
            t.Errorf("Отсутствует ожидаемая метрика '%s'", key)
        }
    }

    expectedCounterKeys := []string{"PollCount"}
    for _, key := range expectedCounterKeys {
        if _, ok := agent.counters[key]; !ok {
            t.Errorf("Отсутствует ожидаемый счетчик '%s'", key)
        }
    }
}

// TestSendMetrics проверяет отправку метрик агентом.
func TestSendMetrics(t *testing.T) {
    // Создаем новый агент
    agent := NewAgent(2*time.Second, 10*time.Second, "example.com")

    // Создаем mock клиента HTTP для тестирования отправки
    mockClient := &MockHTTPClient{}

    // Устанавливаем тестовые данные
    agent.gauges = map[string]float64{
        "TestGauge": 123.45,
    }
    agent.counters = map[string]int64{
        "TestCounter": 67890,
    }

    baseURL := fmt.Sprintf("http://%s/update", agent.addr)

    // Запускаем отправку метрик
    go agent.SendMetrics()

    time.Sleep(100 * time.Millisecond) // Даем немного времени на выполнение отправки

    // Проверяем, что запросы были сделаны
    if len(mockClient.requests) == 0 {
        t.Error("Запросы не были отправлены")
    }

    // Проверяем содержание запросов
    expectedRequests := []string{
        fmt.Sprintf("%s/gauge/TestGauge/123.450000", baseURL),
        fmt.Sprintf("%s/counter/TestCounter/67890", baseURL),
    }

    for i, request := range mockClient.requests {
        if request != expectedRequests[i] {
            t.Errorf("Неправильный запрос #%d: ожидалось '%s', получено '%s'", i+1, expectedRequests[i], request)
        }
    }
}

// MockHTTPClient имитирует поведение HTTP-клиента для тестирования.
type MockHTTPClient struct {
    requests []string
}

func (m *MockHTTPClient) Post(url, contentType string, body interface{}) (*http.Response, error) {
    m.requests = append(m.requests, url)
    return &http.Response{}, nil
}
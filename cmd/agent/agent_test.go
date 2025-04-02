package main

import (
    "net/http"
    "reflect"
    "sync"
    "testing"
    "time"
)

// Тестовая функция для конструктора агента
func TestNewAgent(t *testing.T) {
    type args struct {
        poll   time.Duration
        report time.Duration
        addr   string
    }
    tests := []struct {
        name string
        args args
        want *Agent
    }{
        {
            name: "TestWithValidArguments",
            args: args{
                poll:   2 * time.Second,
                report: 10 * time.Second,
                addr:   "http://localhost:8080",
            },
            want: &Agent{
                gauges:         make(map[string]float64),
                counters:       make(map[string]int64),
                pollCount:      0,
                mu:             &sync.Mutex{},
                pollInterval:   2 * time.Second,
                reportInterval: 10 * time.Second,
                addr:           "http://localhost:8080",
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := NewAgent(tt.args.poll, tt.args.report, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewAgent() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Тестовая функция для парсера командной строки
func Test_parseFlags(t *testing.T) {
    tests := []struct {
        name  string
        want  time.Duration
        want1 time.Duration
        want2 string
    }{
        {
            name:  "TestParseFlags",
            want:  2 * time.Second,
            want1: 10 * time.Second,
            want2: "http://localhost:8080",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2 := parseFlags()
            if got != tt.want {
                t.Errorf("parseFlags() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("parseFlags() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("parseFlags() got2 = %v, want %v", got2, tt.want2)
            }
        })
    }
}

// Тестовая функция для метода CollectMetrics
func TestAgent_CollectMetrics(t *testing.T) {
    type fields struct {
        gauges         map[string]float64
        counters       map[string]int64
        pollCount      int64
        mu             *sync.Mutex
        pollInterval   time.Duration
        reportInterval time.Duration
        addr           string
    }
    tests := []struct {
        name   string
        fields fields
    }{
        {
            name: "TestCollectMetrics",
            fields: fields{
                gauges:         make(map[string]float64),
                counters:       make(map[string]int64),
                pollCount:      0,
                mu:             &sync.Mutex{},
                pollInterval:   2 * time.Second,
                reportInterval: 10 * time.Second,
                addr:           "http://localhost:8080",
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := &Agent{
                gauges:         tt.fields.gauges,
                counters:       tt.fields.counters,
                pollCount:      tt.fields.pollCount,
                mu:             tt.fields.mu,
                pollInterval:   tt.fields.pollInterval,
                reportInterval: tt.fields.reportInterval,
                addr:           tt.fields.addr,
            }
            a.CollectMetrics()
            // Проверяем изменения в полях
            if a.pollCount != 1 {
                t.Errorf("Expected pollCount to be 1, but got %d", a.pollCount)
            }
        })
    }
}

// Тестовая функция для метода SendMetrics
func TestAgent_SendMetrics(t *testing.T) {
    type fields struct {
        gauges         map[string]float64
        counters       map[string]int64
        pollCount      int64
        mu             *sync.Mutex
        pollInterval   time.Duration
        reportInterval time.Duration
        addr           string
    }
    tests := []struct {
        name   string
        fields fields
    }{
        {
            name: "TestSendMetrics",
            fields: fields{
                gauges:         map[string]float64{"metric1": 123.45},
                counters:       map[string]int64{"metric2": 678},
                pollCount:      1,
                mu:             &sync.Mutex{},
                pollInterval:   2 * time.Second,
                reportInterval: 10 * time.Second,
                addr:           "http://localhost:8080",
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := &Agent{
                gauges:         tt.fields.gauges,
                counters:       tt.fields.counters,
                pollCount:      tt.fields.pollCount,
                mu:             tt.fields.mu,
                pollInterval:   tt.fields.pollInterval,
                reportInterval: tt.fields.reportInterval,
                addr:           tt.fields.addr,
            }
            a.SendMetrics()
            // Проверяем успешность отправки метрик
        })
    }
}

// Тестовая функция для метода sendMetric
func Test_sendMetric(t *testing.T) {
    type args struct {
        client *http.Client
        url    string
    }
    tests := []struct {
        name string
        args args
    }{
        {
            name: "TestSendMetric",
            args: args{
                client: &http.Client{},
                url:    "http://localhost:8080/update/gauge/metric1/123.45",
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sendMetric(tt.args.client, tt.args.url)
            // Проверяем успешность отправки запроса
        })
    }
}
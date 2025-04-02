package main

import (
    "net/http"
    "reflect"
    "sync"
    "testing"
    "time"
)

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
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := NewAgent(tt.args.poll, tt.args.report, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewAgent() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_parseFlags(t *testing.T) {
    tests := []struct {
        name  string
        want  time.Duration
        want1 time.Duration
        want2 string
    }{
        // TODO: Add test cases.
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

func TestAgent_CollectMetrics(t *testing.T) {
    type fields struct {
        gauges         map[string]float64
        counters       map[string]int64
        pollCount      int64
        mu             *sync.Mutex // Теперь указатель на Mutex
        pollInterval   time.Duration
        reportInterval time.Duration
        addr           string
    }
    tests := []struct {
        name   string
        fields fields
    }{
        // TODO: Add test cases.
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := &Agent{
                gauges:         tt.fields.gauges,
                counters:       tt.fields.counters,
                pollCount:      tt.fields.pollCount,
                mu:             tt.fields.mu, // Указатель на Mutex
                pollInterval:   tt.fields.pollInterval,
                reportInterval: tt.fields.reportInterval,
                addr:           tt.fields.addr,
            }
            a.CollectMetrics()
        })
    }
}



func Test_sendMetric(t *testing.T) {
    type args struct {
        client *http.Client
        url    string
    }
    tests := []struct {
        name string
        args args
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sendMetric(tt.args.client, tt.args.url)
        })
    }
}
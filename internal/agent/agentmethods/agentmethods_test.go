package agentmethods

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
		// TODO: Add test cases.
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
		})
	}
}

func TestAgent_SendCollectedMetrics(t *testing.T) {
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
		// TODO: Add test cases.
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
			a.SendCollectedMetrics()
		})
	}
}

func TestSendMetricToServer(t *testing.T) {
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
			SendMetricToServer(tt.args.client, tt.args.url)
		})
	}
}

package storage

import (
	"sync"
	"testing"
	"time"
)

func TestMetricStorage_LoadFromFile(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		mu         sync.RWMutex
		saveTicker *time.Ticker
		wg         sync.WaitGroup
		stopCh     chan struct{}
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricStorage{
				kv:         tt.fields.kv,
				gauges:     tt.fields.gauges,
				counters:   tt.fields.counters,
				mu:         tt.fields.mu,
				saveTicker: tt.fields.saveTicker,
				wg:         tt.fields.wg,
				stopCh:     tt.fields.stopCh,
			}
			if err := s.LoadFromFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("MetricStorage.LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricStorage_SaveToFile(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		mu         sync.RWMutex
		saveTicker *time.Ticker
		wg         sync.WaitGroup
		stopCh     chan struct{}
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricStorage{
				kv:         tt.fields.kv,
				gauges:     tt.fields.gauges,
				counters:   tt.fields.counters,
				mu:         tt.fields.mu,
				saveTicker: tt.fields.saveTicker,
				wg:         tt.fields.wg,
				stopCh:     tt.fields.stopCh,
			}
			if err := s.SaveToFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("MetricStorage.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricStorage_StartPeriodicSaving(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		mu         sync.RWMutex
		saveTicker *time.Ticker
		wg         sync.WaitGroup
		stopCh     chan struct{}
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricStorage{
				kv:         tt.fields.kv,
				gauges:     tt.fields.gauges,
				counters:   tt.fields.counters,
				mu:         tt.fields.mu,
				saveTicker: tt.fields.saveTicker,
				wg:         tt.fields.wg,
				stopCh:     tt.fields.stopCh,
			}
			s.StartPeriodicSaving(tt.args.filePath)
		})
	}
}

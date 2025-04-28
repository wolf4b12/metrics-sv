package storage

import (
	"reflect"
	"testing"
	"time"
)

func TestNewMetricStorage(t *testing.T) {
	type args struct {
		restore       bool
		storeInterval time.Duration
		filePath      string
	}
	tests := []struct {
		name    string
		args    args
		want    *MetricStorage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricStorage(tt.args.restore, tt.args.storeInterval, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetricStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		saveTicker *time.Ticker
		stopCh     chan struct{}
	}
	type args struct {
		name  string
		value float64
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
				saveTicker: tt.fields.saveTicker,
				stopCh:     tt.fields.stopCh,
			}
			s.UpdateGauge(tt.args.name, tt.args.value)
		})
	}
}

func TestMetricStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		saveTicker *time.Ticker
		stopCh     chan struct{}
	}
	type args struct {
		name  string
		value int64
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
				saveTicker: tt.fields.saveTicker,
				stopCh:     tt.fields.stopCh,
			}
			s.UpdateCounter(tt.args.name, tt.args.value)
		})
	}
}

func TestMetricStorage_GetGauge(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		saveTicker *time.Ticker
		stopCh     chan struct{}
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
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
				saveTicker: tt.fields.saveTicker,
				stopCh:     tt.fields.stopCh,
			}
			got, err := s.GetGauge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricStorage.GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricStorage.GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricStorage_GetCounter(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		saveTicker *time.Ticker
		stopCh     chan struct{}
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
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
				saveTicker: tt.fields.saveTicker,
				stopCh:     tt.fields.stopCh,
			}
			got, err := s.GetCounter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricStorage.GetCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricStorage.GetCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricStorage_AllMetrics(t *testing.T) {
	type fields struct {
		kv         *KVStorage
		gauges     map[string]float64
		counters   map[string]int64
		saveTicker *time.Ticker
		stopCh     chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]map[string]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricStorage{
				kv:         tt.fields.kv,
				gauges:     tt.fields.gauges,
				counters:   tt.fields.counters,
				saveTicker: tt.fields.saveTicker,
				stopCh:     tt.fields.stopCh,
			}
			if got := s.AllMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricStorage.AllMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

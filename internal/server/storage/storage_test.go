package storage

import (
	"reflect"
	"sync"
	"testing"
)

func TestMemStorage_ErrMetricNotFound(t *testing.T) {
	type fields struct {
		mu       sync.RWMutex
		gauges   map[string]float64
		counters map[string]int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				mu:       tt.fields.mu,
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			if err := s.ErrMetricNotFound(); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.ErrMetricNotFound() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		mu       sync.RWMutex
		gauges   map[string]float64
		counters map[string]int64
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
			s := &MemStorage{
				mu:       tt.fields.mu,
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			s.UpdateGauge(tt.args.name, tt.args.value)
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		mu       sync.RWMutex
		gauges   map[string]float64
		counters map[string]int64
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
			s := &MemStorage{
				mu:       tt.fields.mu,
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			s.UpdateCounter(tt.args.name, tt.args.value)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		mu       sync.RWMutex
		gauges   map[string]float64
		counters map[string]int64
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
			s := &MemStorage{
				mu:       tt.fields.mu,
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			got, err := s.GetGauge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		mu       sync.RWMutex
		gauges   map[string]float64
		counters map[string]int64
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
			s := &MemStorage{
				mu:       tt.fields.mu,
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			got, err := s.GetCounter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.GetCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_AllMetrics(t *testing.T) {
	type fields struct {
		mu       sync.RWMutex
		gauges   map[string]float64
		counters map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				mu:       tt.fields.mu,
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			if got := s.AllMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.AllMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

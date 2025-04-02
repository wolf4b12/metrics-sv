package storage

import (
	"reflect"
//	"sync"
	"testing"
)

// Тесты
func TestMemStorage_ErrMetricNotFound(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Error returned", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{}
			err := s.ErrMetricNotFound()
			if (err != nil) != tt.wantErr {
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
		{"Create new storage", NewMemStorage()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemStorage()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_UpdateGauge(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]float64
		args   struct{ name string; value float64 }
		want   float64
	}{
		{"Update gauge metric", map[string]float64{}, struct{ name string; value float64 }{"testGauge", 123.45}, 123.45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{gauges: tt.fields}
			s.UpdateGauge(tt.args.name, tt.args.value)
			if got := s.gauges[tt.args.name]; got != tt.want {
				t.Errorf("UpdateGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]int64
		args   struct{ name string; value int64 }
		want   int64
	}{
		{"Update counter metric", map[string]int64{}, struct{ name string; value int64 }{"testCounter", 10}, 10},
		{"Update existing counter", map[string]int64{"testCounter": 5}, struct{ name string; value int64 }{"testCounter", 10}, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{counters: tt.fields}
			s.UpdateCounter(tt.args.name, tt.args.value)
			if got := s.counters[tt.args.name]; got != tt.want {
				t.Errorf("UpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	tests := []struct {
		name    string
		fields  map[string]float64
		args    string
		want    float64
		wantErr bool
	}{
		{"Get existing gauge metric", map[string]float64{"testGauge": 123.45}, "testGauge", 123.45, false},
		{"Get non-existing gauge metric", map[string]float64{}, "nonExisting", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{gauges: tt.fields}
			got, err := s.GetGauge(tt.args)
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
	tests := []struct {
		name    string
		fields  map[string]int64
		args    string
		want    int64
		wantErr bool
	}{
		{"Get existing counter metric", map[string]int64{"testCounter": 10}, "testCounter", 10, false},
		{"Get non-existing counter metric", map[string]int64{}, "nonExisting", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{counters: tt.fields}
			got, err := s.GetCounter(tt.args)
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
	tests := []struct {
		name   string
		fields struct {
			gauges   map[string]float64
			counters map[string]int64
		}
		want map[string]map[string]interface{}
	}{
		{"Get all metrics", struct {
			gauges   map[string]float64
			counters map[string]int64
		}{gauges: map[string]float64{"testGauge": 123.45}, counters: map[string]int64{"testCounter": 10}}, map[string]map[string]interface{}{
			"gauges":   {"testGauge": 123.45},
			"counters": {"testCounter": 10},
		}},
		{"Empty storage", struct {
			gauges   map[string]float64
			counters map[string]int64
		}{gauges: map[string]float64{}, counters: map[string]int64{}}, map[string]map[string]interface{}{
			"gauges":   {},
			"counters": {},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			if got := s.AllMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.AllMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
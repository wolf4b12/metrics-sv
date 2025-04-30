package storage

import (
	"testing"
	"time"
)

func TestMetricStorage_LoadFromFile(t *testing.T) {
	type fields struct {
		kv         KVStorageInterface
		metrics    map[string]MetricValue
		saveTicker *time.Ticker
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
				metrics:    tt.fields.metrics,
				saveTicker: tt.fields.saveTicker,
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
		kv         KVStorageInterface
		metrics    map[string]MetricValue
		saveTicker *time.Ticker
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
				metrics:    tt.fields.metrics,
				saveTicker: tt.fields.saveTicker,
				stopCh:     tt.fields.stopCh,
			}
			if err := s.SaveToFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("MetricStorage.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

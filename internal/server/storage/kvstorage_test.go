package storage

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewKVStorage(t *testing.T) {
	tests := []struct {
		name string
		want *KVStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKVStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKVStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Set(t *testing.T) {
	type fields struct {
		data map[string]any
		mu   sync.RWMutex
	}
	type args struct {
		key   string
		value any
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
			s := &KVStorage{
				data: tt.fields.data,
				mu:   tt.fields.mu,
			}
			s.Set(tt.args.key, tt.args.value)
		})
	}
}

func TestKVStorage_Get(t *testing.T) {
	type fields struct {
		data map[string]any
		mu   sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KVStorage{
				data: tt.fields.data,
				mu:   tt.fields.mu,
			}
			got, got1 := s.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KVStorage.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("KVStorage.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKVStorage_Delete(t *testing.T) {
	type fields struct {
		data map[string]any
		mu   sync.RWMutex
	}
	type args struct {
		key string
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
			s := &KVStorage{
				data: tt.fields.data,
				mu:   tt.fields.mu,
			}
			s.Delete(tt.args.key)
		})
	}
}

func TestKVStorage_All(t *testing.T) {
	type fields struct {
		data map[string]any
		mu   sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KVStorage{
				data: tt.fields.data,
				mu:   tt.fields.mu,
			}
			if got := s.All(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KVStorage.All() = %v, want %v", got, tt.want)
			}
		})
	}
}

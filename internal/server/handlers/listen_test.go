package handlers

import (
	"net/http"
	"reflect"
	"testing"

	storage "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
)

func TestListMetricsHandler(t *testing.T) {
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListMetricsHandler(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListMetricsHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

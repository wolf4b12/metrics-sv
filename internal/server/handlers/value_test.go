package handlers

import (
	"net/http"
	"reflect"
	"testing"
)

func TestValueHandler(t *testing.T) {
	type args struct {
		storage GetStorage
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
			if got := ValueHandler(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValueHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

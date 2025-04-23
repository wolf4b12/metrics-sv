package handlers

import (
	"net/http"
	"reflect"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	type args struct {
		storage UpdateStorage
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateHandler(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

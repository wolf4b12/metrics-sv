package main

import (
	"net/http"
	"testing"
)

func TestSendMetric(t *testing.T) {
	type args struct {
		client        *http.Client
		serverAddress string
		metricType    string
		metricName    string
		value         interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMetric(tt.args.client, tt.args.serverAddress, tt.args.metricType, tt.args.metricName, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendAllMetrics(t *testing.T) {
	type args struct {
		client        *http.Client
		serverAddress string
		metrics       *Metrics
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendAllMetrics(tt.args.client, tt.args.serverAddress, tt.args.metrics)
		})
	}
}

func TestRunAgent(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunAgent(tt.args.cfg)
		})
	}
}

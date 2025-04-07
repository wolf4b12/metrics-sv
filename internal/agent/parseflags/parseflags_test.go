package parseflags

import (
	"testing"
	"time"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name  string
		want  time.Duration
		want1 time.Duration
		want2 string
	}{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := ParseFlags()
			if got != tt.want {
				t.Errorf("ParseFlags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseFlags() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ParseFlags() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

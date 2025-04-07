// internal/server/server.go
package srv

import (
	"net/http"
	"testing"
	"github.com/go-chi/chi/v5"
)

func TestServer_Run(t *testing.T) {
	type fields struct {
		router *chi.Mux
		server *http.Server
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				router: tt.fields.router,
				server: tt.fields.server,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package logger

import (
	"net/http"
	"testing"

	"go.uber.org/zap"
    "net/http/httptest"
    "io"
    "log"
    "strings"
)

func TestLoggingMiddleware(t *testing.T) {
    type args struct {
        logger *zap.Logger
        handler http.Handler
        reqMethod string
        reqPath string
    }

    tests := []struct {
        name     string
        args     args
        wantCode int
    }{
        {
            name: "Test logging middleware with simple handler",
            args: args{
                logger: zap.NewExample(),
                handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    w.WriteHeader(http.StatusOK)
                    w.Write([]byte("Hello World"))
                }),
                reqMethod: "GET",
                reqPath:   "/test",
            },
            wantCode: http.StatusOK,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            recorder := httptest.NewRecorder()
            request := httptest.NewRequest(tt.args.reqMethod, tt.args.reqPath, nil)

            mwHandler := LoggingMiddleware(tt.args.logger)(tt.args.handler)
            mwHandler.ServeHTTP(recorder, request)

            respBody, err := io.ReadAll(recorder.Body)
            if err != nil {
                log.Fatalf("failed to read response body: %s\n", err.Error())
            }

            if recorder.Code != tt.wantCode || strings.TrimSpace(string(respBody)) != "Hello World" {
                t.Errorf("Expected code: %d, actual: %d. Expected body: Hello World, actual: %s", tt.wantCode, recorder.Code, respBody)
            }
        })
    }
}
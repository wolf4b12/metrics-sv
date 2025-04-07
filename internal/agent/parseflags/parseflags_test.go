package parseflags
import (
    "flag"
//    "log"
    "os"
//    "strconv"
//    "sync"
    "testing"
    "time"
)

// Функция для тестирования
func TestParseFlags(t *testing.T) {
    tests := []struct {
        name               string
        env                map[string]string
        args               []string
        wantPollInterval   time.Duration
        wantReportInterval time.Duration
        wantAddr           string
        wantErr            bool
    }{
        {
            name: "Read from environment variables",
            env: map[string]string{
                "ADDRESS":          "example.com:80",
                "REPORT_INTERVAL":  "15",
                "POLL_INTERVAL":    "3",
            },
            args:               nil,
            wantPollInterval:   3 * time.Second,
            wantReportInterval: 15 * time.Second,
            wantAddr:           "example.com:80",
            wantErr:            false,
        },
        {
            name: "Use command line flags",
            env:  nil,
            args: []string{"-a", "custom.example.com:9000", "-r", "20", "-p", "5"},
            wantPollInterval:   5 * time.Second,
            wantReportInterval: 20 * time.Second,
            wantAddr:           "custom.example.com:9000",
            wantErr:            false,
        },
        {
            name: "Default values with no arguments or environment variables",
            env:  nil,
            args: nil,
            wantPollInterval:   2 * time.Second,
            wantReportInterval: 10 * time.Second,
            wantAddr:           "localhost:8080",
            wantErr:            false,
        },
        {
            name: "Invalid REPORT_INTERVAL in environment variable",
            env: map[string]string{
                "REPORT_INTERVAL": "abc",
            },
            args:     nil,
            wantErr:  true,
        },
        {
            name: "Unknown flag",
            env:  nil,
            args: []string{"-z", "value"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Устанавливаем временные переменные окружения
            if tt.env != nil {
                for key, value := range tt.env {
                    os.Setenv(key, value)
                }
            }

            // Сбрасываем переменные окружения после каждого теста
            defer func() {
                for key := range tt.env {
                    os.Unsetenv(key)
                }
            }()

            // Устанавливаем аргументы командной строки
            flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
            if tt.args != nil {
                flag.CommandLine.Parse(tt.args)
            }

            var gotErr bool

            _, _, _ = ParseFlags()

            if tt.wantErr {
                if !gotErr {
                    t.Errorf("ParseFlags() did not return an error as expected")
                }
            } else {
                if gotErr {
                    t.Errorf("ParseFlags() returned an unexpected error")
                }
            }
        })
    }
}
package srv

import (
    "bytes"
 //   "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
 //   "github.com/go-chi/chi/v5/middleware"
 "github.com/stretchr/testify/assert"  // Добавляем assert
 "github.com/stretchr/testify/require" // Добавляем require

    storage  "github.com/wolf4b12/metrics-sv.git/internal/server/storage"
    handler "github.com/wolf4b12/metrics-sv.git/internal/server/handlers"
)

// TestNewServer проверяет создание нового сервера с правильными настройками.
func TestNewServer(t *testing.T) {
    addr := ":8080"
    server := NewServer(addr)
    defer server.server.Close()

    assert.NotNil(t, server.router)
    assert.Equal(t, server.server.Addr, addr)
}

// TestUpdateHandler проверяет обработку POST-запросов на обновление метрик.
func TestUpdateHandler(t *testing.T) {
    testCases := []struct {
        name          string
        method        string
        url           string
        body          string
        expectedCode  int
        expectedBody  string
    }{
        {
            name:          "Valid update",
            method:        "POST",
            url:           "/update/gauge/test_gauge/42.0",
            body:          "",
            expectedCode:  http.StatusOK,
            expectedBody:  "",
        },
        {
            name:          "Invalid type",
            method:        "POST",
            url:           "/update/invalid_type/test_gauge/42.0",
            body:          "",
            expectedCode:  http.StatusBadRequest,
            expectedBody:  "Unknown metric type invalid_type\n",
        },
        {
            name:          "Empty metric name",
            method:        "POST",
            url:           "/update/gauge//42.0",
            body:          "",
            expectedCode:  http.StatusBadRequest,
            expectedBody:  "Empty metric name\n",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(tc.body))
            require.NoError(t, err)

            rr := httptest.NewRecorder()

            memStorage := storage.NewMemStorage()
            router := chi.NewRouter()
            router.Post("/update/{metricType}/{metricName}/{metricValue}", handler.UpdateHandler(memStorage))

            router.ServeHTTP(rr, req)

            assert.Equal(t, tc.expectedCode, rr.Code)
            assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
        })
    }
}

// TestValueHandler проверяет обработку GET-запросов на получение значения метрики.
func TestValueHandler(t *testing.T) {
    testCases := []struct {
        name          string
        method        string
        url           string
        expectedCode  int
        expectedBody  string
    }{
        {
            name:          "Existing gauge",
            method:        "GET",
            url:           "/value/gauge/test_gauge",
            expectedCode:  http.StatusOK,
            expectedBody:  `{"data":"42.0"}`,
        },
        {
            name:          "Non-existing gauge",
            method:        "GET",
            url:           "/value/gauge/non_existing_gauge",
            expectedCode:  http.StatusNotFound,
            expectedBody:  "Metric not found\n",
        },
        {
            name:          "Invalid type",
            method:        "GET",
            url:           "/value/invalid_type/test_gauge",
            expectedCode:  http.StatusBadRequest,
            expectedBody:  "Unknown metric type invalid_type\n",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            req, err := http.NewRequest(tc.method, tc.url, nil)
            require.NoError(t, err)

            rr := httptest.NewRecorder()

            memStorage := storage.NewMemStorage()
            router := chi.NewRouter()
            router.Get("/value/{metricType}/{metricName}", handler.ValueHandler(memStorage))

            router.ServeHTTP(rr, req)

            assert.Equal(t, tc.expectedCode, rr.Code)
            assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
        })
    }
}

// TestListMetricsHandler проверяет обработку GET-запросов на получение списка всех метрик.
func TestListMetricsHandler(t *testing.T) {
    testCases := []struct {
        name          string
        method        string
        url           string
        expectedCode  int
        expectedBody  string
    }{
        {
            name:          "Empty list",
            method:        "GET",
            url:           "/",
            expectedCode:  http.StatusOK,
            expectedBody:  "[]",
        },
        {
            name:          "With existing metrics",
            method:        "GET",
            url:           "/",
            expectedCode:  http.StatusOK,
            expectedBody:  `[{"name":"test_gauge","value":"42.0","metricType":"gauge"}]`,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            req, err := http.NewRequest(tc.method, tc.url, nil)
            require.NoError(t, err)

            rr := httptest.NewRecorder()

            memStorage := storage.NewMemStorage()
            router := chi.NewRouter()
            router.Get("/", handler.ListMetricsHandler(memStorage))

            router.ServeHTTP(rr, req)

            assert.Equal(t, tc.expectedCode, rr.Code)
            assert.Equal(t, tc.expectedBody, strings.TrimSpace(rr.Body.String()))
        })
    }
}

// TestRun проверяет успешный запуск сервера.
func TestRun(t *testing.T) {
    addr := ":8000"
    server := NewServer(addr)

    done := make(chan bool)
    go func() {
        err := server.Run()
        require.NoError(t, err)
        close(done)
    }()

    // Ждем, пока сервер запустится
    select {
    case <-done:
        break
    case <-time.After(500 * time.Millisecond):
        t.Fatal("Timed out waiting for the server to start")
    }

    resp, err := http.Get(fmt.Sprintf("http://%s/", addr))
    require.NoError(t, err)
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    require.NoError(t, err)

    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Contains(t, string(body), "Hello from Metrics Server!")
}
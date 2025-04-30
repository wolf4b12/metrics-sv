package handlers_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

   "github.com/wolf4b12/metrics-sv.git/internal/server/handlers"
)

// allMetricsStorageMock is a mock implementation of the allMetricsStorage interface
type allMetricsStorageMock struct{}

func (m *allMetricsStorageMock) AllMetrics() map[string]map[string]interface{} {
    return map[string]map[string]interface{}{
        "gauge": {
            "metric1": 1.23,
            "metric2": 4.56,
        },
        "counter": {
            "metric3": 7,
            "metric4": 8,
        },
    }
}

func TestListMetricsHandler(t *testing.T) {
    // Create a mock storage
    storage := &allMetricsStorageMock{}

    // Create a request to pass to our handler
    req, err := http.NewRequest("GET", "/metrics", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Call the handler with the mock storage and the request
    handler := handlers.ListMetricsHandler(storage)
    handler(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check the response body
    expected := "<html><body><h2>gauge Metrics:</h2><p>metric1: 1.23</p><p>metric2: 4.56</p><h2>counter Metrics:</h2><p>metric3: 7</p><p>metric4: 8</p></body></html>"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

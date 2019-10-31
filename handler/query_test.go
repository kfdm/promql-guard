package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// https://blog.questionable.services/article/testing-http-handlers-go/

func TestQuery(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add Test Query
	q := req.URL.Query()
	q.Add("query", "node_filesystem_free_bytes / node_filesystem_size_bytes")
	q.Add("start", "12345")
	q.Add("end", "54321")
	q.Add("step", "120")
	req.URL.RawQuery = q.Encode()

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Query(logger)
	targetHandler.ServeHTTP(rr, req)
	level.Debug(logger).Log("query", req.URL.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSeries(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/series", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add Test Query
	q := req.URL.Query()
	q.Add("match[]", "node_exporter_build_info")
	req.URL.RawQuery = q.Encode()

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Series(logger)
	targetHandler.ServeHTTP(rr, req)
	level.Debug(logger).Log("query", req.URL.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

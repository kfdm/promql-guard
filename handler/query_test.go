package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kfdm/promql-guard/config"
)

// https://blog.questionable.services/article/testing-http-handlers-go/

func TestMissingHost(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)
	var config, err = config.LoadFile("../example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Query(logger, config)
	targetHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestQuery(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	var config, err = config.LoadFile("../example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = config.VirtualHosts[0].Hostname

	// Add Test Query
	q := req.URL.Query()
	q.Add("query", "node_filesystem_free_bytes / node_filesystem_size_bytes")
	q.Add("start", "12345")
	q.Add("end", "54321")
	q.Add("step", "120")
	req.URL.RawQuery = q.Encode()

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Query(logger, config)
	targetHandler.ServeHTTP(rr, req)
	level.Debug(logger).Log("query", req.URL.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestSeries(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	var config, err = config.LoadFile("../example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/series", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = config.VirtualHosts[0].Hostname

	// Add Test Query
	q := req.URL.Query()
	q.Add("match[]", "node_exporter_build_info")
	req.URL.RawQuery = q.Encode()

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Series(logger, config)
	targetHandler.ServeHTTP(rr, req)
	level.Debug(logger).Log("query", req.URL.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

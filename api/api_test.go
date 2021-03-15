package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/kfdm/promql-guard/config"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/util/testutil"
)

// https://blog.questionable.services/article/testing-http-handlers-go/

func init() {
	// For finding our test configuration files
	os.Chdir("..")
}

func TestMissingAuth(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	var config, err = config.LoadFile("guard.yml")
	testutil.Ok(t, err)

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	testutil.Ok(t, err)

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Query(logger, config)
	targetHandler.ServeHTTP(rr, req)

	testutil.Equals(t, rr.Code, http.StatusUnauthorized)
}

func TestGetQuery(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo())

	var config, err = config.LoadFile("guard.yml")
	testutil.Ok(t, err)

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	testutil.Ok(t, err)
	req.SetBasicAuth("tenantA", "tenantA")

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

	testutil.Equals(t, rr.Code, http.StatusOK)
}

func TestPostQuery(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo())

	var config, err = config.LoadFile("guard.yml")
	testutil.Ok(t, err)

	// Add Test Query
	q := url.Values{}
	q.Add("query", "node_filesystem_free_bytes / node_filesystem_size_bytes")
	q.Add("start", "12345")
	q.Add("end", "54321")
	q.Add("step", "120")
	data := q.Encode()

	// Build Reqeust
	req, err := http.NewRequest("POST", "/api/v1/query", strings.NewReader(data))
	testutil.Ok(t, err)
	req.SetBasicAuth("tenantA", "tenantA")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Query(logger, config)
	targetHandler.ServeHTTP(rr, req)
	level.Debug(logger).Log("query", req.URL.String())

	testutil.Equals(t, rr.Code, http.StatusOK)
}

func TestSeries(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo())

	var config, err = config.LoadFile("guard.yml")
	testutil.Ok(t, err)

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/series", nil)
	testutil.Ok(t, err)
	req.SetBasicAuth("tenantB", "tenantB")

	// Add Test Query
	q := req.URL.Query()
	q.Add("match[]", "node_exporter_build_info")
	req.URL.RawQuery = q.Encode()

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Series(logger, config)
	targetHandler.ServeHTTP(rr, req)
	level.Debug(logger).Log("query", req.URL.String())

	testutil.Equals(t, rr.Code, http.StatusOK)
}

package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/proxy"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/require"
)

// https://blog.questionable.services/article/testing-http-handlers-go/

func init() {
	// For finding our test configuration files
	os.Chdir("..")
}

func TestMissingAuth(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	var config, err = config.LoadFile("guard.yml")
	require.NoError(t, err)

	api := NewAPI(config, logger, nil)

	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	require.NoError(t, err)

	// Test Request
	rr := httptest.NewRecorder()
	api.Query().ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusUnauthorized)
}

func TestGetQuery(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	var config, err = config.LoadFile("guard.yml")
	require.NoError(t, err)

	var mockResult = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "GET", r.Method)
		// Header from config gets added
		require.Equal(t, "Token Foo", r.Header.Get("Authorization"))
		proxy.ExpectedPromql(t,
			r.FormValue("query"),
			"a{service=\"tenantA\"} / b{service=\"tenantA\"}",
		)
	}

	proxy_ := proxy.NewMock(logger, mockResult)
	api := NewAPI(config, logger, proxy_)

	// Build Reqeust
	q := url.Values{}
	q.Add("query", "a / b")

	// https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries
	req, err := proxy.Get("/api/v1/query", q)
	require.NoError(t, err)
	req.SetBasicAuth("tenantA", "tenantA")

	// Test Request
	rr := httptest.NewRecorder()
	api.Query().ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusOK)
}

func TestPostQuery(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	var config, err = config.LoadFile("guard.yml")
	require.NoError(t, err)

	var mockResult = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		// Header from config gets added
		require.Equal(t, "Token Foo", r.Header.Get("Authorization"))
		proxy.ExpectedPromql(t,
			r.FormValue("query"),
			"a{service=\"tenantA\"} / b{service=\"tenantA\"}",
		)
	}

	proxy_ := proxy.NewMock(logger, mockResult)
	api := NewAPI(config, logger, proxy_)

	// Build Reqeust
	q := url.Values{}
	q.Add("query", "a / b")

	// https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries
	req, err := proxy.Post("/api/v1/query", q)
	require.NoError(t, err)
	req.SetBasicAuth("tenantA", "tenantA")

	// Test Request
	rr := httptest.NewRecorder()
	api.Query().ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusOK)
}

func TestPostQueryRange(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)

	var config, err = config.LoadFile("guard.yml")
	require.NoError(t, err)

	var mockResult = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		// Header from config gets added
		require.Equal(t, "Token Foo", r.Header.Get("Authorization"))
		proxy.ExpectedPromql(t,
			r.FormValue("query"),
			"test{service=\"tenantA\"}[1m] offset 1w",
		)
	}

	proxy_ := proxy.NewMock(logger, mockResult)
	api := NewAPI(config, logger, proxy_)

	// Build Reqeust
	q := url.Values{}
	q.Add("query", "test[1m0s] offset 1w")
	q.Add("start", "12345")
	q.Add("end", "54321")
	q.Add("step", "120")

	// https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
	req, err := proxy.Post("/api/v1/query_range", q)
	require.NoError(t, err)
	req.SetBasicAuth("tenantA", "tenantA")

	// Test Request
	rr := httptest.NewRecorder()
	api.QueryRange().ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusOK)
}

func TestGetSeries(t *testing.T) {
	var logger = log.NewJSONLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo())

	var config, err = config.LoadFile("guard.yml")
	require.NoError(t, err)

	var mockResult = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "GET", r.Method)
		// Tenant B doesn't set a specific header so it just gets passed as is
		require.Equal(t, "Basic dGVuYW50Qjp0ZW5hbnRC", r.Header.Get("Authorization"))
		proxy.ExpectedPromql(t,
			r.FormValue("match[]"),
			"node_exporter_build_info{app=~\"appY|appZ\"}",
		)
	}

	api := API{
		config: config,
		logger: logger,
		proxy:  proxy.NewMock(logger, mockResult),
	}

	// Build Reqeust
	q := url.Values{}
	q.Add("match[]", "node_exporter_build_info")

	// https://prometheus.io/docs/prometheus/latest/querying/api/#finding-series-by-label-matchers
	req, err := proxy.Get("/api/v1/series", q)
	require.NoError(t, err)
	req.SetBasicAuth("tenantB", "tenantB")

	// Test Request
	rr := httptest.NewRecorder()
	api.Series().ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusOK)
}

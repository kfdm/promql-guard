package proxy

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/kfdm/promql-guard/config"

	"github.com/go-kit/log"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
)

type MockProxy struct {
	logger log.Logger
	mock   func(rw http.ResponseWriter, req *http.Request)
}

func NewMock(logger log.Logger, mock func(rw http.ResponseWriter, req *http.Request)) *MockProxy {
	return &MockProxy{logger: logger, mock: mock}
}

func (p *MockProxy) ProxyRequest(w http.ResponseWriter, req *http.Request, config *config.VirtualHost) {
	req.Form = nil
	req.PostForm = nil
	config.Prometheus.UpdateRequest(req)
	p.mock(w, req)
}

func Get(path string, q url.Values) (*http.Request, error) {
	// Build Reqeust
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

func Post(path string, q url.Values) (*http.Request, error) {
	data := q.Encode()
	req, err := http.NewRequest("POST", path, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))
	return req, nil
}

func ExpectedPromql(t *testing.T, value string, expected string) {
	expr, err := parser.ParseExpr(value)
	require.NoError(t, err)
	require.Equal(t, expected, expr.String())
}

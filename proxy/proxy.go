package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/kfdm/promql-guard/config"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// RequestProxy for upstream server
type RequestProxy interface {
	ProxyRequest(w http.ResponseWriter, req *http.Request, config *config.VirtualHost)
}

// Proxy model
type Proxy struct {
	logger log.Logger
}

// https://stackoverflow.com/a/53007606

type DebugTransport struct{}

func (DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequestOut(r, false)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))
	return http.DefaultTransport.RoundTrip(r)
}

func NewProxy(logger log.Logger) *Proxy {
	return &Proxy{logger: logger}
}

// ProxyRequest to upstream Prometheus server
func (p *Proxy) ProxyRequest(w http.ResponseWriter, req *http.Request, config *config.VirtualHost) {

	proxy := httputil.NewSingleHostReverseProxy(config.Prometheus.URL())
	config.Prometheus.UpdateRequest(req)

	level.Info(p.logger).Log("msg", "proxying request", "upstream", config.Prometheus.URL(), "query", req.URL.String(), "method", req.Method)

	// proxy.Transport = DebugTransport{}

	proxy.ServeHTTP(w, req)
}

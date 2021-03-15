package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kfdm/promql-guard/config"
)

// Proxy model
type Proxy struct {
	logger log.Logger
	config *config.VirtualHost
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

func NewProxy(config *config.VirtualHost, logger log.Logger) *Proxy {
	return &Proxy{logger: logger, config: config}
}

// ProxyRequest to upstream Prometheus server
func (p *Proxy) ProxyRequest(w http.ResponseWriter, req *http.Request) {

	proxy := httputil.NewSingleHostReverseProxy(p.config.Prometheus.URL())
	req.Host = p.config.Prometheus.Host()
	level.Info(p.logger).Log("msg", "proxying request", "upstream", p.config.Prometheus.URL(), "query", req.URL.String(), "method", req.Method)

	// proxy.Transport = DebugTransport{}

	proxy.ServeHTTP(w, req)
}

package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kfdm/promql-guard/config"
)

// Proxy model
type Proxy struct {
	Logger log.Logger
	Config config.VirtualHost
}

// ProxyRequest to upstream Prometheus server
func (p *Proxy) ProxyRequest(w http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(p.Config.Prometheus.Upstream.URL)
	level.Info(p.Logger).Log("msg", "proxying request", "upstream", p.Config.Prometheus.Upstream, "query", req.URL.String())
	proxy.ServeHTTP(w, req)
}

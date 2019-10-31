package proxy

import (
	"net/http"

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
	level.Debug(p.Logger).Log("proxy", req.URL.String(), "prometheus", p.Config.Prometheus.Upstream)
}

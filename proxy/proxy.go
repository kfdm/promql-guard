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
	Logger log.Logger
	Config config.VirtualHost
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

// ProxyRequest to upstream Prometheus server
func (p *Proxy) ProxyRequest(w http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(p.Config.Prometheus.Upstream.URL)
	req.Host = p.Config.Prometheus.Upstream.Host
	level.Info(p.Logger).Log("msg", "proxying request", "upstream", p.Config.Prometheus.Upstream, "query", req.URL.String())

	// proxy.Transport = DebugTransport{}

	proxy.ServeHTTP(w, req)
}

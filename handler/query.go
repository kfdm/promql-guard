package handler

import (
	"net/http"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/injectproxy"
	"github.com/kfdm/promql-guard/proxy"

	auth "github.com/abbot/go-http-auth"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/promql"
)

type Enforcer struct {
	config *config.Config
	logger log.Logger
	query  string
}

// NewEnforcer returns a Enforcer handler
func NewEnforcer(cfg *config.Config, logger log.Logger, query string) *Enforcer {
	return &Enforcer{config: cfg, logger: logger, query: query}
}

func (h *Enforcer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	htpasswd := auth.HtpasswdFileProvider(h.config.Htpasswd)
	authenticator := auth.NewBasicAuthenticator("Basic Realm", htpasswd)
	user := authenticator.CheckAuth(req)

	if user == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	virtualhost, err := h.config.Find(user)
	if err != nil {
		level.Error(h.logger).Log("msg", "Unable to find virtualhost", "user", user)
		http.Error(w, "No configuration for this host", http.StatusUnauthorized)
		return
	}

	expr, err := promql.ParseExpr(req.FormValue(h.query))
	if err != nil {
		http.Error(w, "Error parsing PromQL", 400)
		return
	}

	// Add our required labels
	level.Debug(h.logger).Log("msg", "Incoming expression", "expression", expr.String(), "user", virtualhost.Username)
	err = injectproxy.InjectMatchers(expr, virtualhost.Prometheus.Matchers)
	if err != nil {
		http.Error(w, "Error enforcing PromQL", 400)
		level.Error(h.logger).Log("msg", "Unable to find virtualhost", "host", req.Host)
		return
	}
	level.Debug(h.logger).Log("msg", "Outgoing expression", "expression", expr.String(), "user", virtualhost.Username)

	// Return updated query
	q := req.URL.Query()
	q.Set(h.query, expr.String())
	req.URL.RawQuery = q.Encode()

	var p = proxy.Proxy{
		Logger: h.logger,
		Config: *virtualhost,
	}
	p.ProxyRequest(w, req)
}

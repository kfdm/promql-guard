package handler

import (
	"net/http"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/injectproxy"
	"github.com/kfdm/promql-guard/proxy"
	"github.com/pkg/errors"

	auth "github.com/abbot/go-http-auth"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/promql"
)

type EnforcerHandler struct {
	config *config.Config
	logger log.Logger
	query  string
}

// NewEnforcer returns a Enforcer handler
func NewEnforcer(cfg *config.Config, logger log.Logger, query string) *EnforcerHandler {
	return &EnforcerHandler{config: cfg, logger: logger, query: query}
}

// BasicAuth enforces our autentication and returns the correct config
func (h *EnforcerHandler) BasicAuth(w http.ResponseWriter, req *http.Request) (*config.VirtualHost, error) {
	htpasswd := auth.HtpasswdFileProvider(h.config.Htpasswd)
	authenticator := auth.NewBasicAuthenticator("Basic Realm", htpasswd)
	user := authenticator.CheckAuth(req)

	if user == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, errors.New("Unauthorized")
	}

	virtualhost, err := h.config.Find(user)
	if err != nil {
		level.Error(h.logger).Log("msg", "Unable to find virtualhost", "user", user)
		http.Error(w, "No configuration for this host", http.StatusUnauthorized)
		return nil, err
	}

	return virtualhost, nil
}

// Error formats and logs our http errors
func (h *EnforcerHandler) Error(w http.ResponseWriter, code int, err error, msg string) {
	http.Error(w, msg, code)
	level.Error(h.logger).Log("msg", msg, "err", err.Error())
}

// ServeHTTP implements our required http.Handler interface
func (h *EnforcerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	virtualhost, err := h.BasicAuth(w, req)
	if err != nil {
		h.Error(w, http.StatusUnauthorized, err, "Unauthorised")
		return
	}

	expr, err := promql.ParseExpr(req.FormValue(h.query))
	if err != nil {
		h.Error(w, http.StatusBadRequest, err, "Error parsing PromQL")
		return
	}

	// Add our required labels
	level.Debug(h.logger).Log("msg", "Incoming expression", "expression", expr.String(), "user", virtualhost.Username)
	err = injectproxy.InjectMatchers(expr, virtualhost.Prometheus.Matchers)
	if err != nil {
		h.Error(w, http.StatusBadRequest, err, "Error enforcing PromQL")
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

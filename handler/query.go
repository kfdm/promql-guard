package handler

import (
	"net/http"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/injectproxy"
	"github.com/kfdm/promql-guard/proxy"

	auth "github.com/abbot/go-http-auth"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/prometheus/promql"
)

var (
	httpCnt = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promqlguard_http_requests_total",
			Help: "Total HTTP requests processed by promql-guard, excluding scrapes.",
		},
		[]string{"handler", "code", "method"},
	)
)

func enforce(query string, w http.ResponseWriter, req *http.Request, cfg *config.Config, logger log.Logger) {
	htpasswd := auth.HtpasswdFileProvider(cfg.Htpasswd)
	authenticator := auth.NewBasicAuthenticator("Basic Realm", htpasswd)
	user := authenticator.CheckAuth(req)

	if user == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	virtualhost, err := cfg.Find(user)
	if err != nil {
		level.Error(logger).Log("msg", "Unable to find virtualhost", "user", user)
		http.Error(w, "No configuration for this host", http.StatusUnauthorized)
		return
	}

	expr, err := promql.ParseExpr(req.FormValue(query))
	if err != nil {
		http.Error(w, "Error parsing PromQL", 400)
		return
	}

	// Add our required labels
	level.Debug(logger).Log("msg", "Incoming expression", "expression", expr.String(), "user", virtualhost.Username)
	err = injectproxy.InjectMatchers(expr, virtualhost.Prometheus.Matchers)
	if err != nil {
		http.Error(w, "Error enforcing PromQL", 400)
		level.Error(logger).Log("msg", "Unable to find virtualhost", "host", req.Host)
		return
	}
	level.Debug(logger).Log("msg", "Outgoing expression", "expression", expr.String(), "user", virtualhost.Username)

	// Return updated query
	q := req.URL.Query()
	q.Set(query, expr.String())
	req.URL.RawQuery = q.Encode()

	var p = proxy.Proxy{
		Logger: logger,
		Config: *virtualhost,
	}
	p.ProxyRequest(w, req)
}

// Query Wraped Prometheus Query
func Query(logger log.Logger, config *config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			enforce("query", w, req, config, logger)
		}),
	)
}

// Series Wraped Prometheus Query
func Series(logger log.Logger, config *config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "series"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			enforce("match[]", w, req, config, logger)
		}),
	)
}

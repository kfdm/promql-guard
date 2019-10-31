package handler

import (
	"net/http"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/injectproxy"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
)

var (
	httpCnt = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promqlguard_http_requests_total",
			Help: "Total HTTP requests processed by the Pushgateway, excluding scrapes.",
		},
		[]string{"handler", "code", "method"},
	)
)

func proxy(w http.ResponseWriter, req *http.Request, vhost *config.VirtualHost, logger log.Logger) {
	level.Debug(logger).Log("proxy", req.URL.String(), "prometheus", vhost.Prometheus.Upstream)
}

func enforce(query string, w http.ResponseWriter, req *http.Request, cfg *config.Config, logger log.Logger) {
	if req.Host == "" {
		http.Error(w, "Missing host header", 400)
		return
	}

	virtualhost, err := cfg.Find(req.Host)
	if err != nil {
		level.Error(logger).Log("msg", "Unable to find virtualhost", "host", req.Host)
		return
	}

	expr, err := promql.ParseExpr(req.FormValue(query))
	if err != nil {
		return
	}

	// Add our required labels
	level.Debug(logger).Log("msg", "Incoming expression", "expression", expr.String())
	err = injectproxy.SetRecursive(expr, []*labels.Matcher{{
		Name:  "key",
		Type:  labels.MatchEqual,
		Value: "value",
	}})
	level.Debug(logger).Log("msg", "Outgoing expression", "expression", expr.String())

	// Return updated query
	q := req.URL.Query()
	q.Set(query, expr.String())
	req.URL.RawQuery = q.Encode()

	proxy(w, req, virtualhost, logger)
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

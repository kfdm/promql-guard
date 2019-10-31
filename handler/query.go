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

// Query Wraped Prometheus Query
func Query(logger log.Logger, config config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Pull out Query
			expr, err := promql.ParseExpr(req.FormValue("query"))
			if err != nil {
				return
			}

			// Add our required labels
			level.Debug(logger).Log("msg", "Incoming query", "query", expr.String())
			err = injectproxy.SetRecursive(expr, []*labels.Matcher{{
				Name:  "key",
				Type:  labels.MatchEqual,
				Value: "value",
			}})
			level.Debug(logger).Log("msg", "Outgoing query", "query", expr.String())

			// Return updated query
			q := req.URL.Query()
			q.Set("query", expr.String())
			req.URL.RawQuery = q.Encode()
		}),
	)
}

// Series Wraped Prometheus Query
func Series(logger log.Logger, config config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "series"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Pull out Query
			expr, err := promql.ParseExpr(req.FormValue("match[]"))
			if err != nil {
				return
			}

			// Add our required labels
			level.Debug(logger).Log("msg", "Incoming query", "query", expr.String())
			err = injectproxy.SetRecursive(expr, []*labels.Matcher{{
				Name:  "key",
				Type:  labels.MatchEqual,
				Value: "value",
			}})
			level.Debug(logger).Log("msg", "Outgoing query", "query", expr.String())

			// Return updated query
			q := req.URL.Query()
			q.Set("match[]", expr.String())
			req.URL.RawQuery = q.Encode()
		}),
	)
}

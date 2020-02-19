package api

import (
	"net/http"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/handler"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

// Query Wraped Prometheus Query
func Query(logger log.Logger, config *config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.Enforce("query", w, req, config, logger)
		}),
	)
}

// QueryRange Wraped Prometheus Query
func QueryRange(logger log.Logger, config *config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query_range"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.Enforce("query", w, req, config, logger)
		}),
	)
}

// Series Wraped Prometheus Query
func Series(logger log.Logger, config *config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "series"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.Enforce("match[]", w, req, config, logger)
		}),
	)
}

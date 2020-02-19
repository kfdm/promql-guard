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
		handler.NewEnforcer(config, logger, "query"),
	)
}

// QueryRange Wraped Prometheus Query
func QueryRange(logger log.Logger, config *config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query_range"}),
		handler.NewEnforcer(config, logger, "query"),
	)
}

// Series Wraped Prometheus Query
func Series(logger log.Logger, config *config.Config) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "series"}),
		handler.NewEnforcer(config, logger, "match[]"),
	)
}

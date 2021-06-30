package api

import (
	"net/http"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/handler"
	"github.com/kfdm/promql-guard/proxy"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promqlguard_http_requests_total",
			Help: "Counter of HTTP requests.",
		},
		[]string{"handler", "code", "method"},
	)
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "promqlguard_http_request_duration_seconds",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: []float64{.1, .2, .4, 1, 3, 8, 20, 60, 120},
		},
		[]string{"handler"},
	)
	responseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "promqlguard_http_response_size_bytes",
			Help:    "Histogram of response size for HTTP requests.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"handler"},
	)
)

type API struct {
	config *config.Config
	logger log.Logger
	proxy  proxy.RequestProxy
}

// NewAPI instance
func NewAPI(config *config.Config, logger log.Logger, proxy proxy.RequestProxy) *API {
	return &API{
		config: config,
		logger: logger,
		proxy:  proxy,
	}
}

// Wrap handler with metrics
func Wrap(handlerName string, handler http.Handler) http.HandlerFunc {
	return promhttp.InstrumentHandlerCounter(
		requestCounter.MustCurryWith(prometheus.Labels{"handler": handlerName}),
		promhttp.InstrumentHandlerDuration(
			requestDuration.MustCurryWith(prometheus.Labels{"handler": handlerName}),
			promhttp.InstrumentHandlerResponseSize(
				responseSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
				handler,
			),
		),
	)
}

// Query Wraped Prometheus Query
func (api *API) Query() http.Handler {
	return Wrap("query", handler.NewEnforcer(api.config, api.logger, "query", api.proxy))
}

// QueryRange Wraped Prometheus Query
func (api *API) QueryRange() http.Handler {
	return Wrap("query_range", handler.NewEnforcer(api.config, api.logger, "query", api.proxy))
}

// Series Wraped Prometheus Query
func (api *API) Series() http.Handler {
	return Wrap("series", handler.NewEnforcer(api.config, api.logger, "match[]", api.proxy))
}

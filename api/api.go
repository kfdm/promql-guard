package api

import (
	"net/http"

	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/handler"
	"github.com/kfdm/promql-guard/proxy"

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

// Query Wraped Prometheus Query
func (api *API) Query() http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query"}),
		handler.NewEnforcer(api.config, api.logger, "query", api.proxy),
	)
}

// QueryRange Wraped Prometheus Query
func (api *API) QueryRange() http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query_range"}),
		handler.NewEnforcer(api.config, api.logger, "query", api.proxy),
	)
}

// Series Wraped Prometheus Query
func (api *API) Series() http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "series"}),
		handler.NewEnforcer(api.config, api.logger, "match[]", api.proxy),
	)
}

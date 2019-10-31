package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kfdm/promql-guard/injectproxy"
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
func Query() http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "query"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			expr, err := promql.ParseExpr(req.FormValue("query"))
			if err != nil {
				return
			}
			fmt.Println(expr.String())

			err = injectproxy.SetRecursive(expr, []*labels.Matcher{{
				Name:  "key",
				Type:  labels.MatchEqual,
				Value: "value",
			}})

			fmt.Println(expr.String())

			q := req.URL.Query()
			q.Set("query", expr.String())
			req.URL.RawQuery = q.Encode()

			io.WriteString(w, "OK")
		}),
	)
}

// Series Wraped Prometheus Query
func Series() http.Handler {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": "series"}),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "OK")
		}),
	)
}

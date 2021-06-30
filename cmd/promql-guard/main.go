package main

import (
	"net"
	"net/http"
	"os"

	"github.com/kfdm/promql-guard/api"
	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/proxy"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/route"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	prometheus.MustRegister(version.NewCollector("promqlguard"))
}

func main() {
	os.Exit(run())
}

func run() int {
	var (
		promlogConfig = promlog.Config{}
		configFile    = kingpin.Flag("config.file", "PromqlGuard configuration file name.").Default("guard.yml").String()
		// Registered at https://github.com/prometheus/prometheus/wiki/Default-port-allocations
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for the web interface and API.").Default(":9218").String()
	)

	flag.AddFlags(kingpin.CommandLine, &promlogConfig)

	kingpin.Version(version.Print("promql-guard"))
	kingpin.CommandLine.GetFlag("help").Short('h')
	kingpin.Parse()

	logger := promlog.New(&promlogConfig)
	level.Info(logger).Log("build_context", version.BuildContext())

	// Load Configuration
	level.Info(logger).Log("msg", "Config", "config", *configFile)
	config, err := config.LoadFile(*configFile)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	// Build Routing Tree
	router := route.New()
	router.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Add in our API Endpoints
	proxy := proxy.NewProxy(logger)
	api := api.NewAPI(config, logger, proxy)
	router.Get("/api/v1/query", api.Query().ServeHTTP)
	router.Post("/api/v1/query", api.Query().ServeHTTP)
	router.Get("/api/v1/query_range", api.QueryRange().ServeHTTP)
	router.Post("/api/v1/query_range", api.QueryRange().ServeHTTP)
	router.Get("/api/v1/series", api.Series().ServeHTTP)
	router.Post("/api/v1/series", api.Series().ServeHTTP)

	// Launch server
	level.Info(logger).Log("listen_address", *listenAddress)
	l, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	err = (&http.Server{Addr: *listenAddress, Handler: router}).Serve(l)
	level.Error(logger).Log("msg", "HTTP server stopped", "err", err)

	return 0
}

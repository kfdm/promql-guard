package main

import (
	"net"
	"net/http"
	"os"
	"path"

	"github.com/go-kit/kit/log/level"
	"github.com/julienschmidt/httprouter"
	"github.com/kfdm/promql-guard/config"
	"github.com/kfdm/promql-guard/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
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
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for the web interface and API.").Default(":9093").String()
	)

	flag.AddFlags(kingpin.CommandLine, &promlogConfig)

	kingpin.Version(version.Print("alertmanager"))
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
	r := httprouter.New()
	r.Handler("GET", path.Join("/metrics"), promhttp.Handler())
	r.Handler("GET", path.Join("/api/v1/query"), handler.Query(logger, config))
	r.Handler("GET", path.Join("/api/v1/query_range"), handler.Query(logger, config))
	r.Handler("GET", path.Join("/api/v1/series"), handler.Series(logger, config))

	// Launch server
	level.Info(logger).Log("listen_address", *listenAddress)
	l, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	err = (&http.Server{Addr: *listenAddress, Handler: r}).Serve(l)
	level.Error(logger).Log("msg", "HTTP server stopped", "err", err)

	return 0
}

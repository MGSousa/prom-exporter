package main

import (
	"flag"
	"fmt"
	"main/exporter"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var (
	internalVersion string

	listenAddress        = flag.String("listen-address", ":19100", "Address to listen on for telemetry.")
	metricsPath          = flag.String("telemetry-path", "/metrics", "Base path under which to expose metrics.")
	serviceName          = flag.String("service-name", "", "Service name to reference")
	serviceUri           = flag.String("service-uri", "http://localhost:5066", "HTTP address of the service.")
	serviceMetricsPath   = flag.String("service-metrics-path", "metrics", "Service path to scrape metrics from.")
	serviceVersionScrape = flag.Bool("service-version-scrape", false, "Enable whether the service will be internally scraped for fetching remote build version or not")
	debugLevel           = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	name := *serviceName
	if strings.Trim(name, " ") == "" {
		log.Fatalln("Service name not known! Specify by -service-name SERVICE")
	}

	if *debugLevel {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Check if target is reachable...")
	if *serviceVersionScrape {
		internalVersion = checkEndpoint(*serviceUri)
	} else {
		checkEndpoint(*serviceUri)
	}
	log.Info("Target endpoint is reachable")

	registry := prometheus.NewRegistry()

	// register current Exporter version metrics
	versionMetric := version.NewCollector(name)
	registry.MustRegister(versionMetric)

	// register remote service metrics
	exporter := exporter.NewCollector(name, fmt.Sprintf("%s/%s", *serviceUri, *serviceMetricsPath), internalVersion)
	registry.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			DisableCompression: false,
			ErrorHandling:      promhttp.ContinueOnError,
		}),
	)
	log.Println("Starting server....")

	srv := &http.Server{
		Addr:              *listenAddress,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("http server quit with error: %v", err)
	}
}

// TODO: improve this method
func checkEndpoint(endpoint string) string {
	stopCh := make(chan bool)
	t := time.NewTicker(2 * time.Second)

	stats := &exporter.HttpStats{}

discovery:
	for {
		select {
		case <-t.C:
			if stats = exporter.FetchStats(endpoint); stats != nil {
				break discovery
			}
			log.Errorln("base endpoint not available, retrying in 2s")
			continue

		case <-stopCh:
			os.Exit(0)
		}
	}
	t.Stop()

	return stats.Version
}

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
	internalVersion, serviceEndpoint string

	// internal settings
	listenAddress = flag.String("listen-address", ":19100", "Address to listen on to be scraped.")
	metricsPath   = flag.String("telemetry-path", "/metrics", "Base path under which to expose metrics.")
	debugLevel    = flag.Bool("debug", false, "Enable debug mode.")

	// remote service related flags
	serviceName          = flag.String("service-name", "", "Remote service name to reference.")
	serviceProtocol      = flag.String("service-protocol", "http", "HTTP Schema of the remote service (http or https).")
	serviceUri           = flag.String("service-uri", "", "Endpoint address of the remote service.")
	servicePort          = flag.String("service-port", "80", "HTTP Port of the remote service.")
	serviceMetricsPath   = flag.String("service-metrics-path", "metrics", "Service path to scrape metrics from.")
	serviceVersionScrape = flag.Bool("service-version-scrape", false, "Enable whether the service will be internally scraped for fetching remote build version or not.")
)

func main() {
	flag.Parse()

	name := *serviceName
	if strings.Trim(name, " ") == "" {
		log.Fatalln("Service name not set! Specify by -service-name SERVICE")
	}

	// build service endpoint
	// if not set try to fetch from env
	serviceEndpoint = fmt.Sprintf("%s://%s:%s", *serviceProtocol, *serviceUri, *servicePort)
	if *serviceUri == "" {
		host, exists := os.LookupEnv("SERVICE_ENDPOINT")
		if exists && host != "" {
			serviceEndpoint = fmt.Sprintf("%s://%s:%s", *serviceProtocol, host, *servicePort)
		} else {
			log.Fatalln("Service URI is not set! Specify either a '-service-uri' flag OR an environment variable 'SERVICE_ENDPOINT'")
		}
	}

	// enable debug mode if required
	if *debugLevel {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Check if target is reachable...")
	if *serviceVersionScrape {
		internalVersion = checkEndpoint(serviceEndpoint)
	} else {
		checkEndpoint(serviceEndpoint)
	}
	log.Info("Target endpoint is reachable")

	registry := prometheus.NewRegistry()

	// register current Exporter version metrics
	versionMetric := version.NewCollector(name)
	registry.MustRegister(versionMetric)

	// register remote service metrics
	exporter := exporter.NewCollector(name, fmt.Sprintf("%s/%s", serviceEndpoint, *serviceMetricsPath), internalVersion)
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
		ReadHeaderTimeout: 30 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("http server quit with error: %v", err)
	}
}

// check if service endpoint is reachable
// increases duration by 5secs for the next ticket if unavailable
func checkEndpoint(endpoint string) string {
	var duration = 5
	stopCh := make(chan bool)
	t := time.NewTicker(time.Duration(duration) * time.Second)

	stats := &exporter.HttpStats{}

discovery:
	for {
		select {
		case <-t.C:
			if stats = exporter.FetchStats(endpoint); stats != nil {
				break discovery
			}
			log.Errorf("service endpoint not available, retrying in %ds", duration)

			t = time.NewTicker(time.Duration(duration) * time.Second)
			duration = duration + 5
			continue

		case <-stopCh:
			os.Exit(0)
		}
	}
	t.Stop()

	return stats.Version
}

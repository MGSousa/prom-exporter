package exporter

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type (
	Data         []string
	CustomLabels map[string]prometheus.Labels

	Exporter struct {
		uri     string
		name    string
		version string
		metrics Metrics

		// extract parsed data
		data Data
	}

	Metrics []struct {
		help  *prometheus.Desc
		value float64
		vtype prometheus.ValueType
	}
)

// NewCollector instantiates the collector
func NewCollector(name, uri, version string) *Exporter {
	return &Exporter{
		name:    name,
		uri:     uri,
		version: version,
		metrics: Metrics{},
	}
}

// Describe describes metrics by setting HELP lines
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range e.metrics {
		ch <- metric.help
	}
}

// Collect collects defined metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.process(ch)

	for _, metric := range e.metrics {
		ch <- prometheus.MustNewConstMetric(metric.help, metric.vtype, metric.value)
	}
}

// process all extracted metrics from remote service
func (e *Exporter) process(ch chan<- prometheus.Metric) {
	var (
		mType prometheus.ValueType
	)

	defer func() {
		e.data = nil
	}()

	body, err := FetchMetrics(e.uri)
	if err != nil {
		// set target down on error
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(fmt.Sprintf("%s_up", e.name), SERVICE_UP_HELP, nil, nil), prometheus.GaugeValue, float64(0))
		log.Debugf("Failed getting metrics endpoint of target: %s ", err.Error())
		return
	}

	log.Debugln("Service scrapped up")
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(fmt.Sprintf("%s_up", e.name), SERVICE_UP_HELP, nil, nil), prometheus.GaugeValue, float64(1))

	// export custom version if there isnt a default metric
	if e.version != "" {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(fmt.Sprintf("%s_version_info", e.name), fmt.Sprintf("%s service info", e.name), nil,
				prometheus.Labels{"version": e.version}), prometheus.GaugeValue, float64(1))
	}

	// parse response body
	e.parse(body)

	// extract strings and process them to be exported
	less := e.extractStrings(ch)

	// sort parsed data
	sort.SliceStable(e.data, func(i, j int) bool { return e.data[i] < e.data[j] })

	// iterate over parsed metrics
	i := 0
	stats := make(Metrics, len(e.data)-less)
	for k, v := range e.data {
		if v != "" {
			m := split(v, RAW_METRIC_DELIM)
			value, err := strconv.ParseFloat(m[2], 64)
			if err != nil {
				log.Debugln(err)
				value = 0
			}

			// TODO: custom filter by a map file
			if len(m) > 3 {
				mType = prometheus.UntypedValue
			} else {
				mType = prometheus.GaugeValue
			}

			// build && extract Full-Qualified Name and Labels
			fqname, labels := e.build(m, k)

			// supply metrics
			stats[i].help = prometheus.NewDesc(
				fqname, "", nil, labels)
			stats[i].value = float64(value)
			stats[i].vtype = mType

			i++
		}
	}

	e.metrics = stats
}

// check if custom string metrics are found
// then extract them as a valid metric to be in Prometheus format
func (e *Exporter) extractStrings(ch chan<- prometheus.Metric) int {
	var (
		metric  []string
		i, less int
	)
	untypedMetrics := make(CustomLabels, 0)
	for _, v := range e.data {
		metric = split(v, RAW_METRIC_DELIM)
		if _, err := strconv.ParseFloat(metric[2], 64); err != nil {
			if len(metric) > 3 {
				if untypedMetrics[metric[0]] == nil {
					untypedMetrics[metric[0]] = make(prometheus.Labels)
				}
				untypedMetrics[metric[0]][metric[1]] = metric[2]

				// unset old metric
				e.data[i] = ""
				less++
			}
		}
		i++
	}

	for m, labels := range untypedMetrics {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(fmt.Sprintf("%s_%s", e.name, m), "", nil,
				labels), prometheus.UntypedValue, float64(1))
	}
	return less
}

package exporter

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type (
	Exporter struct {
		uri     string
		name    string
		version string
		metrics Metrics

		data []string
	}

	Metrics []struct {
		help  *prometheus.Desc
		value float64
		vtype prometheus.ValueType
	}
)

func NewCollector(name, uri, version string) *Exporter {
	return &Exporter{
		name:    name,
		uri:     uri,
		version: version,
		metrics: Metrics{},
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range e.metrics {
		ch <- metric.help
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.process(ch)

	for _, metric := range e.metrics {
		ch <- prometheus.MustNewConstMetric(metric.help, metric.vtype, metric.value)
	}
}

var (
	// s          = []string{"processor_rate_limit_1::dropped::0", "registrar_states::cleanup::0", "registrar_states::update::0", "registrar_writes::success::0", "system_cpu::cores::4", "system_load::1::0.71", "system_load::15::0.47", "system_load::5::0.53", "system_load_norm::1::0.1775", "system_load_norm::15::0.1175", "system_load_norm::5::0.1325"}
	vl, future []string

	fqname string
	labels prometheus.Labels
	mType  prometheus.ValueType
)

func (e *Exporter) process(ch chan<- prometheus.Metric) {
	body, err := FetchMetrics(e.uri)
	if err != nil {
		// set target down on error
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(fmt.Sprintf("%s_up", e.name), "Target up", nil, nil), prometheus.GaugeValue, float64(0))
		log.Debugf("Failed getting metrics endpoint of target: %s ", err.Error())
		return
	}

	log.Debugln("Service scrapped up")
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(fmt.Sprintf("%s_up", e.name), "Target Up", nil, nil), prometheus.GaugeValue, float64(1))

	// inject service version is exists
	if e.version != "" {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(fmt.Sprintf("%s_version_info", e.name), fmt.Sprintf("%s service info", e.name), nil,
				prometheus.Labels{"version": e.version}), prometheus.GaugeValue, float64(1))
	}

	// parse response body
	p := Parser{data: make([]string, 0)}
	p.parse(body)
	e.data = p.data

	sort.Slice(e.data, func(i, j int) bool { return e.data[i] < e.data[j] })

	stats := make(Metrics, len(e.data))
	for k, v := range e.data {
		vl = split(v, "::")

		i, err := strconv.ParseFloat(vl[2], 64)
		if err != nil {
			// a string was found
			log.Debugln(err)
			i = 0
		}

		// check last filed for the specific ValueType
		if contains(vl[1], "total") {
			mType = prometheus.CounterValue
		} else {
			if len(vl) > 3 {
				mType = prometheus.UntypedValue
			} else {
				mType = prometheus.GaugeValue
			}
		}

		// build Full-Qualified Name
		e.build(k)

		// setup metrics
		stats[k].help = prometheus.NewDesc(
			fqname, "", nil, labels)
		stats[k].value = float64(i)
		stats[k].vtype = mType
	}

	e.metrics = stats
}

func (e *Exporter) build(k int) {
	labels = nil
	fp := split(vl[0], "_")

	if len(fp) == 2 {

		if k < (len(e.data) - 1) {
			future = split(e.data[k+1], "::")
			if vl[0] == future[0] {

				fqname = prometheus.BuildFQName(
					e.name, fmt.Sprintf("%s", join(fp[0:len(fp)-1], "_")), fp[len(fp)-1])
				labels = prometheus.Labels{"metric": vl[1]}
			} else {
				future = split(e.data[k-1], "::")
				if vl[0] == future[0] {
					fqname = prometheus.BuildFQName(
						e.name, fmt.Sprintf("%s", join(fp[0:len(fp)-1], "_")), fp[len(fp)-1])
					labels = prometheus.Labels{"metric": vl[1]}
				} else {
					fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s_%s", fp[0], fp[1]), vl[1])
				}
			}
		} else if k == (len(e.data) - 1) {
			future = split(e.data[k-1], "::")
			if vl[0] == future[0] {
				fqname = prometheus.BuildFQName(
					e.name, fmt.Sprintf("%s", join(fp[0:len(fp)-1], "_")), fp[len(fp)-1])
				labels = prometheus.Labels{"metric": vl[1]}
			} else {
				fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s_%s", fp[0], fp[1]), vl[1])
			}
		} else {
			fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s_%s", fp[0], fp[1]), vl[1])
		}

	} else if len(fp) >= 3 {

		if k < (len(e.data) - 1) {
			future = split(e.data[k+1], "::")
			if vl[0] == future[0] {

				fqname = prometheus.BuildFQName(
					e.name, fmt.Sprintf("%s", join(fp[0:len(fp)-1], "_")), fp[len(fp)-1])
				labels = prometheus.Labels{"metric": vl[1]}
			} else {

				if k == 0 {
					fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s_%s", fp[0], fp[1]), vl[1])
				} else {
					future = split(e.data[k-1], "::")
					if vl[0] == future[0] {
						fqname = prometheus.BuildFQName(
							e.name, fmt.Sprintf("%s", join(fp[0:len(fp)-1], "_")), fp[len(fp)-1])
						labels = prometheus.Labels{"metric": vl[1]}
					} else {
						fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s_%s", fp[0], fp[1]), vl[1])
					}
				}
			}
		} else if k == (len(e.data) - 1) {
			future = split(e.data[k-1], "::")
			if vl[0] == future[0] {
				labels = prometheus.Labels{"metric": vl[1]}
			}
		}
		fqname = prometheus.BuildFQName(
			e.name, fmt.Sprintf("%s", join(fp[0:len(fp)-1], "_")), fp[len(fp)-1])
	}
}

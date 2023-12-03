package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	fqname string
	labels prometheus.Labels = nil
)

func (e *Exporter) build(m []string, k int) (string, prometheus.Labels) {
	metricSlices := split(m[0], "_")

	if len(metricSlices) == 1 {
		fqname, labels = e.buildShortMetrics(m, k, metricSlices[0])
	} else if len(metricSlices) >= 2 {
		fqname, labels = e.buildLongMetrics(m, metricSlices, k)
	}

	return fqname, labels
}

func (e *Exporter) buildShortMetrics(m []string, k int, part string) (string, prometheus.Labels) {
	labels = nil

	if k < (len(e.data) - 1) {
		future := split(e.data[k+1], RAW_METRIC_DELIM)
		if m[0] == future[0] {
			fqname = prometheus.BuildFQName(e.name, "", part)
			labels = prometheus.Labels{"metric": m[1]}

		} else {
			if k == 0 {
				fqname = prometheus.BuildFQName(e.name, part, m[1])

			} else {
				future = split(e.data[k-1], RAW_METRIC_DELIM)
				if m[0] == future[0] {
					fqname = prometheus.BuildFQName(e.name, "", part)
					labels = prometheus.Labels{"metric": m[1]}

				} else {
					fqname = prometheus.BuildFQName(e.name, part, m[1])
				}
			}
		}
	} else if k == (len(e.data) - 1) {
		future := split(e.data[k-1], RAW_METRIC_DELIM)
		if m[0] == future[0] {
			fqname = prometheus.BuildFQName(e.name, "", part)
			labels = prometheus.Labels{"metric": m[1]}

		} else {
			fqname = prometheus.BuildFQName(e.name, part, m[1])
		}
	} else {
		fqname = prometheus.BuildFQName(e.name, part, m[1])
	}

	return fqname, labels
}

func (e *Exporter) buildLongMetrics(m, slices []string, k int) (string, prometheus.Labels) {
	labels = nil

	if k < (len(e.data) - 1) {
		future := split(e.data[k+1], RAW_METRIC_DELIM)
		if m[0] == future[0] {
			fqname = prometheus.BuildFQName(
				e.name, fmt.Sprintf("%s", join(slices[0:len(slices)-1], "_")), slices[len(slices)-1])
			labels = prometheus.Labels{"metric": m[1]}

		} else {
			if k == 0 {
				fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s", join(slices[0:], "_")), m[1])

			} else {
				future := split(e.data[k-1], RAW_METRIC_DELIM)
				if m[0] == future[0] {
					fqname = prometheus.BuildFQName(
						e.name, fmt.Sprintf("%s", join(slices[0:len(slices)-1], "_")), slices[len(slices)-1])
					labels = prometheus.Labels{"metric": m[1]}

				} else {
					fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s", join(slices[0:], "_")), m[1])
				}
			}
		}
	} else if k == (len(e.data) - 1) {
		future := split(e.data[k-1], RAW_METRIC_DELIM)
		if m[0] == future[0] {
			fqname = prometheus.BuildFQName(
				e.name, fmt.Sprintf("%s", join(slices[0:len(slices)-1], "_")), slices[len(slices)-1])
			labels = prometheus.Labels{"metric": m[1]}

		} else {
			fqname = prometheus.BuildFQName(e.name, fmt.Sprintf("%s", join(slices[0:], "_")), m[1])
		}
	}

	return fqname, labels
}

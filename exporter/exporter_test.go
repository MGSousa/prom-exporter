package exporter

import (
	"sort"
	"testing"

	"golang.org/x/exp/slices"
)

var (
	m    []string
	data = []string{"processor_rate_limit_1::dropped::0", "registrar_states::cleanup::0", "registrar_states::update::0", "registrar_writes::success::0", "system_cpu::cores::4", "system_load::1::0.71", "system_load::15::0.47", "system_load::5::0.53", "system_load_norm::1::0.1775", "system_load_norm::15::0.1175", "system_load_norm::5::0.1325"}
)

func (d Data) Len() int           { return len(d) }
func (d Data) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Data) Less(i, j int) bool { return d[i] < d[j] }

func BenchmarkBuildSort(t *testing.B) {
	type fields struct {
		uri     string
		name    string
		version string
		metrics Metrics
		data    []string
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				uri:     "http://localhost:5066/stats",
				name:    "test",
				metrics: Metrics{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			e := &Exporter{
				uri:     tt.fields.uri,
				name:    tt.fields.name,
				version: tt.fields.version,
				metrics: tt.fields.metrics,
				data:    data,
			}
			sort.Slice(e.data, func(i, j int) bool { return e.data[i] < e.data[j] })
			for k, v := range e.data {
				m = split(v, "::")
				e.build(m, k)
			}
		})
	}
}

func BenchmarkBuildSortStable(t *testing.B) {
	type fields struct {
		uri     string
		name    string
		version string
		metrics Metrics
		data    []string
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				uri:     "http://localhost:5066/stats",
				name:    "test",
				metrics: Metrics{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			e := &Exporter{
				uri:     tt.fields.uri,
				name:    tt.fields.name,
				version: tt.fields.version,
				metrics: tt.fields.metrics,
				data:    data,
			}
			sort.SliceStable(e.data, func(i, j int) bool { return e.data[i] < e.data[j] })
			for k, v := range e.data {
				m = split(v, "::")
				e.build(m, k)
			}
		})
	}
}

func BenchmarkBuildXSlices(t *testing.B) {
	type fields struct {
		uri     string
		name    string
		version string
		metrics Metrics
		data    []string
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				uri:     "http://localhost:5066/stats",
				name:    "test",
				metrics: Metrics{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			e := &Exporter{
				uri:     tt.fields.uri,
				name:    tt.fields.name,
				version: tt.fields.version,
				metrics: tt.fields.metrics,
				data:    data,
			}

			slices.Sort(e.data)
			for k, v := range e.data {
				m = split(v, "::")
				e.build(m, k)
			}
		})
	}
}

func BenchmarkBuildSortDefault(t *testing.B) {
	type fields struct {
		uri     string
		name    string
		version string
		metrics Metrics
		data    Data
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				uri:     "http://localhost:5066/stats",
				name:    "test",
				metrics: Metrics{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			e := &Exporter{
				uri:     tt.fields.uri,
				name:    tt.fields.name,
				version: tt.fields.version,
				metrics: tt.fields.metrics,
				data:    data,
			}

			sort.Sort(e.data)
			for k, v := range e.data {
				m = split(v, "::")
				e.build(m, k)
			}
		})
	}
}

package exporter

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewCollector(t *testing.T) {
	type args struct {
		name string
		uri  string
	}
	tests := []struct {
		name string
		args args
		want *Exporter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCollector(tt.args.name, tt.args.uri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExporter_Collect(t *testing.T) {
	type fields struct {
		uri     string
		name    string
		metrics Metrics
	}
	type args struct {
		ch chan<- prometheus.Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				uri:     tt.fields.uri,
				name:    tt.fields.name,
				metrics: tt.fields.metrics,
			}
			e.Collect(tt.args.ch)
		})
	}
}

func TestExporter_process(t *testing.T) {
	type fields struct {
		uri     string
		name    string
		metrics Metrics
	}
	type args struct {
		ch chan<- prometheus.Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			// TODO: Add test cases.
			name: "test",
			fields: fields{
				uri:     "http://localhost:5066/stats",
				name:    "test",
				metrics: Metrics{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				uri:     tt.fields.uri,
				name:    tt.fields.name,
				metrics: tt.fields.metrics,
			}
			e.process(tt.args.ch)
		})
	}
}

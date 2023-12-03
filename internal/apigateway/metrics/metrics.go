package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	StatusError = "error"
	StatusOk    = "ok"
)

func mustRegister(collectors ...prometheus.Collector) {
	prometheus.DefaultRegisterer.MustRegister(collectors...)
}

func newHistogramVec(name, help string, buckets []float64, labelValues ...string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "blockchain_crawler_api",
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labelValues,
	)
}

func newCounterVec(name, help string, labelValues ...string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "blockchain_crawler_api",
			Name:      name,
			Help:      help,
		},
		labelValues,
	)
}

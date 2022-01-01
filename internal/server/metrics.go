package server

import (
	"log"

	"github.com/AlfredDobradi/ledgerlog/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	subsystem = "http"

	EndpointCalls         string = "EndpointCalls"
	EndpointErrors        string = "EndpointErrors"
	EndpointCallDurations string = "EndpointCallDurations"
)

type mx struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

var gatherers mx

func init() {
	gatherers = mx{
		counters: map[string]*prometheus.CounterVec{
			EndpointCalls: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: metrics.Namespace,
				Subsystem: subsystem,
				Name:      "endpoint_calls",
				Help:      "Number of endpoint requests received",
			}, []string{"endpoint"}),
			EndpointErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: metrics.Namespace,
				Subsystem: subsystem,
				Name:      "endpoint_errors",
				Help:      "Number of endpoint requests resulted in an error",
			}, []string{"endpoint", "statusCode"}),
		},
		histograms: map[string]*prometheus.HistogramVec{
			EndpointCallDurations: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: metrics.Namespace,
				Subsystem: subsystem,
				Name:      "endpoint_call_durations",
				Help:      "Durations of the endpoint requests",
			}, []string{"endpoint"}),
		},
	}

	for name, collector := range gatherers.counters {
		log.Printf("Registering counter %s...", name)
		prometheus.MustRegister(collector)
	}

	for name, collector := range gatherers.histograms {
		log.Printf("Registering histogram %s...", name)
		prometheus.MustRegister(collector)
	}
}

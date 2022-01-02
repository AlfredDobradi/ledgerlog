package server

import (
	"fmt"
	"log"

	"github.com/AlfredDobradi/ledgerlog/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	subsystem = "http"

	metricEndpointCalls                string = "endpoint_calls"
	metricEndpointErrors               string = "endpoint_errors"
	metricEndpointCallDurations        string = "endpoint_call_durations"
	metricWebsocketConnectionsTotal    string = "websocket_connections_total"
	metricWebsocketConnectionsCurrent  string = "websocket_connections_current"
	metricWebsocketConnectionDurations string = "websocket_connection_durations"
	metricWebsocketErrors              string = "websocket_errors"

	labelEndpoint     string = "endpoint"
	labelStatusCode   string = "status_code"
	labelConnectionID string = "connection_id"
)

var (
	metricsToRegister []*metrics.Opts = []*metrics.Opts{
		{
			Kind:      metrics.KindCounter,
			Namespace: metrics.Namespace,
			Subsystem: subsystem,
			Name:      metricEndpointCalls,
			Help:      "Number of endpoint requests received",
			Labels:    []string{labelEndpoint},
		},
		{
			Kind:      metrics.KindCounter,
			Namespace: metrics.Namespace,
			Subsystem: subsystem,
			Name:      metricEndpointErrors,
			Help:      "Number of endpoint requests resulted in an error",
			Labels:    []string{labelEndpoint, labelStatusCode},
		},
		{
			Kind:      metrics.KindHistogram,
			Namespace: metrics.Namespace,
			Subsystem: subsystem,
			Name:      metricEndpointCallDurations,
			Help:      "Durations of the endpoint requests",
			Labels:    []string{labelEndpoint},
		},
		{
			Kind:      metrics.KindCounter,
			Namespace: metrics.Namespace,
			Subsystem: subsystem,
			Name:      metricWebsocketConnectionsTotal,
			Help:      "Total number of websocket connections established",
			Labels:    []string{labelEndpoint},
		},
		{
			Kind:      metrics.KindGauge,
			Namespace: metrics.Namespace,
			Subsystem: subsystem,
			Name:      metricWebsocketConnectionsCurrent,
			Help:      "Current number of websocket connections established",
			Labels:    []string{labelEndpoint, labelConnectionID},
		},
		{
			Kind:      metrics.KindCounter,
			Namespace: metrics.Namespace,
			Subsystem: subsystem,
			Name:      metricWebsocketErrors,
			Help:      "Number of errors occured during websocket connections",
			Labels:    []string{labelEndpoint, labelStatusCode},
		},
		{
			Kind:      metrics.KindHistogram,
			Namespace: metrics.Namespace,
			Subsystem: subsystem,
			Name:      metricWebsocketConnectionDurations,
			Help:      "Durations of websocket connections",
			Labels:    []string{labelEndpoint},
		},
	}
)

type mx struct {
	counters   map[string]*prometheus.CounterVec
	gauges     map[string]*prometheus.GaugeVec
	histograms map[string]*prometheus.HistogramVec
}

var gatherers mx

func init() {
	gatherers = mx{
		counters:   make(map[string]*prometheus.CounterVec),
		gauges:     make(map[string]*prometheus.GaugeVec),
		histograms: make(map[string]*prometheus.HistogramVec),
	}

	for _, gathererOpts := range metricsToRegister {
		var err error
		switch gathererOpts.Kind {
		case metrics.KindCounter:
			counter := prometheus.NewCounterVec(gathererOpts.AsCounterOpts(), gathererOpts.Labels)
			gatherers.counters[gathererOpts.Name] = counter
			err = prometheus.Register(counter)
		case metrics.KindHistogram:
			histogram := prometheus.NewHistogramVec(gathererOpts.AsHistogramOpts(), gathererOpts.Labels)
			gatherers.histograms[gathererOpts.Name] = histogram
			err = prometheus.Register(histogram)
		case metrics.KindGauge:
			gauge := prometheus.NewGaugeVec(gathererOpts.AsGaugeOpts(), gathererOpts.Labels)
			gatherers.gauges[gathererOpts.Name] = gauge
			err = prometheus.Register(gauge)
		default:
			err = fmt.Errorf("Metric %s not registered. Invalid kind %d", gathererOpts.Name, gathererOpts.Kind)
		}
		if err != nil {
			log.Printf("Error registering %s metric: %v", gathererOpts.Name, err)
		}
	}
}

func gatherEndpointError(url string, code int) {
	errCode := fmt.Sprintf("%d", code)
	gatherers.counters[metricEndpointErrors].WithLabelValues(url, errCode).Inc()
}

func gatherWebsocketError(url string, code int) {
	errCode := fmt.Sprintf("%d", code)
	gatherers.counters[metricEndpointErrors].WithLabelValues(url, errCode).Inc()
}

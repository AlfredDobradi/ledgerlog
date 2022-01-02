package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	Namespace = "ledger"
)

type OptKind byte

const (
	KindCounter OptKind = iota
	KindHistogram
	KindGauge
)

type Opts struct {
	Kind      OptKind
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Buckets   []string
	Labels    []string
}

func (o Opts) AsCounterOpts() prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Namespace: o.Namespace,
		Subsystem: o.Subsystem,
		Name:      o.Name,
		Help:      o.Help,
	}
}

func (o Opts) AsGaugeOpts() prometheus.GaugeOpts {
	return prometheus.GaugeOpts{
		Namespace: o.Namespace,
		Subsystem: o.Subsystem,
		Name:      o.Name,
		Help:      o.Help,
	}
}

func (o Opts) AsHistogramOpts() prometheus.HistogramOpts {
	return prometheus.HistogramOpts{
		Namespace: o.Namespace,
		Subsystem: o.Subsystem,
		Name:      o.Name,
		Help:      o.Help,
	}
}

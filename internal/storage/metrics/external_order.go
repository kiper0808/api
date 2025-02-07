package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	addExternalOrderCounter *prometheus.CounterVec
)

const (
	addExternalOrderCounterName = "add_external_order_total"
)

func registerExternalOrderMetrics() (err error) {
	addExternalOrderCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      addExternalOrderCounterName,
		Help:      "add external order counter",
	}, []string{"status"})
	if err = prometheus.Register(addExternalOrderCounter); err != nil {
		return
	}

	return
}

func IncAddExternalOrderCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	addExternalOrderCounter.With(labels).Inc()
}

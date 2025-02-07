package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	addOrderCounter       *prometheus.CounterVec
	editOrderCounter      *prometheus.CounterVec
	driveeCallbackCounter *prometheus.CounterVec
)

const (
	addOrderCounterName  = "add_order_total"
	editOrderCounterName = "edit_order_total"
	driveeCallbackName   = "drivee_callback_total"
)

func registerOrderMetrics() (err error) {
	addOrderCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      addOrderCounterName,
		Help:      "add order counter",
	}, []string{"status"})
	if err = prometheus.Register(addOrderCounter); err != nil {
		return
	}

	editOrderCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      editOrderCounterName,
		Help:      "edit order counter",
	}, []string{"status"})
	if err = prometheus.Register(editOrderCounter); err != nil {
		return
	}

	driveeCallbackCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeCallbackName,
		Help:      "drivee callback counter",
	}, []string{"drivee_stage", "dostavkee_event"})
	if err = prometheus.Register(driveeCallbackCounter); err != nil {
		return
	}

	return
}

func IncAddOrderCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	addOrderCounter.With(labels).Inc()
}

func IncEditOrderCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	editOrderCounter.With(labels).Inc()
}

func IncDriveeCallbackCounter(driveeStage string, dostavkeeEvent string) {
	labels := prometheus.Labels{
		"drivee_stage":    driveeStage,
		"dostavkee_event": dostavkeeEvent,
	}

	driveeCallbackCounter.With(labels).Inc()
}

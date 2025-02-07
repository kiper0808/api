package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	dgisSuggestsRequestCounter *prometheus.CounterVec
	dgisRoutingRequestCounter  *prometheus.CounterVec
	dgisGeocodeRequestCounter  *prometheus.CounterVec
)

const (
	dgisSuggestsRequestName = "dgis_suggests_request_total"
	dgisRoutingRequestName  = "dgis_routing_request_total"
	dgisGeocodeRequestName  = "dgis_geocode_request_total"
)

func registerDgisMetrics() (err error) {
	dgisSuggestsRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      dgisSuggestsRequestName,
		Help:      "dgis suggests request counter",
	}, []string{"status"})
	if err = prometheus.Register(dgisSuggestsRequestCounter); err != nil {
		return
	}

	dgisRoutingRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      dgisRoutingRequestName,
		Help:      "dgis routing request counter",
	}, []string{"status"})
	if err = prometheus.Register(dgisRoutingRequestCounter); err != nil {
		return
	}

	dgisGeocodeRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      dgisGeocodeRequestName,
		Help:      "dgis geocode request counter",
	}, []string{"status"})
	if err = prometheus.Register(dgisGeocodeRequestCounter); err != nil {
		return
	}

	return
}

func IncDgisSuggestsRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	dgisSuggestsRequestCounter.With(labels).Inc()
}

func IncDgisRoutingRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	dgisRoutingRequestCounter.With(labels).Inc()
}

func IncDgisGeocodeRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	dgisGeocodeRequestCounter.With(labels).Inc()
}

package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	dadataSuggestRequestCounter *prometheus.CounterVec
)

const (
	dadataSuggestRequestName = "dadata_suggest_request_total"
)

func registerDadataMetrics() (err error) {
	dadataSuggestRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      dadataSuggestRequestName,
		Help:      "dadata suggest request counter",
	}, []string{"status"})
	if err = prometheus.Register(dadataSuggestRequestCounter); err != nil {
		return
	}

	return
}

func IncDadataSuggestRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	dadataSuggestRequestCounter.With(labels).Inc()
}

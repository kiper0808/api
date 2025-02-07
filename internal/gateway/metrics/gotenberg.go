package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	gotenbergGenActsLatencyHistogram *prometheus.HistogramVec
)

const (
	gotenbergGenActsLatencyName = "gotenberg_gen_acts_time"
)

func registerGotenbergMetrics() (err error) {
	gotenbergGenActsLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: prometheusNamespace,
		Name:      gotenbergGenActsLatencyName,
		Help:      "gotenberg gen acts time",
	}, []string{"status"})
	if err = prometheus.Register(gotenbergGenActsLatencyHistogram); err != nil {
		return
	}

	return
}

func ObserveGotenbergGenActsLatency(status string, elapsedSeconds float64) {
	labels := prometheus.Labels{
		"status": status,
	}
	gotenbergGenActsLatencyHistogram.With(labels).Observe(elapsedSeconds)
}

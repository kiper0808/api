package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var companyCreateTotal *prometheus.CounterVec

const companyCreateTotalName = "company_create_total"

func registerCompanyMetrics() (err error) {
	companyCreateTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      companyCreateTotalName,
		Help:      "company create total",
	}, []string{"status"})
	if err = prometheus.Register(companyCreateTotal); err != nil {
		return
	}

	return
}

func IncCompanyCreateTotalCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	companyCreateTotal.With(labels).Inc()
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var newPaymentCounter *prometheus.CounterVec

const paymentCounterName = "payment_total"

const (
	PaymentTypeCard    = "card"
	PaymentTypeSbp     = "sbp"
	PaymentTypeInvoice = "invoice"
)

func registerPaymentMetrics() (err error) {
	newPaymentCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      paymentCounterName,
		Help:      "new payment counter",
	}, []string{"type"})
	if err = prometheus.Register(newPaymentCounter); err != nil {
		return
	}

	return
}

func IncNewPaymentCounter(paymentType string) {
	labels := prometheus.Labels{
		"type": paymentType,
	}
	newPaymentCounter.With(labels).Inc()
}

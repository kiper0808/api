package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	paymentCreateInvoiceCounter *prometheus.CounterVec
)

const (
	paymentCreateInvoiceName = "payment_create_invoice_total"
)

func registerInvoiceMetrics() (err error) {
	paymentCreateInvoiceCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      paymentCreateInvoiceName,
		Help:      "payment create invoice counter",
	}, []string{"status"})
	if err = prometheus.Register(paymentCreateInvoiceCounter); err != nil {
		return
	}

	return
}

func IncPaymentCreateInvoiceCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	paymentCreateInvoiceCounter.With(labels).Inc()
}

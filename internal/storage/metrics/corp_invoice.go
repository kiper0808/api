package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	paymentCreateCorpInvoiceCounter *prometheus.CounterVec
)

const (
	paymentCreateCorpInvoiceName = "payment_create_corp_invoice_total"
)

func registerCorpInvoiceMetrics() (err error) {
	paymentCreateCorpInvoiceCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      paymentCreateCorpInvoiceName,
		Help:      "payment create corp invoice counter",
	}, []string{"status"})
	if err = prometheus.Register(paymentCreateCorpInvoiceCounter); err != nil {
		return
	}

	return
}

func IncPaymentCreateCorpInvoiceCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	paymentCreateCorpInvoiceCounter.With(labels).Inc()
}

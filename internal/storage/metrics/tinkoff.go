package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	tinkoffGetStateTotal         *prometheus.CounterVec
	tinkoffGetQrTotal            *prometheus.CounterVec
	tinkoffInitTotal             *prometheus.CounterVec
	tinkoffInvoiceSendTotal      *prometheus.CounterVec
	tinkoffInvoiceInfoTotal      *prometheus.CounterVec
	tinkoffCorpCardBindInitTotal *prometheus.CounterVec
	tinkoffCancelTotal           *prometheus.CounterVec
	tinkoffChargeTotal           *prometheus.CounterVec
)

const (
	tinkoffGetStateName         = "tinkoff_getstate_total"
	tinkoffGetQrName            = "tinkoff_getqr_total"
	tinkoffInitName             = "tinkoff_init_total"
	tinkoffInvoiceSendName      = "tinkoff_invoice_send_total"
	tinkoffInvoiceInfoName      = "tinkoff_invoice_info_total"
	tinkoffCorpCardBindInitName = "tinkoff_corp_card_bind_init_total"
	tinkoffCancelName           = "tinkoff_cancel_total"
	tinkoffChargeName           = "tinkoff_charge_total"
)

func registerTinkoffMetrics() (err error) {
	tinkoffGetStateTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffGetStateName,
		Help:      "tinkoff get state counter",
	}, []string{"status"})
	if err := prometheus.Register(tinkoffGetStateTotal); err != nil {
		return err
	}

	tinkoffGetQrTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffGetQrName,
		Help:      "tinkoff get qr counter",
	}, []string{"status"})
	if err := prometheus.Register(tinkoffGetQrTotal); err != nil {
		return err
	}

	tinkoffInitTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffInitName,
		Help:      "tinkoff init counter",
	}, []string{"status", "type"})
	if err := prometheus.Register(tinkoffInitTotal); err != nil {
		return err
	}

	tinkoffInvoiceSendTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffInvoiceSendName,
		Help:      "tinkoff invoice send counter",
	}, []string{"status"})
	if err := prometheus.Register(tinkoffInvoiceSendTotal); err != nil {
		return err
	}

	tinkoffInvoiceInfoTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffInvoiceInfoName,
		Help:      "tinkoff invoice info counter",
	}, []string{"status"})
	if err := prometheus.Register(tinkoffInvoiceInfoTotal); err != nil {
		return err
	}

	tinkoffCorpCardBindInitTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffCorpCardBindInitName,
		Help:      "tinkoff bind corp card init counter",
	}, []string{"status"})
	if err := prometheus.Register(tinkoffCorpCardBindInitTotal); err != nil {
		return err
	}

	tinkoffCancelTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffCancelName,
		Help:      "tinkoff cancel counter",
	}, []string{"status"})
	if err := prometheus.Register(tinkoffCancelTotal); err != nil {
		return err
	}

	tinkoffChargeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      tinkoffChargeName,
		Help:      "tinkoff charge counter",
	}, []string{"status"})
	if err := prometheus.Register(tinkoffChargeTotal); err != nil {
		return err
	}

	return err
}

func IncTinkoffGetStateCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	tinkoffGetStateTotal.With(labels).Inc()
}

func IncTinkoffGetQrCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	tinkoffGetQrTotal.With(labels).Inc()
}

func IncTinkoffInitCounter(status string, initType string) {
	labels := prometheus.Labels{
		"status": status,
		"type":   initType,
	}
	tinkoffInitTotal.With(labels).Inc()
}

func IncTinkoffInvoiceSendCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	tinkoffInvoiceSendTotal.With(labels).Inc()
}

func IncTinkoffInvoiceInfoCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	tinkoffInvoiceInfoTotal.With(labels).Inc()
}

func IncTinkoffCorpCardBindInitCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	tinkoffCorpCardBindInitTotal.With(labels).Inc()
}

func IncTinkoffCancelCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	tinkoffCancelTotal.With(labels).Inc()
}

func IncTinkoffChargeCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	tinkoffChargeTotal.With(labels).Inc()
}

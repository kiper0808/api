package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	smtpSendCounter *prometheus.CounterVec
)

const (
	smtpSendName = "smtp_send_total"
)

const (
	SMTPSendTypeUserVerification  = "user_verification"
	SMTPSendTypeUserResetPassword = "user_reset_password"
	SMTPSendTypeUserSetupPassword = "user_setup_password"
)

func registerSmtpSendMetrics() (err error) {
	smtpSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      smtpSendName,
		Help:      "smtp send counter",
	}, []string{"status", "type"})
	if err = prometheus.Register(smtpSendCounter); err != nil {
		return
	}

	return
}

func IncSmtpSendCounter(status string, sendType string) {
	labels := prometheus.Labels{
		"status": status,
		"type":   sendType,
	}
	smtpSendCounter.With(labels).Inc()
}

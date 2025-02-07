package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	telegramOrderCancelMessageSendCounter         *prometheus.CounterVec
	telegramOrderAcceptMessageSendCounter         *prometheus.CounterVec
	telegramCompanyRegistrationMessageSendCounter *prometheus.CounterVec
	telegramOrderCreationMessageSendCounter       *prometheus.CounterVec
	telegramCorpOrderCreationMessageSendCounter   *prometheus.CounterVec
	telegramCorpOrderAcceptMessageSendCounter     *prometheus.CounterVec
	telegramCorpOrderCancelMessageSendCounter     *prometheus.CounterVec
)

const (
	telegramOrderCancelMessageSendName         = "telegram_order_cancel_message_send_total"
	telegramOrderAcceptMessageSendName         = "telegram_order_accept_message_send_total"
	telegramCompanyRegistrationMessageSendName = "telegram_company_registration_message_send_total"
	telegramOrderCreationMessageSendName       = "telegram_order_creation_message_send_total"
	telegramCorpOrderCreationMessageSendName   = "telegram_corp_order_creation_message_send_total"
	telegramCorpOrderAcceptMessageSendName     = "telegram_corp_order_accept_message_send_total"
	telegramCorpOrderCancelMessageSendName     = "telegram_corp_order_cancel_message_send_total"
)

const (
	CancelByCourier = "courier"
	CancelByClient  = "client"
)

func registerTelegramMetrics() error {
	telegramOrderCancelMessageSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      telegramOrderCancelMessageSendName,
		Help:      "telegram order cancel message send counter",
	}, []string{"status", "by"})
	if err := prometheus.Register(telegramOrderCancelMessageSendCounter); err != nil {
		return err
	}

	telegramOrderAcceptMessageSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      telegramOrderAcceptMessageSendName,
		Help:      "telegram order accept message send counter",
	}, []string{"status"})
	if err := prometheus.Register(telegramOrderAcceptMessageSendCounter); err != nil {
		return err
	}

	telegramCompanyRegistrationMessageSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      telegramCompanyRegistrationMessageSendName,
		Help:      "telegram company registration message send counter",
	}, []string{"status"})

	telegramOrderCreationMessageSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      telegramOrderCreationMessageSendName,
		Help:      "telegram order creation message send counter",
	}, []string{"status"})
	if err := prometheus.Register(telegramOrderCreationMessageSendCounter); err != nil {
		return err
	}

	telegramCorpOrderCreationMessageSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      telegramCorpOrderCreationMessageSendName,
		Help:      "telegram corp order creation message send counter",
	}, []string{"status"})
	if err := prometheus.Register(telegramCorpOrderCreationMessageSendCounter); err != nil {
		return err
	}

	telegramCorpOrderAcceptMessageSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      telegramCorpOrderAcceptMessageSendName,
		Help:      "telegram corp order accept message send counter",
	}, []string{"status"})
	if err := prometheus.Register(telegramCorpOrderAcceptMessageSendCounter); err != nil {
		return err
	}

	telegramCorpOrderCancelMessageSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      telegramCorpOrderCancelMessageSendName,
		Help:      "telegram corp order cancel message send counter",
	}, []string{"status"})
	if err := prometheus.Register(telegramCorpOrderCancelMessageSendCounter); err != nil {
		return err
	}

	return nil
}

func IncTelegramOrderCancelMessageSendCounter(status string, by string) {
	labels := prometheus.Labels{
		"status": status,
		"by":     by,
	}
	telegramOrderCancelMessageSendCounter.With(labels).Inc()
}

func IncTelegramOrderAcceptMessageSendCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	telegramOrderAcceptMessageSendCounter.With(labels).Inc()
}

func IncTelegramCompanyRegistrationMessageSendCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	telegramCompanyRegistrationMessageSendCounter.With(labels).Inc()
}

func IncTelegramOrderCreationMessageSendCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	telegramOrderCreationMessageSendCounter.With(labels).Inc()
}

func IncTelegramCorpOrderCreationMessageSendCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	telegramCorpOrderCreationMessageSendCounter.With(labels).Inc()
}

func IncTelegramCorpOrderAcceptMessageSendCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	telegramCorpOrderAcceptMessageSendCounter.With(labels).Inc()
}

func IncTelegramCorpOrderCancelMessageSendCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	telegramCorpOrderCancelMessageSendCounter.With(labels).Inc()
}

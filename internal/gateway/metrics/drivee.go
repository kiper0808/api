package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	driveeSetOrderStatusRequestCounter     *prometheus.CounterVec
	driveeAddOrderRequestCounter           *prometheus.CounterVec
	driveeEditOrderRequestCounter          *prometheus.CounterVec
	driveeCreateUserRequestCounter         *prometheus.CounterVec
	driveeUpdateUserRequestCounter         *prometheus.CounterVec
	driveeGetProfileRequestCounter         *prometheus.CounterVec
	driveeCancelOrderRequestCounter        *prometheus.CounterVec
	driveeGetCourierLocationRequestCounter *prometheus.CounterVec
	driveeSendMessageRequestCounter        *prometheus.CounterVec
	driveeCancelOrderMessageRequestCounter *prometheus.CounterVec
)

const (
	driveeSetOrderStatusRequestName     = "drivee_set_order_status_request_total"
	driveeAddOrderRequestName           = "drivee_add_order_request_total"
	driveeEditOrderRequestName          = "drivee_edit_order_request_total"
	driveeCreateUserRequestName         = "drivee_create_user_request_total"
	driveeUpdateUserRequestName         = "drivee_update_user_request_total"
	driveeGetProfileRequestName         = "drivee_get_profile_request_total"
	driveeCancelOrderRequestName        = "drivee_cancel_order_request_total"
	driveeGetCourierLocationRequestName = "drivee_get_courier_location_request_total"
	driveeSendMessageRequestName        = "drivee_send_message_request_total"
	driveeCancelOrderMessageRequestName = "drivee_cancel_order_message_request_total"
)

const (
	DriveeCancelOrderTypeClient = "client"
	DriveeCancelOrderTypeWorker = "worker"
)

func registerDriveeMetrics() error {
	driveeSetOrderStatusRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeSetOrderStatusRequestName,
		Help:      "drivee setOrderStatus request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeSetOrderStatusRequestCounter); err != nil {
		return err
	}

	driveeAddOrderRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeAddOrderRequestName,
		Help:      "drivee addOrder request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeAddOrderRequestCounter); err != nil {
		return err
	}

	driveeEditOrderRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeEditOrderRequestName,
		Help:      "drivee edit order request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeEditOrderRequestCounter); err != nil {
		return err
	}

	driveeCreateUserRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeCreateUserRequestName,
		Help:      "drivee createUser request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeCreateUserRequestCounter); err != nil {
		return err
	}

	driveeUpdateUserRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeUpdateUserRequestName,
		Help:      "drivee updateUser request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeUpdateUserRequestCounter); err != nil {
		return err
	}

	driveeGetProfileRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeGetProfileRequestName,
		Help:      "drivee getProfile request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeGetProfileRequestCounter); err != nil {
		return err
	}

	driveeCancelOrderRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeCancelOrderRequestName,
		Help:      "drivee cancelOrder request counter",
	}, []string{"status", "type"})
	if err := prometheus.Register(driveeCancelOrderRequestCounter); err != nil {
		return err
	}

	driveeGetCourierLocationRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeGetCourierLocationRequestName,
		Help:      "drivee getCourierLocation request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeGetCourierLocationRequestCounter); err != nil {
		return err
	}

	driveeSendMessageRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeSendMessageRequestName,
		Help:      "drivee sendMessage request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeSendMessageRequestCounter); err != nil {
		return err
	}

	driveeCancelOrderMessageRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      driveeCancelOrderMessageRequestName,
		Help:      "drivee cancelOrderMessage request counter",
	}, []string{"status"})
	if err := prometheus.Register(driveeCancelOrderMessageRequestCounter); err != nil {
		return err
	}

	return nil
}

const (
	DriveeUserStatusDuplicate = "duplicate"
)

func IncDriveeSetOrderStatusRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeSetOrderStatusRequestCounter.With(labels).Inc()
}

func IncDriveeAddOrderRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeAddOrderRequestCounter.With(labels).Inc()
}

func IncDriveeEditOrderRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeEditOrderRequestCounter.With(labels).Inc()
}

func IncDriveeCreateUserRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeCreateUserRequestCounter.With(labels).Inc()
}

func IncDriveeUpdateUserRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeUpdateUserRequestCounter.With(labels).Inc()
}

func IncDriveeGetProfileRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeGetProfileRequestCounter.With(labels).Inc()
}

func IncDriveeCancelOrderRequestCounter(status string, cancelType string) {
	labels := prometheus.Labels{
		"status": status,
		"type":   cancelType,
	}
	driveeCancelOrderRequestCounter.With(labels).Inc()
}

func IncDriveeGetCourierLocationRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeGetCourierLocationRequestCounter.With(labels).Inc()
}

func IncDriveeSendMessageRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeSendMessageRequestCounter.With(labels).Inc()
}

func IncDriveeCancelOrderMessageRequestCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	driveeCancelOrderMessageRequestCounter.With(labels).Inc()
}

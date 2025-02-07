package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	fileStorageUploadCounter *prometheus.CounterVec
)

const (
	fileStorageUploadName = "file_storage_upload_total"
)

func registerFileStorageMetrics() (err error) {
	fileStorageUploadCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      fileStorageUploadName,
		Help:      "file storage upload counter",
	}, []string{"status"})
	if err = prometheus.Register(fileStorageUploadCounter); err != nil {
		return
	}

	return
}

func IncFileStorageUploadCounter(status string) {
	labels := prometheus.Labels{
		"status": status,
	}
	fileStorageUploadCounter.With(labels).Inc()
}

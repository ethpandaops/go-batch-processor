package processor

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// DefaultMetrics is the default metrics instance using "batch" namespace.
var DefaultMetrics = NewMetrics("batch")

// Metrics holds Prometheus metrics for the batch processor.
type Metrics struct {
	itemsQueued            *prometheus.GaugeVec
	itemsDropped           *prometheus.CounterVec
	itemsFailed            *prometheus.CounterVec
	itemsExported          *prometheus.CounterVec
	exportDuration         *prometheus.HistogramVec
	batchSize              *prometheus.HistogramVec
	workerCount            *prometheus.GaugeVec
	workerExportInProgress *prometheus.GaugeVec
}

// NewMetrics creates a new Metrics instance with the given namespace.
func NewMetrics(namespace string) *Metrics {
	if namespace != "" {
		namespace += "_"
	}

	namespace += "processor"

	m := &Metrics{
		itemsQueued: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "items_queued",
			Namespace: namespace,
			Help:      "Number of items queued",
		}, []string{"processor"}),
		itemsDropped: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "items_dropped_total",
			Namespace: namespace,
			Help:      "Number of items dropped",
		}, []string{"processor"}),
		itemsFailed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "items_failed_total",
			Namespace: namespace,
			Help:      "Number of items failed",
		}, []string{"processor"}),
		itemsExported: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "items_exported_total",
			Namespace: namespace,
			Help:      "Number of items exported",
		}, []string{"processor"}),
		exportDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "export_duration_seconds",
			Namespace: namespace,
			Help:      "Duration of export operations in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10),
		}, []string{"processor"}),
		batchSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "batch_size",
			Namespace: namespace,
			Help:      "Size of processed batches",
			Buckets:   prometheus.ExponentialBucketsRange(1, 50000, 10),
		}, []string{"processor"}),
		workerCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "worker_count",
			Namespace: namespace,
			Help:      "Number of active workers",
		}, []string{"processor"}),
		workerExportInProgress: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "worker_export_in_progress",
			Namespace: namespace,
			Help:      "Number of workers currently exporting",
		}, []string{"processor"}),
	}

	prometheus.MustRegister(m.itemsQueued)
	prometheus.MustRegister(m.itemsDropped)
	prometheus.MustRegister(m.itemsFailed)
	prometheus.MustRegister(m.itemsExported)
	prometheus.MustRegister(m.exportDuration)
	prometheus.MustRegister(m.batchSize)
	prometheus.MustRegister(m.workerCount)
	prometheus.MustRegister(m.workerExportInProgress)

	return m
}

// SetItemsQueued sets the number of items queued for the given processor.
func (m *Metrics) SetItemsQueued(name string, count float64) {
	m.itemsQueued.WithLabelValues(name).Set(count)
}

// IncItemsDroppedBy increments the number of items dropped by the given count.
func (m *Metrics) IncItemsDroppedBy(name string, count float64) {
	m.itemsDropped.WithLabelValues(name).Add(count)
}

// IncItemsExportedBy increments the number of items exported by the given count.
func (m *Metrics) IncItemsExportedBy(name string, count float64) {
	m.itemsExported.WithLabelValues(name).Add(count)
}

// IncItemsFailedBy increments the number of items failed by the given count.
func (m *Metrics) IncItemsFailedBy(name string, count float64) {
	m.itemsFailed.WithLabelValues(name).Add(count)
}

// ObserveExportDuration records the duration of an export operation.
func (m *Metrics) ObserveExportDuration(name string, duration time.Duration) {
	m.exportDuration.WithLabelValues(name).Observe(duration.Seconds())
}

// ObserveBatchSize records the size of a processed batch.
func (m *Metrics) ObserveBatchSize(name string, size float64) {
	m.batchSize.WithLabelValues(name).Observe(size)
}

// SetWorkerCount sets the number of active workers for the given processor.
func (m *Metrics) SetWorkerCount(name string, count float64) {
	m.workerCount.WithLabelValues(name).Set(count)
}

// IncWorkerExportInProgress increments the number of workers currently exporting.
func (m *Metrics) IncWorkerExportInProgress(name string) {
	m.workerExportInProgress.WithLabelValues(name).Inc()
}

// DecWorkerExportInProgress decrements the number of workers currently exporting.
func (m *Metrics) DecWorkerExportInProgress(name string) {
	m.workerExportInProgress.WithLabelValues(name).Dec()
}

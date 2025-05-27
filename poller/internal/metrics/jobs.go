package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type jobMetrics struct {
	jobScheduled *prometheus.CounterVec
	jobErrorInc  *prometheus.CounterVec
}

func (m *schedulerMetricsHandler) registerJobMetrics() {
	m.job.jobScheduled = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Name:      "offer_scheduler_job_total",
			Help:      "The total number of scheduled jobs",
		},
		[]string{"host", "shard", "job_name"},
	)

	m.job.jobErrorInc = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Name:      "offer_scheduler_job_errors_total",
			Help:      "The total number of failed jobs",
		},
		[]string{"host", "shard", "job_name", "reason"},
	)
}

// ReportScheduleOfJob reports a scheduled job
func (m *schedulerMetricsHandler) ReportScheduleOfJob(jobName string) {
	labels := m.baseLabelsWithValues(jobName)
	m.job.jobScheduled.WithLabelValues(labels...).Inc()
}

// JobErrorInc increases error counter when job fails
func (m *schedulerMetricsHandler) JobErrorInc(jobName string, errMsg string) {
	labels := m.baseLabelsWithValues(jobName, errMsg)
	m.job.jobErrorInc.WithLabelValues(labels...).Inc()
}

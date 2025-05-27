package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type dbMetrics struct {
	dbCall     *prometheus.CounterVec
	dbErrorInc *prometheus.CounterVec
}

func (m *metricsHandler) registerDbMetrics() {
	m.dbMetrics.dbCall = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Name:      "queries_total",
			Help:      "A counter for query",
		},
		[]string{"host", "database", "query"},
	)

	m.dbMetrics.dbErrorInc = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Name:      "db_error_counter",
			Help:      "How many unexpected database error",
		},
		[]string{"method", "reason"},
	)
}

// DBCall increases the number of db calls
func (m *metricsHandler) DBCall(source string) {
	labels := m.baseLabelsWithValues(source)
	m.dbMetrics.dbCall.WithLabelValues(labels...).Inc()
}

// DBErrorInc increases error counter when query fails
func (m *metricsHandler) DBErrorInc(source string, errMsg string) {
	labels := m.baseLabelsWithValues(source, errMsg)
	m.dbMetrics.dbErrorInc.WithLabelValues(labels...).Inc()
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type natsConsumer struct {
	natsErrorInc *prometheus.CounterVec
}

func (m *metricsHandler) registerNatsMetrics() {

	m.natsConsumer.natsErrorInc = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "nats_error_counter",
			Help:      "How many NATS errors received by app.",
		},
		[]string{"host", "shard", "reason"},
	)
}

// NatsErrorInc increases error counter when query fails
func (m *metricsHandler) NatsErrorInc(errMsg string) {
	labels := m.baseLabelsWithValues(errMsg)
	m.natsConsumer.natsErrorInc.WithLabelValues(labels...).Inc()
}

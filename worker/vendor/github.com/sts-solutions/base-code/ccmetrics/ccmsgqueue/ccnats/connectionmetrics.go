package ccnats

import (
	"os"

	"github.com/sts-solutions/base-code/ccmetrics"
	"github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type connectionMetrics struct {
	host          ccmetrics.Host
	shard         string
	errorsCounter *prometheus.CounterVec
}

func NewConnectionMetrics(nameSpace string) ccmsgqueue.ConnectionMetrics {
	hostName, _ := os.Hostname()
	host := ccmetrics.Host(hostName)
	ns := nameSpace

	errorsCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "connection_errors_counter",
		Help:      "The number of NATS connection errors",
	}, []string{"host", "shard", "type"})

	return &connectionMetrics{
		host:          host,
		shard:         host.Shard(),
		errorsCounter: errorsCounter,
	}
}

func (nm *connectionMetrics) Host() string {
	return nm.host.String()
}

func (nm *connectionMetrics) ErrorInc(errorType string) {
	nm.errorsCounter.WithLabelValues(nm.Host(), nm.shard, errorType).Inc()
}

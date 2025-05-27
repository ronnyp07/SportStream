package ccnats

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sts-solutions/base-code/ccmetrics"
	"github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue"
)

type consumerMetrics struct {
	host               ccmetrics.Host
	shard              string
	concurrencyProcess *prometheus.GaugeVec
	messageDuration    *prometheus.HistogramVec
	messagesCounter    *prometheus.CounterVec
	errorCounter       *prometheus.CounterVec
}

func NewConsumerMetrics(nameSpace string) ccmsgqueue.ConsumerMetrics {
	hostName, _ := os.Hostname()
	host := ccmetrics.Host(hostName)
	ns := nameSpace

	concurrencyProcess := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "available_concurrent_processes",
		Help:      "The number of processes currently available for a given consumer.",
	}, []string{"host", "shard", "subject", "name"})

	messagesDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "consumed_messages_lag_seconds",
		Help:      "The duration it took for the message to be picked up by consumer in seconds",
	}, []string{"host", "shard", "subject", "name"})

	messageCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "consumed_messages_counter",
		Help:      "A counter for consumed messages",
	}, []string{"host", "shard", "subject", "name"})

	errorsCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "consumer_error_counter",
		Help:      "The number of Nats errors",
	}, []string{"host", "shard", "type"})

	return &consumerMetrics{
		host:               host,
		shard:              host.Shard(),
		concurrencyProcess: concurrencyProcess,
		messageDuration:    messagesDuration,
		messagesCounter:    messageCounter,
		errorCounter:       errorsCounter,
	}
}

func (nm *consumerMetrics) Host() string {
	return nm.host.String()
}

func (nm *consumerMetrics) ConcurrentProcessInc(quantity int, subject, consumerName string) {
	nm.concurrencyProcess.WithLabelValues(nm.Host(), nm.shard, subject, consumerName).Add(float64(quantity))
}

func (nm *consumerMetrics) MessageLag(duration float64, subject, consumerName string) {
	nm.messageDuration.WithLabelValues(nm.Host(), nm.shard, subject, consumerName).Observe(duration)
}

func (nm *consumerMetrics) MessageCounterInc(subject, consumerName string) {
	nm.messagesCounter.WithLabelValues(nm.Host(), nm.shard, subject, consumerName).Inc()
}

func (nm *consumerMetrics) ErrorInc(errorType string) {
	nm.messagesCounter.WithLabelValues(nm.Host(), nm.shard, errorType).Inc()
}

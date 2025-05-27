package ccnats

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sts-solutions/base-code/ccmetrics"
	"github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue"
)

type publisherMetrics struct {
	host            ccmetrics.Host
	shard           string
	messageDuration *prometheus.HistogramVec
	messagesCounter *prometheus.CounterVec
	errorCounter    *prometheus.CounterVec
}

func NewPublisherMetrics(nameSpace string) ccmsgqueue.PublisherMetrics {
	hostName, _ := os.Hostname()
	host := ccmetrics.Host(hostName)
	ns := nameSpace

	messagesDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "published_messages_lag_seconds",
		Help:      "The duration it took for the message to be published up by publisher in seconds",
	}, []string{"host", "shard", "subject", "name"})

	messageCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "published_messages_counter",
		Help:      "A counter for published messages.",
	}, []string{"host", "shard", "subject"})

	errorsCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: SUBSYSTEM,
		Name:      "consumer_error_counter",
		Help:      "The number of Nats errors",
	}, []string{"host", "shard", "type"})

	return &publisherMetrics{
		host:            host,
		shard:           host.Shard(),
		messageDuration: messagesDuration,
		messagesCounter: messageCounter,
		errorCounter:    errorsCounter,
	}
}

func (nm *publisherMetrics) Host() string {
	return nm.host.String()
}

func (nm *publisherMetrics) MessageLag(duration float64, subject string) {
	nm.messageDuration.WithLabelValues(nm.Host(), nm.shard, subject).Observe(duration)
}

func (nm *publisherMetrics) MessageCounterInc(subject string) {
	nm.messagesCounter.WithLabelValues(nm.Host(), nm.shard, subject).Inc()
}

func (nm *publisherMetrics) ErrorInc(errorType string) {
	nm.messagesCounter.WithLabelValues(nm.Host(), nm.shard, errorType).Inc()
}

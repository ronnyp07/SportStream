package promnats

import (
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	subsystemNatsOutgoing = "nats_outgoing"
)

var (
	duration    *prometheus.HistogramVec
	queries     *prometheus.CounterVec
	initialized bool = false
)

// Init initializes the NATS Prometheus metrics
func Init(namespace string) {
	if initialized {
		return
	}

	duration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystemNatsOutgoing,
			Name:      "publish_duration_histogram_seconds",
			Help:      "A histogram of NATS message publish latencies in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"subject", "success"},
	)

	queries = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystemNatsOutgoing,
			Name:      "messages_total",
			Help:      "Total count of NATS messages published.",
		},
		[]string{"subject", "success"},
	)

	initialized = true
}

// Publish publishes a message to NATS JetStream and records metrics
func Publish(jetstream nats.JetStreamContext, subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	if !initialized {
		return jetstream.Publish(subj, data, opts...)
	}

	start := time.Now()
	ack, err := jetstream.Publish(subj, data, opts...)

	labelValues := []string{subj, strconv.FormatBool(err == nil)}
	duration.WithLabelValues(labelValues...).Observe(time.Since(start).Seconds())
	queries.WithLabelValues(labelValues...).Inc()

	return ack, err
}

// PublishAsync publishes a message asynchronously to NATS JetStream and records metrics
func PublishAsync(jetstream nats.JetStreamContext, subj string, data []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	if !initialized {
		return jetstream.PublishAsync(subj, data, opts...)
	}

	start := time.Now()
	ack, err := jetstream.PublishAsync(subj, data, opts...)

	labelValues := []string{subj, strconv.FormatBool(err == nil)}
	duration.WithLabelValues(labelValues...).Observe(time.Since(start).Seconds())
	queries.WithLabelValues(labelValues...).Inc()

	return ack, err
}

// PublishMsg publishes a NATS message to JetStream and records metrics
func PublishMsg(jetstream nats.JetStreamContext, m *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error) {
	if !initialized {
		return jetstream.PublishMsg(m, opts...)
	}

	start := time.Now()
	ack, err := jetstream.PublishMsg(m, opts...)

	labelValues := []string{m.Subject, strconv.FormatBool(err == nil)}
	duration.WithLabelValues(labelValues...).Observe(time.Since(start).Seconds())
	queries.WithLabelValues(labelValues...).Inc()

	return ack, err
}

// PublishMsgAsync publishes a NATS message asynchronously to JetStream and records metrics
func PublishMsgAsync(jetstream nats.JetStreamContext, m *nats.Msg, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	if !initialized {
		return jetstream.PublishMsgAsync(m, opts...)
	}

	start := time.Now()
	ack, err := jetstream.PublishMsgAsync(m, opts...)

	labelValues := []string{m.Subject, strconv.FormatBool(err == nil)}
	duration.WithLabelValues(labelValues...).Observe(time.Since(start).Seconds())
	queries.WithLabelValues(labelValues...).Inc()

	return ack, err
}

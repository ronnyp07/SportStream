package ccmetrics

import (
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type natsMsgQueueMetrics struct {
	host                Host
	shard               string
	natsMessageLag      *prometheus.HistogramVec
	concurrentProcessor *prometheus.GaugeVec
}

func NewNatsMsgQueueMetrics(nameSpace string) MegQueueMetrics {
	hostName, _ := os.Hostname()
	host := Host(hostName)
	ns := nameSpace
	subsystem := "nats"
	histogramBuckets := []float64{.005, 0.1, 0.2, 0.25, 0.5, .1, .25, .5,
		1, 2.5, 5, 10, 15, 20, 25, 30}

	concurrentProcessor := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ns,
		Subsystem: subsystem,
		Name:      "available_concurrent_processor",
		Help: "The number of processes currently available for a " +
			"given consumer",
	}, []string{"consumername"})

	natsMessageLag := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: subsystem,
		Name:      "nats_message_lag_seconds",
		Help: "The diraction it took for the massage to be picked" +
			"up by consumer in seconds.",
		Buckets: histogramBuckets,
	}, []string{})

	return &natsMsgQueueMetrics{
		host:                host,
		shard:               host.Shard(),
		natsMessageLag:      natsMessageLag,
		concurrentProcessor: concurrentProcessor,
	}
}

func (n *natsMsgQueueMetrics) Host() string {
	return n.host.String()
}

func (n *natsMsgQueueMetrics) MessageLag(start time.Time, labelValues []string) {
	lag := time.Now().UTC().Sub(start)
	n.natsMessageLag.WithLabelValues().Observe(lag.Seconds())
}

func (n *natsMsgQueueMetrics) ConcurrentProcessorInc(counter int, name string) {
	n.concurrentProcessor.WithLabelValues(strings.Replace(name, "-", "_", -1)).Set(float64(counter))
}

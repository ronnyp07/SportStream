package ccmetrics

import "time"

type MegQueueMetrics interface {
	Host() string
	MessageLag(start time.Time, labelValues []string)
	ConcurrentProcessorInc(counter int, name string)
}

package ccmetrics

import "time"

type MetricsNats interface {
	Host() string
	MessageLag(start time.Time, labelValues []string)
	ConcurrentProcessorInc(counter int, name string)
}

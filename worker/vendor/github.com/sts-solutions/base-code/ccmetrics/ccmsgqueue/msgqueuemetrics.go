package ccmsgqueue

type HostProvider interface {
	Host() string
}

type ConnectionMetrics interface {
	HostProvider
	ErrorInc(errorType string)
}

type ConsumerMetrics interface {
	HostProvider
	ConcurrentProcessInc(counter int, subject, consumerName string)
	MessageLag(duration float64, subject, consumerName string)
	MessageCounterInc(subject, consumerName string)
	ErrorInc(errorType string)
}

type PublisherMetrics interface {
	MessageLag(duration float64, subject string)
	MessageCounterInc(subject string)
	ErrorInc(errorType string)
}

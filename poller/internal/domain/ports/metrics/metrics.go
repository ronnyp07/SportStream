package metrics

type AppInfoMetricsHandler interface {
	RegisterVersionInfo()
}

type MessageQueueMetricsHandler interface {
	ReportError(err error)
}

type HTTPMetricsHandler interface {
	RegisterError(method, endpoint, reason string)
}

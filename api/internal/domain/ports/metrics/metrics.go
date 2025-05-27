package metrics

type AppInfoMetricsHandler interface {
	RegisterVersionInfo()
}

type MessageQueueMetricsHandler interface {
	ReportError(err error)
}

type MetricsHandler interface {
	RegisterMetrics()
	DBCall(source string)
	DBErrorInc(source string, errMsg string)
}

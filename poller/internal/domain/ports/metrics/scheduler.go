package metrics

import "time"

type SchedulerMetricsHandler interface {
	RegisterMetrics()
	JobErrorInc(string, string)
	SchedulerOutgoingHttpRequest(startTime time.Time, method string,
		destination string, requestCode int, platformName string)
	SchedulerHTTPClientCall(startTime time.Time, method string,
		destination string, requestCode int, platformName string)
	ReportScheduleOfJob(string)
}

package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type httpClient struct {
	callDuration       *prometheus.HistogramVec
	callRequestCounter *prometheus.CounterVec
	errorCounter       *prometheus.CounterVec
}

// registerHTTPClientMetrics sets up the different metrics used for monitoring httpClient calls.
func (m *schedulerMetricsHandler) registerHTTPClientMetrics() {
	m.httpClient.callDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: m.namespace,
			Name:      "http_client_duration_seconds",
			Help:      "External http calls latencies in seconds",
			Buckets:   m.histogramBuckets,
		},
		[]string{"host", "shard", "method", "destination", "response_code", "platform_name"},
	)

	m.httpClient.callRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Subsystem: "http",
			Name:      "outgoing_requests_total",
			Help:      "Total count of outgoing HTTP requests",
		},
		[]string{"host", "shard", "method", "destination", "response_code", "platform_name"},
	)

	m.httpClient.errorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Name:      "http_client_errors_total",
			Help:      "Total count of errors from external API calls",
		},
		[]string{"host", "shard", "method", "destination", "error_code", "platform_name"},
	)
}

// elapsed calculates duration since start time in seconds
func elapsed(start time.Time) float64 {
	return time.Since(start).Seconds()
}

// SchedulerHTTPClientCall records external request duration and counts successful requests
func (m *schedulerMetricsHandler) SchedulerHTTPClientCall(startTime time.Time, method string,
	destination string, responseCode int, platformName string) {
	if responseCode == 0 {
		responseCode = http.StatusInternalServerError
	}
	labelValues := m.baseLabelsWithValues(
		method,
		destination,
		strconv.Itoa(responseCode),
		platformName,
	)
	m.httpClient.callDuration.WithLabelValues(labelValues...).Observe(elapsed(startTime))
	m.httpClient.callRequestCounter.WithLabelValues(labelValues...).Inc()
}

// SchedulerOutgoingHttpRequest counts outgoing HTTP requests
func (m *schedulerMetricsHandler) SchedulerOutgoingHttpRequest(startTime time.Time, method string,
	destination string, responseCode int, platformName string) {
	if responseCode == 0 {
		responseCode = http.StatusInternalServerError
	}
	labelValues := m.baseLabelsWithValues(
		method,
		destination,
		strconv.Itoa(responseCode),
		platformName,
	)
	m.httpClient.callRequestCounter.WithLabelValues(labelValues...).Inc()
}

// SchedulerHTTPClientErrorInc increases error counter when calling external APIs
func (m *schedulerMetricsHandler) SchedulerHTTPClientErrorInc(httpErrorCode int, method string,
	destination string, platformName string) {
	labelValues := m.baseLabelsWithValues(
		method,
		destination,
		strconv.Itoa(httpErrorCode),
		platformName,
	)
	m.httpClient.errorCounter.WithLabelValues(labelValues...).Inc()
}

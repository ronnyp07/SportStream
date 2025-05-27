package metrics

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/ronnyp07/SportStream/internal/domain/ports/metrics"
)

const (
	twoLayerSkipper   = 2
	threeLayerSkipper = 3
)

type schedulerMetricsHandler struct {
	host             string
	shard            string
	namespace        string
	environment      string
	histogramBuckets []float64
	httpClient       *httpClient
	job              *jobMetrics
}

func NewSchedulerMetricsHandler() metrics.SchedulerMetricsHandler {
	host, _ := os.Hostname()
	m := &schedulerMetricsHandler{
		namespace:        "pooller",
		histogramBuckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 15, 20, 25, 30},
		host:             host,
		shard:            shardFromHostname(host),
		environment:      environmentFromHostname(host),
		httpClient:       &httpClient{},
		job:              &jobMetrics{},
	}

	return m
}

func (m *schedulerMetricsHandler) Host() string {
	return m.host
}

// RegisterMetrics sets up all the metrics used in the application.
func (m *schedulerMetricsHandler) RegisterMetrics() {
	m.registerHTTPClientMetrics()
	m.registerJobMetrics()
}

// baseLabels returns an array with the host and shard
func (m *schedulerMetricsHandler) baseLabels() []string {
	return []string{m.host, m.shard}
}

// baseLabelsWithValues returns an array with the provided values prepending host and shard
func (m *schedulerMetricsHandler) baseLabelsWithValues(values ...string) []string {
	baseLabels := []string{m.host, m.shard}
	return append(baseLabels, values...)
}

// baseLabelsWithLabelsFromError returns an array with the provided values prepending host and shard
func (m *schedulerMetricsHandler) baseLabelsWithLabelsFromError(err error) []string {
	baseLabels := []string{m.host, m.shard}
	lbError := labelsFromError(err, twoLayerSkipper)
	return append(baseLabels, lbError...)
}

// shardFromHostname extracts short instance identifier from hostname e.g. '01' from `coreac01`
func shardFromHostname(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) >= 1 {
		s := regexp.MustCompile(`[0-9]+`).FindString(parts[0])
		if s != "" {
			return s
		}
	}
	return host
}

// environmentFromHostname extracts environment from hostname
func environmentFromHostname(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}

// labelsFromError extracts relevant labels from error
func labelsFromError(err error, skip int) []string {
	if err == nil {
		return []string{"", ""}
	}

	var stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if errors.As(err, &stackTracer) {
		frame := errors.Frame(stackTracer.StackTrace()[skip])
		pc := uintptr(frame) - 1
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			file, line := fn.FileLine(pc)
			return []string{
				path.Base(fn.Name()),
				fmt.Sprintf("%s:%d", path.Base(file), line),
			}
		}
	}

	return []string{"", ""}
}

// Additional methods would be implemented here for:
// - registerAPIMetrics()
// - registerHTTPClientMetrics()
// - registerJobMetrics()
// - The actual metric reporting methods (SchedulerOutgoingHttpRequest, JobErrorInc, etc.)

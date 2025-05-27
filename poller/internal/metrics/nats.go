package metrics

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"emperror.dev/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Host, _ = os.Hostname()
)

var (
	// VersionInfo = prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Namespace: MetricsNamespace,
	// 		Name:      "version_info",
	// 		Help:      "Version information.",
	// 	},
	// 	[]string{"host", "version", "go"},
	// )

	NatsErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "nats_error_counter",
			Help:      "How many NATS errors received by app.",
		},
		[]string{"host", "function", "file", "line"},
	)

	RequestValidationProcessDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  MetricsNamespace,
			Name:       "request_validation_process_duration_seconds",
			Help:       "The Request Validation step process latencies in seconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "response"},
	)

	ErrorRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "error_request_counter",
			Help:      "How many error requests received by app.",
		},
		[]string{"method", "endpoint", "reason"},
	)

	DataBaseErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "error_counter",
			Help:      "How many error requests received in the service.",
		},
		[]string{"method", "reason"},
	)
)

func init() {
	prometheus.MustRegister(
		//VersionInfo,
		NatsErrorCounter,
		RequestValidationProcessDuration,
		ErrorRequestCounter,
		DataBaseErrorCounter,
	)
}

var initializedCounters sync.Map

const (
	sep                     = "_"
	timeAfterInitialization = 15
)

// ReportNatsError increments the NATS error counter with appropriate labels
func ReportNatsError(err error) {
	labels := LabelsFromError(err)
	initializeCounterIfNeeded(NatsErrorCounter, labels)
	NatsErrorCounter.WithLabelValues(labels...).Inc()
}

// isCounterInitialized checks if a counter with given labels was already initialized
func isCounterInitialized(labels []string) bool {
	key := strings.Join(labels, sep)
	_, ok := initializedCounters.Load(key)
	return ok
}

// LabelsFromCaller returns hostname and caller information (function, file, line)
func LabelsFromCaller(skip int) []string {
	pc, file, line, _ := runtime.Caller(skip)
	functionName := runtime.FuncForPC(pc).Name()
	_, function := path.Split(functionName)
	_, fileName := path.Split(file)

	return []string{
		Host,
		function,
		fileName,
		fmt.Sprintf("%d", line),
	}
}

// LabelsFromError extracts labels from an error including stack trace information
func LabelsFromError(err error) []string {
	if s := stack(err); len(s) > 0 {
		var function string
		if v := strings.Split(fmt.Sprintf("%+s", s[0]), "\n\t"); len(v) == 2 {
			function = path.Base(v[0])
		} else {
			function = fmt.Sprintf("%s", s[0])
		}

		return []string{
			Host,
			function,
			fmt.Sprint("%s", s[0]),
			fmt.Sprint("%s", s[0]),
		}
	}

	return LabelsFromCaller(3)
}

// initializeCounterIfNeeded ensures a counter is initialized only once
func initializeCounterIfNeeded(counter *prometheus.CounterVec, labels []string) {
	if isCounterInitialized(labels) {
		return
	}
	key := strings.Join(labels, sep)
	initializedCounters.Store(key, true)
	counter.WithLabelValues(labels...).Add(0)
	time.Sleep(timeAfterInitialization * time.Second)
}

func stack(err error) (stackTrace errors.StackTrace) {
	for {
		stackErr, ok := err.(interface {
			StackTrace() errors.StackTrace
		})

		if ok {
			stackTrace = stackErr.StackTrace()
		}

		u, ok := err.(interface {
			Unwrap() error
		})

		if !ok {
			break
		}

		err = u.Unwrap()
	}

	return stackTrace
}

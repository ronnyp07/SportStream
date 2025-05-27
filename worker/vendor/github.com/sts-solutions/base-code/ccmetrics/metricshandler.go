package ccmetrics

import (
	"emperror.dev/errors"
	"github.com/sts-solutions/base-code/ccerrors/withinnererr"
	"time"
)

// MetricsHandler is the interface that implements the functions to register metrics
type MetricsHandler interface {
	Host() string
	RequestDuration(start time.Time, labelValues []string)
	VersionInfoInc(t string)
	UnexpectedErrorInc(string)
	RegisterMetrics()
	PanicInc(tracker string)
	DBCall(start time.Time, source string)
	DBErrorInc(source string)
	HTTPClientCall(startTime time.Time, destination string, reponseCode int)
	HTTPClientErrorInc(httpErrorCode int, destination string)
}

type Traceable interface {
	RegisterMetric(metricsHandler MetricsHandler)
}

func Register(mHandler MetricsHandler, obj any) {
	registed := false
	for {
		if traceable, ok := obj.(Traceable); ok {
			traceable.RegisterMetric(mHandler)
			registed = true
		}
		if withInnerErr, ok := obj.(withinnererr.WithInnerErr); ok {
			obj = errors.Cause(withInnerErr.GetInnerErr())
			continue
		}
		break
	}

	if !registed {
		mHandler.UnexpectedErrorInc("InternalServerError")
	}
}

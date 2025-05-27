package log

import (
	"emperror.dev/errors"

	"github.com/sts-solutions/base-code/cclogger"
)

var lgr cclogger.Logger

func Logger() cclogger.Logger {
	return lgr
}

func SetupLogger(appName string) error {
	if lgr != nil {
		return errors.New("logger already set")
	}

	logger, err := cclogger.NewBuilder().
		WithLevel(cclogger.Info).
		WithAppName(appName).
		WithCorrelationID().
		Build()

	if err != nil {
		return errors.Wrap(err, "creating logger")
	}

	lgr = logger

	return nil
}
